package gohttp

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
)

var (
	mockupServer = mockServer{
		mocks: make(map[string]*Mock),
	}
)

type mockServer struct {
	enabled      bool
	serverMutext sync.Mutex

	mocks map[string]*Mock
}

func StartMockServer() {
	mockupServer.serverMutext.Lock()
	defer mockupServer.serverMutext.Unlock()

	mockupServer.enabled = true
}

func StopMockServer() {
	mockupServer.serverMutext.Lock()
	defer mockupServer.serverMutext.Unlock()
	mockupServer.enabled = false
}

func AddMock(mock Mock) {
	mockupServer.serverMutext.Lock()
	defer mockupServer.serverMutext.Unlock()

	key := mockupServer.getMockKey(mock.Method, mock.Url, mock.RequestBody)
	mockupServer.mocks[key] = &mock

}

func FlushMocks() {
	mockupServer.serverMutext.Lock()
	defer mockupServer.serverMutext.Unlock()

	mockupServer.mocks = make(map[string]*Mock)
}

func (m *mockServer) getMockKey(method, url, body string) string {
	hasher := md5.New()
	hasher.Write([]byte(method + url + m.cleanBody(body)))
	key := hex.EncodeToString(hasher.Sum(nil))
	return key
}

func (m *mockServer) cleanBody(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}

	body = strings.ReplaceAll(body, "\t", "")
	body = strings.ReplaceAll(body, "\n", "")
	return body

}

func (m *mockServer) getMock(method, url, body string) *Mock {
	if !m.enabled {
		return nil
	}

	if mock := m.mocks[m.getMockKey(method, url, body)]; mock != nil {
		return mock
	}

	return &Mock{
		Error: fmt.Errorf("no mock matching %s from '%s' with given body", method, url),
	}
}
