package events

import (
	"math"
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	dmath "github.com/yohamta/donburi/features/math"
)

type ZoomCameraData struct {
	Delta  float64
	Cursor dmath.Vec2
}

var UpdateCamera = events.NewEventType[dmath.Vec2]()
var ZoomCamera = events.NewEventType[ZoomCameraData]()

func InitCamera(world donburi.World) {
	UpdateCamera.Subscribe(world, updateCamera)
	ZoomCamera.Subscribe(world, zoomCamera)
}

func updateCamera(world donburi.World, delta dmath.Vec2) {
	player := components.GetPlayer(world)
	if player == nil || player.Camera == nil {
		return
	}

	scale := player.Camera.Scale
	if scale <= 0 {
		scale = 1.0
	}

	player.Camera.MovePosition(-delta.X/scale, -delta.Y/scale)
	player.ClampCameraPosition()
}

func zoomCamera(world donburi.World, data ZoomCameraData) {
	player := components.GetPlayer(world)
	if player == nil || player.Camera == nil {
		return
	}

	if data.Delta == 0 {
		return
	}

	factor := math.Pow(1.15, data.Delta)
	newScale := min(components.MaxCameraZoom, max(components.MinCameraZoom, player.Camera.Scale*factor))
	if newScale == player.Camera.Scale {
		return
	}

	cursorWorldBefore := player.ScreenToWorld(data.Cursor)
	player.Camera.SetZoom(newScale)
	cursorWorldAfter := player.ScreenToWorld(data.Cursor)

	diff := cursorWorldBefore.Sub(cursorWorldAfter)
	player.Camera.MovePosition(diff.X, diff.Y)
	player.ClampCameraPosition()
}
