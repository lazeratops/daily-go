package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrFailedAPICall is returned when the call to Daily has failed
	ErrFailedAPICall              = errors.New("the Daily API call has failed")
	ErrFailedBodyRead             = errors.New("failed to read Daily API response body")
	ErrFailedEndpointConstruction = errors.New("failed to deduce Daily API call endpoint")
)

func NewErrFailedAPICall(statusCode int, body string) error {
	return fmt.Errorf("status code: %d; body: %s: %w", statusCode, body, ErrFailedAPICall)
}

func NewErrFailedBodyRead(err error) error {
	return fmt.Errorf("%s: %w", err, ErrFailedBodyRead)
}

func NewErrFailedEndpointConstruction(err error) error {
	return fmt.Errorf("%s: %w", err, ErrFailedEndpointConstruction)
}
