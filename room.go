package daily

import (
	"fmt"
	"golang/auth"
	"golang/room"
	"regexp"
	"time"
)

type RoomCreateParams struct {
	Name            string
	Prefix          string
	IsPrivate       bool
	Props           room.RoomProps
	AdditionalProps map[string]interface{}
}

// CreateRoom creates a Daily room using Daily's REST API
func (d *Daily) CreateRoom(params RoomCreateParams) (*room.Room, error) {
	creds := auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}
	if params.Props.Exp == 0 {
		params.Props.SetExpiry(time.Now().Add(d.defaultRoomExp))
	}
	if params.Prefix != "" {
		return room.CreateWithPrefix(room.CreateParams{
			Creds:           creds,
			IsPrivate:       params.IsPrivate,
			Props:           params.Props,
			AdditionalProps: params.AdditionalProps,
		}, params.Prefix)
	}
	return room.Create(room.CreateParams{
		Creds:           creds,
		Name:            params.Name,
		IsPrivate:       params.IsPrivate,
		Props:           params.Props,
		AdditionalProps: params.AdditionalProps,
	})
}

// GetRooms returns multiple Daily rooms matching the given
// limits, if any
func (d *Daily) GetRooms(params *room.GetManyParams) ([]room.Room, error) {
	return room.GetMany(auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}, params)
}

func (d *Daily) GetRoomsWithRegexStr(params *room.GetManyParams, nameRegexStr string) ([]room.Room, error) {
	reg, err := regexp.Compile(nameRegexStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse regex: %w", err)
	}
	return d.GetRoomsWithRegex(params, reg)
}

func (d *Daily) GetRoomsWithRegex(params *room.GetManyParams, nameRegex *regexp.Regexp) ([]room.Room, error) {
	rooms, err := room.GetManyWithRegex(auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}, params, nameRegex)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetRoom returns a single Daily room matching the given name
func (d *Daily) GetRoom(name string) (*room.Room, error) {
	return room.GetOne(auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}, name)
}

// DeleteRoom deletes the given Daily room
func (d *Daily) DeleteRoom(roomName string) error {
	return room.Delete(auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}, roomName)
}
