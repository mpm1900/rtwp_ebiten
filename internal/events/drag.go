package events

import (
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
)

const dragClickThreshold = 40

var StartDrag = events.NewEventType[math.Vec2]()
var UpdateDrag = events.NewEventType[math.Vec2]()
var EndDrag = events.NewEventType[math.Vec2]()

func InitDrag(world donburi.World) {
	StartDrag.Subscribe(world, startDrag)
	UpdateDrag.Subscribe(world, updateDrag)
	EndDrag.Subscribe(world, endDrag)
}

func startDrag(world donburi.World, point math.Vec2) {
	player := components.GetPlayer(world)
	player.StartDrag(point)
}

func updateDrag(world donburi.World, point math.Vec2) {
	player := components.GetPlayer(world)
	player.UpdateDrag(point)
}

func endDrag(world donburi.World, point math.Vec2) {
	player := components.GetPlayer(world)
	if player.DragStart == nil {
		return
	}

	if player.DragDistance() <= dragClickThreshold {
		selectAt(world, *player.DragStart)
	} else {
		selectInRect(world, player.DragRect())
	}

	player.ClearDrag()
}
