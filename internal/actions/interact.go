package actions

import (
	"fmt"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/events"
	"rtwp_ebitengine/internal/util"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type InteractAction struct {
	components.ActionData
}

func (a InteractAction) Data() components.ActionData {
	return a.ActionData
}
func (a InteractAction) Publish(world donburi.World, event components.ActionEvent) {
	shift := slices.Contains(event.Keys, ebiten.KeyShift)
	for selected := range components.SelectedActorsQuery.Iter(world) {
		actor := components.Actor.Get(selected)
		point := event.Point

		// just add point to current path
		if shift && pushActiveMove(world, selected, actor, point, false) {
			continue
		}

		event := components.ActionEvent{
			Action: a,
			Source: selected.Entity(),
			Point:  point,
			Loop:   false,
		}

		if actor.QueueActionEvent(world, event, shift) {
			events.Actions.Publish(world, event)
		}
	}
}
func (a InteractAction) Handle(world donburi.World, event components.ActionEvent) {
	if !components.IsInWorld(event.Point) {
		return
	}

	if !world.Valid(event.Source) {
		return
	}

	point := util.ToPoint(event.Point)
	target, ok := components.FirstInteractableAtPoint(world, point)
	if !ok {
		return
	}

	interact_point := components.Interactable.Get(target).Point(components.Center(target))
	interact_range := components.DEFAULT_STOP_DISTANCE
	if target.HasComponent(components.Range) {
		interact_range = *components.Range.Get(target)
	}

	source := world.Entry(event.Source)
	source_center := components.Center(source)
	distance := interact_point.Distance(source_center)
	if distance > interact_range {
		moveTo(world, event.Source, interact_point, components.DEFAULT_STOP_DISTANCE, false)
		components.Actor.Get(source).PushNextActionEvent(event)
		return
	}

	fmt.Println("dong the interact")
}
func (a InteractAction) IsComplete(world donburi.World, source donburi.Entity) bool {
	if !world.Valid(source) {
		return true
	}

	return !world.Entry(source).HasComponent(components.Movement) // interacting?
}
func (a InteractAction) Cancel(world donburi.World, source donburi.Entity) {
	if !world.Valid(source) {
		return
	}

	entry := world.Entry(source)
	if entry.HasComponent(components.Movement) {
		entry.RemoveComponent(components.Movement)
	}
}
func (a InteractAction) Valid(world donburi.World, point math.Vec2) bool {
	_, ok := components.FirstInteractableAtPoint(world, util.ToPoint(point))
	return ok

}

var Interact = InteractAction{
	Key:  ebiten.Key1,
	Name: "Interact",
}
