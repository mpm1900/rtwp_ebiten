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
		was_blocked := actor.ActionCooldown > 0 || actor.ActionDelay > 0
		if actor.ActionCooldown > 0 {
			actor.ActionCooldown--
		}
		if actor.ActionDelay > 0 {
			actor.ActionDelay--
		}
		if actor.ActionCooldown > 0 || actor.ActionDelay > 0 {
			continue
		}

		if was_blocked || !actor.StartedAction {
			events.HandleActionQueue(ecs.World, entry.Entity())
		}
	}
}
