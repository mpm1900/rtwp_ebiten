package systems

import (
	"rtwp_ebitengine/internal/effects"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi/ecs"
)

func ResolveModifiers(frame *util.Frame) ecs.System {
	return func(ecs *ecs.ECS) {
		effects.ResolveModifiers(ecs, frame)
	}
}
