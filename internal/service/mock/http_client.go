package mock

import (
	"errors"
	"net/http"
)

type MockHTTPClient struct{}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("failed to make request")
}
