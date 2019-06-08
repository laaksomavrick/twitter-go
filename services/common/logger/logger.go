package logger

import (
	"twitter-go/services/common/env"

	log "github.com/sirupsen/logrus"
)

var logLevelMappingTable = map[string]int{
	"debug":   0,
	"info":    1,
	"warning": 2,
	"error":   3,
}

// Loggable defines the shape of a message output to std
type Loggable struct {
	Message string
	Data    map[string]interface{}
}

// Init initializes the logger configuration
func Init() {}

// Debug logs loggable data with a level of debug
func Debug(loggable Loggable) {
	ok := logForLevel("Debug", loggable)
	if ok {
		log.WithFields(loggable.Data).Debug(loggable.Message)
	}
}

// Info logs loggable data with a level of info
func Info(loggable Loggable) {
	ok := logForLevel("Info", loggable)
	if ok {
		log.WithFields(loggable.Data).Info(loggable.Message)
	}
}

// Warning logs loggable data with a level of warning
func Warning(loggable Loggable) {
	ok := logForLevel("Warning", loggable)
	if ok {
		log.WithFields(loggable.Data).Warn(loggable.Message)
	}
}

// Error logs loggable data with a level of error
func Error(loggable Loggable) {
	ok := logForLevel("Error", loggable)
	if ok {
		log.WithFields(loggable.Data).Error(loggable.Message)
	}
}

func logForLevel(logLevel string, loggable Loggable) bool {
	envLogLevel := env.GetEnv("LOG_LEVEL", "debug")
	envLogLevelInt := logLevelMappingTable[envLogLevel]
	logLevelInt := logLevelMappingTable[logLevel]
	return envLogLevelInt <= logLevelInt
}
