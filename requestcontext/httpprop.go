package requestcontext

import (
	"context"
	"net/http"
)

type roundTripper struct {
	base http.RoundTripper
}

// RoundTripper wraps base so every outbound request carries HeaderName when
// its context has a request id and the header is not already set. If base
// is nil, http.DefaultTransport is used. The wrapper clones the request
// before mutating headers — never modifies the caller's *http.Request.
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

// WrapClient sets client.Transport = RoundTripper(client.Transport) and
// returns the same *http.Client. Use at HTTP-client construction time.
func WrapClient(client *http.Client) *http.Client {
	client.Transport = RoundTripper(client.Transport)

	return client
}

// InjectHTTPHeader sets HeaderName on h when ctx has a request id. Escape
// hatch for code paths that build a raw *http.Request without going
// through a wrapped client.
func InjectHTTPHeader(ctx context.Context, h http.Header) {
	if id := ID(ctx); id != "" {
		h.Set(HeaderName, id)
	}
}
