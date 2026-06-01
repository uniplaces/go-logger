package requestcontext

import (
	"context"
	"net/http"
)

type roundTripper struct {
	base http.RoundTripper
}

// RoundTripper wraps base so outbound requests carry HeaderName from ctx; nil base uses http.DefaultTransport.
func RoundTripper(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}

	return roundTripper{base: base}
}

// RoundTrip implements http.RoundTripper.
func (r roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if id := ID(req.Context()); id != "" && req.Header.Get(HeaderName) == "" {
		req = req.Clone(req.Context())
		req.Header.Set(HeaderName, id)
	}

	return r.base.RoundTrip(req)
}

// WrapClient returns a shallow copy of client with its Transport wrapped; the input is never mutated.
func WrapClient(client *http.Client) *http.Client {
	if client == nil {
		client = http.DefaultClient
	}

	wrapped := *client
	wrapped.Transport = RoundTripper(client.Transport)

	return &wrapped
}

// InjectHTTPHeader sets HeaderName on header when ctx has a request id.
func InjectHTTPHeader(ctx context.Context, header http.Header) {
	if id := ID(ctx); id != "" {
		header.Set(HeaderName, id)
	}
}
