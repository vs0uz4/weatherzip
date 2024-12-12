package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"weatherzip/internal/domain"
	"weatherzip/internal/service/mock"
)

const (
	weatherServiceBaseURL  = "http://example.com/weather?apikey=%s&q=%s&lang=%s"
	weatherServiceApiKey   = "dummy-api-key"
	weatherServiceLanguage = "en"
)

func TestWeatherServiceCreateRequest(t *testing.T) {
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
			s := &WeatherService{}
			_, err := s.GetWeather(tt.inputURL)

			if err == nil || !strings.Contains(err.Error(), tt.expectErr) {
				t.Errorf("Expected error containing %q, got %v", tt.expectErr, err)
			}
		})
	}
}

func TestWeatherServiceExecuteRequest(t *testing.T) {
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

			service := WeatherService{
				HttpClient: mockClient,
				BaseURL:    weatherServiceBaseURL,
				ApiKey:     weatherServiceApiKey,
				Language:   weatherServiceLanguage,
			}

			_, err := service.GetWeather("valid-location")

			if err == nil || !strings.Contains(err.Error(), tt.expectErr) {
				t.Errorf("Expected error containing %q, got %q", tt.expectErr, err.Error())
			}
		})
	}
}

func TestWeatherServiceStatusCodeHandling(t *testing.T) {
	tests := []struct {
		name        string
		inputStatus int
		expectErr   string
	}{
		{
			name:        "Unexpected Status Code",
			inputStatus: 500,
			expectErr:   "unexpected status code: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.inputStatus != http.StatusOK && tt.inputStatus != http.StatusBadRequest {
				err = domain.NewUnexpectedStatusCodeError(tt.inputStatus)
			}

			if err == nil || !strings.Contains(err.Error(), tt.expectErr) {
				t.Errorf("Expected error containing %q, got %q", tt.expectErr, err.Error())
			}
		})
	}
}

func TestWeatherServiceDecodeResponse(t *testing.T) {
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
			var response domain.WeatherResponse
			err := json.Unmarshal([]byte(tt.inputResponse), &response)

			if err == nil || !strings.Contains(err.Error(), tt.expectErr) {
				t.Errorf("Expected error containing %q, got %q", tt.expectErr, err.Error())
			}
		})
	}
}

func TestWeatherServiceBadRequestHandling(t *testing.T) {
	tests := []struct {
		name      string
		inputCode int
		expectErr error
	}{
		{
			name:      "Unexpected Bad Request",
			inputCode: 0001,
			expectErr: domain.ErrUnexpectedBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mock.MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       io.NopCloser(strings.NewReader(fmt.Sprintf(`{"error":{"code":%d}}`, tt.inputCode))),
					}, nil
				},
			}

			service := WeatherService{
				HttpClient: mockClient,
				BaseURL:    weatherServiceBaseURL,
				ApiKey:     weatherServiceApiKey,
				Language:   weatherServiceLanguage,
			}
			_, err := service.GetWeather("invalid-location")

			if err == nil || err != domain.ErrUnexpectedBadRequest {
				t.Errorf("Expected error %v, got %v", domain.ErrUnexpectedBadRequest, err)
			}
		})
	}
}

func TestWeatherServiceUnexpectedStatusCode(t *testing.T) {
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

			service := WeatherService{
				HttpClient: mockClient,
				BaseURL:    weatherServiceBaseURL,
				ApiKey:     weatherServiceApiKey,
				Language:   weatherServiceLanguage,
			}
			_, err := service.GetWeather("valid-location")

			if err == nil || err.Error() != tt.expectErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectErr, err)
			}
		})
	}
}

func TestWeatherServiceGetWeatherData(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		inputLocation  string
		expectErr      error
		expectOutput   domain.WeatherResponse
	}{
		{
			name:           "Valid Location",
			mockResponse:   `{"location": {"name": "Cidade C", "region": "Região R", "country": "País P"}, "current": {"temp_c": 25.0, "temp_f": 77.0, "condition": {"text": "Sunny", "icon": "icon_url"}}}`,
			mockStatusCode: http.StatusOK,
			inputLocation:  "Cidade C",
			expectErr:      nil,
			expectOutput:   domain.WeatherResponse{Location: domain.LocationData{Name: "Cidade C", Region: "Região R", Country: "País P"}, Current: domain.CurrentWeather{TempC: 25.0, TempF: 77.0, TempK: 298.15, Condition: domain.WeatherCondition{Text: "Sunny", Icon: "icon_url"}}},
		},
		{
			name:           "Location Not Found",
			mockResponse:   `{"error": {"code": 1006,	"message": "No matching location found."}}`,
			mockStatusCode: http.StatusBadRequest,
			inputLocation:  "Unknown",
			expectErr:      domain.ErrLocationNotFound,
			expectOutput:   domain.WeatherResponse{},
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

			weatherService := NewWeatherService(mockServer.Client(), mockServer.URL+"?key=%s&q=%s&lang=%s&aqi=no", "APIKEY", "pt")
			encodedInputLocation := url.QueryEscape(tt.inputLocation)
			result, err := weatherService.GetWeather(encodedInputLocation)

			if !errors.Is(err, tt.expectErr) {
				t.Errorf("Expected error %v, got %v", tt.expectErr, err)
			}

			if result != tt.expectOutput {
				t.Errorf("Expected output %+v, got %+v", tt.expectOutput, result)
			}
		})
	}
}
