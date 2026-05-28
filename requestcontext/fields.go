package requestcontext

import (
	"context"
	"sync"
)

type fieldsCtxKey struct{}

type fieldsBag struct {
	mu sync.RWMutex
	m  map[string]any
}

// WithFields attaches an empty mutable field bag to ctx. Call once at the
// request boundary, before code paths that may add fields. Subsequent Set
// calls on ctx (or any context derived from it) mutate the same bag.
func WithFields(ctx context.Context) context.Context {
	return context.WithValue(ctx, fieldsCtxKey{}, &fieldsBag{m: make(map[string]any)})
}

// Set assigns key=value on the field bag attached to ctx. Safe for
// concurrent use. No-op (silent) if WithFields was not called on ctx.
func Set(ctx context.Context, key string, value any) {
	bag, ok := ctx.Value(fieldsCtxKey{}).(*fieldsBag)
	if !ok {
		return
	}

	bag.mu.Lock()
	defer bag.mu.Unlock()

	bag.m[key] = value
}

// Snapshot returns a shallow copy of the field bag attached to ctx, or an
// empty map (never nil) if WithFields was not called.
func Snapshot(ctx context.Context) map[string]any {
	bag, ok := ctx.Value(fieldsCtxKey{}).(*fieldsBag)
	if !ok {
		return map[string]any{}
	}

	bag.mu.RLock()
	defer bag.mu.RUnlock()

	result := make(map[string]any, len(bag.m))
	for k, v := range bag.m {
		result[k] = v
	}

	return result
}
