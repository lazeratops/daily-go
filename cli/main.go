package main

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/kong"
	"go.uber.org/zap"
	daily "golang"
	"golang/room"
	"regexp"
	"time"
)

type RoomCreateCmd struct {
	Name      string                 `name:"name" help:"Room name"`
	Prefix    string                 `name:"prefix" help:"Prefix to use for otherwise randomly generated room name"`
	IsPrivate bool                   `name:"isprivate" help:"Whether the room should be private" default:"false"`
	Props     map[string]interface{} `name:"props" help:"Room properties"`
}

type RoomGetCmd struct {
	Name          string    `name:"name" help:"Name of room to get"`
	Interactive   bool      `name:"interactive" help:"Show results in interactive format"`
	Limit         int       `name:"limit" help:"Maximum number of rooms to retrieve"`
	CreatedBefore time.Time `name:"createdbefore" help:"Latest creation date"`
	CreatedAfter  time.Time `name:"createdafter" help:"Earliest creation date"`
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
		if err := roomGet(sugar, cli.APIKey, cli.Room.Get); err != nil {
			sugar.Fatal("failed to get room(s)", err)
		}
	default:
		panic(ctx.Command())
	}
}

func roomCreate(logger *zap.SugaredLogger, apiKey string, cmd RoomCreateCmd) error {
	checkNamePrefix(cmd)

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
	r, err := d.CreateRoom(daily.RoomCreateParams{
		Name:            cmd.Name,
		Prefix:          cmd.Prefix,
		IsPrivate:       cmd.IsPrivate,
		Props:           rp,
		AdditionalProps: cmd.Props,
	})
	if err != nil {
		return err
	}
	roomData, err := json.Marshal(r)
	if err != nil {
		return err
	}
	logger.Infof("created room: %s", string(roomData))
	return nil
}

func checkNamePrefix(cmd RoomCreateCmd) {
	n := cmd.Name
	p := cmd.Prefix
	if n != "" && p != "" {
		fmt.Println("Arguments contain both name and prefix. Name will be ignored.")
	}
}

func roomGet(logger *zap.SugaredLogger, apiKey string, cmd RoomGetCmd) error {
	d, err := daily.NewDaily(apiKey)
	if err != nil {
		return err
	}
	name := cmd.Name
	if name != "" {
		reg, err := regexp.Compile(name)
		fmt.Println("reg:", reg, err)
		if err != nil {
			r, err := d.GetRoom(name)
			if err != nil {
				return err
			}
			roomData, err := json.Marshal(r)
			if err != nil {
				return err
			}
			logger.Infof("got room: %s", string(roomData))
			return nil

		}

		var endingBefore, startingAfter int64
		if !cmd.CreatedBefore.IsZero() {
			endingBefore = cmd.CreatedBefore.Unix()
		}
		if !cmd.CreatedAfter.IsZero() {
			startingAfter = cmd.CreatedAfter.Unix()
		}
		params := &room.GetManyParams{
			Limit:         cmd.Limit,
			EndingBefore:  endingBefore,
			StartingAfter: startingAfter,
		}

		r, err := d.GetRoomsWithRegex(params, reg)
		roomData, err := json.Marshal(r)
		if err != nil {
			return err
		}
		logger.Infof("got rooms: %s", string(roomData))
	}

	return nil
}
