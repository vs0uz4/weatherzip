package contracts

import "weatherzip/internal/domain"

type WeatherService interface {
	GetWeather(location string) (domain.WeatherResponse, error)
}
