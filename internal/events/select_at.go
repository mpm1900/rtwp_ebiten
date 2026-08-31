package events

import (
	"rtwp_ebitengine/internal/abilities"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
)

var SelectAt = events.NewEventType[math.Vec2]()

func InitSelectAt(world donburi.World) {
	SelectAt.Subscribe(world, selectAt)
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
	player.Ability = nil

	components.EachActorAtPoint(world, util.ToPoint(at), func(entry *donburi.Entry) {
		selectActor(world, entry.Entity())
		player.Ability = &abilities.Move
	})
}
