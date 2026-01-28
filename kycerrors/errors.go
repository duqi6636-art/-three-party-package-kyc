package kycerrors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidConfig  = errors.New("kyc-sdk: invalid config")
	ErrUnauthorized   = errors.New("kyc-sdk: unauthorized")
	ErrRateLimited    = errors.New("kyc-sdk: rate limited")
	ErrBadRequest     = errors.New("kyc-sdk: bad request")
	ErrServerInternal = errors.New("kyc-sdk: server internal")
	ErrUnexpectedHTTP = errors.New("kyc-sdk: unexpected http error")
)

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	if e == nil {
		return "kyc-sdk: http error"
	}
	if e.Body == "" {
		return fmt.Sprintf("kyc-sdk: http %d", e.StatusCode)
	}
	return fmt.Sprintf("kyc-sdk: http %d: %s", e.StatusCode, e.Body)
}

func (e *HTTPError) Is(target error) bool {
	if e == nil {
		return false
	}

	switch target {
	case ErrBadRequest:
		return e.StatusCode == 400
	case ErrUnauthorized:
		return e.StatusCode == 401 || e.StatusCode == 403
	case ErrRateLimited:
		return e.StatusCode == 429
	case ErrServerInternal:
		return e.StatusCode >= 500 && e.StatusCode <= 599
	case ErrUnexpectedHTTP:
		return e.StatusCode >= 400
	default:
		return false
	}
}
