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
	event, ok := actor.PeekActionQueue()
	if !ok {
		return false
	}
	if event.Action == nil {
		actor.NextActionEvent()
		return true
	}

	active_event := *event

	if !actor.ActionStarted {
		if !actor.StartAction(world, entry, active_event) {
			return false
		}
	}

	switch active_event.Action.Update(world, active_event) {
	case ActionRunning:
		return false
	case ActionComplete, ActionCanceled:
		active_event.Action.Cancel(world, active_event)
		actor.EndAction(world, entry, &active_event)
		actor.NextActionEvent()
		return true
	default:
		return false
	}
}
