package room

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lazeratops/daily-go/daily/auth"
	"github.com/lazeratops/daily-go/daily/errors"
	"io"
	"net/http"
)

type SendAppMessageParams struct {
	RoomName  string
	Data      string
	Recipient string
}
type sendAppMessageBody struct {
	Data      string `json:"privacy,omitempty"`
	Recipient string `json:"properties,omitempty"`
}

func SendAppMessage(creds auth.Creds, params SendAppMessageParams) error {
	endpoint, err := roomsEndpoint(creds.APIURL, params.RoomName)
	if err != nil {
		return err
	}

	body := sendAppMessageBody{Data: params.Data, Recipient: params.Recipient}
	bodyBlob, err := json.Marshal(body)
	reqBody := bytes.NewBuffer(bodyBlob)
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
	return nil
}
