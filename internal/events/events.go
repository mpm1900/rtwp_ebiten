package events

import "github.com/yohamta/donburi"

func Load(world donburi.World) {
	InitDrag(world)
	InitClearSelected(world)
	InitSelectAt(world)
	InitSelectInRect(world)
	InitCamera(world)
	InitMinimap(world)
}

func ProcessEvents(world donburi.World) {
	StartDrag.ProcessEvents(world)
	UpdateDrag.ProcessEvents(world)
	EndDrag.ProcessEvents(world)
	ClearSelected.ProcessEvents(world)
	SelectInRect.ProcessEvents(world)
	SelectAt.ProcessEvents(world)
	UpdateCamera.ProcessEvents(world)
	LeftClickMinimap.ProcessEvents(world)
	RightClickMinimap.ProcessEvents(world)
}
