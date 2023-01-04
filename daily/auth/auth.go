package auth

import (
	"fmt"
	"net/http"
)

type Creds struct {
	APIKey string
	APIURL string
}

// SetAPIKeyAuthHeaders sets a Daily API key as a Bearer token
// on the given request.
func SetAPIKeyAuthHeaders(req *http.Request, apiKey string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
}
