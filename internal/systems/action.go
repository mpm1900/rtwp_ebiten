package systems

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/events"

	"github.com/yohamta/donburi/ecs"
)

func TickActorActions(ecs *ecs.ECS) {
	if ecs.IsPaused() {
		return
	}

	for entry := range components.ActorQuery.Iter(ecs.World) {
		actor := components.Actor.Get(entry)

		active_event, ok := actor.PeekActionQueue()
		if !ok {
			actor.TickActionTimers()
			continue
		}

		was_blocked := actor.ActionDelay > 0 || actor.CooldownForAction(active_event.Action) > 0
		actor.TickActionTimers()
		if actor.ActionDelay > 0 || actor.CooldownForAction(active_event.Action) > 0 {
			continue
		}

		if was_blocked || !actor.ActionStarted {
			events.HandleActionQueue(ecs.World, entry.Entity())
		}
	}
}
