// Package requestcontext carries per-request id and fields through context.Context and propagates them over HTTP.
package requestcontext

import (
	"context"

	"github.com/google/uuid"
)

// HeaderName is the inbound/outbound HTTP header carrying the request id.
const HeaderName = "X-Request-Id"

type idCtxKey struct{}

// WithID returns ctx carrying requestID; an empty requestID is ignored and ctx is returned unchanged.
func WithID(ctx context.Context, requestID string) context.Context {
	if requestID == "" {
		return ctx
	}

	return context.WithValue(ctx, idCtxKey{}, requestID)
}

// ID returns the request id stashed on ctx, or "" if absent.
func ID(ctx context.Context) string {
	if id, ok := ctx.Value(idCtxKey{}).(string); ok {
		return id
	}

	return ""
}

// Ensure returns ctx with a request id, minting a UUIDv4 when absent; idempotent when present.
func Ensure(ctx context.Context) (context.Context, string) {
	if id := ID(ctx); id != "" {
		return ctx, id
	}

	id := uuid.NewString()

	return WithID(ctx, id), id
}
