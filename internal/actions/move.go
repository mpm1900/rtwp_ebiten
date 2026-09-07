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
	first_anchor := moveAnchor(first, shift)

	_, interact := components.FirstInteractableAtPoint(world, util.ToPoint(event.Position))
	for selected := range components.SelectedActorsQuery.Iter(world) {
		actor := components.Actor.Get(selected)
		point := event.Position

		if ctrl {
			anchor := moveAnchor(selected, shift)
			point = point.Add(anchor.Sub(first_anchor))
		}

		if shift && pushActiveMove(world, selected, actor, point, loop) {
			continue
		}

		event := components.ActionEvent{
			Action:   action,
			Source:   selected.Entity(),
			Position: point,
			Loop:     loop,
		}

		if interact {
			event.Action = Interact
		}

		actor.QueueActionEvent(world, event, shift)
	}
}
func (b MoveBehavior) Start(action *components.Action, world donburi.World, event components.ActionEvent) {
	if !components.IsInWorld(event.Position) {
		return
	}

	if !world.Valid(event.Source) {
		return
	}

	if f, ok := components.FirstActorAtPoint(world, event.Position); ok {
		follow := f.Entity()
		if follow == event.Source {
			return
		}
		components.WithMovement(world.Entry(event.Source), components.NewPathFollow(follow))
		return
	}

	if _, ok := components.FirstColliderAtPoint(world, event.Position); ok {
		return
	}

	if len(event.Path) > 0 {
		setMoveThroughPath(world, event.Source, event.Path, event.Position, event.Loop)
		return
	}

	setMoveTo(world, event.Source, event.Position, event.Loop)
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

func moveAnchor(entry *donburi.Entry, queued bool) math.Vec2 {
	if queued {
		if entry.HasComponent(components.Movement) {
			movement := components.Movement.Get(entry)
			if last, ok := movement.Last(); ok {
				return last
			}
		}
		if entry.HasComponent(components.Actor) {
			actor := components.Actor.Get(entry)
			if active, ok := actor.PeekActionQueue(); ok && active.Action == Move {
				if len(active.Path) > 0 {
					return active.Path[len(active.Path)-1]
				}
				return active.Position
			}
		}
	}

	return components.Center(entry)
}
func getPath(world donburi.World, start, point math.Vec2) []math.Vec2 {
	path, ok := pathing.FindPath(world, start, point)
	if !ok || len(path) == 0 {
		return []math.Vec2{}
	}

	return path
}
func setMoveTo(world donburi.World, source donburi.Entity, point math.Vec2, loop bool) {
	entry := world.Entry(source)
	start := components.Center(entry)

	path := getPath(world, start, point)
	if len(path) == 0 {
		return
	}
	components.WithMovement(entry, components.NewPathMovement(entry, path, loop))
}
func pushMoveTo(world donburi.World, source donburi.Entity, point math.Vec2, loop bool) {
	entry := world.Entry(source)
	start := components.Center(entry)
	movement := components.Movement.Get(entry)
	if len(movement.Path) > 0 {
		start, _ = movement.Last()
	}

	path := getPath(world, start, point)
	components.PushMovementList(entry, path, movement.StopDistance, loop)
}
func setMoveThroughPath(world donburi.World, source donburi.Entity, path []math.Vec2, final math.Vec2, loop bool) {
	entry := world.Entry(source)
	start := components.Center(entry)

	var full_path []math.Vec2
	for _, wp := range path {
		full_path = append(full_path, getPath(world, start, wp)...)
		start = wp
	}
	full_path = append(full_path, getPath(world, start, final)...)

	if len(full_path) == 0 {
		return
	}
	components.WithMovement(entry, components.NewPathMovement(entry, full_path, loop))
}
func pushActiveMove(world donburi.World, entry *donburi.Entry, actor *components.ActorData, point math.Vec2, loop bool) bool {
	active_event, ok := actor.PeekActionQueue()
	if !ok || actor.ActionQueueLen() != 1 || active_event.Action != Move {
		return false
	}

	if actor.ActionStarted {
		if !entry.HasComponent(components.Movement) {
			return false
		}
		pushMoveTo(world, entry.Entity(), point, loop)
		return true
	}

	active_event.Path = append(active_event.Path, point)
	return true
}
