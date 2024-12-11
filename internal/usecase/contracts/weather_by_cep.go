package contracts

import "weatherzip/internal/domain"

type WeatherByCepUsecase interface {
	GetWeatherByCep(cep string) (domain.WeatherResponse, error)
}
