package gohttp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/maxkopysow/go-httpclient/core"
	"github.com/maxkopysow/go-httpclient/gohttp_mock"

	"github.com/maxkopysow/go-httpclient/gomime"
)

const (
	defaultMaxIdleConnenctions = 5
	defaultResponseTimeout     = 5 * time.Second
	defaultConnectionTimeout   = 1 * time.Second
)

func (c *httpClient) getRequestBody(contentType string, body interface{}) ([]byte, error) {
	if body == nil {
		return nil, nil
	}

	switch strings.ToLower(contentType) {
	case gomime.ContentTypeJson:
		return json.Marshal(body)
	case gomime.ContentTypeXML:
		return xml.Marshal(body)
	default:
		return json.Marshal(body)
	}
}

func (c *httpClient) getHttpClient() core.HttpClient {

	if gohttp_mock.MockupServer.IsEnabled() {
		return gohttp_mock.MockupServer.GetMockedClient()
	}

	c.clientOnce.Do(func() {
		if c.builder.client != nil {
			c.client = c.builder.client
			return
		}

		c.client = &http.Client{
			Timeout: c.getConnectionTimeout() + c.getResponseTimeout(),
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   c.getMaxIdleConnections(),
				ResponseHeaderTimeout: c.getResponseTimeout(),
				DialContext: (&net.Dialer{
					Timeout: c.getConnectionTimeout(),
				}).DialContext,
			},
		}

	})

	return c.client

}

func (c *httpClient) do(method string, url string, headers http.Header, body interface{}) (*core.Response, error) {

	fullHeaders := c.getRequestHeaders(headers)

	requestBody, err := c.getRequestBody(fullHeaders.Get("Content-Type"), body)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))

	if err != nil {
		errorText := "unable to create a new request"
		return nil, errors.New(errorText)

	}

	request.Header = fullHeaders

	client := c.getHttpClient()

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err

	}

	finalResponse := core.Response{
		Status:     response.Status,
		StatusCode: response.StatusCode,
		Headers:    response.Header,
		Body:       responseBody,
	}

	return &finalResponse, nil
}

func (c *httpClient) getMaxIdleConnections() int {
	if c.builder.connectionTimeout > 0 {
		return c.builder.maxIdleConnections
	}

	return defaultMaxIdleConnenctions
}

func (c *httpClient) getResponseTimeout() time.Duration {
	if c.builder.responseTimeout > 0 {
		return c.builder.responseTimeout
	}

	if c.builder.disableTimeouts {
		return 0
	}

	return defaultResponseTimeout
}
func (c *httpClient) getConnectionTimeout() time.Duration {
	if c.builder.connectionTimeout > 0 {
		return c.builder.connectionTimeout
	}

	if c.builder.disableTimeouts {
		return 0
	}

	return defaultConnectionTimeout
}
