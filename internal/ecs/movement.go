package ecs

import (
	"math"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

const (
	DefaultMovementSpeed  = 2.0
	DEFAULT_STOP_DISTANCE = 1.0
)

type MovementData struct {
	Target       Point
	Speed        float64
	StopDistance float64
}

var Movement = donburi.NewComponentType[MovementData]()
var MovementQuery = donburi.NewQuery(
	filter.Contains(Movement, Position),
)

func WithMovement(entry *donburi.Entry, target Point) {
	WithMovementSpeed(entry, target, DefaultMovementSpeed, DEFAULT_STOP_DISTANCE)
}

func WithMovementSpeed(entry *donburi.Entry, target Point, speed float64, stopDistance float64) {
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
		position := Position.Get(entry)
		movement := Movement.Get(entry)

		dx := movement.Target.X - position.X
		dy := movement.Target.Y - position.Y
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
		position.X += dx / distance * step
		position.Y += dy / distance * step

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
