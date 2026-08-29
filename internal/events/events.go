package events

import "github.com/yohamta/donburi"

func Load(world donburi.World) {
	InitClearSelected(world)
	InitSelectAt(world)
	InitSelectInRect(world)
}

func ProcessEvents(world donburi.World) {
	ClearSelected.ProcessEvents(world)
	SelectInRect.ProcessEvents(world)
	SelectAt.ProcessEvents(world)
}
