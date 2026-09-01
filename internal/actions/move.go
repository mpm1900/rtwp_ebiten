package actions

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/events"
	"rtwp_ebitengine/internal/pathing"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type MoveAction struct {
	components.ActionData
}

func (a MoveAction) Data() components.ActionData {
	return a.ActionData
}
func (a MoveAction) Publish(world donburi.World, point math.Vec2) {
	for selected := range components.Selected.Iter(world) {
		events.Actions.Publish(world, components.ActionEvent{
			Action: a,
			Source: selected.Entity(),
			Point:  point,
		})
	}
}
func (a MoveAction) Handle(world donburi.World, source donburi.Entity, point math.Vec2) {
	if !components.IsInWorld(point) {
		return
	}

	if _, has_selection := components.Selected.First(world); !has_selection {
		return
	}

	first, ok := components.FirstActorAtPoint(world, util.ToPoint(point))
	if ok {
		moveFollow(world, source, first.Entity(), components.DEFAULT_STOP_DISTANCE)
	} else {
		moveTo(world, source, point, components.DEFAULT_STOP_DISTANCE)
	}

}
func (a MoveAction) Valid(world donburi.World, point math.Vec2) bool {
	return components.IsInWorld(point)
}

var Move = MoveAction{
	Key:          ebiten.Key1,
	Name:         "Move",
	CursorOffset: math.NewVec2(-8, -8),
}

func moveTo(world donburi.World, source donburi.Entity, point math.Vec2, stopDistance float64) {
	push := ebiten.IsKeyPressed(ebiten.KeyShift)

	entry := world.Entry(source)
	start := components.Center(entry)
	if push && entry.HasComponent(components.Movement) {
		movement := components.Movement.Get(entry)
		if len(movement.Targets) > 0 {
			start = movement.Targets[len(movement.Targets)-1]
		}
	}

	path, ok := pathing.FindPath(world, start, point)
	if !ok || len(path) == 0 {
		path = []math.Vec2{point}
	}

	if push {
		components.PushMovementList(entry, path, stopDistance)
	} else {
		components.WithMovementList(entry, path, stopDistance)
	}
}

func moveFollow(world donburi.World, source donburi.Entity, follow donburi.Entity, stopDistance float64) {
	if source == follow {
		return
	}

	components.WithMovementFollow(world.Entry(source), follow, stopDistance)
}
