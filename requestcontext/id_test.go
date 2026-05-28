package requestcontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := WithID(context.Background(), "abc-123")

	assert.Equal(t, "abc-123", ID(ctx))
}

func TestIDAbsentReturnsEmpty(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "", ID(context.Background()))
}

func TestEnsureMintsUUIDWhenAbsent(t *testing.T) {
	t.Parallel()

	ctx, id := Ensure(context.Background())

	assert.NotEmpty(t, id)
	assert.Equal(t, id, ID(ctx))
}

func TestEnsurePreservesExistingID(t *testing.T) {
	t.Parallel()

	in := WithID(context.Background(), "existing-id")

	out, id := Ensure(in)

	assert.Equal(t, in, out)
	assert.Equal(t, "existing-id", id)
}

func TestHeaderNameConstant(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "X-Request-Id", HeaderName)
}
