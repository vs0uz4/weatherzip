package mock

import "weatherzip/internal/domain"

type MockWeatherService struct {
	GetWeatherFunc func(string) (domain.WeatherResponse, error)
}

func (m *MockWeatherService) GetWeather(location string) (domain.WeatherResponse, error) {
	return m.GetWeatherFunc(location)
}
