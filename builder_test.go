package go_logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	initializedFields := make(map[string]interface{})

	builder := Builder()
	assert.Equal(t, builder.fields, initializedFields)
	assert.Equal(t, builder.contextFields, initializedFields)
}

func TestBuilder_AddField(t *testing.T) {
	expectedFields := map[string]interface{}{
		"test": "value",
	}

	builder := Builder()
	builder.AddField("test", "value")

	assert.Equal(t, expectedFields, builder.fields)
}

func TestBuilder_AddContextField(t *testing.T) {
	expectedContextFields := map[string]interface{}{
		"test": "value",
	}

	builder := Builder()
	builder.AddContextField("test", "value")

	assert.Equal(t, expectedContextFields, builder.contextFields)
}

func TestBuilder_GetFields(t *testing.T) {
	expectedFields := map[string]interface{}{
		"test": "value",
		"context": map[string]interface{}{
			"foo": "bar",
		},
	}

	builder := Builder()
	builder.
		AddField("test", "value").
		AddContextField("foo", "bar")

	assert.Equal(t, expectedFields, builder.getFields())
}
