package events

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
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
	if player == nil {
		return
	}

	player.StartDrag(point)
}

func updateDrag(world donburi.World, point math.Vec2) {
	player := components.GetPlayer(world)
	if player == nil {
		return
	}

	player.UpdateDrag(point)
}

func endDrag(world donburi.World, point math.Vec2) {
	player := components.GetPlayer(world)
	if player == nil || player.DragStart == nil {
		return
	}

	shift := ebiten.IsKeyPressed(ebiten.KeyShift)
	if player.DragStart.Distance(point) <= dragClickThreshold {
		selectAt(world, SelectAtEvent{
			Point: player.ScreenToWorld(*player.DragStart),
			Shift: shift,
		})
	} else {
		start := player.ScreenToWorld(*player.DragStart)
		end := player.ScreenToWorld(point)
		selectInRect(world, SelectInRectEvent{
			Rect:  util.ToRect(start, end),
			Shift: shift,
		})
	}

	player.ClearDrag()
}
