package gohttp

import (
	"net/http"
	"time"
)

type clientBuilder struct {
	headers            http.Header
	maxIdleConnections int
	connectionTimeout  time.Duration
	responseTimeout    time.Duration
	disableTimeouts    bool
	client             *http.Client
	userAgent          string
}

type ClientBuilder interface {
	SetHeaders(headers http.Header) ClientBuilder
	SetConnectionTimeout(timeout time.Duration) ClientBuilder
	SetResponseTimeout(timeout time.Duration) ClientBuilder
	SetMaxIdleConnections(connections int) ClientBuilder
	DisableTimeouts(b bool) ClientBuilder
	SetHttpClient(c *http.Client) ClientBuilder
	SetUserAgent(userAgent string) ClientBuilder
	Build() Client
}

func NewBuilder() ClientBuilder {

	builder := &clientBuilder{}
	return builder
}

func (c *clientBuilder) Build() Client {

	client := httpClient{
		builder: c,
	}

	return &client
}

func (c *clientBuilder) DisableTimeouts(disable bool) ClientBuilder {
	c.disableTimeouts = disable
	return c
}
func (c *clientBuilder) SetHeaders(headers http.Header) ClientBuilder {
	c.headers = headers
	return c
}

func (c *clientBuilder) SetConnectionTimeout(timeout time.Duration) ClientBuilder {
	c.connectionTimeout = timeout
	return c
}

func (c *clientBuilder) SetResponseTimeout(timeout time.Duration) ClientBuilder {
	c.responseTimeout = timeout
	return c
}

func (c *clientBuilder) SetMaxIdleConnections(connections int) ClientBuilder {
	c.maxIdleConnections = connections
	return c
}

func (c *clientBuilder) SetHttpClient(client *http.Client) ClientBuilder {
	c.client = client
	return c
}

func (c *clientBuilder) SetUserAgent(userAgent string) ClientBuilder {
	c.userAgent = userAgent
	return c
}
