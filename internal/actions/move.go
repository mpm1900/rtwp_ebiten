package actions

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/pathing"
	"rtwp_ebitengine/internal/util"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type MoveBehavior struct{}

func (b MoveBehavior) Publish(action *components.Action, world donburi.World, event components.ActionEvent) {
	first, ok := components.SelectedActorsQuery.First(world)
	if !ok {
		return
	}

	shift := slices.Contains(event.Keys, ebiten.KeyShift)
	ctrl := slices.Contains(event.Keys, ebiten.KeyControl)
	loop := slices.Contains(event.Keys, ebiten.KeyZ)

	_, interact := components.FirstInteractableAtPoint(world, util.ToPoint(event.Point))
	for selected := range components.SelectedActorsQuery.Iter(world) {
		actor := components.Actor.Get(selected)
		point := event.Point

		if ctrl {
			first_center := components.Center(first)
			center := components.Center(selected)
			point = point.Add(center.Sub(first_center))
		}

		if shift && pushActiveMove(world, selected, actor, point, loop) {
			continue
		}

		event := components.ActionEvent{
			Action: action,
			Source: selected.Entity(),
			Point:  point,
			Loop:   loop,
		}

		if interact {
			event.Action = Interact
		}

		actor.QueueAndAdvance(world, selected, event, shift)
	}
}
func (b MoveBehavior) Start(action *components.Action, world donburi.World, event components.ActionEvent) {
	if !components.IsInWorld(event.Point) {
		return
	}

	if !world.Valid(event.Source) {
		return
	}

	if f, ok := components.FirstActorAtPoint(world, event.Point); ok {
		follow := f.Entity()
		if follow == event.Source {
			return
		}
		components.WithMovement(world.Entry(event.Source), components.NewPathFollow(follow))
		return
	}

	if _, ok := components.FirstColliderAtPoint(world, event.Point); ok {
		return
	}

	moveTo(world, event.Source, event.Point, components.DEFAULT_STOP_DISTANCE, event.Loop)
}
func (b MoveBehavior) Update(action *components.Action, world donburi.World, event components.ActionEvent) components.ActionStatus {
	if !world.Valid(event.Source) {
		return components.ActionComplete
	}

	if world.Entry(event.Source).HasComponent(components.Movement) {
		return components.ActionRunning
	}

	return components.ActionComplete
}
func (b MoveBehavior) Cancel(action *components.Action, world donburi.World, event components.ActionEvent) {
	if !world.Valid(event.Source) {
		return
	}

	entry := world.Entry(event.Source)
	if entry.HasComponent(components.Movement) {
		entry.RemoveComponent(components.Movement)
	}
}
func (b MoveBehavior) Valid(action *components.Action, world donburi.World, point math.Vec2) bool {
	return components.IsInWorld(point)
}

var Move = &components.Action{
	Key:          ebiten.Key1,
	Name:         "Move",
	CursorOffset: math.NewVec2(-8, -8),
	Behavior:     MoveBehavior{},
}

func pushActiveMove(world donburi.World, entry *donburi.Entry, actor *components.ActorData, point math.Vec2, loop bool) bool {
	active_event, ok := actor.PeekActionQueue()
	if !ok || actor.ActionQueueLen() != 1 || !actor.ActionStarted {
		return false
	}
	if active_event.Action != Move {
		return false
	}
	if !entry.HasComponent(components.Movement) {
		return false
	}

	pushMoveTo(world, entry.Entity(), point, components.DEFAULT_STOP_DISTANCE, loop)
	return true
}
func moveTo(world donburi.World, source donburi.Entity, point math.Vec2, stopDistance float64, loop bool) {
	entry := world.Entry(source)
	start := components.Center(entry)

	path, ok := pathing.FindPath(world, start, point)
	if !ok || len(path) == 0 {
		path = []math.Vec2{point}
	}
	if loop {
		path = components.AppendLoopOrigin(path, start)
	}

	components.WithMovement(entry, components.NewPathMovement(entry, path, loop))
}
func pushMoveTo(world donburi.World, source donburi.Entity, point math.Vec2, stopDistance float64, loop bool) {
	entry := world.Entry(source)
	start := components.Center(entry)
	if entry.HasComponent(components.Movement) {
		movement := components.Movement.Get(entry)
		if len(movement.Path) > 0 {
			start, _ = movement.Last()
		}
	}

	path, ok := pathing.FindPath(world, start, point)
	if !ok || len(path) == 0 {
		path = []math.Vec2{point}
	}
	if loop {
		path = components.AppendLoopOrigin(path, components.LoopOriginForEntry(entry))
	}

	components.PushMovementList(entry, path, stopDistance, loop)
}
