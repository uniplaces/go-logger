package go_logger

import (
	"errors"
	"fmt"
)

const contextFieldsKey = "context"

type builder struct {
	fields        map[string]interface{}
	contextFields map[string]interface{}
}

// LogBuilder is an exported type alias for the per-log field-builder returned by Builder().
// Use this in function signatures in external packages that need to name the type
// (e.g. github.com/uniplaces/go-logger/requestcontext.Logger).
//
// The alias is named LogBuilder rather than Builder because go-logger already exports
// a Builder() function and Go uses a single namespace for types and functions within a
// package.
type LogBuilder = builder

// Builder initializes logger fields builder
func Builder() builder {
	return builder{
		fields:        make(map[string]interface{}),
		contextFields: make(map[string]interface{}),
	}
}

// AddField adds a new field
func (builder builder) AddField(key string, value interface{}) builder {
	builder.fields[key] = value

	return builder
}

func (builder builder) AddFields(fields map[string]interface{}) builder {
	for k, v := range fields {
		builder.AddField(k, v)
	}

	return builder
}

// AddContextField adds a new context field
func (builder builder) AddContextField(key string, value interface{}) builder {
	builder.contextFields[key] = value

	return builder
}

func (builder builder) getDefaultFields(fields []defaultField) map[string]interface{} {
	for _, field := range fields {
		if field.isContextField {
			builder.AddContextField(field.key, field.value)

			continue
		}

		builder.AddField(field.key, field.value)
	}

	return builder.getFields()
}

func (builder builder) getFields() map[string]interface{} {
	fields := builder.fields
	if len(builder.contextFields) > 0 {
		fields[contextFieldsKey] = builder.contextFields
	}

	if len(fields) == 0 {
		return nil
	}

	return fields
}

// EmitFailure terminates the builder at error level, logging the wrapped err or just reason when err is nil.
func (builder builder) EmitFailure(reason string, err error) {
	if err != nil {
		builder.AddField("error_message", err.Error()).Error(fmt.Errorf("%s: %w", reason, err))

		return
	}

	builder.Error(errors.New(reason))
}
