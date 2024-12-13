package contracts

import "github.com/vs0uz4/weatherzip/internal/domain"

type CepService interface {
	GetLocation(cep string) (domain.CepResponse, error)
}
