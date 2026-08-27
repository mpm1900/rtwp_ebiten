package ecs

import "github.com/yohamta/donburi"

var Targets = donburi.NewComponentType[[]donburi.Entity]()

func WithTargets(entry *donburi.Entry, targets ...donburi.Entity) {
	entry.AddComponent(Targets)
	Targets.SetValue(entry, targets)
}
