package requestcontext

import (
	"context"

	logger "github.com/uniplaces/go-logger"
)

// Logger returns a logger.LogBuilder pre-populated with request_id (when
// ID(ctx) != "") and every field currently in the bag (when WithFields was
// called on ctx). Callers chain .AddField / .AddContextField for call-site
// fields and terminate with .Info / .Warning / .Error / .Debug.
func Logger(ctx context.Context) logger.LogBuilder {
	b := logger.Builder()

	if id := ID(ctx); id != "" {
		b = b.AddField("request_id", id)
	}

	for k, v := range Snapshot(ctx) {
		b = b.AddField(k, v)
	}

	return b
}
