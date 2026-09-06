package events

import (
	"image"
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
)

type SelectAtEvent struct {
	Point math.Vec2
	Shift bool
}
type SelectInRectEvent struct {
	Rect  image.Rectangle
	Shift bool
}

var ClearSelected = events.NewEventType[struct{}]()
var SelectAt = events.NewEventType[SelectAtEvent]()
var SelectInRect = events.NewEventType[SelectInRectEvent]()

func InitSelection(world donburi.World) {
	ClearSelected.Subscribe(world, clearSelected)
	SelectAt.Subscribe(world, selectAt)
	SelectInRect.Subscribe(world, selectInRect)
}

func clearSelected(world donburi.World, _ struct{}) {
	cursor := components.GetCursor(world)
	player := components.GetPlayer(world)
	cursor.ClearDrag()
	player.SelectedAction = nil

	for selected := range components.Selected.Iter(world) {
		selected.RemoveComponent(components.Selected)
	}
}

func selectActor(world donburi.World, entity donburi.Entity) {
	if !world.Valid(entity) {
		return
	}

	entry := world.Entry(entity)
	actor := components.Actor.Get(entry)
	player := components.GetPlayerEntity(world)
	if actor.Player != player {
		return
	}

	entry.AddComponent(components.Selected)
}

func selectAt(world donburi.World, event SelectAtEvent) {
	clicked_entity := singleSelectedAtPoint(world, event.Point)
	if clicked_entity != donburi.Null {
		centerCameraOn(world, clicked_entity)
	}

	if !event.Shift {
		clearSelected(world, struct{}{})
	}

	player := components.GetPlayer(world)
	player.SelectedAction = nil

	components.EachActorAtPoint(world, event.Point, func(entry *donburi.Entry) {
		selectActor(world, entry.Entity())
		actor := components.Actor.Get(entry)
		player.SelectedAction = actor.Actions[0]
	})
}

func singleSelectedAtPoint(world donburi.World, point math.Vec2) donburi.Entity {
	var selected_entity donburi.Entity
	components.EachActorAtPoint(world, point, func(entry *donburi.Entry) {
		if entry.HasComponent(components.Selected) {
			selected_entity = entry.Entity()
		}
	})

	count := components.SelectedActorsQuery.Count(world)
	if count == 1 && selected_entity != donburi.Null {
		return selected_entity
	}

	return donburi.Null
}

func centerCameraOn(world donburi.World, entity donburi.Entity) {
	if !world.Valid(entity) {
		return
	}

	camera := components.GetCamera(world)
	camera.CenterOn(components.Center(world.Entry(entity)))
}

func selectInRect(world donburi.World, event SelectInRectEvent) {
	if !event.Shift {
		clearSelected(world, struct{}{})
	}

	player := components.GetPlayer(world)
	player.SelectedAction = nil

	for entry := range components.ActorQuery.Iter(world) {
		actorRect, ok := components.Rect(entry)
		if !ok {
			continue
		}

		if event.Rect.Overlaps(actorRect) {
			selectActor(world, entry.Entity())
			actor := components.Actor.Get(entry)
			player.SelectedAction = actor.Actions[0]
		}
	}
}
