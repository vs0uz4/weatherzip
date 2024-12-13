package service

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vs0uz4/weatherzip/internal/domain"
	"github.com/vs0uz4/weatherzip/internal/service/mock"
)

const (
	cepServiceBaseURL = "https://example.com/ws/%s/json/"
)

func TestCepServiceCreateRequest(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		expectErr string
	}{
		{
			name:      "Request Creation Error",
			inputURL:  "",
			expectErr: "failed to create request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CepService{}
			_, err := s.GetLocation(tt.inputURL)

			if err == nil || !strings.Contains(err.Error(), tt.expectErr) {
				t.Errorf("Expected error containing %q, got %v", tt.expectErr, err)
			}
		})
	}
}

func TestCepServiceExecuteRequest(t *testing.T) {
	tests := []struct {
		name      string
		expectErr string
	}{
		{
			name:      "Request Execution Error",
			expectErr: "failed to make request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mock.MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("network error")
				},
			}

			service := CepService{
				HttpClient: mockClient,
				BaseURL:    cepServiceBaseURL,
			}

			_, err := service.GetLocation("12345678")

			if err == nil || !strings.Contains(err.Error(), tt.expectErr) {
				t.Errorf("Expected error containing %q, got %q", tt.expectErr, err.Error())
			}
		})
	}
}

func TestCepServiceUnexpectedStatusCode(t *testing.T) {
	tests := []struct {
		name      string
		inputCode int
		expectErr error
	}{
		{
			name:      "Unexpected Status Code",
			inputCode: 500,
			expectErr: domain.NewUnexpectedStatusCodeError(500),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mock.MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(strings.NewReader("")),
					}, nil
				},
			}

			service := CepService{
				HttpClient: mockClient,
				BaseURL:    cepServiceBaseURL,
			}
			_, err := service.GetLocation("12345678")

			if err == nil || err.Error() != tt.expectErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectErr, err)
			}
		})
	}
}

func TestCepServiceDecodeResponse(t *testing.T) {
	tests := []struct {
		name          string
		inputResponse string
		expectErr     string
	}{
		{
			name:          "Failed to Decode Response",
			inputResponse: "invalid_json",
			expectErr:     "invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response domain.CepResponse
			err := json.Unmarshal([]byte(tt.inputResponse), &response)

			if err == nil || !strings.Contains(err.Error(), tt.expectErr) {
				t.Errorf("Expected error containing %q, got %q", tt.expectErr, err.Error())
			}
		})
	}
}

func TestCepServiceDecodeResponseError(t *testing.T) {
	mockClient := &mock.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`invalid_json`)),
			}, nil
		},
	}

	service := CepService{
		HttpClient: mockClient,
		BaseURL:    cepServiceBaseURL,
	}

	_, err := service.GetLocation("12345678")

	if err == nil || !strings.Contains(err.Error(), "failed to decode response") {
		t.Errorf("Expected error containing %q, got %q", "failed to decode response", err.Error())
	}
}

func TestCepServicePopulateMapError(t *testing.T) {
	mockClient := &mock.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"cep": "21831200", "logradouro": "Rua A", "bairro": "Bairro C", "localidade": null, "uf": "SP", "estado": "São Paulo"}`)),
			}, nil
		},
	}

	service := CepService{
		HttpClient: mockClient,
		BaseURL:    cepServiceBaseURL,
	}
	_, err := service.GetLocation("12345678")

	expectedError := "failed to map response"
	if err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing %q, got %q", expectedError, err.Error())
	}
}

func TestCepServiceGetCepData(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		inputCep       string
		expectErr      error
		expectOutput   domain.CepResponse
	}{
		{
			name:           "Valid CEP",
			mockResponse:   `{"cep": "12345678", "logradouro": "Rua A", "bairro": "Bairro B", "localidade": "Cidade C", "uf": "SP"}`,
			mockStatusCode: http.StatusOK,
			inputCep:       "12345678",
			expectErr:      nil,
			expectOutput:   domain.CepResponse{Cep: "12345678", Logradouro: "Rua A", Bairro: "Bairro B", Localidade: "Cidade C", Uf: "SP"},
		},
		{
			name:           "CEP Not Found",
			mockResponse:   `{"erro": "true"}`,
			mockStatusCode: http.StatusOK,
			inputCep:       "00000000",
			expectErr:      domain.ErrZipcodeNotFound,
			expectOutput:   domain.CepResponse{},
		},
		{
			name:           "Invalid CEP Format",
			mockResponse:   "",
			mockStatusCode: http.StatusBadRequest,
			inputCep:       "1234",
			expectErr:      domain.ErrInvalidZipcode,
			expectOutput:   domain.CepResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				if _, err := w.Write([]byte(tt.mockResponse)); err != nil {
					t.Fatalf("Failed to write mock response: %v", err)
				}
			}))
			defer mockServer.Close()

			cepService := NewCepService(mockServer.Client(), mockServer.URL+"/%s")
			result, err := cepService.GetLocation(tt.inputCep)

			if !errors.Is(err, tt.expectErr) {
				t.Errorf("Expected error %v, got %v", tt.expectErr, err)
			}

			if result != tt.expectOutput {
				t.Errorf("Expected output %+v, got %+v", tt.expectOutput, result)
			}
		})
	}
}
