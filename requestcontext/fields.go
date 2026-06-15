package requestcontext

import (
	"context"
	"sync"
)

type fieldsCtxKey struct{}

type fieldsBag struct {
	mutex  sync.RWMutex
	fields map[string]any
}

// WithFields attaches an empty mutable field bag to ctx; call once at the request boundary.
func WithFields(ctx context.Context) context.Context {
	return context.WithValue(ctx, fieldsCtxKey{}, &fieldsBag{fields: make(map[string]any)})
}

// Set assigns key=value on the field bag attached to ctx; concurrency-safe, no-op without WithFields.
func Set(ctx context.Context, key string, value any) {
	bag, ok := ctx.Value(fieldsCtxKey{}).(*fieldsBag)
	if !ok {
		return
	}

	bag.mutex.Lock()
	defer bag.mutex.Unlock()

	bag.fields[key] = value
}

// IterateFields calls fn for each field on ctx under the read lock; allocation-free, no-op without WithFields.
func IterateFields(ctx context.Context, fn func(key string, value any)) {
	bag, ok := ctx.Value(fieldsCtxKey{}).(*fieldsBag)
	if !ok {
		return
	}

	bag.mutex.RLock()
	defer bag.mutex.RUnlock()

	for key, value := range bag.fields {
		fn(key, value)
	}
}

// Snapshot returns a shallow copy of the field bag on ctx, or an empty map (never nil).
func Snapshot(ctx context.Context) map[string]any {
	bag, ok := ctx.Value(fieldsCtxKey{}).(*fieldsBag)
	if !ok {
		return map[string]any{}
	}

	bag.mutex.RLock()
	defer bag.mutex.RUnlock()

	result := make(map[string]any, len(bag.fields))
	for key, value := range bag.fields {
		result[key] = value
	}

	return result
}
