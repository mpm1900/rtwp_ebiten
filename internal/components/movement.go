package components

import (
	"image/color"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

const (
	DEFAULT_STOP_DISTANCE = 1.0
)

type MovementData struct {
	Targets      []dmath.Vec2
	Follow       donburi.Entity
	StopDistance float64
}

var Movement = donburi.NewComponentType[MovementData]()
var MovementQuery = donburi.NewQuery(
	filter.Contains(Movement, transform.Transform),
)

func WithMovementTo(entry *donburi.Entry, target dmath.Vec2, stopDistance float64) {
	WithMovementList(entry, []dmath.Vec2{target}, stopDistance)
}

func WithMovementFollow(entry *donburi.Entry, follow donburi.Entity, stopDistance float64) {
	if !entry.HasComponent(Movement) {
		entry.AddComponent(Movement)
	}

	Movement.SetValue(entry, MovementData{
		Follow:       follow,
		Targets:      nil,
		StopDistance: stopDistance,
	})
}

func PushMovement(entry *donburi.Entry, target dmath.Vec2, stopDistance float64) {
	if !entry.HasComponent(Movement) {
		WithMovementTo(entry, target, stopDistance)
		return
	}

	movement := Movement.Get(entry)
	movement.Targets = append(movement.Targets, target)
	if movement.StopDistance <= 0 {
		movement.StopDistance = DEFAULT_STOP_DISTANCE
	}
}

func WithMovementList(entry *donburi.Entry, targets []dmath.Vec2, stopDistance float64) {
	if !entry.HasComponent(Movement) {
		entry.AddComponent(Movement)
	}

	Movement.SetValue(entry, MovementData{
		Follow:       donburi.Null,
		Targets:      targets,
		StopDistance: stopDistance,
	})
}

func MoveSelectedTo(world donburi.World, point dmath.Vec2, stopDistance float64) {
	push := ebiten.IsKeyPressed(ebiten.KeyShift)

	for selected := range Selected.Iter(world) {
		if push {
			PushMovement(selected, point, stopDistance)
			continue
		}

		WithMovementTo(selected, point, stopDistance)
	}
}

func MoveSelectedFollow(world donburi.World, follow donburi.Entity, stopDistance float64) {
	for selected := range Selected.Iter(world) {
		if selected.Entity() == follow {
			continue
		}

		WithMovementFollow(selected, follow, stopDistance)
	}
}

func GetSpeed(entry *donburi.Entry) float64 {
	speed := 0.0

	if entry.HasComponent(Stats) {
		stats := Stats.Get(entry)
		speed = stats.Base[StatSpeed]
	}

	return speed
}

func nextTarget(movement *MovementData) bool {
	movement.Targets = movement.Targets[1:]
	return len(movement.Targets) == 0
}

func movementTarget(world donburi.World, movement *MovementData) (dmath.Vec2, bool) {
	if movement.Follow != donburi.Null {
		if !world.Valid(movement.Follow) {
			return dmath.Vec2{}, false
		}

		follow := world.Entry(movement.Follow)
		if !follow.HasComponent(transform.Transform) {
			return dmath.Vec2{}, false
		}

		return Center(follow), true
	}

	if len(movement.Targets) == 0 {
		return dmath.Vec2{}, false
	}

	return movement.Targets[0], true
}

func finishMovementTarget(movement *MovementData) bool {
	if movement.Follow != donburi.Null {
		return false
	}

	return nextTarget(movement)
}

func collisionStopDistance(entry *donburi.Entry, stopDistance float64) float64 {
	if !entry.HasComponent(Collision) || !entry.HasComponent(transform.Transform) {
		return stopDistance
	}

	scale := transform.Transform.Get(entry).LocalScale
	if scale.X <= 0 || scale.Y <= 0 {
		return stopDistance
	}

	return max(stopDistance, scale.Magnitude())
}

func MoveEntities(world donburi.World) {
	completed := []donburi.Entity{}

	for entry := range MovementQuery.Iter(world) {
		movement := Movement.Get(entry)
		target, ok := movementTarget(world, movement)
		if !ok {
			completed = append(completed, entry.Entity())
			continue
		}

		position := Center(entry)
		direction := target.Sub(position)
		distance := direction.Magnitude()
		stopDistance := movement.StopDistance
		if stopDistance <= 0 {
			stopDistance = DEFAULT_STOP_DISTANCE
		}

		remainingDistance := distance
		if distance > stopDistance {
			step := min(GetSpeed(entry), distance)
			moveResult := MoveWithCollision(world, entry, direction.Normalized().MulScalar(step))
			remainingDistance = Center(entry).Distance(target)
			if moveResult.Collided && remainingDistance <= collisionStopDistance(entry, stopDistance) {
				if finishMovementTarget(movement) {
					completed = append(completed, entry.Entity())
				}
				continue
			}
		}

		if remainingDistance <= stopDistance {
			if finishMovementTarget(movement) {
				completed = append(completed, entry.Entity())
			}
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

func DrawMovement(screen *ebiten.Image, world donburi.World) {
	lineColor := color.RGBA{0xff, 0xff, 0xff, 0xff}

	for entry := range MovementQuery.Iter(world) {
		movement := Movement.Get(entry)
		from := Center(entry)
		if movement.Follow != donburi.Null {
			to, ok := movementTarget(world, movement)
			if !ok {
				continue
			}

			util.DrawPoints(screen, from, to, 1, lineColor)
			continue
		}

		for _, to := range movement.Targets {
			util.DrawPoints(screen, from, to, 1, lineColor)
			from = to
		}
	}
}
