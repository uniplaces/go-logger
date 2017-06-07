package go_logger

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestCreateExtraField(t *testing.T) {
	t.Parallel()

	createdExtraField := CreateExtraField("test", 123)

	assert.Equal(t, createdExtraField, extraField{key: "test", value: 123})
}

func TestCreateFields(t *testing.T) {
	t.Parallel()

	fieldsArg := map[string]interface{}{"test": 123}
	extraField := extraField{key: "extra", value: 321}

	expected := fields{
		"test": 123,
		"extra_info": map[string]interface{}{
			"extra": 321,
		},
	}

	fields := CreateFields(fieldsArg, extraField)

	assert.Equal(t, fields, expected)
}
