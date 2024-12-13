package mock

import "github.com/vs0uz4/weatherzip/internal/domain"

type MockCepService struct {
	GetLocationFunc func(string) (domain.CepResponse, error)
}

func (m *MockCepService) GetLocation(cep string) (domain.CepResponse, error) {
	return m.GetLocationFunc(cep)
}
