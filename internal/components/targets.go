package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
)

var Targets = donburi.NewComponentType[[]donburi.Entity]()
var TargetsWhere = donburi.NewComponentType[func(donburi.Entity) bool]()

func WithTargets(entry *donburi.Entry, targets ...donburi.Entity) {
	entry.AddComponent(Targets)
	Targets.SetValue(entry, targets)
}

func WithTargetsWhere(entry *donburi.Entry, where func(donburi.Entity) bool) {
	entry.AddComponent(TargetsWhere)
	TargetsWhere.SetValue(entry, where)
}

func EachTarget(world donburi.World, entry *donburi.Entry, yield func(*donburi.Entry)) {
	if !entry.HasComponent(Targets) {
		return
	}

	targets := Targets.Get(entry)
	for _, target := range *targets {
		if !world.Valid(target) {
			continue
		}

		target_entry := world.Entry(target)
		if !target_entry.HasComponent(transform.Transform) {
			continue
		}

		yield(target_entry)
	}
}

func FirstTarget(world donburi.World, entry *donburi.Entry) (*donburi.Entry, bool) {
	var first *donburi.Entry
	EachTarget(world, entry, func(target *donburi.Entry) {
		if first == nil {
			first = target
		}
	})

	return first, first != nil
}
