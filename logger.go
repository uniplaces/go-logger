package go_logger

import (
	"errors"
	"os"
	"sync"

	"github.com/uniplaces/go-logger/internal"
)

const logType = "app"

var instance Logger
var once sync.Once
var defaultFields []defaultField

type defaultField struct {
	key            string
	value          interface{}
	isContextField bool
}

// Init initializes logger instance
func Init(config Config) error {
	if instance != nil {
		return errors.New("logger cannot be initialized more than once")
	}

	once.Do(func() {
		// todo use implementation according to the env
		instance = internal.NewLogrusLogger(config.level, os.Stdout)
	})

	addMandatoryDefaultFields()

	return nil
}

// InitWithInstance sets logger to an instance (for testing purposes)
func InitWithInstance(newInstance Logger) error {
	if instance != nil {
		return errors.New("logger cannot be initialized more than once")
	}

	instance = newInstance

	addMandatoryDefaultFields()

	return nil
}

func AddDefaultField(key string, value interface{}, isContextField bool) {
	defaultFields = append(defaultFields, defaultField{key: key, value: value, isContextField: isContextField})
}

// Error logs a error message
func Error(err error) {
	Builder().Error(err)
}

// Error logs a error message with fields
func (builder builder) Error(err error) {
	if instance == nil {
		return
	}

	instance.ErrorWithFields(err, builder.getDefaultFields(defaultFields))
}

// Warning logs a warning message
func Warning(message string) {
	Builder().Warning(message)
}

// Warning logs a warning message with fields
func (builder builder) Warning(message string) {
	if instance == nil {
		return
	}

	instance.WarningWithFields(message, builder.getDefaultFields(defaultFields))
}

// Info logs a info message
func Info(message string) {
	Builder().Info(message)
}

// Info logs a info message with fields
func (builder builder) Info(message string) {
	if instance == nil {
		return
	}

	instance.InfoWithFields(message, builder.getDefaultFields(defaultFields))
}

// Debug logs a debug message
func Debug(message string) {
	Builder().Debug(message)
}

// Debug logs a debug message with fields
func (builder builder) Debug(message string) {
	if instance == nil {
		return
	}

	instance.DebugWithFields(message, builder.getDefaultFields(defaultFields))
}

func addMandatoryDefaultFields() {
	AddDefaultField("type", logType, false)
	AddDefaultField("app-id", os.Getenv("APPID"), false)
	AddDefaultField("env", os.Getenv("GOENV"), false)
	AddDefaultField("git-hash", os.Getenv("GITHASH"), false)
}
