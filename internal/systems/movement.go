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
		budget := getSpeed(entry)

		for budget > 0 {
			distance := movement.TargetDistance(ecs.World, entry)
			if distance <= movement.StopDistance {
				if movement.Next() {
					completed = append(completed, entry.Entity())
					break
				}
				continue
			}

			step := min(budget, distance)
			delta, ok := movement.Delta(ecs.World, entry, step)
			if !ok {
				completed = append(completed, entry.Entity())
				break
			}

			budget -= step
			if !moveWithCollision(ecs.World, entry, delta) {
				continue
			}

			if movement.TargetDistance(ecs.World, entry) <= components.CollisionStopDistance(entry, movement.StopDistance) {
				if movement.Next() {
					completed = append(completed, entry.Entity())
					break
				}
				continue
			}
			break
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

	// Fast path: try the full delta in one collision check.
	full_position := components.ClampWorldPosition(start_pos.Add(delta))
	if isFreeAt(world, entry, full_position) {
		trans.LocalPosition = full_position
		return false
	}

	// Blocked: sub-step to find how far we can get before the collision.
	magnitude := delta.Magnitude()
	direction := delta.Normalized()
	position := start_pos
	traveled := 0.0
	for traveled < magnitude {
		step := min(1.0, magnitude-traveled)
		next := components.ClampWorldPosition(position.Add(direction.MulScalar(step)))
		if !isFreeAt(world, entry, next) {
			break
		}
		position = next
		traveled += step
	}

	// Slide along axes for the remaining distance.
	remaining := magnitude - traveled
	if remaining > 0 {
		axes := [2]dmath.Vec2{{X: direction.X * remaining}, {Y: direction.Y * remaining}}
		for _, axis_delta := range axes {
			if axis_delta.IsZero() {
				continue
			}
			next := components.ClampWorldPosition(position.Add(axis_delta))
			if isFreeAt(world, entry, next) {
				position = next
			}
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

	return speed / 5
}
