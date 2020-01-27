package guppy

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// Client is an HTTP Client.
//
// It exposes a single Do method that allows to perform HTTP requests.
type Client struct {
	// name is the name that will be used for circuit-breaking
	name string
	// userAgent is the userAgent that will be passed as a header to all requests.
	userAgent string
	// client is the http.Client that will be used for HTTP transport
	client *http.Client

	// tracer is the (optional) opentracing.Tracer that will be used to inject tracing information
	// in outgoing http.Request headers
	tracer *opentracing.Tracer
}

// Do executes a given http.Request and returns the http.Response or an error.
//
// This method ony decorates the request before executing it using Go's internal http.HTTPClient.
//
// For general documentation of response and error behavior, see http.HTTPClient.Do
func (c *Client) Do(ctx context.Context, operationName string, r *http.Request) (*http.Response, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		err := opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			return nil, errors.Wrap(err, "unable to inject opentracing headers")
		}
	}

	r.Header.Add("User-Agent", c.userAgent)

	return c.client.Do(r)
}