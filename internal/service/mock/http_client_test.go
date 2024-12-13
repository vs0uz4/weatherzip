package mock

import (
	"net/http"
	"testing"

	"github.com/vs0uz4/weatherzip/internal/domain"
)

func TestMockHTTPClientDo(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == "http://example.com" {
				return &http.Response{
					StatusCode: http.StatusOK,
				}, nil
			}
			return nil, domain.ErrUnexpectedUrl
		},
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	res, err := mockClient.Do(req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}
}
