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
	c := components.GetCamera(world)
	if c.Camera == nil {
		return
	}

	scale := c.Camera.Scale
	if scale <= 0 {
		scale = 1.0
	}

	c.Camera.MovePosition(-delta.X/scale, -delta.Y/scale)
	c.ClampCameraPosition()
}

func zoomCamera(world donburi.World, data ZoomCameraData) {
	c := components.GetCamera(world)
	if c.Camera == nil {
		return
	}

	if data.Delta == 0 {
		return
	}

	factor := math.Pow(1.15, data.Delta)
	newScale := min(components.MAX_CAMERA_ZOOM, max(components.MIN_CAMERA_ZOOM, c.Camera.Scale*factor))
	if newScale == c.Camera.Scale {
		return
	}

	cursorWorldBefore := c.ScreenToWorld(data.Cursor)
	c.Camera.SetZoom(newScale)
	cursorWorldAfter := c.ScreenToWorld(data.Cursor)

	diff := cursorWorldBefore.Sub(cursorWorldAfter)
	c.Camera.MovePosition(diff.X, diff.Y)
	c.ClampCameraPosition()
}
