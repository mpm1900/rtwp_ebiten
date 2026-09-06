package components

import (
	"github.com/yohamta/donburi"
)

const maxActionSteps = 64

// AdvanceActionQueue runs the action state machine for one actor.
//
// Pipeline:
//  1. Publish (input) - queue ActionEvents on ActorData
//  2. Advance (here) - start, update, cleanup, pop, repeat
//  3. TickActorActions (system) - tick cooldowns and call Advance each frame
func AdvanceActionQueue(world donburi.World, entry *donburi.Entry) {
	for range maxActionSteps {
		if !advanceActionStep(world, entry) {
			return
		}
	}
}

func advanceActionStep(world donburi.World, entry *donburi.Entry) bool {
	if !entry.HasComponent(Actor) {
		return false
	}

	actor := Actor.Get(entry)
	active_event, ok := actor.PeekActionQueue()
	if !ok {
		return false
	}
	if active_event.Action == nil {
		actor.NextActionEvent()
		return true
	}

	event := *active_event
	if !actor.ActionStarted {
		if !startAction(world, entry, actor, event) {
			return false
		}
	}

	switch event.Action.Update(world, event) {
	case ActionRunning:
		return false
	case ActionComplete, ActionCanceled:
		event.Action.Cancel(world, event)
		actor.NextActionEvent()
		return true
	default:
		return false
	}
}

func startAction(world donburi.World, entry *donburi.Entry, actor *ActorData, event ActionEvent) bool {
	if actor.CooldownForAction(event.Action) > 0 {
		return false
	}

	if entry.HasComponent(Delay) {
		if delay := *Delay.Get(entry); delay > 0 {
			return false
		}

		entry.RemoveComponent(Delay)
	} else if delay := event.Action.Delay; delay > 0 {
		WithDelay(entry, delay)
		return false
	}

	actor.ActionStarted = true
	event.Action.Start(world, event)
	actor.SetActionCooldown(event.Action)
	return true
}
