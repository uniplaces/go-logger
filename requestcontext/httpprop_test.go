package requestcontext_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/uniplaces/go-logger/requestcontext"
)

func TestRoundTripperInjectsHeader(t *testing.T) {
	t.Parallel()

	var got string

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("X-Request-Id")
	}))
	defer srv.Close()

	client := requestcontext.WrapClient(&http.Client{})

	req, err := http.NewRequestWithContext(
		requestcontext.WithID(context.Background(), "req-out-1"),
		http.MethodGet,
		srv.URL,
		nil,
	)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, "req-out-1", got)
}

func TestRoundTripperNoIDNoHeader(t *testing.T) {
	t.Parallel()

	var got string

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("X-Request-Id")
	}))
	defer srv.Close()

	client := requestcontext.WrapClient(&http.Client{})

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, "", got)
}

func TestRoundTripperPreservesExistingHeader(t *testing.T) {
	t.Parallel()

	var got string

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("X-Request-Id")
	}))
	defer srv.Close()

	client := requestcontext.WrapClient(&http.Client{})

	req, err := http.NewRequestWithContext(
		requestcontext.WithID(context.Background(), "from-ctx"),
		http.MethodGet,
		srv.URL,
		nil,
	)
	require.NoError(t, err)
	req.Header.Set("X-Request-Id", "explicit-override")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, "explicit-override", got)
}

func TestRoundTripperDoesNotMutateRequest(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	defer srv.Close()

	client := requestcontext.WrapClient(&http.Client{})

	req, err := http.NewRequestWithContext(
		requestcontext.WithID(context.Background(), "req-no-mutate"),
		http.MethodGet,
		srv.URL,
		nil,
	)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, "", req.Header.Get("X-Request-Id"),
		"the caller's request header must not be mutated")
}

func TestWrapClientReturnsSameClient(t *testing.T) {
	t.Parallel()

	c := &http.Client{}

	got := requestcontext.WrapClient(c)

	assert.Same(t, c, got)
	assert.NotNil(t, c.Transport, "WrapClient must set Transport")
}

func TestInjectHTTPHeaderSetsWhenIDPresent(t *testing.T) {
	t.Parallel()

	h := http.Header{}

	requestcontext.InjectHTTPHeader(requestcontext.WithID(context.Background(), "abc"), h)

	assert.Equal(t, "abc", h.Get("X-Request-Id"))
}

func TestInjectHTTPHeaderNoopWhenIDAbsent(t *testing.T) {
	t.Parallel()

	h := http.Header{}

	requestcontext.InjectHTTPHeader(context.Background(), h)

	assert.Empty(t, h.Get("X-Request-Id"))
}

func TestRoundTripperNilBaseUsesDefaultTransport(t *testing.T) {
	t.Parallel()

	rt := requestcontext.RoundTripper(nil)

	assert.NotNil(t, rt)
}
