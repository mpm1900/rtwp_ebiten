package systems

import (
	"math"
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

func HandleMovement(ecs *ecs.ECS) {
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
			if moveWithCollision(ecs.World, entry, delta) && distance <= components.CollisionStopDistance(entry, movement.StopDistance) {
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
	}
}

func moveWithCollision(world donburi.World, entry *donburi.Entry, delta dmath.Vec2) bool {
	if delta.IsZero() {
		return false
	}

	trans := transform.Transform.Get(entry)
	start_pos := trans.LocalPosition
	trans.LocalRotation = dmath.ToDegrees(math.Atan2(delta.Y, delta.X))

	full_position := components.ClampWorldPosition(start_pos.Add(delta))
	if isFreeAt(world, entry, full_position) {
		trans.LocalPosition = full_position
		return false
	}

	position := start_pos
	axes := [2]dmath.Vec2{{X: delta.X}, {Y: delta.Y}}
	for _, axis_delta := range axes {
		if axis_delta.IsZero() {
			continue
		}

		next_position := components.ClampWorldPosition(position.Add(axis_delta))
		if isFreeAt(world, entry, next_position) {
			position = next_position
		}
	}

	trans.LocalPosition = position
	return true
}

func isFreeAt(world donburi.World, entry *donburi.Entry, position dmath.Vec2) bool {
	_, colliding := components.CollidesAt(world, entry, position)
	return !colliding
}

func getSpeed(entry *donburi.Entry) float64 {
	speed := 0.0

	if entry.HasComponent(components.Stats) {
		stats := components.Stats.Get(entry)
		speed = stats.Stats[components.StatSpeed]
	}

	return speed
}
