package components

import (
	"github.com/google/uuid"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type ModifierData struct {
	EffectID uuid.UUID
	Priority int
}

func (mod ModifierData) Order() int {
	return mod.Priority
}

func (mod ModifierData) Modifier() ModifierData {
	return mod
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
		entities := Targets.Get(modifier)
		for _, entity := range *entities {
			yield(world.Entry(entity))
		}
	}

	if modifier.HasComponent(Range) {
		EachActorsInRange(world, modifier, yield)
	}
}
