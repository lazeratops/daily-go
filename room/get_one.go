package room

import (
	"encoding/json"
	"fmt"
	"golang/auth"
	"golang/errors"
	"io"
	"net/http"
)

func GetOne(roomName string, creds auth.Creds) (*Room, error) {
	endpoint, err := roomsEndpoint(creds.APIURL, roomName)
	if err != nil {
		return nil, err
	}

	// Make the actual HTTP request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request to room endpoint: %w", err)
	}

	// Prepare auth and content-type headers for request
	auth.SetAPIKeyAuthHeaders(req, creds.APIKey)

	// Do the thing!!!
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	// Parse the response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.NewErrFailedBodyRead(err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.NewErrFailedAPICall(res.StatusCode, string(resBody))
	}

	var room Room
	if err := json.Unmarshal(resBody, &room); err != nil {
		return nil, NewErrFailUnmarshal(err)
	}

	return &room, nil
}
