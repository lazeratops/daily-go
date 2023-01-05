package main

import (
	"context"
	"github.com/alecthomas/kong"
	"go.uber.org/zap"
	"time"
)

type RoomCreateCmd struct {
	Name      string                 `help:"Room Name"`
	Prefix    string                 `help:"Prefix to use for otherwise randomly generated room Name"`
	IsPrivate bool                   `help:"Whether the room should be private" default:"false"`
	Props     map[string]interface{} `help:"Room properties"`
}

type RoomGetCmd struct {
	Name        string `help:"Name of room to get"`
	Regex       string `help:"Regex to filter room names by"`
	Interactive bool   `help:"Show results in interactive format"`
	Limit       int    `help:"Maximum number of rooms to retrieve"`
	//	CreatedBefore  time.Time `help:"Latest creation date"`
	//	CreatedAfter   time.Time `help:"Earliest creation date"`
	IncludeExpired bool `help:"Include expired rooms" default:"true" negatable:""`
}

var cli struct {
	APIKey string `short:"a" help:"Daily API key" type:"string" env:"DAILY_API_KEY" required:""`
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
