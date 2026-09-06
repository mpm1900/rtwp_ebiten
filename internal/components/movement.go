package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

const (
	DEFAULT_STOP_DISTANCE = 1.0
)

type MovementData struct {
	Follow       donburi.Entity
	Loop         bool
	Path         []math.Vec2
	PathIndex    int
	StopDistance float64
}

func (m *MovementData) PushPath(path ...math.Vec2) {
	m.Path = append(m.Path, path...)
	m.Follow = donburi.Null
}
func (m *MovementData) Next() bool {
	if m.Follow != donburi.Null {
		return false
	}

	if m.Loop && len(m.Path) > 0 {
		m.PathIndex = (m.PathIndex + 1) % len(m.Path)
		return false
	}
	m.Path = m.Path[1:]
	return len(m.Path) == 0
}
func (m *MovementData) Target(world donburi.World) (math.Vec2, bool) {
	if m.Follow != donburi.Null {
		if !world.Valid(m.Follow) {
			return math.Vec2{}, false
		}

		follow := world.Entry(m.Follow)
		if !follow.HasComponent(transform.Transform) {
			return math.Vec2{}, false
		}

		return Center(follow), true
	}

	if len(m.Path) == 0 {
		return math.Vec2{}, false
	}

	target := m.Path[0]
	if m.Loop {
		if m.PathIndex >= len(m.Path) {
			m.PathIndex = 0
		}
		target = m.Path[m.PathIndex]
	}

	return target, true
}
func (m *MovementData) TowardsTarget(world donburi.World, entry *donburi.Entry) (math.Vec2, bool) {
	target, ok := m.Target(world)
	if !ok {
		return math.Vec2{}, false
	}
	vec := target.Sub(Center(entry))
	return vec, true
}
func (m *MovementData) Delta(world donburi.World, entry *donburi.Entry, scalar float64) (math.Vec2, bool) {
	vec, ok := m.TowardsTarget(world, entry)
	if !ok {
		return math.Vec2{}, false
	}

	return vec.Normalized().MulScalar(scalar), true
}
func (m *MovementData) TargetDistance(world donburi.World, entry *donburi.Entry) float64 {
	target, ok := m.Target(world)
	if !ok {
		return 0
	}

	return Center(entry).Distance(target)
}

func NewPathMovement(parent *donburi.Entry, path []math.Vec2, loop bool) MovementData {
	return MovementData{
		Follow:       donburi.Null,
		Loop:         loop,
		Path:         path,
		PathIndex:    initialLoopIndex(path, loop, Center(parent), DEFAULT_STOP_DISTANCE),
		StopDistance: DEFAULT_STOP_DISTANCE,
	}
}
func NewPathFollow(follow donburi.Entity) MovementData {
	return MovementData{
		Follow:       follow,
		StopDistance: DEFAULT_STOP_DISTANCE,
	}
}

var Movement = donburi.NewComponentType[MovementData]()
var MovementQuery = donburi.NewQuery(filter.And(
	filter.Contains(Movement, transform.Transform),
))

func WithMovement(entry *donburi.Entry, data MovementData) {
	if !entry.HasComponent(Movement) {
		entry.AddComponent(Movement)
	}

	Movement.SetValue(entry, data)
}

func PushMovementList(entry *donburi.Entry, path []math.Vec2, stopDistance float64, loop bool) {
	if len(path) == 0 {
		return
	}
	if !entry.HasComponent(Movement) {
		WithMovement(entry, NewPathMovement(entry, path, loop))
		return
	}

	movement := Movement.Get(entry)
	movement.Follow = donburi.Null
	movement.Loop = movement.Loop || loop
	movement.StopDistance = stopDistance
	movement.Path = append(movement.Path, path...)
}

func initialLoopIndex(path []math.Vec2, loop bool, current math.Vec2, stopDistance float64) int {
	if !loop || len(path) == 0 {
		return 0
	}

	if stopDistance <= 0 {
		stopDistance = DEFAULT_STOP_DISTANCE
	}
	if path[0].Distance(current) <= stopDistance && len(path) > 1 {
		return 1
	}

	return 0
}

func LoopOriginForEntry(entry *donburi.Entry) math.Vec2 {
	if entry.HasComponent(Movement) {
		movement := Movement.Get(entry)
		if movement.Loop && len(movement.Path) > 0 {
			return movement.Path[len(movement.Path)-1]
		}
	}

	return Center(entry)
}

func AppendLoopOrigin(targets []math.Vec2, origin math.Vec2) []math.Vec2 {
	if len(targets) == 0 {
		return []math.Vec2{origin}
	}

	last_target := targets[len(targets)-1]
	if last_target.Distance(origin) <= DEFAULT_STOP_DISTANCE {
		return targets
	}

	return append(targets, origin)
}
