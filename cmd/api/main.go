package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	"github.com/jahidhimon/greenlight.git/internal/data"
	"github.com/jahidhimon/greenlight.git/internal/greenlog"
	_ "github.com/lib/pq"
)

// Declare a sting containing the application version number
const version = "0.0.1"

// TODO: Define more configs
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

// application struct to hold the dependencies for our
// HTTP handlers, helpers and middleware
type application struct {
	config config
	logger *greenlog.Greenlog
	models data.Models
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development",
		"Environment (development/staging/production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"),
		"PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25,
		"PostgresQL max open connection")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25,
		"PostgreSQL max idle connection")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m",
		"PostgreSQL max connection idle time")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2,
		"Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4,
		"Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()

	logger := greenlog.New(os.Stdout, greenlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	// create a context witha 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use pingcontext to establish a new connection to the databse, passing
	// in the context we created above as a parameter. If the connection
	// couldn't be established successfully within 5 second deadline, then this
	// will return an error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
