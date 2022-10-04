// Package room handles Daily room creation, retrieval, and modification
package room

import (
	"errors"
	"golang/room"
)

const (
	dailyURL = "https://api.daily.co/v1/"
)

var (
	// ErrInvalidTokenExpiry is returned when the caller attempts to create
	// a meeting token without a valid expiry time.
	ErrInvalidTokenExpiry = errors.New("expiry cannot be empty or in the past")
	// ErrInvalidAPIKey is returned when the caller attempts to provide
	// an invalid Daily API key.
	ErrInvalidAPIKey = errors.New("API key is invalid")
)

// Daily communicates with Daily's REST API
type Daily struct {
	apiKey string
	apiURL string
}

// NewDaily returns a new instance of Daily
func NewDaily(apiKey string) (*Daily, error) {
	// Check that user passed in what at least COULD be a valid
	// API key. In a prod implementation you probably want to
	// have additional validity checks here.
	if apiKey == "" {
		return nil, ErrInvalidAPIKey
	}
	return &Daily{
		apiKey: apiKey,
		// This is set on the struct instead of just reusing the
		// const to enable overriding for unit tests.
		apiURL: dailyURL,
	}, nil
}

// CreateRoom creates a Daily room using Daily's REST API
func (d *Daily) CreateRoom(name string, isPrivate bool, props room.RoomProps, additionalProps map[string]interface{}) (*room.Room, error) {
	return room.Create(room.CreateParams{
		APIKey:          d.apiKey,
		APIURL:          d.apiURL,
		Name:            name,
		IsPrivate:       isPrivate,
		Props:           props,
		AdditionalProps: additionalProps,
	})
}
