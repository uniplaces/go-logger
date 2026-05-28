package requestcontext_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	logger "github.com/uniplaces/go-logger"
	"github.com/uniplaces/go-logger/internal"
	"github.com/uniplaces/go-logger/requestcontext"
)

var testBuf bytes.Buffer

func TestMain(m *testing.M) {
	if err := logger.InitWithInstance(internal.NewLogrusLogger("info", "test", &testBuf, nil)); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func decodeLastLine(t *testing.T) map[string]any {
	t.Helper()

	line := testBuf.Bytes()
	if i := bytes.LastIndexByte(bytes.TrimRight(line, "\n"), '\n'); i >= 0 {
		line = line[i+1:]
	}

	var got map[string]any

	require.NoError(t, json.Unmarshal(bytes.TrimSpace(line), &got))

	return got
}

func TestLoggerPrePopulatesRequestIDAndFields(t *testing.T) {
	testBuf.Reset()

	ctx := requestcontext.WithID(context.Background(), "req-1")
	ctx = requestcontext.WithFields(ctx)
	requestcontext.Set(ctx, "graphql_operation", "query")
	requestcontext.Set(ctx, "status", 200)

	requestcontext.Logger(ctx).Info("hello")

	got := decodeLastLine(t)
	assert.Equal(t, "req-1", got["request_id"])
	assert.Equal(t, "query", got["graphql_operation"])
	assert.Equal(t, float64(200), got["status"]) // JSON numbers decode as float64
}

func TestLoggerEmptyContextOmitsRequestID(t *testing.T) {
	testBuf.Reset()

	requestcontext.Logger(context.Background()).Info("hello")

	got := decodeLastLine(t)

	_, hasID := got["request_id"]
	assert.False(t, hasID, "request_id should not be present when ctx has no id")
}

func TestLoggerWithIDOnlyNoBag(t *testing.T) {
	testBuf.Reset()

	ctx := requestcontext.WithID(context.Background(), "req-2")

	requestcontext.Logger(ctx).Info("hello")

	got := decodeLastLine(t)
	assert.Equal(t, "req-2", got["request_id"])
}

func TestLoggerCallSiteFieldsCanChain(t *testing.T) {
	testBuf.Reset()

	ctx := requestcontext.WithID(context.Background(), "req-3")

	requestcontext.Logger(ctx).
		AddField("call_site", "specific").
		Info("hello")

	got := decodeLastLine(t)
	assert.Equal(t, "req-3", got["request_id"])
	assert.Equal(t, "specific", got["call_site"])
}
