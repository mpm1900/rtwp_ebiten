package ecs

import "github.com/yohamta/donburi"

var Position = donburi.NewComponentType[Point]()

func WithPosition(entry *donburi.Entry, position Point) {
	entry.AddComponent(Position)
	Position.SetValue(entry, position)
}
