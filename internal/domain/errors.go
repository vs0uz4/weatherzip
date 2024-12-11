package domain

import "errors"

var ErrCepNotFound = errors.New("cep not found")
var ErrCepIsInvalid = errors.New("invalid cep")
var ErrLocationNotFound = errors.New("location not found")
var ErrUnexpectedBadRequest = errors.New("unexpected bad request error")
var ErrParameterNotProvided = errors.New("parameter 'q' not provided")
var ErrApiUrlIsInvalid = errors.New("API request URL is invalid")
var ErrJsonBodyIsInvalid = errors.New("invalid JSON body in bulk request")
var ErrTooManyLocations = errors.New("too many locations in bulk request")
var ErrInternalApplication = errors.New("internal application error")
