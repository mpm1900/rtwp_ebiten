package events

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

var ActorDeath = events.NewEventType[donburi.Entity]()

func InitDeath(world donburi.World) {
	ActorDeath.Subscribe(world, handleActorDeath)
}

func handleActorDeath(world donburi.World, entity donburi.Entity) {
	world.Remove(entity)
}
