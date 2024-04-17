package gohttp_mock

import (
	"fmt"
	"net/http"

	"github.com/maxkopysow/go-httpclient/core"
)

// Provides a clean way to counfigure HTTP mocks
// based on combination between request method, URL and request body.
type Mock struct {
	Method      string
	Url         string
	RequestBody string

	ResponseBody       string
	ResponseStatusCode int
	Error              error
}

// Returns a Response object based on the mock configuration
func (m *Mock) GetResponse() (*core.Response, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	response := core.Response{
		StatusCode: m.ResponseStatusCode,
		Body:       []byte(m.ResponseBody),
		Status:     fmt.Sprintf("%d %s", m.ResponseStatusCode, http.StatusText(m.ResponseStatusCode)),
	}

	return &response, nil
}
