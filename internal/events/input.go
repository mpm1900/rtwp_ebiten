package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
)

type ActionClickEvent struct {
	Point math.Vec2
	Shift bool
}

var ActionClick = events.NewEventType[ActionClickEvent]()

func InitInput(world donburi.World) {
	ActionClick.Subscribe(world, handleActionClick)
}

func handleActionClick(world donburi.World, event ActionClickEvent) {
	player := components.GetPlayer(world)
	if player.SelectedAction == nil {
		return
	}

	if !player.SelectedAction.Valid(world, event.Point) {
		return
	}

	player.SelectedAction.Publish(world, event.Point, event.Shift)
}
