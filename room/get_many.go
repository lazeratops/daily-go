package room

import (
	"encoding/json"
	"fmt"
	"golang/auth"
	"golang/errors"
	"io"
	"net/http"
	"time"
)

type roomsResponse struct {
	TotalCount int    `json:"total_count"`
	Data       []Room `json:"data"`
}

type GetManyParams struct {
	Limit         int32 `json:"limit"`
	EndingBefore  int64 `json:"ending_before"`
	StartingAfter int64 `json:"starting_after"`
}

// SetEndingBefore sets EndingBefore as a Unix timestamp
func (p *GetManyParams) SetEndingBefore(endingBefore time.Time) {
	p.EndingBefore = endingBefore.Unix()
}

// GetEndingBefore retrieves the room EndingBefore
func (p *GetManyParams) GetEndingBefore() time.Time {
	return time.Unix(p.EndingBefore, 0)
}

// SetStartingAfter sets StartingAfter as a Unix timestamp
func (p *GetManyParams) SetStartingAfter(startingAfter time.Time) {
	p.StartingAfter = startingAfter.Unix()
}

// GetStartingAfter retrieves the room StartingAfter
func (p *GetManyParams) GetStartingAfter() time.Time {
	return time.Unix(p.StartingAfter, 0)
}

func GetMany(creds auth.Creds, params *GetManyParams) ([]Room, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rooms retrieval params: %w", err)
	}

	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON params to map: %w", err)
	}

	endpoint, err := roomsEndpointWithParams(creds.APIURL, m)
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

	var rooms roomsResponse
	if err := json.Unmarshal(resBody, &rooms); err != nil {
		return nil, NewErrFailUnmarshal(err)
	}

	return rooms.Data, nil
}
