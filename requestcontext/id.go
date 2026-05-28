// Package requestcontext carries per-request correlation id and structured
// fields through context.Context, and provides origin-level failure helpers
// + outbound HTTP propagation for Uniplaces Go services.
package requestcontext

import (
	"context"

	"github.com/google/uuid"
)

// HeaderName is the inbound/outbound HTTP header carrying the request id.
const HeaderName = "X-Request-Id"

type idCtxKey struct{}

// WithID returns a new context carrying the given request id. Use at the
// request boundary once the id is known (header read, SQS message id, etc.).
func WithID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, idCtxKey{}, requestID)
}

// ID returns the request id stashed on ctx, or "" if absent.
func ID(ctx context.Context) string {
	if id, ok := ctx.Value(idCtxKey{}).(string); ok {
		return id
	}

	return ""
}

// Ensure returns a context with a request id attached. If ctx already has
// one, ctx is returned unchanged. Otherwise a UUIDv4 is minted and attached.
// Use at the entry point of a non-HTTP request path (SQS handler, CLI item
// loop, scheduled job) where there is no inbound header to honour.
func Ensure(ctx context.Context) (context.Context, string) {
	if id := ID(ctx); id != "" {
		return ctx, id
	}

	id := uuid.NewString()

	return WithID(ctx, id), id
}
