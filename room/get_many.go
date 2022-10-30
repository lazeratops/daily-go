package room

import (
	"encoding/json"
	"fmt"
	"golang/auth"
	"golang/errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type getManyResponse struct {
	TotalCount int    `json:"total_count"`
	Data       []Room `json:"data"`
}

type GetManyParams struct {
	Limit         int   `json:"limit"`
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
	// If no params given, find all rooms
	if params == nil {
		return getAllRooms(creds, nil)
	}

	rooms, err := doGetRooms(creds, params)
	if err != nil {
		return nil, err
	}

	// Make another request that starts with the
	// oldest room returned, to confirm that we didn't
	// miss any rooms
	l := len(rooms.Data)
	if l < params.Limit {
		// Get the remaining rooms
		lastRoom := rooms.Data[l-1]
		newParams := GetManyParams{
			StartingAfter: params.StartingAfter,
			EndingBefore:  lastRoom.CreatedAt.Unix(),
		}
		moreRooms, err := GetMany(creds, &newParams)
		if err != nil {
			return nil, err
		}
		rooms.Data = append(rooms.Data, moreRooms...)
	}

	return rooms.Data, nil
}

func GetManyWithRegex(creds auth.Creds, params *GetManyParams, regex *regexp.Regexp) ([]Room, error) {
	rooms, err := GetMany(creds, params)
	if err != nil {
		return nil, err
	}
	var matchedRooms []Room
	for _, r := range rooms {
		if regex.MatchString(r.Name) {
			matchedRooms = append(matchedRooms, r)
		}
	}
	return matchedRooms, nil
}

func getAllRooms(creds auth.Creds, params *GetManyParams) ([]Room, error) {
	rooms, err := doGetRooms(creds, params)
	if err != nil {
		return nil, err
	}
	l := len(rooms.Data)
	// If there are more rooms to retrieve,
	// do so now
	if rooms.TotalCount > l {
		lastRoom := rooms.Data[l-1]
		newParams := GetManyParams{
			EndingBefore: lastRoom.CreatedAt.Unix(),
		}
		moreRooms, err := getAllRooms(creds, &newParams)
		if err != nil {
			return nil, err
		}
		rooms.Data = append(rooms.Data, moreRooms...)
	}

	return rooms.Data, nil
}

func doGetRooms(creds auth.Creds, params *GetManyParams) (*getManyResponse, error) {
	var endpoint string
	if params == nil {
		var err error
		endpoint, err = roomsEndpoint(creds.APIURL)
		if err != nil {
			return nil, err
		}
	} else {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal rooms retrieval params: %w", err)
		}

		var paramsMap map[string]string
		var m map[string]int
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON params to map: %w", err)
		}

		paramsMap = make(map[string]string)
		for k, v := range m {
			if v == 0 {
				continue
			}
			s := strconv.Itoa(v)
			paramsMap[k] = s
		}

		endpoint, err = roomsEndpointWithParams(creds.APIURL, paramsMap)
		if err != nil {
			return nil, err
		}
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

	var rooms getManyResponse
	if err := json.Unmarshal(resBody, &rooms); err != nil {
		return nil, NewErrFailUnmarshal(err)
	}
	return &rooms, nil
}
