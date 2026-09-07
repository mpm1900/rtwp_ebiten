package systems

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi/ecs"
)

// TickActorActions ticks cooldown timers and advances queued actions for every actor.
func TickActorActions(ecs *ecs.ECS) {
	components.RebuildCollisionIndex(ecs.World)
	if ecs.IsPaused() {
		return
	}

	for entry := range components.ActorQuery.Iter(ecs.World) {
		actor := components.Actor.Get(entry)
		actor.TickActionTimers()
		if !actor.HasAction() {
			continue
		}

		components.AdvanceActionQueue(ecs.World, entry)
	}
}
