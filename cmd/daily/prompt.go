package main

import (
	"context"
	"fmt"
	"github.com/lazeratops/daily-go/daily"
	"github.com/lazeratops/daily-go/daily/room"
	"github.com/manifoldco/promptui"
	"go.uber.org/zap"
	"strings"
)

type roomItem struct {
	ID         string
	Name       string
	IsSelected bool
}

func selectRoomNames(selectedPos int, selection []*roomItem, allRooms []room.Room) ([]*roomItem, error) {
	var items = []*roomItem{
		{
			ID:   "Done",
			Name: "Finished",
		},
	}
	for _, r := range allRooms {
		var isSelected bool
		for _, selectedR := range selection {
			if r.ID == selectedR.ID {
				isSelected = true
				break
			}
		}
		items = append(items, &roomItem{
			ID:         r.ID,
			Name:       r.Name,
			IsSelected: isSelected,
		})
	}

	templates := &promptui.SelectTemplates{
		Label: `{{if .IsSelected}}
					✔
				{{end}} {{ .Name }} - label`,
		Active:   "→ {{if .IsSelected}}✔ {{end}}{{ .Name | cyan }} ({{ .ID | red }})",
		Inactive: "{{if .IsSelected}}✔ {{end}}{{ .Name | cyan }} ({{ .ID | red }})",
		Details: `
--------- Room ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Selected:" | faint }}	{{ .IsSelected }}`,
	}

	searcher := func(input string, index int) bool {
		roomItem := items[index]
		name := strings.Replace(strings.ToLower(roomItem.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:        "Room",
		Items:        items,
		Templates:    templates,
		Size:         25,
		Searcher:     searcher,
		CursorPos:    selectedPos,
		HideSelected: true,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	chosenItem := items[i]
	if chosenItem.Name == "Finished" {
		return selection, nil
	}
	selected := chosenItem.IsSelected
	chosenItem.IsSelected = !selected

	var newSelection []*roomItem
	for _, item := range items {
		if item.IsSelected {
			newSelection = append(newSelection, item)
		}
	}
	return selectRoomNames(i, newSelection, allRooms)
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

func roomItemsToRooms(items []*roomItem, rooms []room.Room) []room.Room {
	var retRooms []room.Room
	for _, i := range items {
		for _, r := range rooms {
			if i.ID == r.ID {
				retRooms = append(retRooms, r)
				break
			}
		}
	}
	return retRooms
}
