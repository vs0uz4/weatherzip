package mock

import (
	"testing"

	"github.com/vs0uz4/weatherzip/internal/domain"
)

func TestMockCepService(t *testing.T) {
	mock := MockCepService{
		GetLocationFunc: func(cep string) (domain.CepResponse, error) {
			if cep == "12345678" {
				return domain.CepResponse{Cep: "12345678"}, nil
			}
			return domain.CepResponse{}, domain.ErrZipcodeNotFound
		},
	}

	response, err := mock.GetLocation("12345678")
	if response.Cep != "12345678" || err != nil {
		t.Errorf("Expected Cep: 12345678, got: %v, err: %v", response.Cep, err)
	}

	_, err = mock.GetLocation("00000000")
	if err != domain.ErrZipcodeNotFound {
		t.Errorf("Expected error: %v, got: %v", domain.ErrZipcodeNotFound, err)
	}
}
