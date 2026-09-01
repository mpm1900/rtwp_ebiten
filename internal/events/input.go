package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

var ActionClick = events.NewEventType[components.ActionEvent]()

func InitInput(world donburi.World) {
	ActionClick.Subscribe(world, handleActionClick)
}

func handleActionClick(world donburi.World, event components.ActionEvent) {
	player := components.GetPlayer(world)
	if player.SelectedAction == nil {
		return
	}

	if !player.SelectedAction.Valid(world, event.Point) {
		return
	}

	player.SelectedAction.Publish(world, event)
}
