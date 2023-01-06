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

func TestCreate(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		params  room.CreateParams
		retCode int
		retBody string
		wantErr error
	}{
		{
			name:    "success",
			params:  room.CreateParams{},
			retCode: 200,
			retBody: `
				{
				  "id": "987b5eb5-d116-4a4e-8e2c-14fcb5710966",
				  "name": "getting-started-webinar",
				  "api_created": true,
				  "privacy":"private",
				  "url":"https://api-demo.daily.co/getting-started-webinar",
				  "created_at":"2019-01-26T09:01:22.000Z",
				  "config":{
					"start_audio_off": true,
					"start_video_off": true
				  }
				}
        `,
		},
		{
			name:    "failure",
			params:  room.CreateParams{},
			retCode: 400,
			wantErr: errors.ErrFailedAPICall,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.retCode)
				_, err := w.Write([]byte(tc.retBody))
				require.NoError(t, err)
			}))
			p := tc.params
			p.Creds = auth.Creds{
				APIKey: "somekey",
				APIURL: testServer.URL,
			}
			_, gotErr := room.Create(p)
			require.ErrorIs(t, gotErr, tc.wantErr)
			defer testServer.Close()
		})

	}
}
