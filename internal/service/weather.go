package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"weatherzip/internal/domain"
	"weatherzip/internal/service/contracts"
)

var _ contracts.WeatherService = (*WeatherService)(nil)

var weatherErrorCodes = map[int]error{
	1003: domain.ErrParameterNotProvided,
	1005: domain.ErrApiUrlIsInvalid,
	1006: domain.ErrLocationNotFound,
	9000: domain.ErrJsonBodyIsInvalid,
	9001: domain.ErrTooManyLocations,
	9999: domain.ErrInternalApplication,
}

type WeatherService struct {
	HttpClient contracts.HttpClient
	BaseURL    string
	ApiKey     string
	Language   string
}

func NewWeatherService(client *http.Client, baseURL, apiKey, language string) *WeatherService {
	return &WeatherService{
		HttpClient: client,
		BaseURL:    baseURL,
		ApiKey:     apiKey,
		Language:   language,
	}
}

func (s *WeatherService) GetWeather(location string) (domain.WeatherResponse, error) {
	var response domain.WeatherResponse

	encodedLocation := url.QueryEscape(location)
	url := fmt.Sprintf(s.BaseURL, s.ApiKey, encodedLocation, s.Language)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, domain.NewFailedToCreateRequestError(err)
	}

	res, err := s.HttpClient.Do(req)
	if err != nil {
		return response, domain.NewFailedToMakeRequestError(err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		var errorResponse struct {
			Error struct {
				Code int `json:"code"`
			} `json:"error"`
		}
		if err := json.NewDecoder(res.Body).Decode(&errorResponse); err == nil {
			if apiErr, exists := weatherErrorCodes[errorResponse.Error.Code]; exists {
				return response, apiErr
			}
		}
		return response, domain.ErrUnexpectedBadRequest
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusBadRequest {
		return response, domain.NewUnexpectedStatusCodeError(res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return response, domain.NewFailedToDecodeResponseError(err)
	}

	response.Current.TempK = response.Current.TempC + 273.15

	return response, nil
}
