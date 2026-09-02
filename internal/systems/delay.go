package systems

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi/ecs"
)

func DecrementDelays(ecs *ecs.ECS) {
	if ecs.IsPaused() {
		return
	}

	for entry := range components.Delay.Iter(ecs.World) {
		delay := components.Delay.Get(entry)
		if *delay > 0 {
			*delay--
		}
	}
}

func RemoveCompletedDelays(ecs *ecs.ECS) {
	for entry := range components.Delay.Iter(ecs.World) {
		delay := *components.Delay.Get(entry)
		if delay == 0 {
			if entry.HasComponent(components.Actor) {
				actor := components.Actor.Get(entry)
				if actor.HasAction() && !actor.ActionStarted {
					continue
				}
			}

			entry.RemoveComponent(components.Delay)
		}
	}
}
