package guppy

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
)

var (
	ErrCommandNameNotSpecified = errors.New("must specify operationName and/or serviceName/upstreamServiceName")
)

type Option func(c *Configuration) error

type Configuration struct {
	serviceName         string
	upstreamServiceName string
	operationName       string
	userAgent           string

	// HTTP transport timeouts
	idleConnTimeout time.Duration
	tlsHandshakeTimeout time.Duration
	responseHeaderTimeout time.Duration

	// global HTTP timeout
	httpTimeout time.Duration

	// circuit breaker
	circuitTimeout time.Duration
	circuitOpenTimeout time.Duration

	httpClient *http.Client

	tracer *opentracing.Tracer
}

func ServiceName(serviceName string) Option {
	return func(c *Configuration) error {
		c.serviceName = serviceName
		return nil
	}
}

func UpstreamServiceName(upstreamServiceName string) Option {
	return func(c *Configuration) error {
		c.upstreamServiceName = upstreamServiceName
		return nil
	}
}

func OperationName(commandName string) Option {
	return func(c *Configuration) error {
		c.operationName = commandName
		return nil
	}
}

func UserAgent(userAgent string) Option {
	return func(c *Configuration) error {
		c.userAgent = userAgent
		return nil
	}
}

func IdleConnTimeout(idleConnTimeout time.Duration) Option {
	return func(c *Configuration) error {
		c.idleConnTimeout = idleConnTimeout
		return nil
	}
}

func TlsHandshakeTimeout(tlsHandshakeTimeout time.Duration) Option {
	return func(c *Configuration) error {
		c.tlsHandshakeTimeout = tlsHandshakeTimeout
		return nil
	}
}

func ResponseHeaderTimeout(responseHeaderTimeout time.Duration) Option {
	return func(c *Configuration) error {
		c.responseHeaderTimeout = responseHeaderTimeout
		return nil
	}
}

func HttpTimeout(httpTimeout time.Duration) Option {
	return func(c *Configuration) error {
		c.httpTimeout = httpTimeout
		return nil
	}
}

func CircuitTimeout(circuitTimeout time.Duration) Option {
	return func(c *Configuration) error {
		c.circuitTimeout = circuitTimeout
		return nil
	}
}

func CircuitOpenTimeout(circuitOpenTimeout time.Duration) Option {
	return func(c *Configuration) error {
		c.circuitOpenTimeout = circuitOpenTimeout
		return nil
	}
}

func HTTPClient(client *http.Client) Option {
	return func(c *Configuration) error {
		c.httpClient = client
		return nil
	}
}

func Tracer(tracer *opentracing.Tracer) Option {
	return func(c *Configuration) error {
		c.tracer = tracer
		return nil
	}
}

func defaults() *Configuration {
	return &Configuration{
		httpTimeout:    500 * time.Millisecond,
		circuitTimeout: 500 * time.Millisecond,
		circuitOpenTimeout: time.Minute,
	}
}

func NewConfiguration(opts ...Option) (*Configuration, error) {
	o := defaults()
	
	for _, option := range opts {
		option(o)
	}

	if o.operationName == "" {
		if o.serviceName == "" || o.upstreamServiceName == "" {
			return nil, ErrCommandNameNotSpecified
		}
		o.operationName = fmt.Sprintf("%v/%v", o.serviceName, o.upstreamServiceName)
	}
	if o.serviceName == "" || o.upstreamServiceName == "" {
		return nil, ErrCommandNameNotSpecified
	}

	return o, nil
}

func (o *Configuration) Build() (*Client, error) {
	return nil, nil
}