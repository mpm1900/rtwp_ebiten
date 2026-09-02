package events

import (
	"image"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
)

var ClearSelected = events.NewEventType[struct{}]()
var SelectAt = events.NewEventType[math.Vec2]()
var SelectInRect = events.NewEventType[image.Rectangle]()

func InitSelection(world donburi.World) {
	ClearSelected.Subscribe(world, clearSelected)
	SelectAt.Subscribe(world, selectAt)
	SelectInRect.Subscribe(world, selectInRect)
}

func clearSelected(world donburi.World, _ struct{}) {
	player := components.GetPlayer(world)
	player.ClearDrag()
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

func selectAt(world donburi.World, at math.Vec2) {
	clearSelected(world, struct{}{})
	player := components.GetPlayer(world)
	player.SelectedAction = nil

	components.EachActorAtPoint(world, util.ToPoint(at), func(entry *donburi.Entry) {
		selectActor(world, entry.Entity())
		actor := components.Actor.Get(entry)
		player.SelectedAction = actor.Actions[0]
	})
}

func selectInRect(world donburi.World, rect image.Rectangle) {
	clearSelected(world, struct{}{})
	player := components.GetPlayer(world)
	player.SelectedAction = nil

	for entry := range components.ActorQuery.Iter(world) {
		actorRect, ok := components.Rect(entry)
		if !ok {
			continue
		}

		if rect.Overlaps(actorRect) {
			selectActor(world, entry.Entity())
			actor := components.Actor.Get(entry)
			player.SelectedAction = actor.Actions[0]
		}
	}
}
