package internal

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"

	"github.com/uniplaces/logrus"
	"github.com/pkg/errors"
)

// defines which levels log stack trace
var stackTraceLevels = map[logrus.Level]bool{
	logrus.ErrorLevel: true,
	logrus.WarnLevel:  false,
	logrus.InfoLevel:  false,
	logrus.DebugLevel: false,
}

type logrusLogger struct {
	*logrus.Logger
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type causer interface {
	Cause() error
}

const (
	stackTraceKey    = "stack_trace"
	stackTraceFormat = "%s:%d"
	stackTraceErrorPkgFormat = "%+v"
)

func NewLogrusLogger(level string, writer io.Writer) logrusLogger {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		panic(fmt.Sprintf("invalid log level (%s)", level))
	}

	instance := logrusLogger{
		Logger: logrus.New(),
	}

	instance.Formatter = &logrus.JSONFormatter{
		EnableIntLogLevels:true,
	}
	instance.Level = logrusLevel
	instance.Out = writer

	log.SetOutput(instance.Writer())

	return instance
}

func (logger logrusLogger) ErrorWithFields(err error, fields map[string]interface{}) {
	stackTrace := logger.getErrorLevelStackTrace(err)

	entry := logger.entry(fields, stackTrace)
	entry.Error(err.Error())
}

func (logger logrusLogger) WarningWithFields(message string, fields map[string]interface{}) {
	var stackTrace []string
	if stackTraceLevels[logrus.WarnLevel] {
		stackTrace = buildStackTrace()
	}

	entry := logger.entry(fields, stackTrace)
	entry.Warning(message)
}

func (logger logrusLogger) InfoWithFields(message string, fields map[string]interface{}) {
	var stackTrace []string
	if stackTraceLevels[logrus.InfoLevel] {
		stackTrace = buildStackTrace()
	}

	entry := logger.entry(fields, stackTrace)
	entry.Info(message)
}

func (logger logrusLogger) DebugWithFields(message string, fields map[string]interface{}) {
	var stackTrace []string
	if stackTraceLevels[logrus.DebugLevel] {
		stackTrace = buildStackTrace()
	}

	entry := logger.entry(fields, stackTrace)
	entry.Debug(message)
}

func (logger logrusLogger) getErrorLevelStackTrace(err error) []string {
	var stackTrace []string
	if stackTraceLevels[logrus.ErrorLevel] {
		switch err := firstStackTracerInErrorChain(err).(type) {
		case stackTracer:
			for _, frame := range err.StackTrace() {
				stackTrace = append(stackTrace, fmt.Sprintf(stackTraceErrorPkgFormat, frame))
			}
		default:
			stackTrace = buildStackTrace()
		}
	}

	return stackTrace
}

func (logger logrusLogger) entry(fields map[string]interface{}, stackTrace []string) *logrus.Entry {
	if len(fields) == 0 {
		logFields := logrus.Fields{}
		if len(stackTrace) > 0 {
			logFields[stackTraceKey] = stackTrace
		}

		return logger.WithFields(logFields)
	}

	logFields := logrus.Fields(fields)
	if len(stackTrace) > 0 {
		logFields[stackTraceKey] = stackTrace
	}

	return logger.WithFields(logFields)
}

func buildStackTrace() []string {
	var stacktrace []string

	skip := 0
	for {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		skip++

		if shouldSkipFile(file) {
			continue
		}

		stacktrace = append(stacktrace, fmt.Sprintf(stackTraceFormat, file, line))
	}

	return stacktrace
}

func shouldSkipFile(file string) bool {
	skipIfContainsList := []string{
		"github.com/uniplaces/go-logger",
		"github.com/gin-gonic",
		"autogenerated", // https://github.com/golang/go/issues/16723
		"go/src/runtime/asm_amd64.s",
		"go/src/net/http/server.go",
	}

	for _, s := range skipIfContainsList {
		if strings.Contains(file, s) {
			return true
		}
	}

	return false
}

func firstStackTracerInErrorChain(err error) error {
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}

		if _, ok := err.(stackTracer); ok {
			break
		}

		err = cause.Cause()
	}

	return err
}
