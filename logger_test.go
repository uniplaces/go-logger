package go_logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Nil(t, instance)

	config := NewConfig("env", "warning")
	Init(config)

	assert.NotNil(t, instance)

	instance = nil
}
