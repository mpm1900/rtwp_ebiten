package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
)

var UpdateCamera = events.NewEventType[math.Vec2]()

func InitCamera(world donburi.World) {
	UpdateCamera.Subscribe(world, updateCamera)
}

func updateCamera(world donburi.World, delta math.Vec2) {
	player := components.GetPlayer(world)
	player.Camera.MovePosition(-delta.X, -delta.Y)
}
