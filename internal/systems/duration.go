package systems

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi/ecs"
)

func DecrementDurations(ecs *ecs.ECS) {
	if ecs.IsPaused() {
		return
	}

	for entry := range components.DurationQuery.Iter(ecs.World) {
		duration := components.Duration.Get(entry)
		if *duration > 0 {
			*duration--
		}
	}
}

func RemoveCompleted(ecs *ecs.ECS) {
	for entry := range components.DurationQuery.Iter(ecs.World) {
		duration := *components.Duration.Get(entry)
		if duration == 0 {
			ecs.World.Remove(entry.Entity())
		}
	}
}
