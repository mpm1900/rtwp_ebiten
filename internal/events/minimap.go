package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
)

var LeftClickMinimap = events.NewEventType[math.Vec2]()

func InitMinimap(world donburi.World) {
	LeftClickMinimap.Subscribe(world, leftClickMinimap)
}

func leftClickMinimap(world donburi.World, point math.Vec2) {
	camera := components.GetCamera(world)
	if camera.Camera != nil {
		camera.Camera.SetPosition(point.X, point.Y)
		camera.ClampCameraPosition()
	}
}
