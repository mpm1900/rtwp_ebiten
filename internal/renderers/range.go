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
	for entry := range components.RangeQuery.Iter(ecs.World) {
		t := transform.GetTransform(entry)
		r := components.Range.Get(entry)
		vector.StrokeCircle(
			screen,
			float32(t.LocalPosition.X),
			float32(t.LocalPosition.Y),
			float32(*r),
			2,
			color.RGBA{0xff, 0xff, 0, 0xff},
			false,
		)
	}
}
