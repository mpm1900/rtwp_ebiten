package components

import (
	"image/color"
	"math"

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

func WithMovement(entry *donburi.Entry, target dmath.Vec2) {
	WithMovementList(entry, []dmath.Vec2{target}, DEFAULT_STOP_DISTANCE)
}

func PushMovement(entry *donburi.Entry, target dmath.Vec2) {
	if !entry.HasComponent(Movement) {
		WithMovement(entry, target)
		return
	}

	movement := Movement.Get(entry)
	movement.Targets = append(movement.Targets, target)
	WithMovementList(entry, append(movement.Targets, target), movement.StopDistance)
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

func GetSpeed(entry *donburi.Entry) float64 {
	speed := 0.0

	if entry.HasComponent(Stats) {
		stats := Stats.Get(entry)
		speed = stats.Base[StatSpeed]
	}

	return speed
}

func advanceMovementTarget(movement *MovementData) bool {
	movement.Targets = movement.Targets[1:]
	return len(movement.Targets) == 0
}

func MoveEntities(world donburi.World) {
	completed := []donburi.Entity{}

	for entry := range MovementQuery.Iter(world) {
		trans := transform.Transform.Get(entry)
		movement := Movement.Get(entry)
		if len(movement.Targets) == 0 {
			completed = append(completed, entry.Entity())
			continue
		}

		target := movement.Targets[0]
		dx := target.X - trans.LocalPosition.X
		dy := target.Y - trans.LocalPosition.Y
		distance := math.Hypot(dx, dy)
		stopDistance := movement.StopDistance
		if stopDistance <= 0 {
			stopDistance = DEFAULT_STOP_DISTANCE
		}

		remainingDistance := distance
		if distance > stopDistance {
			step := math.Min(GetSpeed(entry), distance)
			trans.LocalPosition.X += dx / distance * step
			trans.LocalPosition.Y += dy / distance * step
			remainingDistance -= step
		}

		if remainingDistance <= stopDistance {
			if advanceMovementTarget(movement) {
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

		from := transform.GetTransform(entry).LocalPosition
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
