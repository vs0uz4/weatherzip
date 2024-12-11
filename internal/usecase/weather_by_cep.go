package usecase

import (
	"errors"
	"weatherzip/internal/domain"
	"weatherzip/internal/service/contracts"
)

type weatherByCepUsecase struct {
	CepService     contracts.CepService
	WeatherService contracts.WeatherService
}

func NewWeatherByCepUsecase(cepService contracts.CepService, weatherService contracts.WeatherService) *weatherByCepUsecase {
	return &weatherByCepUsecase{
		CepService:     cepService,
		WeatherService: weatherService,
	}
}

func (uc *weatherByCepUsecase) GetWeatherByCep(cep string) (domain.WeatherResponse, error) {
	if len(cep) != 8 || !isNumeric(cep) {
		return domain.WeatherResponse{}, errors.New("invalid zipcode")
	}

	location, err := uc.CepService.GetLocation(cep)
	if err != nil {
		return domain.WeatherResponse{}, err
	}

	weather, err := uc.WeatherService.GetWeather(location.Localidade)
	if err != nil {
		return domain.WeatherResponse{}, err
	}

	return weather, nil
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
