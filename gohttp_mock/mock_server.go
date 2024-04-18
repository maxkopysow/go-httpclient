package gohttp_mock

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"sync"

	"github.com/maxkopysow/go-httpclient/core"
)

var (
	MockupServer = mockServer{
		mocks:      make(map[string]*Mock),
		httpClient: &httpClientMock{},
	}
)

type mockServer struct {
	enabled      bool
	serverMutext sync.Mutex

	httpClient core.HttpClient

	mocks map[string]*Mock
}

func (m *mockServer) Start() {
	m.serverMutext.Lock()
	defer m.serverMutext.Unlock()

	m.enabled = true
}

func (m *mockServer) Stop() {
	m.serverMutext.Lock()
	defer m.serverMutext.Unlock()
	m.enabled = false
}

func (m *mockServer) AddMock(mock Mock) {
	m.serverMutext.Lock()
	defer m.serverMutext.Unlock()

	key := m.getMockKey(mock.Method, mock.Url, mock.RequestBody)
	m.mocks[key] = &mock

}

func (m *mockServer) DeleteMocks() {
	m.serverMutext.Lock()
	defer m.serverMutext.Unlock()

	m.mocks = make(map[string]*Mock)
}

func (m *mockServer) IsEnabled() bool {
	return m.enabled
}

func (m *mockServer) GetMockedClient() core.HttpClient {
	return m.httpClient
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
