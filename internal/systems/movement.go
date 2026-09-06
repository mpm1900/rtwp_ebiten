package systems

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/events"

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
		distance := movement.TargetDistance(ecs.World, entry)
		if distance > movement.StopDistance {
			step := min(getSpeed(entry), distance)
			delta, ok := movement.Delta(ecs.World, entry, step)
			if !ok {
				completed = append(completed, entry.Entity())
				continue
			}

			distance = movement.TargetDistance(ecs.World, entry)
			result := MoveWithCollision(ecs.World, entry, delta)
			if result.Collided && distance <= components.CollisionStopDistance(entry, movement.StopDistance) {
				if movement.Next() {
					completed = append(completed, entry.Entity())
				}
				continue
			}
		} else {
			if movement.Next() {
				completed = append(completed, entry.Entity())
			}
			continue
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
		if entry.HasComponent(components.Actor) {
			events.HandleActionQueue(ecs.World, entity)
		}
	}
}

func getSpeed(entry *donburi.Entry) float64 {
	speed := 0.0

	if entry.HasComponent(components.Stats) {
		stats := components.Stats.Get(entry)
		speed = stats.Stats[components.StatSpeed]
	}

	return speed
}
