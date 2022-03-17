package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jahidhimon/greenlight.git/internal/validator"
	"github.com/lib/pq"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genre,omitempty"`
	Version    int32     `json:"version,omitempty"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes")
	
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "mus tbe a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicates")
}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(mv *Movie) error {
	query := `
Insert INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version`
	
	args := []interface{}{mv.Title, mv.Year, mv.Runtime, pq.Array(mv.Genres)}

	return m.DB.QueryRow(query, args...).Scan(&mv.ID, &mv.CreatedAt, &mv.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	// The psql bigserial data type that wer're using to store id
	// is auto incrementing from 1 by default, less than 1 is not possible
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE id = $1`
	var movie Movie
	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	return nil
}

func (m MovieModel) Delete(id int64) error {
	return nil
}
