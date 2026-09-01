package events

import "github.com/yohamta/donburi"

func Load(world donburi.World) {
	InitInput(world)
	InitDrag(world)
	InitClearSelected(world)
	InitSelectAt(world)
	InitSelectInRect(world)
	InitCamera(world)
	InitMinimap(world)
	InitDamage(world)
	InitActions(world)
}

func ProcessEvents(world donburi.World) {
	// pre events
	ActionClick.ProcessEvents(world)

	// normal events
	StartDrag.ProcessEvents(world)
	UpdateDrag.ProcessEvents(world)
	EndDrag.ProcessEvents(world)
	ClearSelected.ProcessEvents(world)
	SelectInRect.ProcessEvents(world)
	SelectAt.ProcessEvents(world)
	UpdateCamera.ProcessEvents(world)
	ZoomCamera.ProcessEvents(world)
	LeftClickMinimap.ProcessEvents(world)
	DamageAt.ProcessEvents(world)
	Actions.ProcessEvents(world)
}
