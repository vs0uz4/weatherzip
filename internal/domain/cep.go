package domain

type CepResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
	Estado     string `json:"estado,omitempty"`
	Regiao     string `json:"regiao,omitempty"`
}

func (c *CepResponse) PopulateFromMap(data map[string]interface{}) error {
	if erro, ok := data["erro"].(string); ok && erro == "true" {
		return ErrCepNotFound
	}

	if cep, ok := data["cep"].(string); ok {
		c.Cep = cep
	} else {
		return ErrInvalidZipCodeData
	}

	if logradouro, ok := data["logradouro"].(string); ok {
		c.Logradouro = logradouro
	} else {
		return ErrInvalidStreetData
	}

	if bairro, ok := data["bairro"].(string); ok {
		c.Bairro = bairro
	} else {
		return ErrInvalidNeighborhoodData
	}

	if localidade, ok := data["localidade"].(string); ok {
		c.Localidade = localidade
	} else {
		return ErrInvalidLocationData
	}

	if uf, ok := data["uf"].(string); ok {
		c.Uf = uf
	}

	return nil
}
