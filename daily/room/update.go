package room

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang/daily/auth"
	"golang/daily/errors"
	"io"
	"net/http"
)

type UpdateParams struct {
	Creds           auth.Creds
	Name            string
	Privacy         *Privacy
	Props           RoomProps
	AdditionalProps map[string]interface{}
}

type updateRoomBody struct {
	Privacy    Privacy                `json:"privacy,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

func Update(params UpdateParams) error {
	creds := params.Creds
	endpoint, err := roomsEndpoint(creds.APIURL, params.Name)
	if err != nil {
		return err
	}

	reqBody, err := makeUpdateRoomBody(params.Privacy, params.Props, params.AdditionalProps)
	if err != nil {
		return fmt.Errorf("failed to make room update request body: %w", err)
	}

	// Make the actual HTTP request
	req, err := http.NewRequest("POST", endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create POST request to room endpoint: %w", err)
	}

	// Prepare auth and content-type headers for request
	auth.SetAPIKeyAuthHeaders(req, creds.APIKey)

	// Do the thing!!!
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	// Parse the response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.NewErrFailedBodyRead(err)
	}

	if res.StatusCode != http.StatusOK {
		return errors.NewErrFailedAPICall(res.StatusCode, string(resBody))
	}

	var room Room
	if err := json.Unmarshal(resBody, &room); err != nil {
		return NewErrFailUnmarshal(err)
	}

	return nil
}

func makeUpdateRoomBody(privacy *Privacy, props RoomProps, additionalProps map[string]interface{}) (*bytes.Buffer, error) {
	// Concatenate original and additional properties into a JSON blob
	propsData, err := concatRoomProperties(props, additionalProps)
	if err != nil {
		return nil, fmt.Errorf("failed to build room props JSON: %w", err)
	}

	// Prep request body
	reqBody := updateRoomBody{
		Properties: propsData,
	}
	if privacy != nil {
		reqBody.Privacy = *privacy
	}

	bodyBlob, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return bytes.NewBuffer(bodyBlob), nil
}
