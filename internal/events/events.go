package events

import "github.com/yohamta/donburi"

func Load(world donburi.World) {
	InitInput(world)
	InitDrag(world)
	InitSelection(world)
	InitCamera(world)
	InitMinimap(world)
	InitDamage(world)
	InitActions(world)
	InitDeath(world)
}

func ProcessEvents(world donburi.World) {
	// pre events
	ActionClick.ProcessEvents(world)

	// normal events
	StartDrag.ProcessEvents(world)
	UpdateDrag.ProcessEvents(world)
	EndDrag.ProcessEvents(world)
	ClearActions.ProcessEvents(world)
	ClearSelected.ProcessEvents(world)
	SelectInRect.ProcessEvents(world)
	SelectAt.ProcessEvents(world)
	UpdateCamera.ProcessEvents(world)
	ZoomCamera.ProcessEvents(world)
	LeftClickMinimap.ProcessEvents(world)
	DamageAt.ProcessEvents(world)
	Actions.ProcessEvents(world)

	// post events
	ActorDeath.ProcessEvents(world)
}
