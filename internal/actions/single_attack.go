package actions

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/events"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type SingleAttackBehavior struct{}

func (b SingleAttackBehavior) PublishSingle(action *components.Action, world donburi.World, event components.ActionEvent, entry *donburi.Entry) {
	actor := components.Actor.Get(entry)
	action_event := components.ActionEvent{
		Action:   action,
		Source:   entry.Entity(),
		Position: event.Position,
	}

	shift := slices.Contains(event.Keys, ebiten.KeyShift)
	if action.Range != nil {
		center := components.Center(entry)
		distance := center.Distance(event.Position)
		if distance < action.Range[0] {
			return
		}
		if distance > action.Range[1] {
			direction := center.Sub(event.Position).Normalized()
			move_point := event.Position.Add(direction.MulScalar(action.Range[1]))
			move_event := components.ActionEvent{
				Action:   Move,
				Source:   entry.Entity(),
				Position: move_point,
			}
			if shift {
				actor.QueueActionEvent(world, move_event, true)
				actor.QueueActionEvent(world, action_event, true)
			} else {
				actor.SetActionEvent(world, move_event)
				actor.PushNextActionEvent(action_event)
			}
			return
		}
	}

	actor.QueueActionEvent(world, action_event, shift)
}

func (b SingleAttackBehavior) Publish(action *components.Action, world donburi.World, event components.ActionEvent) {
	for selected := range components.Selected.Iter(world) {
		b.PublishSingle(action, world, event, selected)
	}
}

func (b SingleAttackBehavior) Start(action *components.Action, world donburi.World, event components.ActionEvent) {
	if !world.Valid(event.Source) {
		return
	}

	source := world.Entry(event.Source)
	target, ok := b.resolveAttackTarget(world, source, event.Position)
	if !ok {
		return
	}

	components.WithTargets(source, target.Entity())
	events.DamageAt.Publish(world, events.DamageEvent{
		Source:   event.Source,
		Position: components.Center(target),
		Action:   action,
	})
}

func (b SingleAttackBehavior) Update(action *components.Action, world donburi.World, event components.ActionEvent) components.ActionStatus {
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

	b.chainNextAttack(action, world, event, entry)
	return components.ActionComplete
}

func (b SingleAttackBehavior) Cancel(action *components.Action, world donburi.World, event components.ActionEvent) {
	if !world.Valid(event.Source) {
		return
	}

	entry := world.Entry(event.Source)
	if entry.HasComponent(components.Targets) {
		entry.RemoveComponent(components.Targets)
	}
}

func (b SingleAttackBehavior) Valid(action *components.Action, world donburi.World, point math.Vec2) bool {
	return components.IsInWorld(point)
}

func (b SingleAttackBehavior) chainNextAttack(
	action *components.Action,
	world donburi.World,
	event components.ActionEvent,
	entry *donburi.Entry,
) {
	actor := components.Actor.Get(entry)
	if actor.ActionQueueLen() != 1 {
		return
	}

	components.PruneTargets(world, entry)

	if target, ok := components.FirstTarget(world, entry); ok && b.isValidAttackTarget(world, entry, target, action) {
		b.pushAttackOnTarget(action, world, event, entry, target)
		return
	}

	target, ok := components.NearestLivingActorFromEntity(world, entry)
	if !ok {
		return
	}

	b.pushAttackOnTarget(action, world, event, entry, target)
}

func (b SingleAttackBehavior) pushAttackOnTarget(
	action *components.Action,
	world donburi.World,
	event components.ActionEvent,
	entry *donburi.Entry,
	target *donburi.Entry,
) {
	actor := components.Actor.Get(entry)
	target_center := components.Center(target)
	action_event := components.ActionEvent{
		Action:   action,
		Source:   event.Source,
		Position: target_center,
	}

	if action.Range != nil {
		distance := components.Center(entry).Distance(target_center)
		if distance < action.Range[0] {
			return
		}
		if distance > action.Range[1] {
			direction := components.Center(entry).Sub(target_center).Normalized()
			move_point := target_center.Add(direction.MulScalar(action.Range[1]))
			actor.PushActionEvent(components.ActionEvent{
				Action:   Move,
				Source:   event.Source,
				Position: move_point,
			})
		}
	}

	components.WithTargets(entry, target.Entity())
	actor.PushActionEvent(action_event)
}

func (b SingleAttackBehavior) resolveAttackTarget(
	world donburi.World,
	source *donburi.Entry,
	position math.Vec2,
) (*donburi.Entry, bool) {
	if target, ok := components.FirstActorAtPoint(world, position); ok && components.IsAlive(target) {
		return target, true
	}

	return components.NearestLivingActorFromEntity(world, source)
}

func (b SingleAttackBehavior) isValidAttackTarget(
	world donburi.World,
	source *donburi.Entry,
	target *donburi.Entry,
	action *components.Action,
) bool {
	if target == nil || target.Entity() == source.Entity() {
		return false
	}
	if !world.Valid(target.Entity()) || !target.HasComponent(components.Actor) {
		return false
	}
	if !components.IsAlive(target) {
		return false
	}
	if action.Range == nil {
		return true
	}

	distance := components.Center(source).Distance(components.Center(target))
	return distance >= action.Range[0] && distance <= action.Range[1]
}

var SingleAttack = &components.Action{
	Behavior: SingleAttackBehavior{},
	Key:      ebiten.Key2,
	Name:     "Single Attack",
	Accuracy: 100,
	Cooldown: 60,
	Delay:    10,
	Power:    50,
	Range:    []float64{5, 55},
}
