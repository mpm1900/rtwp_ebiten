package components

import (
	"rtwp_ebitengine/internal/util"

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
	filter.And(filter.Contains(Modifier), filter.Not(filter.Contains(Delay))),
)

func ResolveModifiers(world donburi.World, frame *util.Frame, effectRegistry map[uuid.UUID]Effect) {
	for modifier := range ModifierQuery.IterOrdered(world, Modifier) {
		instance := Modifier.Get(modifier)
		effect, ok := effectRegistry[instance.EffectID]
		if ok {
			if effect.Active(world, modifier) {
				effect.Apply(world, frame, modifier)
			}
		}
	}
}

func CreateModifier(world donburi.World, mod ModifierData) *donburi.Entry {
	entity := world.Create(Modifier)
	entry := world.Entry(entity)
	Modifier.SetValue(entry, mod)
	return entry
}

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
