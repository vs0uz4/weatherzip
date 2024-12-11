package contracts

import "weatherzip/internal/domain"

type CepService interface {
	GetLocation(cep string) (domain.CepResponse, error)
}
