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
func (a MoveAction) Publish(world donburi.World, point math.Vec2, push bool) {
	for selected := range components.Selected.Iter(world) {
		actor := components.Actor.Get(selected)
		if push && a.pushActiveMove(world, selected, actor, point) {
			continue
		}

		event := components.ActionEvent{
			Action: a,
			Source: selected.Entity(),
			Point:  point,
		}

		if actor.QueueActionEvent(world, event, push) {
			events.Actions.Publish(world, event)
		}
	}
}
func (a MoveAction) Handle(world donburi.World, source donburi.Entity, point math.Vec2) {
	if !components.IsInWorld(point) {
		return
	}

	if !world.Valid(source) {
		return
	}

	f, ok := components.FirstActorAtPoint(world, util.ToPoint(point))
	if ok {
		follow := f.Entity()
		if follow == source {
			return
		}
		components.WithMovementFollow(world.Entry(source), follow, components.DEFAULT_STOP_DISTANCE)
	} else {
		moveTo(world, source, point, components.DEFAULT_STOP_DISTANCE)
	}

}
func (a MoveAction) IsComplete(world donburi.World, source donburi.Entity) bool {
	if !world.Valid(source) {
		return true
	}

	return !world.Entry(source).HasComponent(components.Movement)
}
func (a MoveAction) Cancel(world donburi.World, source donburi.Entity) {
	if !world.Valid(source) {
		return
	}

	entry := world.Entry(source)
	if entry.HasComponent(components.Movement) {
		entry.RemoveComponent(components.Movement)
	}
}
func (a MoveAction) pushActiveMove(world donburi.World, entry *donburi.Entry, actor *components.ActorData, point math.Vec2) bool {
	active_event, ok := actor.PeekActionQueue()
	if !ok || actor.ActionQueueLen() != 1 || !active_event.Started {
		return false
	}
	if _, ok := active_event.Action.(MoveAction); !ok {
		return false
	}
	if !entry.HasComponent(components.Movement) {
		return false
	}

	pushMoveTo(world, entry.Entity(), point, components.DEFAULT_STOP_DISTANCE)
	return true
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
	entry := world.Entry(source)
	start := components.Center(entry)

	path, ok := pathing.FindPath(world, start, point)
	if !ok || len(path) == 0 {
		path = []math.Vec2{point}
	}

	components.WithMovementList(entry, path, stopDistance)
}

func pushMoveTo(world donburi.World, source donburi.Entity, point math.Vec2, stopDistance float64) {
	entry := world.Entry(source)
	start := components.Center(entry)
	if entry.HasComponent(components.Movement) {
		movement := components.Movement.Get(entry)
		if len(movement.Targets) > 0 {
			start = movement.Targets[len(movement.Targets)-1]
		}
	}

	path, ok := pathing.FindPath(world, start, point)
	if !ok || len(path) == 0 {
		path = []math.Vec2{point}
	}

	components.PushMovementList(entry, path, stopDistance)
}
