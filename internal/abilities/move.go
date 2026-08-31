package abilities

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/pathing"
	"rtwp_ebitengine/internal/util"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

var Move = components.Ability{
	AbilityID:    uuid.New(),
	Key:          ebiten.Key1,
	Name:         "Move",
	CursorOffset: math.NewVec2(-8, -8),
	Handle: func(ecs *ecs.ECS, screenPoint math.Vec2) {
		if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			return
		}

		if components.IsOverMinimap(screenPoint, components.MinimapRect()) {
			return
		}

		player := components.GetPlayer(ecs.World)
		if player == nil {
			return
		}
		worldPoint := player.ScreenToWorld(screenPoint)
		if !components.IsInWorld(worldPoint) {
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
	Valid: func(ecs *ecs.ECS, screenPoint math.Vec2) bool {
		if components.IsOverMinimap(screenPoint, components.MinimapRect()) {
			return false
		}
		player := components.GetPlayer(ecs.World)
		if player == nil {
			return false
		}
		worldPoint := player.ScreenToWorld(screenPoint)
		return components.IsInWorld(worldPoint)
	},
}

func moveSelectedTo(world donburi.World, point math.Vec2, stopDistance float64) {
	push := ebiten.IsKeyPressed(ebiten.KeyShift)

	for selected := range components.Selected.Iter(world) {
		start := components.Center(selected)
		if push && selected.HasComponent(components.Movement) {
			movement := components.Movement.Get(selected)
			if len(movement.Targets) > 0 {
				start = movement.Targets[len(movement.Targets)-1]
			}
		}

		path, ok := pathing.FindPath(world, start, point)
		if !ok || len(path) == 0 {
			path = []math.Vec2{point}
		}

		if push {
			components.PushMovementList(selected, path, stopDistance)
		} else {
			components.WithMovementList(selected, path, stopDistance)
		}
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
