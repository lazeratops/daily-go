package tests

import (
	"github.com/lazeratops/daily-go/daily/auth"
	"github.com/lazeratops/daily-go/daily/room"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteRoom(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name               string
		roomName           string
		dailyResStatusCode int
		dailyResBody       string
		wantErr            error
	}{
		{
			name:               "success",
			roomName:           "room-0253",
			dailyResStatusCode: 200,
			dailyResBody: `
			{
			  "deleted": true,
			  "name": "room-0253"
			}`,
		},
		{
			name:               "wrong-room-name",
			roomName:           "room-0253",
			dailyResStatusCode: 200,
			dailyResBody: `
			{
			  "deleted": true,
			  "name": "room-0251"
			}`,
			wantErr: room.ErrFailRoomDelete,
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

			gotErr := room.Delete(auth.Creds{
				APIKey: "someKey",
				APIURL: testServer.URL,
			}, tc.roomName)
			require.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}
