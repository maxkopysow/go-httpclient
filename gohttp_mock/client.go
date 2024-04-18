package gohttp_mock

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type httpClientMock struct{}

func (c *httpClientMock) Do(request *http.Request) (*http.Response, error) {
	requestBody, err := request.GetBody()
	if err != nil {
		return nil, err
	}

	defer requestBody.Close()

	body, err := io.ReadAll(requestBody)

	if err != nil {
		return nil, err
	}

	mock := MockupServer.mocks[MockupServer.getMockKey(request.Method, request.URL.String(), string(body))]

	if mock != nil {

		if mock.Error != nil {
			return nil, mock.Error
		}
		response := http.Response{
			StatusCode:    mock.ResponseStatusCode,
			Body:          io.NopCloser(strings.NewReader(mock.ResponseBody)),
			ContentLength: int64(len(mock.ResponseBody)),
		}

		return &response, nil
	}

	return nil, errors.New(fmt.Sprintf("no mock matching %s from '%s' with given body", request.Method, request.URL.String()))
}
