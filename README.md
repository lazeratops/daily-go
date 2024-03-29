## daily-go

![Build and test](https://github.com/lazeratops/daily-go/actions/workflows/build-and-test.yml/badge.svg)
[![codecov](https://codecov.io/gh/lazeratops/daily-go/branch/main/graph/badge.svg?token=LM3KMNT197)](https://codecov.io/gh/lazeratops/daily-go)

`daily-go` is a Go library to communicate with [Daily's REST API](https://docs.daily.co/reference/rest-api).


### Usage

Run `go get github.com/lazeratops/daily-go`

```go
import (
    "fmt"
    "github.com/lazeratops/daily-go/daily"
    "github.com/lazeratops/daily-go/daily/room"
)

func main() {
    d, err := daily.NewDaily("YOUR_DAILY_API_KEY")
    if err != nil {
        panic(err)
    }
	
    // Create a room
    r, err := d.CreateRoom(room.CreateParams{
        Name:            "roomName",
        IsPrivate:       true,
        Props:           room.RoomProps{
            MaxParticipants: 2,
        },
    })
	
    if err != nil {
        panic(err)
    }

    // Get existing rooms
    rooms, err := d.GetRooms(&room.GetManyParams{
        Limit: 5,
    })
	
    if err != nil {
        panic(err)
    }
    
    fmt.Println(len(rooms))
}
```