package internal

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestLogrusLoggerLevel(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	l := NewLogrusLogger("warning", "dev", &buffer, nil)
	l.Info("test info")
	l.Debug("test debug")

	assert.Empty(t, buffer.String())

	l.Error("test error")

	assert.Contains(t, buffer.String(), "level=2")
	assert.Contains(t, buffer.String(), "msg=\"test error\"")
}

func TestLogrusLoggerStackTrace(t *testing.T) {
	t.Parallel()

	var zeroValueFields map[string]interface{}
	var buffer bytes.Buffer
	l := NewLogrusLogger("debug", "dev", &buffer, nil)

	l.DebugWithFields("debug", zeroValueFields)
	assert.NotContains(t, buffer.String(), "stack_trace")
	buffer.Reset()

	l.InfoWithFields("info", zeroValueFields)
	assert.NotContains(t, buffer.String(), "stack_trace")
	buffer.Reset()

	l.WarningWithFields("warning", zeroValueFields)
	assert.NotContains(t, buffer.String(), "stack_trace")
	buffer.Reset()

	l.ErrorWithFields(errors.New("error"), zeroValueFields)
	assert.Contains(t, buffer.String(), "stack_trace")
}

func TestLogrusLoggerStackTraceShouldSkip(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	var buffer bytes.Buffer
	l := NewLogrusLogger("warning", "dev", &buffer, map[string]interface{}{})

	l.DebugWithFields("debug", map[string]interface{}{"test": 123})
	assert.Empty(t, buffer.String())

	l.WarningWithFields("warning", map[string]interface{}{"abc": 321})
	assert.Contains(t, buffer.String(), "abc=321")
}

func TestLogrusLoggerInvalidConfig(t *testing.T) {
	t.Parallel()

	defer func() {
		assert.NotNil(t, recover())
	}()

	NewLogrusLogger("invalid level", "dev", nil, nil)
}

func TestFirstStackTracerInErrorChainWithOneError(t *testing.T) {
	rootError := errors.New("root error")

	_, ok := rootError.(stackTracer)
	assert.True(t, ok)

	gotError := firstStackTracerInErrorChain(rootError)

	assert.EqualError(t, gotError, rootError.Error())
}

func TestFirstStackTracerInErrorChainWithMultipleErrors(t *testing.T) {
	rootError := errors.New("root error")
	stackTracerError := errors.Wrap(rootError, "wrap root error")
	newContextError := errors.WithMessage(stackTracerError, "new context added")

	gotError := firstStackTracerInErrorChain(newContextError)

	_, ok := gotError.(stackTracer)
	assert.True(t, ok)

	assert.EqualError(t, gotError, stackTracerError.Error())
}
