package internal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogrusLoggerLevel(t *testing.T) {
	var buffer bytes.Buffer
	l := NewLogrusLogger("warning", &buffer)
	l.Info("test info")
	l.Debug("test debug")

	assert.Empty(t, buffer.String())

	l.Error("test error")

	assert.Contains(t, buffer.String(), "\"level\":\"error\"")
	assert.Contains(t, buffer.String(), "\"msg\":\"test error\"")
}

func TestLogrusLoggerStackTrace(t *testing.T) {
	var buffer bytes.Buffer
	l := NewLogrusLogger("debug", &buffer)

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

func TestLogrusLoggerStackTraceShouldSkip(t *testing.T) {
	testData := map[string]bool{
		"should not skip test":                             false,
		"/usr/local/go/src/net/http/server.go":             true,
		"/vendor/github.com/uniplaces/go-logger/logger.go": true,
		"/vendor/github.com/gin-gonic/gin/gin.go":          true,
		"delivery/api/handlers/ping/ping.go":               false,
	}

	for file, expectedToSkip := range testData {
		assert.Equal(t, expectedToSkip, shouldSkipFile(file))
	}
}

func TestLogrusLoggerWithFields(t *testing.T) {
	var buffer bytes.Buffer
	l := NewLogrusLogger("warning", &buffer)

	l.DebugWithFields("debug", map[string]interface{}{"test": 123})
	assert.Empty(t, buffer.String())

	l.WarningWithFields("warning", map[string]interface{}{"abc": 321})
	assert.Contains(t, buffer.String(), "\"abc\":321")
}

func TestLogrusLoggerInvalidConfig(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	NewLogrusLogger("invalid level", nil)
}
