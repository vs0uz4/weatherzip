package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"weatherzip/internal/domain"
	"weatherzip/internal/service/contracts"
)

var _ contracts.CepService = (*CepService)(nil)

type CepService struct {
	HttpClient contracts.HttpClient
	BaseURL    string
}

func NewCepService(client *http.Client, baseURL string) *CepService {
	return &CepService{
		HttpClient: client,
		BaseURL:    baseURL,
	}
}

func (s *CepService) GetLocation(cep string) (domain.CepResponse, error) {
	var response domain.CepResponse
	var raw map[string]interface{}

	url := fmt.Sprintf(s.BaseURL, cep)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := s.HttpClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		return response, domain.ErrCepIsInvalid
	}

	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return response, fmt.Errorf("failed to decode response: %w", err)
	}

	if erro, ok := raw["erro"].(bool); ok && erro {
		return response, domain.ErrCepNotFound
	}

	if err := response.PopulateFromMap(raw); err != nil {
		if err.Error() == "cep not found" {
			return response, domain.ErrCepNotFound
		}
		return response, fmt.Errorf("failed to map response: %w", err)
	}

	return response, nil
}
