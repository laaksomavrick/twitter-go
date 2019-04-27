package logger

import (
	"encoding/json"
	"log"
)

// Loggable defines the shape of a message output to std
type Loggable struct {
	Caller string
	Data   map[string]interface{}
}

// Debug logs loggable data with a level of debug
func Debug(loggable Loggable) {
	logForLevel("Debug", loggable)
}

// Info logs loggable data with a level of info
func Info(loggable Loggable) {
	logForLevel("Info", loggable)
}

// Warning logs loggable data with a level of warning
func Warning(loggable Loggable) {
	logForLevel("Warning", loggable)
}

// Error logs loggable data with a level of error
func Error(loggable Loggable) {
	logForLevel("Error", loggable)
}

func logForLevel(level string, loggable Loggable) {
	// TODO: os.getEnv loglevel; log level depending on env
	// LOG_LEVEL debug: all
	// LOG_LEVEL info: all except debug
	// LOG_LEVEL warning: warning and error
	// LOG_LEVEL error: only error
	json, _ := json.MarshalIndent(loggable.Data, "", "\t")
	log.Printf("%s\nCaller: %s\n%+v\n", level, loggable.Caller, string(json))
}
