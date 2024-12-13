package web

import (
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"weatherzip/internal/domain"
	"weatherzip/internal/usecase/mock"
)

func TestWeatherHandler(t *testing.T) {
	tests := []struct {
		name           string
		inputCEP       string
		mockUsecase    func() *mock.MockWeatherByCepUsecase
		expectedStatus int
		expectedBody   string
		expectedError  string
	}{
		{
			name:     "CEP Inválido",
			inputCEP: "123",
			mockUsecase: func() *mock.MockWeatherByCepUsecase {
				return &mock.MockWeatherByCepUsecase{
					GetWeatherByCepFunc: func(cep string) (domain.WeatherResponse, error) {
						return domain.WeatherResponse{}, domain.ErrInvalidZipcode
					},
				}
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "invalid zipcode",
			expectedError:  "Invalid zipcode",
		},
		{
			name:     "CEP Não Encontrado",
			inputCEP: "99999999",
			mockUsecase: func() *mock.MockWeatherByCepUsecase {
				return &mock.MockWeatherByCepUsecase{
					GetWeatherByCepFunc: func(cep string) (domain.WeatherResponse, error) {
						return domain.WeatherResponse{}, domain.ErrZipcodeNotFound
					},
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "can not find zipcode",
			expectedError:  "Zipcode not found",
		},
		{
			name:     "Erro no Serviço de Clima",
			inputCEP: "12345678",
			mockUsecase: func() *mock.MockWeatherByCepUsecase {
				return &mock.MockWeatherByCepUsecase{
					GetWeatherByCepFunc: func(cep string) (domain.WeatherResponse, error) {
						return domain.WeatherResponse{}, domain.ErrWeatherService
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
			expectedError:  "Internal server error",
		},
		{
			name:     "Sucesso",
			inputCEP: "12345678",
			mockUsecase: func() *mock.MockWeatherByCepUsecase {
				return &mock.MockWeatherByCepUsecase{
					GetWeatherByCepFunc: func(cep string) (domain.WeatherResponse, error) {
						return domain.WeatherResponse{
							Current: domain.CurrentWeather{
								TempC: 25.0,
								TempF: 77.0,
								TempK: 298.15,
							},
						}, nil
					},
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"temp_C":25,"temp_F":77,"temp_K":298.15}`,
		},
		{
			name:     "Erro no JSON Encode",
			inputCEP: "12345678",
			mockUsecase: func() *mock.MockWeatherByCepUsecase {
				return &mock.MockWeatherByCepUsecase{
					GetWeatherByCepFunc: func(cep string) (domain.WeatherResponse, error) {
						return domain.WeatherResponse{
							Current: domain.CurrentWeather{
								TempC: math.NaN(),
							},
						}, nil
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := tt.mockUsecase()
			handler := NewWeatherHandler(mockUsecase)

			req := httptest.NewRequest(http.MethodGet, "/weather/"+tt.inputCEP, nil)
			w := httptest.NewRecorder()
			handler.GetWeatherByCep(w, req)

			resp := w.Result()
			body := strings.TrimSpace(w.Body.String())

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if body != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

func TestNewWeatherHandlerInitialization(t *testing.T) {
	mockUsecase := &mock.MockWeatherByCepUsecase{}
	handler := NewWeatherHandler(mockUsecase)

	if handler.Usecase != mockUsecase {
		t.Errorf("Expected usecase %v, got %v", mockUsecase, handler.Usecase)
	}
}
