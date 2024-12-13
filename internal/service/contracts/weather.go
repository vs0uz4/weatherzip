package contracts

import "github.com/vs0uz4/weatherzip/internal/domain"

type WeatherService interface {
	GetWeather(location string) (domain.WeatherResponse, error)
}
