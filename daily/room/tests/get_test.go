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

func TestGetOne(t *testing.T) {
	const rfc3339 = "2006-01-02T15:04:05Z07:00"
	creationTime, gotErr := time.Parse(rfc3339, "2019-01-26T09:01:22.000Z")
	require.NoError(t, gotErr)
	testCases := []struct {
		name               string
		dailyResStatusCode int
		dailyResBody       string
		wantRoom           room.Room
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
			wantRoom: room.Room{
				ID:        "d61cd7b2-a273-42b4-89bd-be763fd562c1",
				Name:      "w2pp2cf4kltgFACPKXmX",
				Privacy:   room.PrivacyPublic,
				Url:       "https://api-demo.daily.co/w2pp2cf4kltgFACPKXmX",
				CreatedAt: creationTime,
				Config: room.Props{
					StartVideoOff: true,
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

			gotRoom, gotErr := room.GetOne(auth.Creds{
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
		getWantRooms       func(t *testing.T) []room.Room
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
			getWantRooms: func(t *testing.T) []room.Room {
				room1CreationTime, gotErr := time.Parse(rfc3339, "2019-01-26T09:01:22.000Z")
				require.NoError(t, gotErr)

				room1 := room.Room{
					ID:        "5e3cf703-5547-47d6-a371-37b1f0b4427f",
					Name:      "w2pp2cf4kltgFACPKXmX",
					Privacy:   room.PrivacyPublic,
					Url:       "https://api-demo.daily.co/w2pp2cf4kltgFACPKXmX",
					CreatedAt: room1CreationTime,
					Config: room.Props{
						StartVideoOff: true,
					},
				}

				room2CreationTime, gotErr := time.Parse(rfc3339, "2019-01-25T23:49:42.000Z")
				require.NoError(t, gotErr)
				room2 := room.Room{
					ID:        "d61cd7b2-a273-42b4-89bd-be763fd562c1",
					Name:      "hello",
					Privacy:   room.PrivacyPublic,
					Url:       "https://your-domain.daily.co/hello",
					CreatedAt: room2CreationTime,
				}

				return []room.Room{room1, room2}
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

			gotRooms, gotErr := room.GetMany(auth.Creds{
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
