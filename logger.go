package go_logger

import (
	"errors"
	"os"
	"sync"

	"github.com/uniplaces/go-logger/internal"
)

var instance Logger
var once sync.Once

// Init initializes logger instance
func Init(config Config) error {
	if instance != nil {
		return errors.New("logger cannot be initialized more than once")
	}

	once.Do(func() {
		// todo use implementation according to the env
		instance = internal.NewLogrusLogger(config.level, os.Stdout)
	})

	return nil
}

// InitWithInstance sets logger to an instance (for testing purposes)
func InitWithInstance(newInstance Logger) error {
	if instance != nil {
		return errors.New("logger cannot be initialized more than once")
	}

	instance = newInstance

	return nil
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

	instance.ErrorWithFields(err.Error(), builder.getFields())
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

	instance.WarningWithFields(message, builder.getFields())
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

	instance.InfoWithFields(message, builder.getFields())
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

	instance.DebugWithFields(message, builder.getFields())
}
