package room

import (
	"errors"
	"fmt"
)

var (
	ErrFailUnmarshal = errors.New("failed to unmarshal response body into Room")
)

func NewErrFailUnmarshal(unmarshalErr error) error {
	return fmt.Errorf("%s: %w", unmarshalErr, ErrFailUnmarshal)
}
