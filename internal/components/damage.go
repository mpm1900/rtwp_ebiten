package components

import (
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
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
	if !entry.HasComponent(Damage) {
		entry.AddComponent(Damage)
	}
	Damage.SetValue(entry, damage)
}

func DamageAt(world donburi.World, point math.Vec2, amount float64) (*donburi.Entry, bool) {
	entry, ok := FirstActorAtPoint(world, util.ToPoint(point))
	if !ok {
		return nil, false
	}

	if !entry.HasComponent(Damage) {
		WithDamage(entry, amount)
		return nil, false
	}

	damage := Damage.Get(entry)
	Damage.SetValue(entry, *damage+amount)
	return entry, true
}
