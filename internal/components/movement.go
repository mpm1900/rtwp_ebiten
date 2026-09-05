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
	Targets         []dmath.Vec2
	TargetLoopIndex int
	Loop            bool
	Follow          donburi.Entity
	StopDistance    float64
}

func (m *MovementData) NextTarget() bool {
	if m.Follow != donburi.Null {
		return false
	}

	if m.Loop && len(m.Targets) > 0 {
		m.TargetLoopIndex = (m.TargetLoopIndex + 1) % len(m.Targets)
		return false
	}
	m.Targets = m.Targets[1:]
	return len(m.Targets) == 0
}

var Movement = donburi.NewComponentType[MovementData]()
var MovementQuery = donburi.NewQuery(filter.And(
	filter.Contains(Movement, transform.Transform),
	filter.Not(filter.Contains(Delay)),
))

func WithMovementTo(entry *donburi.Entry, target dmath.Vec2, stopDistance float64) {
	WithMovementList(entry, []dmath.Vec2{target}, stopDistance, false)
}
func WithMovementLoopTo(entry *donburi.Entry, target dmath.Vec2, stopDistance float64) {
	WithMovementList(entry, []dmath.Vec2{target}, stopDistance, true)
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
	PushMovementList(entry, []dmath.Vec2{target}, stopDistance, false)
}
func PushMovementLoop(entry *donburi.Entry, target dmath.Vec2, stopDistance float64) {
	PushMovementList(entry, []dmath.Vec2{target}, stopDistance, true)
}

func PushMovementList(entry *donburi.Entry, targets []dmath.Vec2, stopDistance float64, loop bool) {
	if len(targets) == 0 {
		return
	}
	if !entry.HasComponent(Movement) {
		WithMovementList(entry, targets, stopDistance, loop)
		return
	}

	movement := Movement.Get(entry)
	movement.Follow = donburi.Null
	movement.Loop = movement.Loop || loop
	movement.StopDistance = stopDistance
	movement.Targets = append(movement.Targets, targets...)
}

func WithMovementList(entry *donburi.Entry, targets []dmath.Vec2, stopDistance float64, loop bool) {
	if !entry.HasComponent(Movement) {
		entry.AddComponent(Movement)
	}

	current_center := Center(entry)
	Movement.SetValue(entry, MovementData{
		Follow:          donburi.Null,
		Targets:         targets,
		StopDistance:    stopDistance,
		Loop:            loop,
		TargetLoopIndex: initialTargetLoopIndex(targets, loop, current_center, stopDistance),
	})
}

func initialTargetLoopIndex(targets []dmath.Vec2, loop bool, current dmath.Vec2, stopDistance float64) int {
	if !loop || len(targets) == 0 {
		return 0
	}

	if stopDistance <= 0 {
		stopDistance = DEFAULT_STOP_DISTANCE
	}
	if targets[0].Distance(current) <= stopDistance && len(targets) > 1 {
		return 1
	}

	return 0
}

func LoopOriginForEntry(entry *donburi.Entry) dmath.Vec2 {
	if entry.HasComponent(Movement) {
		movement := Movement.Get(entry)
		if movement.Loop && len(movement.Targets) > 0 {
			return movement.Targets[len(movement.Targets)-1]
		}
	}

	return Center(entry)
}

func AppendLoopOrigin(targets []dmath.Vec2, origin dmath.Vec2) []dmath.Vec2 {
	if len(targets) == 0 {
		return []dmath.Vec2{origin}
	}

	last_target := targets[len(targets)-1]
	if last_target.Distance(origin) <= DEFAULT_STOP_DISTANCE {
		return targets
	}

	return append(targets, origin)
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

	target := movement.Targets[0]
	if movement.Loop {
		if movement.TargetLoopIndex >= len(movement.Targets) {
			movement.TargetLoopIndex = 0
		}
		target = movement.Targets[movement.TargetLoopIndex]
	}

	return ClampWorldPosition(target), true
}
