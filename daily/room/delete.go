package room

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lazeratops/daily-go/daily/auth"
	errors2 "github.com/lazeratops/daily-go/daily/errors"
	"io"
	"net/http"
)

type deleteResponse struct {
	Deleted  bool   `json:"deleted"`
	RoomName string `json:"name"`
}

func Delete(creds auth.Creds, roomName string) error {
	endpoint, err := roomsEndpoint(creds.APIURL, roomName)
	if err != nil {
		return err
	}

	// Make the actual HTTP request
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request to room endpoint: %w", err)
	}

	// Prepare auth and content-type headers for request
	auth.SetAPIKeyAuthHeaders(req, creds.APIKey)

	// Do the thing!!!
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}

	// Parse the response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return errors2.NewErrFailedBodyRead(err)
	}

	if res.StatusCode != http.StatusOK {
		return errors2.NewErrFailedAPICall(res.StatusCode, string(resBody))
	}

	var dr deleteResponse
	if err := json.Unmarshal(resBody, &dr); err != nil {
		return NewErrFailUnmarshal(err)
	}
	if dr.Deleted {
		if dr.RoomName != roomName {
			err := fmt.Errorf("requested deletion was of room name '%s', but room reported deleted was '%s'", roomName, dr.RoomName)
			return NewErrFailRoomDelete(err)
		}
		return nil
	}
	return NewErrFailRoomDelete(errors.New("room not deleted"))
}
