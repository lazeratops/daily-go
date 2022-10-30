package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"go.uber.org/zap"
	daily "golang"
	"golang.org/x/sync/errgroup"
	"golang/room"
	"os"
	"regexp"
)

// roomCreate() creates a Daily room
func roomCreate(logger *zap.SugaredLogger, apiKey string, cmd RoomCreateCmd) error {
	n := cmd.Name
	p := cmd.Prefix
	if n != "" && p != "" {
		logger.Warnf("Arguments contain both name (%s) and prefix (%s). Name will be ignored.", n, p)
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

// roomGet() retrieves the relevant room(s) in table or interactive mode
func roomGet(ctx context.Context, logger *zap.SugaredLogger, apiKey string, cmd RoomGetCmd) error {
	// Init Daily with given API key
	d, err := daily.NewDaily(apiKey)
	if err != nil {
		return err
	}

	// If name is provided, just get single room
	// by that name
	if cmd.Name != "" {
		return roomGetSingle(ctx, logger, cmd, d)
	}

	// Prep get many params
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
		return showInteractive(ctx, logger, rooms, d)
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

func showInteractive(ctx context.Context, logger *zap.SugaredLogger, rooms []room.Room, daily *daily.Daily) error {
	var items []string
	for _, r := range rooms {
		items = append(items, r.Name)
	}
	prompt := promptui.Select{
		Label: "Select room",
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	remainingRooms, err := showRoomDetails(ctx, logger, result, rooms, daily)
	if len(remainingRooms) > 0 {
		return showInteractive(ctx, logger, rooms, daily)
	}
	// Remove selected room from list and rerender
	for i, room := range rooms {
		if room.Name == result {
			rooms = removeRoom(rooms, i)
		}
	}
	return showInteractive(ctx, logger, rooms, daily)
}

// showRoomDetails() shows details about the selected rooms
func showRoomDetails(ctx context.Context, logger *zap.SugaredLogger, name string, rooms []room.Room, daily *daily.Daily) ([]room.Room, error) {
	for _, r := range rooms {
		if r.Name == name {
			deleted, err := showWithControls(ctx, logger, []room.Room{r}, daily)
			if err != nil {
				return nil, err
			}
			if deleted {
				// Remove this room from list and rerender
				for i, room := range rooms {
					if room.Name == name {
						rooms = removeRoom(rooms, i)
					}
				}
			}
			return rooms, nil
		}
	}
	return nil, fmt.Errorf("room by name %s not found", name)
}

// showWithControls() shows room details with control options
func showWithControls(ctx context.Context, logger *zap.SugaredLogger, rooms []room.Room, daily *daily.Daily) (bool, error) {
	if err := showInTable(rooms); err != nil {
		return false, err
	}
	prompt := promptui.Select{
		Label: "Action",
		Items: []string{
			"Delete",
			"Update",
			"Back",
		},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return false, fmt.Errorf("prompt failed: %w", err)
	}
	switch result {
	case "Delete":
		if err := deleteRooms(ctx, logger, rooms, daily); err != nil {
			return false, err
		}
		return true, nil
	case "Back":
		return false, nil
	default:
		return false, fmt.Errorf("invalid action choice: %s", result)
	}
}

// showInTable() shows rooms in a non-interactive ASCII table view
func showInTable(rooms []room.Room) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "URL", "Private", "Created at", "Additional props"})

	for _, r := range rooms {
		var private string
		if r.Privacy == room.PrivacyPrivate {
			private = "\u2713"
		}
		propsData, err := json.Marshal(r.AdditionalProps)
		if err != nil {
			return err
		}
		table.Append([]string{r.ID, r.Url, private, r.CreatedAt.String(), string(propsData)})
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
