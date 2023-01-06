package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lazeratops/daily-go/daily"
	"github.com/lazeratops/daily-go/daily/room"
	"github.com/olekukonko/tablewriter"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
	"regexp"
	"strings"
)

// roomCreate() creates a Daily room
func roomCreate(logger *zap.SugaredLogger, apiKey string, cmd RoomCreateCmd) error {
	n := cmd.Name
	p := cmd.Prefix
	if n != "" && p != "" {
		logger.Warnf("Arguments contain both Name (%s) and prefix (%s). Name will be ignored.", n, p)
	}

	// Init Daily with given API key
	d, err := daily.NewDaily(apiKey)
	if err != nil {
		return err
	}

	// Prepare room properties
	var rp room.RoomProps
	data, err := json.Marshal(cmd.Props)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &rp); err != nil {
		return err
	}
	r, err := d.CreateRoom(room.CreateParams{
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

// roomGet() retrieves the relevant room(s) in table or interactive mode
func roomGet(ctx context.Context, logger *zap.SugaredLogger, apiKey string, cmd RoomGetCmd) error {
	// Init Daily with given API key
	d, err := daily.NewDaily(apiKey)
	if err != nil {
		return err
	}

	// If Name is provided, just get single room
	// by that Name
	if cmd.Name != "" {
		return roomGetSingle(ctx, logger, cmd, d)
	}

	// Prep get many params
	/*	var endingBefore, startingAfter int64
		if !cmd.CreatedBefore.IsZero() {
			endingBefore = cmd.CreatedBefore.Unix()
		}
		if !cmd.CreatedAfter.IsZero() {
			startingAfter = cmd.CreatedAfter.Unix()
		} */
	params := &room.GetManyParams{
		Limit: cmd.Limit,
		//	EndingBefore:  endingBefore,
		//	StartingAfter: startingAfter,
	}
	var rooms []room.Room

	// If regex is provided, get rooms with regex
	if cmd.Regex != "" {
		reg, err := regexp.Compile(cmd.Regex)
		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}
		rooms, err = d.GetRoomsWithRegex(params, reg)
		if err != nil {
			return err
		}
	} else {
		// If regex is not provided, just get rooms with
		// given params
		rooms, err = d.GetRooms(params)
		if err != nil {
			return err
		}
	}
	// Show retrieved rooms in either interactive or
	// ASCII table mode
	if cmd.Interactive {
		selectedNames, err := selectRoomNames(0, make([]*roomItem, 0), rooms)
		if err != nil {
			return err
		}
		rooms := roomItemsToRooms(selectedNames, rooms)
		_, err = showWithControls(ctx, logger, rooms, d)
		if err != nil {
			return err
		}
		return roomGet(ctx, logger, apiKey, cmd)
	}
	return showInTable(rooms)
}

func roomGetSingle(ctx context.Context, logger *zap.SugaredLogger, cmd RoomGetCmd, d *daily.Daily) error {
	r, err := d.GetRoom(cmd.Name)
	if err != nil {
		return err
	}
	// Show room in non-interactive ASCII table
	if !cmd.Interactive {
		return showInTable([]room.Room{*r})
	}
	// Show room in interactive mode
	deleted, err := showWithControls(ctx, logger, []room.Room{*r}, d)
	if err != nil {
		return err
	}
	if deleted {
		logger.Infof("Room '%s' has been deleted", cmd.Name)
	}
	return nil
}

// showInTable() shows rooms in a non-interactive ASCII table view
func showInTable(rooms []room.Room) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(true)
	table.SetHeader([]string{"URL", "Private", "Created at", "Additional props"})
	hc := tablewriter.Colors{tablewriter.Bold, tablewriter.BgHiCyanColor}
	table.SetHeaderColor(hc, hc, hc, hc)

	w1 := tablewriter.Colors{tablewriter.FgWhiteColor}
	w2 := tablewriter.Colors{tablewriter.FgHiWhiteColor}

	for i, r := range rooms {
		var private string
		if r.Privacy == room.PrivacyPrivate {
			private = "\u2713"
		}

		// Marshal properties to string, add some spaces for wrapping
		propsData, err := json.Marshal(r.AdditionalProps)
		if err != nil {
			return err
		}
		ps := strings.ReplaceAll(string(propsData), ",", ", ")

		// Set color to use for row
		c := w1
		if i%2 == 0 {
			c = w2
		}

		table.Rich([]string{r.Url, private, r.CreatedAt.String(), ps}, []tablewriter.Colors{c, c, c})
	}
	table.Render()
	return nil
}

// deleteRooms() deletes the given rooms
func deleteRooms(ctx context.Context, logger *zap.SugaredLogger, rooms []room.Room, daily *daily.Daily) error {
	errs, ctx := errgroup.WithContext(ctx)

	for _, r := range rooms {
		r := r
		errs.Go(func() error {
			logger.Debugf("Deleting room '%s'", r.Name)
			return daily.DeleteRoom(r.Name)
		})

	}
	return errs.Wait()
}

// removeRoom() removes the given index from a slice of rooms
func removeRoom(rooms []room.Room, idx int) []room.Room {
	return append(rooms[:idx], rooms[idx+1:]...)
}
