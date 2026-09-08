package events

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

var Destroy = events.NewEventType[donburi.Entity]()

func InitDestroy(world donburi.World) {
	Destroy.Subscribe(world, handleDestroy)
}

func handleDestroy(world donburi.World, entity donburi.Entity) {
	world.Remove(entity)
}
