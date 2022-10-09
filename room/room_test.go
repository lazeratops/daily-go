package room

import (
	"github.com/stretchr/testify/require"
	"golang/auth"
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
				Privacy:   PrivacyPrivate,
				CreatedAt: creationTime,
				AdditionalProps: map[string]interface{}{
					"start_audio_off": true,
					"start_video_off": true,
				},
			},
		},
		{
			name: "room created with exp",
			createParams: CreateParams{
				Name:      "getting-started-webinar",
				IsPrivate: true,
			},
			getRoomProps: func() RoomProps {
				return RoomProps{
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
			wantRoom: Room{
				ID:        "987b5eb5-d116-4a4e-8e2c-14fcb5710966",
				Name:      "ePR84NQ1bPigp79dDezz",
				Url:       "https://api-demo.daily.co/ePR84NQ1bPigp79dDezz",
				Privacy:   PrivacyPublic,
				CreatedAt: creationTime,
				Config: RoomProps{
					Exp: 1548709695,
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
			createParams.Creds = auth.Creds{
				APIKey: "someKey",
				APIURL: testServer.URL,
			}

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

func TestGetOne(t *testing.T) {
	const rfc3339 = "2006-01-02T15:04:05Z07:00"
	creationTime, gotErr := time.Parse(rfc3339, "2019-01-26T09:01:22.000Z")
	require.NoError(t, gotErr)
	testCases := []struct {
		name               string
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
			name:               "room retrieved",
			dailyResStatusCode: http.StatusOK,
			// This data is retrieved from API examples:
			// https://docs.daily.co/reference/rest-api/rooms/get-room-config#example-request
			dailyResBody: `
				{
					"id":"d61cd7b2-a273-42b4-89bd-be763fd562c1",
					"name":"w2pp2cf4kltgFACPKXmX",
					"api_created":false,
					"privacy":"public",
					"url":"https://api-demo.daily.co/w2pp2cf4kltgFACPKXmX",
					"created_at":"2019-01-26T09:01:22.000Z",
					"config":{"start_video_off":true}
				}
			`,
			wantRoom: Room{
				ID:        "d61cd7b2-a273-42b4-89bd-be763fd562c1",
				Name:      "w2pp2cf4kltgFACPKXmX",
				Privacy:   PrivacyPublic,
				Url:       "https://api-demo.daily.co/w2pp2cf4kltgFACPKXmX",
				CreatedAt: creationTime,
				AdditionalProps: map[string]interface{}{
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

			gotRoom, gotErr := GetOne(auth.Creds{
				APIKey: "someKey",
				APIURL: testServer.URL,
			}, "someName")
			require.ErrorIs(t, gotErr, tc.wantErr)
			if gotErr == nil {
				require.EqualValues(t, tc.wantRoom, *gotRoom)
			}
		})
	}
}

func TestGetMany(t *testing.T) {
	const rfc3339 = "2006-01-02T15:04:05Z07:00"
	testCases := []struct {
		name               string
		dailyResStatusCode int
		dailyResBody       string
		wantErr            error
		getWantRooms       func(t *testing.T) []Room
	}{
		{
			name:               "bad status code",
			dailyResStatusCode: http.StatusBadRequest,
			dailyResBody:       "{}",
			wantErr:            errors.ErrFailedAPICall,
		},
		{
			name:               "room retrieved",
			dailyResStatusCode: http.StatusOK,
			// This data is retrieved from API examples:
			// https://docs.daily.co/reference/rest-api/rooms/list-rooms#example-request
			dailyResBody: `
				{
					"total_count":2,
					"data":[
						{
							"id":"5e3cf703-5547-47d6-a371-37b1f0b4427f",
							"name":"w2pp2cf4kltgFACPKXmX",
							"api_created":false,
							"privacy":"public",
							"url":"https://api-demo.daily.co/w2pp2cf4kltgFACPKXmX",
							"created_at":"2019-01-26T09:01:22.000Z",
							"config":{"start_video_off":true}
						},
						{
							"id":"d61cd7b2-a273-42b4-89bd-be763fd562c1",
							"name":"hello",
							"api_created":false,
							"privacy":"public",
							"url":"https://your-domain.daily.co/hello",
							"created_at":"2019-01-25T23:49:42.000Z",
							"config":{}
						}
					]
				}
			`,
			getWantRooms: func(t *testing.T) []Room {
				room1CreationTime, gotErr := time.Parse(rfc3339, "2019-01-26T09:01:22.000Z")
				require.NoError(t, gotErr)

				room1 := Room{
					ID:        "5e3cf703-5547-47d6-a371-37b1f0b4427f",
					Name:      "w2pp2cf4kltgFACPKXmX",
					Privacy:   PrivacyPublic,
					Url:       "https://api-demo.daily.co/w2pp2cf4kltgFACPKXmX",
					CreatedAt: room1CreationTime,
					AdditionalProps: map[string]interface{}{
						"start_video_off": true,
					},
				}

				room2CreationTime, gotErr := time.Parse(rfc3339, "2019-01-25T23:49:42.000Z")
				require.NoError(t, gotErr)
				room2 := Room{
					ID:        "d61cd7b2-a273-42b4-89bd-be763fd562c1",
					Name:      "hello",
					Privacy:   PrivacyPublic,
					Url:       "https://your-domain.daily.co/hello",
					CreatedAt: room2CreationTime,
				}

				return []Room{room1, room2}
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

			gotRooms, gotErr := GetMany(auth.Creds{
				APIKey: "someKey",
				APIURL: testServer.URL,
			}, nil)
			require.ErrorIs(t, gotErr, tc.wantErr)
			if gotErr == nil && tc.getWantRooms != nil {
				wantRooms := tc.getWantRooms(t)
				require.EqualValues(t, wantRooms, gotRooms)
			}
		})
	}
}
