package go_logger

const contextFieldsKey = "context"

type builder struct {
	fields        map[string]interface{}
	contextFields map[string]interface{}
}

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

func (builder builder) getDefaultFields(fields []DefaultField) map[string]interface{} {
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
