package usecase

import (
	"errors"
	"testing"

	"github.com/vs0uz4/weatherzip/internal/domain"
	"github.com/vs0uz4/weatherzip/internal/service/mock"
)

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{input: "123456", expected: true},
		{input: "123abc", expected: false},
		{input: "", expected: true},
		{input: "123 456", expected: false},
		{input: "!@#$%", expected: false},
	}

	for _, tt := range tests {
		result := isNumeric(tt.input)
		if result != tt.expected {
			t.Errorf("For input %q, expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestNewWeatherByCepUsecase(t *testing.T) {
	mockCepSvc := &mock.MockCepService{}
	mockWeatherSvc := &mock.MockWeatherService{}

	usecase := NewWeatherByCepUsecase(mockCepSvc, mockWeatherSvc)

	if usecase.CepService != mockCepSvc {
		t.Errorf("Expected CepService to be %v, got %v", mockCepSvc, usecase.CepService)
	}
	if usecase.WeatherService != mockWeatherSvc {
		t.Errorf("Expected WeatherService to be %v, got %v", mockWeatherSvc, usecase.WeatherService)
	}
}

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
