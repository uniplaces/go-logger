package requestcontext

import (
	"context"

	logger "github.com/uniplaces/go-logger"
)

// HTTPFailure emits an error-level log for an HTTP-adapter failure (status=0 if no response, err=nil for non-2xx).
func HTTPFailure(
	ctx context.Context,
	component, requestURL string,
	status int,
	err error,
	reason string,
) {
	logger.Builder().
		AddField("request_id", ID(ctx)).
		AddField("component", component).
		AddField("url", requestURL).
		AddField("status", status).
		AddField("reason", reason).
		EmitFailure(reason, err)
}

// QueryFailure emits an error-level log for a DB query failure; queryID must not contain bound parameter values.
func QueryFailure(
	ctx context.Context,
	component, queryID string,
	err error,
	reason string,
) {
	logger.Builder().
		AddField("request_id", ID(ctx)).
		AddField("component", component).
		AddField("query", queryID).
		AddField("reason", reason).
		EmitFailure(reason, err)
}
