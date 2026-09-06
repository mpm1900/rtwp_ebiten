package components

import (
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/filter"
)

type ModifierConfig struct {
	Priority int
}

type ModifierBehavior interface {
	Active(*ModifierData, donburi.World, *donburi.Entry) bool
	Apply(*ModifierData, donburi.World, *util.Frame, *donburi.Entry)
	Spawn(*ModifierData, *ecs.ECS, math.Vec2) donburi.Entity
}

type ModifierData struct {
	ModifierConfig
	Behavior ModifierBehavior
}

func (mod *ModifierData) Active(world donburi.World, modifier *donburi.Entry) bool {
	if mod == nil || mod.Behavior == nil {
		return false
	}

	return mod.Behavior.Active(mod, world, modifier)
}

func (mod *ModifierData) Apply(world donburi.World, frame *util.Frame, modifier *donburi.Entry) {
	if mod == nil || mod.Behavior == nil {
		return
	}

	mod.Behavior.Apply(mod, world, frame, modifier)
}

func (mod *ModifierData) Spawn(ecs *ecs.ECS, position math.Vec2) donburi.Entity {
	if mod == nil || mod.Behavior == nil {
		return donburi.Null
	}

	return mod.Behavior.Spawn(mod, ecs, position)
}

func NewModifierInstance(mod *ModifierData) ModifierData {
	if mod == nil {
		return ModifierData{}
	}

	return *mod
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
