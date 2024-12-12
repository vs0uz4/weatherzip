package domain

import (
	"errors"
	"fmt"
)

var (
	ErrCepNotFound          = errors.New("cep not found")
	ErrCepIsInvalid         = errors.New("invalid cep")
	ErrLocationNotFound     = errors.New("location not found")
	ErrUnexpectedBadRequest = errors.New("unexpected bad request error")
	ErrParameterNotProvided = errors.New("parameter 'q' not provided")
	ErrApiUrlIsInvalid      = errors.New("API request URL is invalid")
	ErrJsonBodyIsInvalid    = errors.New("invalid JSON body in bulk request")
	ErrTooManyLocations     = errors.New("too many locations in bulk request")
	ErrInternalApplication  = errors.New("internal application error")
	ErrInvalidLocationData  = errors.New("invalid location data")
	ErrInvalidCurrentData   = errors.New("invalid current weather data")
)

func NewUnexpectedStatusCodeError(statusCode int) error {
	return fmt.Errorf("unexpected status code: %d", statusCode)
}

func NewFailedToCreateRequestError(err error) error {
	return fmt.Errorf("failed to create request: %w", err)
}

func NewFailedToMakeRequestError(err error) error {
	return fmt.Errorf("failed to make request: %w", err)
}

func NewFailedToDecodeResponseError(err error) error {
	return fmt.Errorf("failed to decode response: %w", err)
}
