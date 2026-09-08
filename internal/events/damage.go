package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
)

type DamageEvent struct {
	Position math.Vec2
	Source   donburi.Entity
	Action   *components.Action
}

var DamageAt = events.NewEventType[DamageEvent]()

func InitDamage(world donburi.World) {
	DamageAt.Subscribe(world, damageAt)
}

func damageAt(world donburi.World, event DamageEvent) {
	entry, ok := components.DamageAt(world, event.Source, event.Position, event.Action.Accuracy, event.Action.Power)
	if ok {
		damage := components.Damage.Get(entry)
		stats := components.Stats.Get(entry)
		if *damage >= stats.Stats[components.StatHealth] {
			Destroy.Publish(world, entry.Entity())
		}
	}
}
