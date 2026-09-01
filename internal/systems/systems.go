package systems

import (
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi/ecs"
)

func Load(ecs *ecs.ECS, frame *util.Frame) {
	ecs.AddSystem(HandlePause)
	ecs.AddSystem(HandleInput)

	ecs.AddSystem(DecrementDelays)
	ecs.AddSystem(RemoveCompletedDelays)
	ecs.AddSystem(DecrementDurations)
	ecs.AddSystem(RemoveCompleted)
	ecs.AddSystem(ResolveModifiers(frame))

	ecs.AddSystem(TickActorActions)
	ecs.AddSystem(MoveEntities)
}
