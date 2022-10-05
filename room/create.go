package room

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang/auth"
	"golang/errors"
	"io"
	"net/http"
)

type CreateParams struct {
	APIKey          string
	APIURL          string
	Name            string
	IsPrivate       bool
	Props           RoomProps
	AdditionalProps map[string]interface{}
}

type createRoomBody struct {
	Name       string                 `json:"name,omitempty"`
	Privacy    string                 `json:"privacy,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

func Create(params CreateParams) (*Room, error) {
	// Make the request body for room creation
	reqBody, err := makeCreateRoomBody(params.Name, params.IsPrivate, params.Props, params.AdditionalProps)
	if err != nil {
		return nil, fmt.Errorf("failed to make room creation request body: %w", err)
	}

	endpoint, err := roomsEndpoint(params.APIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain rooms endpoint: %w", err)
	}
	// Make the actual HTTP request
	req, err := http.NewRequest("POST", endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request to rooms endpoint: %w", err)
	}

	// Prepare auth and content-type headers for request
	auth.SetAPIKeyAuthHeaders(req, params.APIKey)

	// Do the thing!!!
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	// Parse the response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.NewErrFailedAPICall(res.StatusCode, string(resBody))
	}

	var room Room
	if err := json.Unmarshal(resBody, &room); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body into Room: %w", err)
	}

	return &room, nil
}

func makeCreateRoomBody(name string, isPrivate bool, props RoomProps, additionalProps map[string]interface{}) (*bytes.Buffer, error) {
	// Concatenate original and additional properties into a JSON blob
	propsData, err := concatRoomProperties(props, additionalProps)
	if err != nil {
		return nil, fmt.Errorf("failed to build room props JSON: %w", err)
	}

	// Prep request body
	reqBody := createRoomBody{
		Name:       name,
		Properties: propsData,
	}

	// Rooms are public by default
	if isPrivate {
		reqBody.Privacy = "private"
	}

	bodyBlob, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return bytes.NewBuffer(bodyBlob), nil
}

func concatRoomProperties(props RoomProps, additionalProps map[string]interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(&props)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal room props: %w", err)
	}

	// Unmarshal all the original props into a map for us to work with
	var mProps map[string]interface{}
	if err := json.Unmarshal(data, &mProps); err != nil {
		return nil, fmt.Errorf("failed to unmarshal props: %w", err)
	}

	// Add additional props to prop map, but only if given key
	// does not already exist in original props.
	for k, v := range additionalProps {
		if _, ok := mProps[k]; ok {
			// This key already exists, skip it
			continue
		}
		mProps[k] = v
	}

	return mProps, nil
}
