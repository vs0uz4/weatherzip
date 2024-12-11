package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"weatherzip/internal/domain"
)

func TestCepServiceGetCepData(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		inputCep       string
		expectErr      error
		expectOutput   domain.CepResponse
	}{
		{
			name:           "Valid CEP",
			mockResponse:   `{"cep": "12345678", "logradouro": "Rua A", "bairro": "Bairro B", "localidade": "Cidade C", "uf": "SP"}`,
			mockStatusCode: http.StatusOK,
			inputCep:       "12345678",
			expectErr:      nil,
			expectOutput:   domain.CepResponse{Cep: "12345678", Logradouro: "Rua A", Bairro: "Bairro B", Localidade: "Cidade C", Uf: "SP"},
		},
		{
			name:           "CEP Not Found",
			mockResponse:   `{"erro": "true"}`,
			mockStatusCode: http.StatusOK,
			inputCep:       "00000000",
			expectErr:      domain.ErrCepNotFound,
			expectOutput:   domain.CepResponse{},
		},
		{
			name:           "Invalid CEP Format",
			mockResponse:   "",
			mockStatusCode: http.StatusBadRequest,
			inputCep:       "1234",
			expectErr:      domain.ErrCepIsInvalid,
			expectOutput:   domain.CepResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()

			cepService := NewCepService(mockServer.Client(), mockServer.URL+"/%s")
			result, err := cepService.GetLocation(tt.inputCep)

			if !errors.Is(err, tt.expectErr) {
				t.Errorf("Expected error %v, got %v", tt.expectErr, err)
			}

			if result != tt.expectOutput {
				t.Errorf("Expected output %+v, got %+v", tt.expectOutput, result)
			}
		})
	}
}
