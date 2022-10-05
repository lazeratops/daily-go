package room

import (
	"github.com/stretchr/testify/require"
	"golang/errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	const rfc3339 = "2006-01-02T15:04:05Z07:00"
	creationTime, gotErr := time.Parse(rfc3339, "2019-01-26T09:01:22.000Z")
	require.NoError(t, gotErr)
	testCases := []struct {
		name               string
		createParams       CreateParams
		getRoomProps       func() RoomProps
		getAdditionalProps func() map[string]interface{}
		dailyResStatusCode int
		dailyResBody       string
		wantRoom           Room
		wantErr            error
	}{
		{
			name:               "bad status code",
			dailyResStatusCode: http.StatusBadRequest,
			dailyResBody:       "{}",
			wantErr:            errors.ErrFailedAPICall,
		},
		{
			name: "room created",
			createParams: CreateParams{
				APIKey:    "some-key",
				Name:      "getting-started-webinar",
				IsPrivate: true,
			},
			getRoomProps: func() RoomProps {
				return RoomProps{}
			},
			getAdditionalProps: func() map[string]interface{} {
				return map[string]interface{}{
					"start_audio_off": true,
					"start_video_off": true,
				}
			},
			dailyResStatusCode: http.StatusOK,
			// This data is retrieved from API examples:
			// https://docs.daily.co/reference/rest-api/rooms/create-room#example-requests
			dailyResBody: `
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
			wantRoom: Room{
				ID:        "987b5eb5-d116-4a4e-8e2c-14fcb5710966",
				Name:      "getting-started-webinar",
				Url:       "https://api-demo.daily.co/getting-started-webinar",
				CreatedAt: creationTime,
				AdditionalProps: map[string]interface{}{
					"start_audio_off": true,
					"start_video_off": true,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.dailyResStatusCode)
				_, err := w.Write([]byte(tc.dailyResBody))
				require.NoError(t, err)
			}))

			defer testServer.Close()

			createParams := tc.createParams
			createParams.APIURL = testServer.URL

			if tc.getRoomProps != nil {
				createParams.Props = tc.getRoomProps()
			}
			if tc.getAdditionalProps != nil {
				createParams.AdditionalProps = tc.getAdditionalProps()
			}

			gotRoom, gotErr := Create(createParams)
			require.ErrorIs(t, gotErr, tc.wantErr)
			if gotErr == nil {
				require.EqualValues(t, tc.wantRoom, *gotRoom)
			}
		})
	}
}
