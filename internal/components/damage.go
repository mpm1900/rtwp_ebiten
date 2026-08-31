package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var Damage = donburi.NewComponentType[float64]()
var DamageQuery = donburi.NewQuery(filter.Contains(Damage, Stats))

func GetHealth(entry *donburi.Entry) (health float64, damage float64) {
	stats := Stats.Get(entry)
	damage = 0.0
	if entry.HasComponent(Damage) {
		damage = *Damage.Get(entry)
	}

	return stats.Stats[StatHealth], damage
}

func WithDamage(entry *donburi.Entry, damage float64) {
	entry.AddComponent(Damage)
	Damage.SetValue(entry, damage)
}
