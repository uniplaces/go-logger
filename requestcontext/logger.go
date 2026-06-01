package requestcontext

import (
	"context"

	logger "github.com/uniplaces/go-logger"
)

// Logger returns a LogBuilder pre-populated with request_id and any fields set on ctx.
func Logger(ctx context.Context) logger.LogBuilder {
	logBuilder := logger.Builder()

	if id := ID(ctx); id != "" {
		logBuilder = logBuilder.AddField("request_id", id)
	}

	IterateFields(ctx, func(key string, value any) {
		logBuilder = logBuilder.AddField(key, value)
	})

	return logBuilder
}
