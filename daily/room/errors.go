package room

import (
	"errors"
	"fmt"
)

var (
	ErrFailUnmarshal  = errors.New("failed to unmarshal response body into Room")
	ErrFailRoomDelete = errors.New("failed to delete room")
)

func NewErrFailUnmarshal(unmarshalErr error) error {
	return fmt.Errorf("%s: %w", unmarshalErr, ErrFailUnmarshal)
}

func NewErrFailRoomDelete(deleteErr error) error {
	return fmt.Errorf("%s: %w", deleteErr, ErrFailRoomDelete)
}
