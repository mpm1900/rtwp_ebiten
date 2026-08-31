package components

import (
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

func (m *MovementData) NextTarget() bool {
	if m.Follow != donburi.Null {
		return false
	}

	m.Targets = m.Targets[1:]
	return len(m.Targets) == 0
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
	PushMovementList(entry, []dmath.Vec2{target}, stopDistance)
}

func PushMovementList(entry *donburi.Entry, targets []dmath.Vec2, stopDistance float64) {
	if len(targets) == 0 {
		return
	}
	if !entry.HasComponent(Movement) {
		WithMovementList(entry, targets, stopDistance)
		return
	}

	movement := Movement.Get(entry)
	movement.Follow = donburi.Null
	movement.StopDistance = stopDistance
	movement.Targets = append(movement.Targets, targets...)
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

func GetSpeed(entry *donburi.Entry) float64 {
	speed := 0.0

	if entry.HasComponent(Stats) {
		stats := Stats.Get(entry)
		speed = stats.Stats[StatSpeed]
	}

	return speed
}

func MovementPosition(world donburi.World, movement *MovementData) (dmath.Vec2, bool) {
	if movement.Follow != donburi.Null {
		if !world.Valid(movement.Follow) {
			return dmath.Vec2{}, false
		}

		follow := world.Entry(movement.Follow)
		if !follow.HasComponent(transform.Transform) {
			return dmath.Vec2{}, false
		}

		return ClampWorldPosition(Center(follow)), true
	}

	if len(movement.Targets) == 0 {
		return dmath.Vec2{}, false
	}

	return ClampWorldPosition(movement.Targets[0]), true
}
