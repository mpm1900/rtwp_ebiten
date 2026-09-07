package renderers

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi/ecs"
)

func RenderRanges(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)

	for entry := range components.RangeQuery.Iter(ecs.World) {
		r := components.Range.Get(entry)
		center := view.Point(components.Center(entry))
		vector.StrokeCircle(
			screen,
			float32(center.X),
			float32(center.Y),
			float32(*r),
			2,
			assets.ColorRange,
			true,
		)
	}

	player := components.GetPlayer(ecs.World)
	if player.SelectedAction != nil && player.SelectedAction.Range != nil {
		for selected := range components.SelectedActorsQuery.Iter(ecs.World) {
			inside := player.SelectedAction.Range[0]
			outside := player.SelectedAction.Range[1]
			center := view.Point(components.Center(selected))
			vector.StrokeCircle(
				screen,
				float32(center.X),
				float32(center.Y),
				float32(inside),
				2,
				assets.ColorRange,
				true,
			)
			vector.StrokeCircle(
				screen,
				float32(center.X),
				float32(center.Y),
				float32(outside),
				2,
				assets.ColorRange,
				true,
			)
		}
	}
}
