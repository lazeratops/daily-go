package room

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/lazeratops/daily-go/daily/auth"
	"github.com/lazeratops/daily-go/daily/errors"
	"io"
	"math/big"
	"net/http"
)

type CreateParams struct {
	Name            string
	IsPrivate       bool
	Props           RoomProps
	AdditionalProps map[string]interface{}
	Prefix          string
}

type createRoomBody struct {
	Name       string                 `json:"name,omitempty"`
	Privacy    Privacy                `json:"privacy,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// CreateWithPrefix creates a room with the name containing the specified
// prefix. The rest of the name is randomized.
func CreateWithPrefix(creds auth.Creds, params CreateParams) (*Room, error) {
	if len(params.Prefix) > 10 {
		return nil, fmt.Errorf("prefix too long, must be up to 10 characters")
	}
	name, err := generateNameWithPrefix(params.Prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to generate room name: %w", err)
	}
	params.Name = name
	return Create(creds, params)
}

func Create(creds auth.Creds, params CreateParams) (*Room, error) {
	// Make the request body for room creation
	reqBody, err := makeCreateRoomBody(params.Name, params.IsPrivate, params.Props, params.AdditionalProps)
	if err != nil {
		return nil, fmt.Errorf("failed to make room creation request body: %w", err)
	}

	endpoint, err := roomsEndpoint(creds.APIURL)
	if err != nil {
		return nil, err
	}
	// Make the actual HTTP request
	req, err := http.NewRequest("POST", endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request to rooms endpoint: %w", err)
	}

	// Prepare auth and content-type headers for request
	auth.SetAPIKeyAuthHeaders(req, creds.APIKey)

	// Do the thing!!!
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
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

	if isPrivate {
		reqBody.Privacy = PrivacyPrivate
	} else {
		reqBody.Privacy = PrivacyPublic
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

func generateNameWithPrefix(prefix string) (string, error) {
	s, err := generateRandStr(20)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", prefix, s), nil
}

func generateRandStr(length int) (string, error) {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[num.Int64()]
	}

	return string(result), nil
}
