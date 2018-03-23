package go_logger

import (
	"bytes"
	"errors"
	"os"
	"testing"

	errorsPkg "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uniplaces/go-logger/internal"
)

func TestNew(t *testing.T) {
	assert.Nil(t, instance)

	config := NewConfig("env", "warning")
	Init(config)

	assert.NotNil(t, instance)

	resetInstance()
}

func TestLogWithFields(t *testing.T) {
	var buffer bytes.Buffer
	err := InitWithInstance(internal.NewLogrusLogger("error", &buffer))
	require.Nil(t, err)

	Builder().
		AddField("key", "value").
		AddContextField("foo", "bar").
		Error(errors.New("error test"))

	assert.Contains(t, buffer.String(), "\"context\":{\"foo\":\"bar\"}")
	assert.Contains(t, buffer.String(), "\"key\":\"value\"")
	assert.Contains(t, buffer.String(), "\"stack_trace\"")

	resetInstance()
}

func TestLogWithDefaultFields(t *testing.T) {
	os.Setenv("APPID", "app_id")
	os.Setenv("GOENV", "go_env")
	os.Setenv("GITHASH", "git_hash")

	var buffer bytes.Buffer
	err := InitWithInstance(internal.NewLogrusLogger("info", &buffer))
	require.Nil(t, err)

	AddDefaultField("test-field", "field_value", false)
	AddDefaultField("test-context-field", "context_field_value", true)

	expectedFields := map[string]interface{}{
		"type":       "app",
		"env":        "go_env",
		"git-hash":   "git_hash",
		"app-id":     "app_id",
		"test-field": "field_value",
		"key":        "value",
		"context": map[string]interface{}{
			"foo":                "bar",
			"test-context-field": "context_field_value",
		},
	}

	builder := Builder()
	builder.
		AddField("key", "value").
		AddContextField("foo", "bar").
		Info("info test")

	assert.Equal(t, expectedFields, builder.getFields())

	resetInstance()
}

func TestLogWithFieldsAndStacktrace(t *testing.T) {
	var buffer bytes.Buffer
	err := InitWithInstance(internal.NewLogrusLogger("error", &buffer))
	require.Nil(t, err)

	errorWithStackTrace := justToShowUpInStackTrace()

	Builder().
		AddField("key", "value").
		AddContextField("foo", "bar").
		Error(errorWithStackTrace)

	assert.Contains(t, buffer.String(), "\"context\":{\"foo\":\"bar\"}")
	assert.Contains(t, buffer.String(), "\"key\":\"value\"")
	assert.Contains(t, buffer.String(), "\"stack_trace\":")
	// test stack trace strings
	assert.Contains(t, buffer.String(), "github.com/uniplaces/go-logger.justToShowUpInStackTrace")
	assert.Contains(t, buffer.String(), "github.com/uniplaces/go-logger.TestLogWithFieldsAndStacktrace")
	assert.Contains(t, buffer.String(), "github.com/uniplaces/go-logger/logger_test.go:122")

	resetInstance()
}

func TestLog(t *testing.T) {
	var buffer bytes.Buffer
	err := InitWithInstance(internal.NewLogrusLogger("error", &buffer))
	require.Nil(t, err)

	Error(errors.New("test error"))

	assert.NotContains(t, buffer.String(), "\"context\"")
	assert.NotContains(t, buffer.String(), "\"key\":\"value\"")
	assert.Contains(t, buffer.String(), "\"stack_trace\"")

	resetInstance()
}

func resetInstance() {
	instance = nil
	defaultFields = []defaultField{}
}

func justToShowUpInStackTrace() error {
	return errorsPkg.New("error test")
}
