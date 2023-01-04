// Package daily handles communication with the Daily REST API
package daily

import (
	"errors"
	"time"
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
	apiKey         string
	apiURL         string
	defaultRoomExp time.Duration
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
		apiURL:         dailyURL,
		defaultRoomExp: time.Hour * 24,
	}, nil
}

func (d *Daily) WithDefaultRoomExpiry(duration time.Duration) {
	d.defaultRoomExp = duration
}
