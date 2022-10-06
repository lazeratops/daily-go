package room

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"time"
)

type Privacy string

const (
	PrivacyPrivate Privacy = "private"
	PrivacyPublic          = "public"
)

// Room represents a Daily room
type Room struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Url             string    `json:"url"`
	Privacy         Privacy   `json:"privacy"`
	CreatedAt       time.Time `json:"created_at"`
	Config          RoomProps `json:"config"`
	AdditionalProps map[string]interface{}
}

func (r *Room) UnmarshalJSON(data []byte) error {
	rm := struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		Privacy   Privacy   `json:"privacy"`
		CreatedAt time.Time `json:"created_at"`
		Config    RoomProps `json:"config"`
	}{}

	if err := json.Unmarshal(data, &rm); err != nil {
		return err
	}

	r.ID = rm.ID
	r.Name = rm.Name
	r.Url = rm.Url
	r.Privacy = rm.Privacy
	r.CreatedAt = rm.CreatedAt
	r.Config = rm.Config
	r.AdditionalProps = make(map[string]interface{})

	// Check config values that are not in RoomProps
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to unmarshal body to map: %w", err)
	}

	if config, ok := m["config"].(map[string]interface{}); ok {
		// Get all room properties keys that should NOT go into additionalProps.
		// (Opted for this vs reflection for now)
		roomPropsKeys := GetRoomPropsKeys()
		// Iterate over all config values and, if the keys are not
		// in existing RoomProps keys retrieved above, add these
		// config keys and values into AdditionalProps
		for k, v := range config {
			if !isInSlice(k, roomPropsKeys) {
				r.AdditionalProps[k] = v
			}
		}
	}

	return nil
}

func isInSlice(ele string, s []string) bool {
	for _, propsKey := range s {
		if propsKey == ele {
			return true
		}
	}
	return false
}

func roomsEndpoint(apiURL string, paths ...string) (string, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return "", err
	}

	allPaths := append([]string{u.Path}, paths...)
	u.Path = path.Join(allPaths...)
	return u.String(), nil
}
