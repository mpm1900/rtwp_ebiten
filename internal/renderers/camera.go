package renderers

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

type cameraView struct {
	offset math.Vec2
}

func newCameraView(ecs *ecs.ECS) cameraView {
	c := components.GetCamera(ecs.World)
	if c.Camera == nil {
		return cameraView{}
	}

	rect := c.Camera.Surface.Bounds()
	return cameraView{
		offset: math.NewVec2(
			float64(rect.Dx())/2-c.Camera.X,
			float64(rect.Dy())/2-c.Camera.Y,
		),
	}
}

func (view cameraView) Point(point math.Vec2) math.Vec2 {
	return point.Add(view.offset)
}
