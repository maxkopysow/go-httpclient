package gohttp

import (
	"fmt"
	"net/http"
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
func (m *Mock) GetResponse() (*Response, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	response := Response{
		statusCode: m.ResponseStatusCode,
		body:       []byte(m.ResponseBody),
		status:     fmt.Sprintf("%d %s", m.ResponseStatusCode, http.StatusText(m.ResponseStatusCode)),
	}

	return &response, nil
}
