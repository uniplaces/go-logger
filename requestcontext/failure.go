package requestcontext

import (
	"context"
	"errors"
	"fmt"

	logger "github.com/uniplaces/go-logger"
)

// HTTPFailure emits an error-level log for a failure originating in an
// HTTP-style infrastructure adapter. Tags: request_id, component, url,
// status, reason, error_message (when err != nil).
//
// Pass status=0 when no response was received (network or request-build
// error) and err=nil when the response itself is the failure (non-2xx).
// 404 is not inspected by this helper — callers decide whether a 404 is a
// benign business miss or a true failure.
func HTTPFailure(
	ctx context.Context,
	component, requestURL string,
	status int,
	err error,
	reason string,
) {
	b := logger.Builder().
		AddField("request_id", ID(ctx)).
		AddField("component", component).
		AddField("url", requestURL).
		AddField("status", status).
		AddField("reason", reason)

	if err != nil {
		b.AddField("error_message", err.Error()).
			Error(fmt.Errorf("%s: %w", reason, err))

		return
	}

	b.Error(errors.New(reason))
}

// QueryFailure emits an error-level log for a database query failure.
// queryID identifies the query — never include bound parameter values.
// Tags: request_id, component, query, reason, error_message.
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
		AddField("error_message", err.Error()).
		Error(fmt.Errorf("%s: %w", reason, err))
}
