package requestcontext_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniplaces/go-logger/requestcontext"
)

func TestHTTPFailureWithError(t *testing.T) {
	testBuf.Reset()

	ctx := requestcontext.WithID(context.Background(), "req-fail-1")

	requestcontext.HTTPFailure(
		ctx,
		"aggregator_http",
		"https://aggregator/api/v1/offer/1",
		0,
		errors.New("dial tcp: timeout"),
		"http call failed",
	)

	got := decodeLastLine(t)
	assert.Equal(t, "error", got["level"])
	assert.Equal(t, "req-fail-1", got["request_id"])
	assert.Equal(t, "aggregator_http", got["component"])
	assert.Equal(t, "https://aggregator/api/v1/offer/1", got["url"])
	assert.Equal(t, float64(0), got["status"])
	assert.Equal(t, "http call failed", got["reason"])
	assert.Equal(t, "dial tcp: timeout", got["error_message"])
	assert.Equal(t, "http call failed: dial tcp: timeout", got["msg"])
}

func TestHTTPFailureWithoutError(t *testing.T) {
	testBuf.Reset()

	ctx := requestcontext.WithID(context.Background(), "req-fail-2")

	requestcontext.HTTPFailure(
		ctx,
		"supplyclassifier_http",
		"https://supply/classify",
		500,
		nil,
		"non-2xx response",
	)

	got := decodeLastLine(t)
	assert.Equal(t, "error", got["level"])
	assert.Equal(t, "req-fail-2", got["request_id"])
	assert.Equal(t, "supplyclassifier_http", got["component"])
	assert.Equal(t, float64(500), got["status"])
	assert.Equal(t, "non-2xx response", got["reason"])

	_, hasErrMsg := got["error_message"]
	assert.False(t, hasErrMsg, "error_message should be omitted when err is nil")

	assert.Equal(t, "non-2xx response", got["msg"])
}

func TestQueryFailure(t *testing.T) {
	testBuf.Reset()

	ctx := requestcontext.WithID(context.Background(), "req-fail-3")

	requestcontext.QueryFailure(
		ctx,
		"mysql",
		"open_api_excluded_offers.select",
		errors.New("connection refused"),
		"mysql query failed",
	)

	got := decodeLastLine(t)
	assert.Equal(t, "error", got["level"])
	assert.Equal(t, "req-fail-3", got["request_id"])
	assert.Equal(t, "mysql", got["component"])
	assert.Equal(t, "open_api_excluded_offers.select", got["query"])
	assert.Equal(t, "mysql query failed", got["reason"])
	assert.Equal(t, "connection refused", got["error_message"])
	assert.Equal(t, "mysql query failed: connection refused", got["msg"])
}
