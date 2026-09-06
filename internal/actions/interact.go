package actions

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type InteractBehavior struct{}

func (b InteractBehavior) Publish(action *components.Action, world donburi.World, event components.ActionEvent) {
	shift := slices.Contains(event.Keys, ebiten.KeyShift)
	for selected := range components.SelectedActorsQuery.Iter(world) {
		actor := components.Actor.Get(selected)
		point := event.Position

		// just add point to current path
		if shift && pushActiveMove(world, selected, actor, point, false) {
			continue
		}

		event := components.ActionEvent{
			Action:   action,
			Source:   selected.Entity(),
			Position: point,
			Loop:     false,
		}

		actor.QueueActionEvent(world, event, shift)
	}
}
func (b InteractBehavior) Start(action *components.Action, world donburi.World, event components.ActionEvent) {
	if !components.IsInWorld(event.Position) {
		return
	}

	if !world.Valid(event.Source) {
		return
	}

	point := util.ToPoint(event.Position)
	target, ok := components.FirstInteractableAtPoint(world, point)
	if !ok {
		return
	}

	entry := components.Interactable.Get(target)
	interact_point := entry.Point(components.Center(target))
	interact_range := components.DEFAULT_STOP_DISTANCE
	if target.HasComponent(components.Range) {
		interact_range = *components.Range.Get(target)
	}

	source := world.Entry(event.Source)
	source_center := components.Center(source)
	distance := interact_point.Distance(source_center)
	if distance > interact_range {
		setMoveTo(world, event.Source, interact_point, false)
		components.Actor.Get(source).PushNextActionEvent(event)
		return
	}

	entry.OnInteract(world, event.Source)
}
func (b InteractBehavior) Update(action *components.Action, world donburi.World, event components.ActionEvent) components.ActionStatus {
	if !world.Valid(event.Source) {
		return components.ActionComplete
	}

	if world.Entry(event.Source).HasComponent(components.Movement) {
		return components.ActionRunning
	}

	return components.ActionComplete
}
func (b InteractBehavior) Cancel(action *components.Action, world donburi.World, event components.ActionEvent) {
	if !world.Valid(event.Source) {
		return
	}

	entry := world.Entry(event.Source)
	if entry.HasComponent(components.Movement) {
		entry.RemoveComponent(components.Movement)
	}
}
func (b InteractBehavior) Valid(action *components.Action, world donburi.World, point math.Vec2) bool {
	_, ok := components.FirstInteractableAtPoint(world, util.ToPoint(point))
	return ok

}

var Interact = &components.Action{
	Key:      ebiten.Key1,
	Name:     "Interact",
	Behavior: InteractBehavior{},
}
