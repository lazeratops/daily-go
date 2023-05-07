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

func TestUpdate(t *testing.T) {
	t.Parallel()
	pp := room.PrivacyPrivate

	testCases := []struct {
		name               string
		updateParams       room.UpdateParams
		getRoomProps       func() room.Props
		getAdditionalProps func() map[string]interface{}
		dailyResStatusCode int
		dailyResBody       string
		wantErr            error
	}{
		{
			name:               "bad status code",
			dailyResStatusCode: http.StatusBadRequest,
			dailyResBody:       "{}",
			wantErr:            errors.ErrFailedAPICall,
		},
		{
			name: "room updated",
			updateParams: room.UpdateParams{
				Privacy: &pp,
			},
			getRoomProps: func() room.Props {
				return room.Props{}
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
		},
		{
			name:         "room updated with exp",
			updateParams: room.UpdateParams{},
			getRoomProps: func() room.Props {
				return room.Props{
					Exp: 1548709695,
				}
			},
			dailyResStatusCode: http.StatusOK,
			// This data is retrieved from API examples:
			// https://docs.daily.co/reference/rest-api/rooms/create-room#example-requests
			dailyResBody: `
				{
				  "id": "987b5eb5-d116-4a4e-8e2c-14fcb5710966",
				  "name": "ePR84NQ1bPigp79dDezz",
				  "api_created": true,
				  "privacy": "public",
				  "url": "https://api-demo.daily.co/ePR84NQ1bPigp79dDezz",
				  "created_at": "2019-01-26T09:01:22.000Z",
				  "config": {
					"exp": 1548709695
				  }
				}
			`,
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

			updateParams := tc.updateParams
			updateParams.Creds = auth.Creds{
				APIKey: "someKey",
				APIURL: testServer.URL,
			}

			if tc.getRoomProps != nil {
				updateParams.Props = tc.getRoomProps()
			}
			if tc.getAdditionalProps != nil {
				updateParams.AdditionalProps = tc.getAdditionalProps()
			}

			gotErr := room.Update(updateParams)
			require.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}
