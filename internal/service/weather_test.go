package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"weatherzip/internal/domain"
)

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
