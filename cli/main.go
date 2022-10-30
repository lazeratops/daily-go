package main

import (
	"context"
	"github.com/alecthomas/kong"
	"go.uber.org/zap"
	"time"
)

type RoomCreateCmd struct {
	Name      string                 `name:"name" help:"Room name"`
	Prefix    string                 `name:"prefix" help:"Prefix to use for otherwise randomly generated room name"`
	IsPrivate bool                   `name:"is-private" help:"Whether the room should be private" default:"false"`
	Props     map[string]interface{} `name:"props" help:"Room properties"`
}

type RoomGetCmd struct {
	Name           string    `name:"name" help:"Name of room to get"`
	Regex          string    `name:"regex" help:"Regex to filter room names by"`
	Interactive    bool      `name:"interactive" help:"Show results in interactive format"`
	Limit          int       `name:"limit" help:"Maximum number of rooms to retrieve"`
	CreatedBefore  time.Time `name:"created-before" help:"Latest creation date"`
	CreatedAfter   time.Time `name:"created-after" help:"Earliest creation date"`
	IncludeExpired bool      `name:"include-expired" help:"Include expired rooms" default:"true" negatable:""`
}

var cli struct {
	APIKey string `name:"apikey" short:"a" help:"Daily API key" type:"string" env:"DAILY_API_KEY" required:""`
	Room   struct {
		Create RoomCreateCmd `cmd:"" help:"Create a Daily room."`
		Get    RoomGetCmd    `cmd:"" help:"Get rooms."`
	} `cmd:"" help:"Daily room operations."`
}

func main() {
	ctx := kong.Parse(&cli)
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	defer logger.Sync() // flushes buffer, if any
	switch ctx.Command() {
	case "room create":
		if err := roomCreate(sugar, cli.APIKey, cli.Room.Create); err != nil {
			sugar.Fatal("failed to create room", err)
		}
	case "room get":
		getCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if err := roomGet(getCtx, sugar, cli.APIKey, cli.Room.Get); err != nil {
			sugar.Fatal("failed to get room(s): %v", err)
		}
	default:
		panic(ctx.Command())
	}
}
