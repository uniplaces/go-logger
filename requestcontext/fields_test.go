package requestcontext

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldsRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := WithFields(context.Background())

	Set(ctx, "key", "value")
	Set(ctx, "num", 42)

	snap := Snapshot(ctx)

	assert.Equal(t, "value", snap["key"])
	assert.Equal(t, 42, snap["num"])
}

func TestFieldsSetWithoutBagIsNoop(t *testing.T) {
	t.Parallel()

	// must not panic
	Set(context.Background(), "key", "value")

	assert.Empty(t, Snapshot(context.Background()))
}

func TestFieldsSnapshotEmptyWhenAbsent(t *testing.T) {
	t.Parallel()

	snap := Snapshot(context.Background())

	assert.NotNil(t, snap)
	assert.Empty(t, snap)
}

func TestFieldsSharedAcrossDerivedContexts(t *testing.T) {
	t.Parallel()

	parent := WithFields(context.Background())
	child, cancel := context.WithCancel(parent)
	defer cancel()

	Set(child, "from-child", true)

	assert.Equal(t, true, Snapshot(parent)["from-child"])
}

func TestIterateFields(t *testing.T) {
	t.Parallel()

	ctx := WithFields(context.Background())
	Set(ctx, "a", 1)
	Set(ctx, "b", "two")

	got := map[string]any{}

	IterateFields(ctx, func(k string, v any) {
		got[k] = v
	})

	assert.Equal(t, 1, got["a"])
	assert.Equal(t, "two", got["b"])
}

func TestIterateFieldsNoopWhenAbsent(t *testing.T) {
	t.Parallel()

	called := false

	IterateFields(context.Background(), func(_ string, _ any) {
		called = true
	})

	assert.False(t, called)
}

func TestFieldsSnapshotIsACopy(t *testing.T) {
	t.Parallel()

	ctx := WithFields(context.Background())
	Set(ctx, "k", "v")

	snap := Snapshot(ctx)
	snap["k"] = "mutated"

	assert.Equal(t, "v", Snapshot(ctx)["k"])
}

func TestFieldsConcurrentAccess(t *testing.T) {
	t.Parallel()

	ctx := WithFields(context.Background())

	var wg sync.WaitGroup

	for i := 0; i < 200; i++ {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			Set(ctx, "n", n)
			_ = Snapshot(ctx)
		}(i)
	}

	wg.Wait()
}
