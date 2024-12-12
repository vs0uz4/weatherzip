package domain

import "errors"

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
