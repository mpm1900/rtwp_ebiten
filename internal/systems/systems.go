package systems

import (
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi/ecs"
)

func LoadSystems(ecs *ecs.ECS, frame *util.Frame) {
	ecs.AddSystem(DecrementDelays)
	ecs.AddSystem(RemoveCompletedDelays)
	ecs.AddSystem(DecrementDurations)
	ecs.AddSystem(RemoveCompleted)
	ecs.AddSystem(ResolveModifiers(frame))

	ecs.AddSystem(MoveEntities)
}
