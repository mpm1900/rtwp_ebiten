package actions

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/events"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type AttackBehavior struct{}

func (b AttackBehavior) Publish(action *components.Action, world donburi.World, event components.ActionEvent) {
	shift := slices.Contains(event.Keys, ebiten.KeyShift)
	for selected := range components.Selected.Iter(world) {
		actor := components.Actor.Get(selected)
		action_event := components.ActionEvent{
			Action:   action,
			Source:   selected.Entity(),
			Position: event.Position,
		}

		if action.Range != nil {
			center := components.Center(selected)
			distance := center.Distance(event.Position)
			if distance < action.Range[0] {
				return
			}
			if distance > action.Range[1] {
				direction := center.Sub(event.Position).Normalized()
				move_point := event.Position.Add(direction.MulScalar(action.Range[1]))
				move_event := components.ActionEvent{
					Action:   Move,
					Source:   selected.Entity(),
					Position: move_point,
				}
				if shift {
					actor.QueueActionEvent(world, move_event, true)
					actor.QueueActionEvent(world, action_event, true)
				} else {
					actor.SetActionEvent(world, move_event)
					actor.PushNextActionEvent(action_event)
				}
				continue
			}
		}

		actor.QueueActionEvent(world, action_event, shift)
	}
}
func (b AttackBehavior) Start(action *components.Action, world donburi.World, event components.ActionEvent) {
	entry, ok := components.FirstActorAtPoint(world, event.Position)
	if !ok {
		return
	}
	if !world.Valid(event.Source) {
		return
	}

	source := world.Entry(event.Source)
	components.WithTargets(source, entry.Entity())
	events.DamageAt.Publish(world, events.DamageEvent{
		Source:   event.Source,
		Position: event.Position,
		Action:   action,
	})
}

func (b AttackBehavior) Update(action *components.Action, world donburi.World, event components.ActionEvent) components.ActionStatus {
	if !world.Valid(event.Source) {
		return components.ActionComplete
	}

	entry := world.Entry(event.Source)
	if !entry.HasComponent(components.Actor) {
		return components.ActionComplete
	}

	actor := components.Actor.Get(entry)
	if actor.CooldownForAction(action) > 0 {
		return components.ActionRunning
	}

	if actor.ActionQueueLen() == 1 {
		if target, ok := components.FirstTarget(world, entry); ok {
			actor.PushActionEvent(components.ActionEvent{
				Action:   action,
				Source:   event.Source,
				Position: components.Center(target),
			})
		}
	}
	return components.ActionComplete
}
func (b AttackBehavior) Cancel(action *components.Action, world donburi.World, event components.ActionEvent) {
	if !world.Valid(event.Source) {
		return
	}

	entry := world.Entry(event.Source)
	if entry.HasComponent(components.Targets) {
		entry.RemoveComponent(components.Targets)
	}
}
func (b AttackBehavior) Valid(action *components.Action, world donburi.World, point math.Vec2) bool {
	return components.IsInWorld(point)
}

var Attack = &components.Action{
	Behavior: AttackBehavior{},
	Key:      ebiten.Key2,
	Name:     "Attack",
	Cooldown: 60,
	Delay:    10,
	Power:    50,
	Range:    []float64{5, 55},
}
