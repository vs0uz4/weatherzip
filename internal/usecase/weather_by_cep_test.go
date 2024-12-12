package usecase

import (
	"errors"
	"testing"

	"weatherzip/internal/domain"
	"weatherzip/internal/service/mock"
)

func TestGetWeatherByCep(t *testing.T) {
	tests := []struct {
		name           string
		inputCep       string
		mockCepSvc     func() *mock.MockCepService
		mockWeatherSvc func() *mock.MockWeatherService
		expectErr      error
		expectOutput   domain.WeatherResponse
	}{
		{
			name:     "Invalid CEP",
			inputCep: "123",
			mockCepSvc: func() *mock.MockCepService {
				return &mock.MockCepService{}
			},
			mockWeatherSvc: func() *mock.MockWeatherService {
				return &mock.MockWeatherService{}
			},
			expectErr: domain.ErrInvalidZipcode,
		},
		{
			name:     "CEP Not Found",
			inputCep: "99999999",
			mockCepSvc: func() *mock.MockCepService {
				return &mock.MockCepService{
					GetLocationFunc: func(cep string) (domain.CepResponse, error) {
						return domain.CepResponse{}, domain.ErrZipcodeNotFound
					},
				}
			},
			mockWeatherSvc: func() *mock.MockWeatherService {
				return &mock.MockWeatherService{}
			},
			expectErr: domain.ErrZipcodeNotFound,
		},
		{
			name:     "Weather Service Error",
			inputCep: "12345678",
			mockCepSvc: func() *mock.MockCepService {
				return &mock.MockCepService{
					GetLocationFunc: func(cep string) (domain.CepResponse, error) {
						return domain.CepResponse{
							Localidade: "City",
							Uf:         "State",
						}, nil
					},
				}
			},
			mockWeatherSvc: func() *mock.MockWeatherService {
				return &mock.MockWeatherService{
					GetWeatherFunc: func(location string) (domain.WeatherResponse, error) {
						return domain.WeatherResponse{}, domain.ErrWeatherService
					},
				}
			},
			expectErr: domain.ErrWeatherService,
		},
		{
			name:     "Success",
			inputCep: "12345678",
			mockCepSvc: func() *mock.MockCepService {
				return &mock.MockCepService{
					GetLocationFunc: func(cep string) (domain.CepResponse, error) {
						return domain.CepResponse{
							Localidade: "City",
							Uf:         "State",
						}, nil
					},
				}
			},
			mockWeatherSvc: func() *mock.MockWeatherService {
				return &mock.MockWeatherService{
					GetWeatherFunc: func(location string) (domain.WeatherResponse, error) {
						return domain.WeatherResponse{
							Current: domain.CurrentWeather{
								TempC: 25.0,
							},
						}, nil
					},
				}
			},
			expectOutput: domain.WeatherResponse{
				Current: domain.CurrentWeather{
					TempC: 25.0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := weatherByCepUsecase{
				CepService:     tt.mockCepSvc(),
				WeatherService: tt.mockWeatherSvc(),
			}

			result, err := usecase.GetWeatherByCep(tt.inputCep)

			if !errors.Is(err, tt.expectErr) {
				t.Errorf("Expected error %v, got %v", tt.expectErr, err)
			}

			if result != tt.expectOutput {
				t.Errorf("Expected output %+v, got %+v", tt.expectOutput, result)
			}
		})
	}
}
