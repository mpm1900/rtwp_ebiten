package systems

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func MoveEntities(ecs *ecs.ECS) {
	completed := []donburi.Entity{}
	if ecs.IsPaused() {
		return
	}

	for entry := range components.MovementQuery.Iter(ecs.World) {
		movement := components.Movement.Get(entry)
		target, ok := components.MovementPosition(ecs.World, movement)
		if !ok {
			completed = append(completed, entry.Entity())
			continue
		}

		position := components.Center(entry)
		direction := target.Sub(position)
		distance := direction.Magnitude()
		stopDistance := movement.StopDistance
		if stopDistance <= 0 {
			stopDistance = components.DEFAULT_STOP_DISTANCE
		}

		remainingDistance := distance
		if distance > stopDistance {
			step := min(components.GetSpeed(entry), distance)
			moveResult := MoveWithCollision(ecs.World, entry, direction.Normalized().MulScalar(step))
			remainingDistance = components.Center(entry).Distance(target)
			if moveResult.Collided && remainingDistance <= components.CollisionStopDistance(entry, stopDistance) {
				if movement.NextTarget() {
					completed = append(completed, entry.Entity())
				}
				continue
			}
		}

		if remainingDistance <= stopDistance {
			if movement.NextTarget() {
				completed = append(completed, entry.Entity())
			}
		}
	}

	for _, entity := range completed {
		if !ecs.World.Valid(entity) {
			continue
		}

		entry := ecs.World.Entry(entity)
		if entry.HasComponent(components.Movement) {
			entry.RemoveComponent(components.Movement)
		}
	}
}
