package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	StopDistance float64
}

var Movement = donburi.NewComponentType[MovementData]()
var MovementQuery = donburi.NewQuery(
	filter.Contains(Movement, transform.Transform),
)

func WithMovement(entry *donburi.Entry, target dmath.Vec2, stopDistance float64) {
	WithMovementList(entry, []dmath.Vec2{target}, stopDistance)
}

func PushMovement(entry *donburi.Entry, target dmath.Vec2, stopDistance float64) {
	if !entry.HasComponent(Movement) {
		WithMovement(entry, target, stopDistance)
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

		WithMovement(selected, point, stopDistance)
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

func movementPosition(entry *donburi.Entry) dmath.Vec2 {
	if center, ok := CollisionCenter(entry); ok {
		return center
	}

	return transform.Transform.Get(entry).LocalPosition
}

func nextTarget(movement *MovementData) bool {
	movement.Targets = movement.Targets[1:]
	return len(movement.Targets) == 0
}

func collisionStopDistance(entry *donburi.Entry, stopDistance float64) float64 {
	if !entry.HasComponent(Collision) {
		return stopDistance
	}

	collision := Collision.Get(entry)
	if collision.Size.X <= 0 || collision.Size.Y <= 0 {
		return stopDistance
	}

	return max(stopDistance, collision.Size.Magnitude())
}

func MoveEntities(world donburi.World) {
	completed := []donburi.Entity{}

	for entry := range MovementQuery.Iter(world) {
		movement := Movement.Get(entry)
		if len(movement.Targets) == 0 {
			completed = append(completed, entry.Entity())
			continue
		}

		target := movement.Targets[0]
		position := movementPosition(entry)
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
			remainingDistance = movementPosition(entry).Distance(target)
			if moveResult.Collided && remainingDistance <= collisionStopDistance(entry, stopDistance) {
				if nextTarget(movement) {
					completed = append(completed, entry.Entity())
				}
				continue
			}
		}

		if remainingDistance <= stopDistance {
			if nextTarget(movement) {
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
		if len(movement.Targets) == 0 {
			continue
		}

		from := movementPosition(entry)
		for _, to := range movement.Targets {
			vector.StrokeLine(
				screen,
				float32(from.X),
				float32(from.Y),
				float32(to.X),
				float32(to.Y),
				1,
				lineColor,
				false,
			)

			from = to
		}
	}
}
