package domain

import (
	"errors"
	"testing"
)

func TestCepResponsePopulateFromMap(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		expectErr error
		output    CepResponse
	}{
		{
			name:      "Valid Data",
			input:     map[string]interface{}{"cep": "12345678", "logradouro": "Rua A", "bairro": "Bairro B", "localidade": "Cidade C", "uf": "SP"},
			output:    CepResponse{Cep: "12345678", Logradouro: "Rua A", Bairro: "Bairro B", Localidade: "Cidade C", Uf: "SP"},
			expectErr: nil,
		},
		{
			name:      "CEP Not Found",
			input:     map[string]interface{}{"erro": "true"},
			expectErr: ErrCepNotFound,
			output:    CepResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result CepResponse
			err := result.PopulateFromMap(tt.input)

			if !errors.Is(err, tt.expectErr) {
				t.Errorf("Expected error %v, got %v", tt.expectErr, err)
			}

			if result != tt.output {
				t.Errorf("Expected output %+v, got %+v", tt.output, result)
			}
		})
	}
}
