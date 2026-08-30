package components

import "github.com/yohamta/donburi"

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
