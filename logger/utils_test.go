package logger

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestCreateExtraField(t *testing.T) {
	createdExtraField := CreateExtraField("test", 123)

	assert.Equal(t, createdExtraField, extraField{key: "test", value: 123})
}

func TestCreateFields(t *testing.T) {
	fieldsArg := map[string]interface{}{"test": 123}
	extraField := extraField{key: "extra", value: 321}

	expected := Fields{
		"test": 123,
		"extra_info": map[string]interface{}{
			"extra": 321,
		},
	}

	fields := CreateFields(fieldsArg, extraField)

	assert.Equal(t, fields, expected)
}
