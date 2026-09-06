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
			Action: action,
			Source: selected.Entity(),
			Point:  event.Point,
		}

		actor.QueueAndAdvance(world, selected, action_event, shift)
	}
}
func (b AttackBehavior) Start(action *components.Action, world donburi.World, event components.ActionEvent) {
	entry, ok := components.FirstActorAtPoint(world, event.Point)
	if !ok {
		return
	}
	if !world.Valid(event.Source) {
		return
	}

	source := world.Entry(event.Source)
	components.WithTargets(source, entry.Entity())
	events.DamageAt.Publish(world, events.DamageEvent{
		Point:  event.Point,
		Amount: 10,
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
				Action: action,
				Source: event.Source,
				Point:  components.Center(target),
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
	Key:      ebiten.Key2,
	Name:     "Attack",
	Delay:    10,
	Cooldown: 60,
	Behavior: AttackBehavior{},
}
