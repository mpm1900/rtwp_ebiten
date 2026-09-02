package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
)

type DamageEvent struct {
	Point  math.Vec2
	Amount float64
}

var DamageAt = events.NewEventType[DamageEvent]()

func InitDamage(world donburi.World) {
	DamageAt.Subscribe(world, damageAt)
}

func damageAt(world donburi.World, event DamageEvent) {
	entry, ok := components.DamageAt(world, event.Point, event.Amount)
	if ok {
		damage := components.Damage.Get(entry)
		stats := components.Stats.Get(entry)
		if *damage >= stats.Stats[components.StatHealth] {
			ActorDeath.Publish(world, entry.Entity())
		}
	}
}
