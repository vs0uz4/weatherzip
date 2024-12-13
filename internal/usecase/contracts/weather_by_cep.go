package contracts

import "github.com/vs0uz4/weatherzip/internal/domain"

type WeatherByCepUsecase interface {
	GetWeatherByCep(cep string) (domain.WeatherResponse, error)
}
