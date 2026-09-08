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
	components.RebuildCollisionIndex(ecs.World)
	if ecs.IsPaused() {
		return
	}

	for entry := range components.MovementQuery.Iter(ecs.World) {
		entity := entry.Entity()
		if !ecs.World.Valid(entity) {
			continue
		}

		movement := components.Movement.Get(entry)
		budget := getSpeed(entry)

		for budget > 0 {
			if !ecs.World.Valid(entity) {
				break
			}
			entry = ecs.World.Entry(entity)

			distance := movement.TargetDistance(ecs.World, entry)
			if movement.Follow == donburi.Null && distance <= movement.StopDistance {
				if movement.Next() {
					completed = append(completed, entity)
					break
				}
				continue
			}

			step := min(budget, distance)
			delta, ok := movement.Delta(ecs.World, entry, step)
			if !ok {
				completed = append(completed, entity)
				break
			}

			budget -= step
			if !moveWithCollision(ecs.World, entity, delta) {
				continue
			}

			if !ecs.World.Valid(entity) {
				break
			}
			entry = ecs.World.Entry(entity)

			if movement.TargetDistance(ecs.World, entry) <= components.CollisionStopDistance(entry, movement.StopDistance) {
				if movement.Next() {
					completed = append(completed, entity)
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

func moveWithCollision(world donburi.World, entity donburi.Entity, delta dmath.Vec2) bool {
	if delta.IsZero() || !world.Valid(entity) {
		return false
	}

	entry := world.Entry(entity)
	trans := transform.Transform.Get(entry)
	start_pos := trans.LocalPosition
	trans.LocalRotation = dmath.ToDegrees(math.Atan2(delta.Y, delta.X))

	full_position := components.ClampWorldPosition(start_pos.Add(delta))
	if isFreeAt(world, entity, full_position) {
		if !world.Valid(entity) {
			return false
		}

		transform.Transform.Get(world.Entry(entity)).LocalPosition = full_position
		return false
	}

	magnitude := delta.Magnitude()
	direction := delta.Normalized()
	position := start_pos
	traveled := 0.0
	for traveled < magnitude {
		if !world.Valid(entity) {
			return false
		}

		step := min(1.0, magnitude-traveled)
		next := components.ClampWorldPosition(position.Add(direction.MulScalar(step)))
		if !isFreeAt(world, entity, next) {
			break
		}
		if !world.Valid(entity) {
			return false
		}

		position = next
		traveled += step
	}

	remaining := magnitude - traveled
	if remaining > 0 && world.Valid(entity) {
		axes := [2]dmath.Vec2{{X: direction.X * remaining}, {Y: direction.Y * remaining}}
		for _, axis_delta := range axes {
			if axis_delta.IsZero() || !world.Valid(entity) {
				continue
			}

			next := components.ClampWorldPosition(position.Add(axis_delta))
			if isFreeAt(world, entity, next) {
				if !world.Valid(entity) {
					return false
				}
				position = next
			}
		}
	}

	if !world.Valid(entity) {
		return false
	}

	transform.Transform.Get(world.Entry(entity)).LocalPosition = position
	return true
}

func isFreeAt(world donburi.World, entity donburi.Entity, position dmath.Vec2) bool {
	if !world.Valid(entity) {
		return true
	}

	entry := world.Entry(entity)
	follow := donburi.Null
	if entry.HasComponent(components.Movement) {
		follow = components.Movement.Get(entry).Follow
	}

	other, colliding := components.CollidesAt(world, entry, position)
	if colliding && other != nil {
		if entry.HasComponent(components.Collision) {
			collision := components.Collision.Get(entry)
			if collision.OnCollide != nil {
				collision.OnCollide(world, entry, other.Entity())
			}
		}
		if world.Valid(other.Entity()) && other.HasComponent(components.Collision) {
			collision := components.Collision.Get(other)
			if collision.OnCollide != nil {
				collision.OnCollide(world, other, entity)
			}
		}
	}

	if !world.Valid(entity) {
		return true
	}

	if follow != donburi.Null && colliding && other != nil {
		if other.HasComponent(components.Actor) {
			return false
		}
		return true
	}

	return !colliding
}

func getSpeed(entry *donburi.Entry) float64 {
	speed := 0.0

	if entry.HasComponent(components.Movement) {
		movement := components.Movement.Get(entry)
		if movement.Speed > 0 {
			speed = movement.Speed
		}
	}

	if entry.HasComponent(components.Stats) {
		stats := components.Stats.Get(entry)
		speed = stats.Stats[components.StatSpeed]
	}

	return speed / 10.0
}
