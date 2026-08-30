package renderers

import (
	"image/color"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

func RenderRanges(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)

	for entry := range components.RangeQuery.Iter(ecs.World) {
		t := transform.GetTransform(entry)
		r := components.Range.Get(entry)
		center := view.Point(t.LocalPosition)
		vector.StrokeCircle(
			screen,
			float32(center.X),
			float32(center.Y),
			float32(*r),
			2,
			color.RGBA{0xff, 0xff, 0, 0xff},
			false,
		)
	}
}
