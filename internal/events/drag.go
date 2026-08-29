package events

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

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
	player.DragStart = &point
	player.DragEnd = nil
}

func updateDrag(world donburi.World, point math.Vec2) {
	player := components.GetPlayer(world)
	player.DragEnd = &point
}

func endDrag(world donburi.World, point math.Vec2) {
	player := components.GetPlayer(world)
	if player.DragStart == nil {
		return
	}

	distance := player.DragStart.Distance(point)
	if distance <= dragClickThreshold {
		selectAt(world, *player.DragStart)
	} else {
		rect := util.ToRect(*player.DragStart, point)
		selectInRect(world, rect)
	}
	player.DragStart = nil
	player.DragEnd = nil
}
