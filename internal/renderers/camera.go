package renderers

import (
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

type cameraView struct {
	offset math.Vec2
}

func newCameraView(ecs *ecs.ECS) cameraView {
	player := components.GetPlayer(ecs.World)
	if player == nil || player.Camera == nil {
		return cameraView{}
	}

	rect := player.Camera.Surface.Bounds()
	return cameraView{
		offset: math.NewVec2(
			float64(rect.Dx())/2-player.Camera.X,
			float64(rect.Dy())/2-player.Camera.Y,
		),
	}
}

func (view cameraView) Point(point math.Vec2) math.Vec2 {
	return point.Add(view.offset)
}

func (view cameraView) Translate(options *ebiten.DrawImageOptions, position math.Vec2) {
	point := view.Point(position)
	options.GeoM.Translate(point.X, point.Y)
}
