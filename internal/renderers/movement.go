package renderers

import (
	"image/color"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var renderMovementQuery = donburi.NewQuery(
	filter.Contains(components.Movement, transform.Transform, components.Selected),
)

func RenderMovement(ecs *ecs.ECS, screen *ebiten.Image) {
	lineColor := color.RGBA{0xff, 0xff, 0xff, 0xff}
	view := newCameraView(ecs)

	for entry := range renderMovementQuery.Iter(ecs.World) {
		movement := components.Movement.Get(entry)
		from := components.Center(entry)
		if movement.Follow != donburi.Null {
			to, ok := components.MovementPosition(ecs.World, movement)
			if !ok {
				continue
			}

			util.DrawPoints(screen, view.Point(from), view.Point(to), 1, lineColor)
			continue
		}

		for _, to := range movement.Targets {
			util.DrawPoints(screen, view.Point(from), view.Point(to), 1, lineColor)
			from = to
		}
	}
}
