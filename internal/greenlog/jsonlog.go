package greenlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo  Level = iota // Has the value 0
	LevelError              // Has the value 1
	LevelFatal              // Has the value 2
	LevelOff                // Has the value 3
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Greenlog struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Greenlog {
	return &Greenlog{
		out:      out,
		minLevel: minLevel,
	}
}

// Helper methods

func (gl *Greenlog) PrintInfo(message string, properties map[string]string) {
	gl.print(LevelInfo, message, properties)
}

func (gl *Greenlog) PrintError(err error, properties map[string]string) {
	gl.print(LevelError, err.Error(), properties)
}

func (gl *Greenlog) PrintFatal(err error, properties map[string]string) {
	gl.print(LevelFatal, err.Error(), properties)
	os.Exit(1) // For entries at the FATAL level, we also terminate the application
}

func (glog *Greenlog) print(level Level, message string, properties map[string]string) (int, error) {
	if level < glog.minLevel {
		return 0, nil
	}
	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}
	// Include a stack trace for entries at the ERROR and FATAL levels.
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}
	var line []byte
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message" + err.Error())
	}
	glog.mu.Lock()
	defer glog.mu.Unlock()

	return glog.out.Write(append(line, '\n'))
}

func (glog *Greenlog) Write(message []byte) (n int, err error) {
	return glog.print(LevelError, string(message), nil)
}
