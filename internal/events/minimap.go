package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
)

var LeftClickMinimap = events.NewEventType[math.Vec2]()
var RightClickMinimap = events.NewEventType[math.Vec2]()

func InitMinimap(world donburi.World) {
	LeftClickMinimap.Subscribe(world, leftClickMinimap)
	RightClickMinimap.Subscribe(world, rightClickMinimap)
}

func leftClickMinimap(world donburi.World, point math.Vec2) {
	player := components.GetPlayer(world)
	if player != nil && player.Camera != nil {
		player.Camera.SetPosition(point.X, point.Y)
		player.ClampCameraPosition()
	}
}

func rightClickMinimap(world donburi.World, point math.Vec2) {
	if _, hasSelection := components.Selected.First(world); hasSelection {
		for selected := range components.Selected.Iter(world) {
			components.WithMovementTo(selected, point, components.DEFAULT_STOP_DISTANCE)
		}
	}
}
