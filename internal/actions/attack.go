package actions

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/events"
	"rtwp_ebitengine/internal/util"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type AttackAction struct {
	components.ActionData
}

func (a AttackAction) Data() components.ActionData {
	return a.ActionData
}
func (a AttackAction) Publish(world donburi.World, event components.ActionEvent) {
	shift := slices.Contains(event.Keys, ebiten.KeyShift)
	for selected := range components.Selected.Iter(world) {
		actor := components.Actor.Get(selected)
		action_event := components.ActionEvent{
			Action: a,
			Source: selected.Entity(),
			Point:  event.Point,
		}

		if actor.QueueActionEvent(world, action_event, shift) {
			events.Actions.Publish(world, action_event)
		}
	}
}
func (a AttackAction) Handle(world donburi.World, event components.ActionEvent) {
	entry, ok := components.FirstActorAtPoint(world, util.ToPoint(event.Point))
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

func (a AttackAction) IsComplete(world donburi.World, source donburi.Entity) bool {
	if !world.Valid(source) {
		return true
	}

	entry := world.Entry(source)
	if !entry.HasComponent(components.Actor) {
		return true
	}

	actor := components.Actor.Get(entry)
	if actor.CooldownForAction(a) > 0 {
		return false
	}

	if actor.ActionQueueLen() == 1 {
		if target, ok := components.FirstTarget(world, entry); ok {
			actor.PushActionEvent(components.ActionEvent{
				Action: a,
				Source: source,
				Point:  components.Center(target),
			})
		}
	}
	return true
}
func (a AttackAction) Cancel(world donburi.World, source donburi.Entity) {
	if !world.Valid(source) {
		return
	}

	entry := world.Entry(source)
	if entry.HasComponent(components.Targets) {
		entry.RemoveComponent(components.Targets)
	}
}
func (a AttackAction) Valid(world donburi.World, point math.Vec2) bool {
	return components.IsInWorld(point)
}

var Attack = AttackAction{
	Key:      ebiten.Key2,
	Name:     "Attack",
	Delay:    10,
	Cooldown: 60,
}
