package main

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/kong"
	daily "golang"
	"golang/room"
)

type RoomCreateCmd struct {
	Name      string                 `name:"name" help:"Room name"`
	IsPrivate bool                   `name:"isprivate" help:"Whether the room should be private" default:"false"`
	Props     map[string]interface{} `name:"props" help:"Room properties"`
}

var cli struct {
	APIKey string `name:"apikey" help:"Daily API key" type:"string" env:"DAILY_API_KEY" required:""`
	Room   struct {
		Create RoomCreateCmd `cmd:"" help:"Create a Daily room."`
	} `cmd:"" help:"Daily room operations."`
}

func main() {
	ctx := kong.Parse(&cli)
	cmd := ctx.Command()
	fmt.Println(cmd)
	p := ctx.Path
	fmt.Println(p)

	switch ctx.Command() {
	case "room create":
		if err := roomCreate(cli.APIKey, cli.Room.Create); err != nil {
			panic(err)
		}
	default:
		panic(ctx.Command())
	}
}

func roomCreate(apiKey string, cmd RoomCreateCmd) error {
	d, err := daily.NewDaily(apiKey)
	if err != nil {
		return err
	}
	var rp room.RoomProps
	data, err := json.Marshal(cmd.Props)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &rp); err != nil {
		return err
	}
	r, err := d.CreateRoom(cmd.Name, cmd.IsPrivate, rp, cmd.Props)
	if err != nil {
		return err
	}
	roomData, err := json.Marshal(r)
	if err != nil {
		return err
	}
	fmt.Println(string(roomData))
	return nil
}
