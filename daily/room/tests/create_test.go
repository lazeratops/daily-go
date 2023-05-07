package tests

import (
	"github.com/lazeratops/daily-go/daily/auth"
	"github.com/lazeratops/daily-go/daily/errors"
	"github.com/lazeratops/daily-go/daily/room"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name               string
		params             room.CreateParams
		retCode            int
		retBody            string
		wantErr            error
		wantRoom           room.Room
		getWantedCreatedAt func() time.Time
	}{
		{
			name:    "success without params",
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
			wantRoom: room.Room{
				ID:      "987b5eb5-d116-4a4e-8e2c-14fcb5710966",
				Name:    "getting-started-webinar",
				Url:     "https://api-demo.daily.co/getting-started-webinar",
				Privacy: room.PrivacyPrivate,
				Config: room.Props{
					StartAudioOff: true,
					StartVideoOff: true,
				},
			},
			getWantedCreatedAt: func() time.Time {
				gotTime, gotErr := time.Parse(time.RFC3339, "2019-01-26T09:01:22.000Z")
				require.NoError(t, gotErr)
				return gotTime
			},
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
			creds := auth.Creds{
				APIKey: "somekey",
				APIURL: testServer.URL,
			}
			gotRoom, gotErr := room.Create(creds, p)
			require.ErrorIs(t, gotErr, tc.wantErr)
			if tc.wantErr == nil {
				tc.wantRoom.CreatedAt = tc.getWantedCreatedAt()
				require.EqualValues(t, &tc.wantRoom, gotRoom)
			}
			defer testServer.Close()
		})

	}
}
