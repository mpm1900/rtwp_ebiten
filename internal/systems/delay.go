package systems

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi/ecs"
)

func DecrementDelays(ecs *ecs.ECS) {
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
			entry.RemoveComponent(components.Delay)
		}
	}
}
