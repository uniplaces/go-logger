package go_logger

import "os"

const (
	contextFieldsKey = "context"
	logType          = "app"
)

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

func (builder builder) getFieldsWithMandatoryKeys() map[string]interface{} {
	builder.
		AddField("type", logType).
		AddField("app-id", os.Getenv("APPID")).
		AddField("env", os.Getenv("GOENV")).
		AddField("git-hash", os.Getenv("GITHASH"))

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
