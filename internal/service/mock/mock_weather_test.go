package mock

import (
	"testing"

	"github.com/vs0uz4/weatherzip/internal/domain"
)

func TestMockWeatherService(t *testing.T) {
	mock := MockWeatherService{
		GetWeatherFunc: func(location string) (domain.WeatherResponse, error) {
			if location == "Valid Location" {
				return domain.WeatherResponse{Current: domain.CurrentWeather{TempC: 25.0}}, nil
			}
			return domain.WeatherResponse{}, domain.ErrUnexpectedBadRequest
		},
	}

	response, err := mock.GetWeather("Valid Location")
	if response.Current.TempC != 25.0 || err != nil {
		t.Errorf("Expected TempC: 25.0, got: %v, err: %v", response.Current.TempC, err)
	}

	_, err = mock.GetWeather("Invalid Location")
	if err != domain.ErrUnexpectedBadRequest {
		t.Errorf("Expected error: %v, got: %v", domain.ErrUnexpectedBadRequest, err)
	}
}
