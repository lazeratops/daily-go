package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrFailedAPICall is returned when the call to Daily has failed
	ErrFailedAPICall = errors.New("the Daily API call has failed")
)

func NewErrFailedAPICall(statusCode int, body string) error {
	return fmt.Errorf("status code: %d; body: %s: %w", statusCode, body, ErrFailedAPICall)
}
