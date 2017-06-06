package internal_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uniplaces/go-logger/internal"
)

func TestLogrusLoggerLevel(t *testing.T) {
	var buffer bytes.Buffer
	l := internal.NewLogrusLogger("warning", &buffer)
	l.Info("test info")
	l.Debug("test debug")

	assert.Empty(t, buffer.String())

	l.Error("test error")

	assert.Contains(t, buffer.String(), "\"level\":\"error\"")
	assert.Contains(t, buffer.String(), "\"msg\":\"test error\"")
}

func TestLogrusLoggerStackTrace(t *testing.T) {
	var buffer bytes.Buffer
	l := internal.NewLogrusLogger("debug", &buffer)

	l.Debug("debug")
	assert.NotContains(t, buffer.String(), "stacktrace")
	buffer.Reset()

	l.Info("info")
	assert.NotContains(t, buffer.String(), "stacktrace")
	buffer.Reset()

	l.Warning("warning")
	assert.NotContains(t, buffer.String(), "stacktrace")
	buffer.Reset()

	l.Error("error")
	assert.Contains(t, buffer.String(), "stacktrace")
}

func TestLogrusLoggerWithFields(t *testing.T) {
	var buffer bytes.Buffer
	l := internal.NewLogrusLogger("warning", &buffer)

	l.DebugWithFields("debug", map[string]interface{}{"test": 123})
	assert.Empty(t, buffer.String())

	l.WarningWithFields("warning", map[string]interface{}{"abc": 321})
	assert.Contains(t, buffer.String(), "\"abc\":321")
}
