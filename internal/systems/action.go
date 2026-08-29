package systems

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

func HandleActions(ecs *ecs.ECS) {
	mousePoint := cursorPoint()
	player := components.GetPlayer(ecs.World)

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		_, has_selected := components.Selected.First(ecs.World)
		if has_selected {
			player.ActionName = "Move"
		}
	}

	switch player.ActionName {
	case "Move":
		{
			if _, has_selection := components.Selected.First(ecs.World); has_selection {
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
					first, ok := components.FirstActorAtPoint(ecs.World, util.ToPoint(mousePoint))
					if ok {
						moveSelectedFollow(ecs.World, first.Entity(), DEFAULT_STOP_DISTANCE)
					} else {
						moveSelectedTo(ecs.World, mousePoint, DEFAULT_STOP_DISTANCE)
					}
				}
			}
		}
	}
}

func moveSelectedTo(world donburi.World, point math.Vec2, stopDistance float64) {
	push := ebiten.IsKeyPressed(ebiten.KeyShift)

	for selected := range components.Selected.Iter(world) {
		if push {
			components.PushMovement(selected, point, stopDistance)
			continue
		}

		components.WithMovementTo(selected, point, stopDistance)
	}
}

func moveSelectedFollow(world donburi.World, follow donburi.Entity, stopDistance float64) {
	for selected := range components.Selected.Iter(world) {
		if selected.Entity() == follow {
			continue
		}

		components.WithMovementFollow(selected, follow, stopDistance)
	}
}
