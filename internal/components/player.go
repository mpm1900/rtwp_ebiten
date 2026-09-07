package components

import (
	"github.com/yohamta/donburi"
)

type PlayerData struct {
	SelectedAction *Action
}

var Player = donburi.NewComponentType[PlayerData]()

func GetPlayerEntity(world donburi.World) donburi.Entity {
	entry := Player.MustFirst(world)
	return entry.Entity()
}
func GetPlayer(world donburi.World) *PlayerData {
	entry := Player.MustFirst(world)
	return Player.Get(entry)
}
