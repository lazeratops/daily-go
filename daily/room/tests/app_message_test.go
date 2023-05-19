package tests

import (
	"github.com/lazeratops/daily-go/daily/auth"
	"github.com/lazeratops/daily-go/daily/errors"
	"github.com/lazeratops/daily-go/daily/room"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendAppMessage(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		params  room.SendAppMessageParams
		retCode int
		retBody string
		wantErr error
	}{
		{
			name: "success",
			params: room.SendAppMessageParams{
				Data:      "my data",
				Recipient: "*",
			},
			retCode: 200,
			retBody: `
				{
				  "sent": "true"
				}`,
		},
		{
			name: "failure",
			params: room.SendAppMessageParams{
				Data:      "my data",
				Recipient: "*",
			},
			retCode: 400,
			wantErr: errors.ErrFailedAPICall,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.retCode)
				_, err := w.Write([]byte(tc.retBody))
				require.NoError(t, err)
			}))
			p := tc.params
			creds := auth.Creds{
				APIKey: "somekey",
				APIURL: testServer.URL,
			}
			gotErr := room.SendAppMessage(creds, p)
			require.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}
