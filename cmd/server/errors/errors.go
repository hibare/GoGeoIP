package errors

import "errors"

var (
	ErrAuthenticationRequired = errors.New("authentication required")
	ErrInvalidAuthToken       = errors.New("invalid authentication token")
	ErrSomethingWentWrong     = errors.New("something went wrong")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrReadingPayload         = errors.New("unable to read payload")
)
