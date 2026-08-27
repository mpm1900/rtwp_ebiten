package ecs

import (
	"math"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

const (
	DefaultMovementSpeed  = 2.0
	DEFAULT_STOP_DISTANCE = 1.0
)

type MovementData struct {
	Target       dmath.Vec2
	Speed        float64
	StopDistance float64
}

var Movement = donburi.NewComponentType[MovementData]()
var MovementQuery = donburi.NewQuery(
	filter.Contains(Movement, transform.Transform),
)

func WithMovement(entry *donburi.Entry, target dmath.Vec2) {
	WithMovementSpeed(entry, target, DefaultMovementSpeed, DEFAULT_STOP_DISTANCE)
}

func WithMovementSpeed(entry *donburi.Entry, target dmath.Vec2, speed float64, stopDistance float64) {
	entry.AddComponent(Movement)
	Movement.SetValue(entry, MovementData{
		Target:       target,
		Speed:        speed,
		StopDistance: stopDistance,
	})
}

func MoveEntities(world donburi.World) {
	completed := []donburi.Entity{}

	for entry := range MovementQuery.Iter(world) {
		trans := transform.Transform.Get(entry)
		movement := Movement.Get(entry)

		dx := movement.Target.X - trans.LocalPosition.X
		dy := movement.Target.Y - trans.LocalPosition.Y
		distance := math.Hypot(dx, dy)
		stopDistance := movement.StopDistance
		if stopDistance <= 0 {
			stopDistance = DEFAULT_STOP_DISTANCE
		}

		if distance <= stopDistance {
			completed = append(completed, entry.Entity())
			continue
		}

		speed := movement.Speed
		if speed <= 0 {
			speed = DefaultMovementSpeed
		}

		step := math.Min(speed, distance)
		trans.LocalPosition.X += dx / distance * step
		trans.LocalPosition.Y += dy / distance * step

		if distance-step <= stopDistance {
			completed = append(completed, entry.Entity())
		}
	}

	for _, entity := range completed {
		if !world.Valid(entity) {
			continue
		}

		entry := world.Entry(entity)
		if entry.HasComponent(Movement) {
			entry.RemoveComponent(Movement)
		}
	}
}
