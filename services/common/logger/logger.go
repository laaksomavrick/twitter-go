package logger

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"
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
	Data    interface{}
}

// Init initializes the logger configuration
func Init() {}

// Debug logs loggable data with a level of debug
func Debug(loggable Loggable) {
	ok := logForLevel("Debug", loggable)
	if ok {
		data := convertDataToMap(loggable.Data)
		log.WithFields(data).Debug(loggable.Message)
	}
}

// Info logs loggable data with a level of info
func Info(loggable Loggable) {
	ok := logForLevel("Info", loggable)
	if ok {
		data := convertDataToMap(loggable.Data)
		log.WithFields(data).Info(loggable.Message)
	}
}

// Warning logs loggable data with a level of warning
func Warning(loggable Loggable) {
	ok := logForLevel("Warning", loggable)
	if ok {
		data := convertDataToMap(loggable.Data)
		log.WithFields(data).Warn(loggable.Message)
	}
}

// Error logs loggable data with a level of error
func Error(loggable Loggable) {
	ok := logForLevel("Error", loggable)
	if ok {
		data := convertDataToMap(loggable.Data)
		log.WithFields(data).Error(loggable.Message)
	}
}

func logForLevel(logLevel string, loggable Loggable) bool {
	envLogLevel := env.GetEnv("LOG_LEVEL", "debug")
	envLogLevelInt := logLevelMappingTable[envLogLevel]
	logLevelInt := logLevelMappingTable[logLevel]
	return envLogLevelInt <= logLevelInt
}

func convertDataToMap(data interface{}) map[string]interface{} {
	if data == nil {
		return map[string]interface{}{}
	}

	klass := reflect.TypeOf(data).Kind()

	if klass != reflect.Map && klass != reflect.Struct && klass != reflect.Slice && klass != reflect.String {
		log.Printf("Unrecognized type in convertDataToMap: %s", klass)
		return map[string]interface{}{}
	}

	if klass == reflect.String {
		strings.Replace(data.(string), "\n", "", -1)
		data = map[string]interface{}{"data": data}
	}

	var serialized map[string]interface{}

	jsonString, err := json.Marshal(data)

	if err != nil {
		log.Printf("Something went wrong in convertDataToMap for: %s", klass)
		log.Printf(err.Error())
		return map[string]interface{}{}
	}

	json.Unmarshal(jsonString, &serialized)

	if serialized != nil {
		serialized["timestamp"] = time.Now().UTC()
	}

	return serialized
}
