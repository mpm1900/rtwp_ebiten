package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type PlayerData struct {
	DragStart *math.Vec2
	DragEnd   *math.Vec2
}

var Player = donburi.NewComponentType[PlayerData]()

func WithSelection(entry *donburi.Entry, start math.Vec2) {
	entry.AddComponent(Player)
	Player.SetValue(entry, PlayerData{
		DragStart: &start,
		DragEnd:   nil,
	})
}

func UpdateDrag(entry *donburi.Entry, end math.Vec2) {
	selection := Player.Get(entry)
	selection.DragEnd = &end
}

func GetPlayer(world donburi.World) *PlayerData {
	entry, ok := Player.First(world)
	if !ok {
		return nil
	}

	return Player.Get(entry)
}
