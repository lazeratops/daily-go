package daily

import (
	"fmt"
	"golang/daily/auth"
	room2 "golang/daily/room"
	"regexp"
	"time"
)

// CreateRoom creates a Daily room using Daily's REST API
func (d *Daily) CreateRoom(params room2.CreateParams) (*room2.Room, error) {
	creds := auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}
	if params.Props.Exp == 0 {
		params.Props.SetExpiry(time.Now().Add(d.defaultRoomExp))
	}
	if params.Prefix != "" {
		return room2.CreateWithPrefix(room2.CreateParams{
			Creds:           creds,
			IsPrivate:       params.IsPrivate,
			Props:           params.Props,
			AdditionalProps: params.AdditionalProps,
			Prefix:          params.Prefix,
		})
	}
	return room2.Create(room2.CreateParams{
		Creds:           creds,
		Name:            params.Name,
		IsPrivate:       params.IsPrivate,
		Props:           params.Props,
		AdditionalProps: params.AdditionalProps,
	})
}

// GetRooms returns multiple Daily rooms matching the given
// limits, if any
func (d *Daily) GetRooms(params *room2.GetManyParams) ([]room2.Room, error) {
	return room2.GetMany(auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}, params)
}

func (d *Daily) GetRoomsWithRegexStr(params *room2.GetManyParams, nameRegexStr string) ([]room2.Room, error) {
	reg, err := regexp.Compile(nameRegexStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse regex: %w", err)
	}
	return d.GetRoomsWithRegex(params, reg)
}

func (d *Daily) GetRoomsWithRegex(params *room2.GetManyParams, nameRegex *regexp.Regexp) ([]room2.Room, error) {
	rooms, err := room2.GetManyWithRegex(auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}, params, nameRegex)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetRoom returns a single Daily room matching the given name
func (d *Daily) GetRoom(name string) (*room2.Room, error) {
	return room2.GetOne(auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}, name)
}

// DeleteRoom deletes the given Daily room
func (d *Daily) DeleteRoom(roomName string) error {
	return room2.Delete(auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}, roomName)
}

// SendAppMessage sends an "app-message" event to the given room
func (d *Daily) SendAppMessage(roomName string, data string, recipient *string) error {
	r := "*"
	if recipient != nil {
		r = *recipient
	}
	return room2.SendAppMessage(auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}, room2.SendAppMessageParams{
		RoomName:  roomName,
		Data:      data,
		Recipient: r,
	})
}
