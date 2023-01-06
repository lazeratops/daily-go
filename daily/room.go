package daily

import (
	"fmt"
	"github.com/lazeratops/daily-go/daily/auth"
	"github.com/lazeratops/daily-go/daily/room"
	"regexp"
	"time"
)

// CreateRoom creates a Daily room using Daily's REST API
func (d *Daily) CreateRoom(params room.CreateParams) (*room.Room, error) {
	creds := auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}
	if params.Props.Exp == 0 {
		params.Props.SetExpiry(time.Now().Add(d.defaultRoomExp))
	}
	if params.Prefix != "" {
		return room.CreateWithPrefix(creds, room.CreateParams{
			IsPrivate:       params.IsPrivate,
			Props:           params.Props,
			AdditionalProps: params.AdditionalProps,
			Prefix:          params.Prefix,
		})
	}
	return room.Create(creds, room.CreateParams{
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

// SendAppMessage sends an "app-message" event to the given room
func (d *Daily) SendAppMessage(roomName string, data string, recipient *string) error {
	r := "*"
	if recipient != nil {
		r = *recipient
	}
	return room.SendAppMessage(auth.Creds{
		APIKey: d.apiKey,
		APIURL: d.apiURL,
	}, room.SendAppMessageParams{
		RoomName:  roomName,
		Data:      data,
		Recipient: r,
	})
}
