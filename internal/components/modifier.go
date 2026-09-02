package components

import (
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/filter"
)

type Effect interface {
	Active(world donburi.World, modifier *donburi.Entry) bool
	Apply(world donburi.World, frame *util.Frame, modifier *donburi.Entry)
	Spawn(ecs *ecs.ECS, position math.Vec2) donburi.Entity
}

type ModifierConfig struct {
	Priority int
}

type ModifierData struct {
	ModifierConfig
	Effect Effect
}

func NewModifierInstance(effect Effect, mod ModifierConfig) ModifierData {
	return ModifierData{
		ModifierConfig: mod,
		Effect:         effect,
	}
}

func (mod ModifierData) Order() int {
	return mod.Priority
}

var Modifier = donburi.NewComponentType[ModifierData]()
var ModifierQuery = donburi.NewOrderedQuery[ModifierData](
	filter.And(
		filter.Contains(Modifier),
		filter.Not(filter.Contains(Delay)),
	),
)

func EachDependent(world donburi.World, modifier *donburi.Entry, yield func(*donburi.Entry)) {
	if modifier.HasComponent(Targets) {
		EachTarget(world, modifier, yield)
	}

	if modifier.HasComponent(TargetsWhere) {
		where := *TargetsWhere.Get(modifier)
		for actor := range Actor.Iter(world) {
			if where(actor.Entity()) {
				yield(actor)
			}
		}
	}

	if modifier.HasComponent(Range) {
		EachActorsInRange(world, modifier, yield)
	}
}
