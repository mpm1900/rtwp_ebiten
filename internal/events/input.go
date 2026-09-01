package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
)

var ActionClick = events.NewEventType[math.Vec2]()

func InitInput(world donburi.World) {
	ActionClick.Subscribe(world, handleActionClick)
}

func handleActionClick(world donburi.World, worldPoint math.Vec2) {
	player := components.GetPlayer(world)
	if player.Action == nil {
		return
	}

	if !player.Action.Valid(world, worldPoint) {
		return
	}

	player.Action.Publish(world, worldPoint)
}
