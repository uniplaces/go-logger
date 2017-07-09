package go_logger

import (
	"bytes"
	"errors"
	"testing"

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
}
