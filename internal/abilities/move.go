package abilities

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

var Move = components.Ability{
	AbilityID: uuid.New(),
	Key:       ebiten.Key1,
	Name:      "Move",
	Handle: func(ecs *ecs.ECS) {
		screenPoint := util.CursorPoint()
		player := components.GetPlayer(ecs.World)
		if player == nil {
			return
		}
		worldPoint := player.ScreenToWorld(screenPoint)

		if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			return
		}

		if _, has_selection := components.Selected.First(ecs.World); !has_selection {
			return
		}

		first, ok := components.FirstActorAtPoint(ecs.World, util.ToPoint(worldPoint))
		if ok {
			moveSelectedFollow(ecs.World, first.Entity(), components.DEFAULT_STOP_DISTANCE)
		} else {
			moveSelectedTo(ecs.World, worldPoint, components.DEFAULT_STOP_DISTANCE)
		}
	},
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
