package logger

import (
	"encoding/json"
	"log"
	"twitter-go/services/common/env"
)

var logLevelMappingTable = map[string]int{
	"debug":   0,
	"info":    1,
	"warning": 2,
	"error":   3,
}

// Loggable defines the shape of a message output to std
type Loggable struct {
	Caller  string
	Message string
	Data    map[string]interface{}
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

// TODO-12: make this similar to cassandra logger
// Give each request a uuid for tracing?
func logForLevel(logLevel string, loggable Loggable) {
	envLogLevel := env.GetEnv("LOG_LEVEL", "debug")
	envLogLevelInt := logLevelMappingTable[envLogLevel]
	logLevelInt := logLevelMappingTable[logLevel]
	if envLogLevelInt <= logLevelInt {
		json, _ := json.MarshalIndent(loggable.Data, "", "\t")
		log.Printf(
			"%s\nCaller: %s\nMessage: %s\n%+v\n",
			logLevel,
			loggable.Caller,
			loggable.Message,
			string(json),
		)
	}
}
