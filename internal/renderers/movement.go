package renderers

import (
	"image/color"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func RenderMovement(ecs *ecs.ECS, screen *ebiten.Image) {
	lineColor := color.RGBA{0xff, 0xff, 0xff, 0xff}

	for entry := range components.MovementQuery.Iter(ecs.World) {
		movement := components.Movement.Get(entry)
		from := components.Center(entry)
		if movement.Follow != donburi.Null {
			to, ok := components.MovementTarget(ecs.World, movement)
			if !ok {
				continue
			}

			util.DrawPoints(screen, from, to, 1, lineColor)
			continue
		}

		for _, to := range movement.Targets {
			util.DrawPoints(screen, from, to, 1, lineColor)
			from = to
		}
	}
}
