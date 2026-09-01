package actions

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/events"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type AttackAction struct {
	components.ActionData
}

func (a AttackAction) Data() components.ActionData {
	return a.ActionData
}
func (a AttackAction) Handle(world donburi.World, point math.Vec2) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		return
	}

	if components.IsOverMinimap(point, components.MinimapRect()) {
		return
	}

	player := components.GetPlayer(world)
	if player == nil {
		return
	}
	worldPoint := player.ScreenToWorld(point)

	for _ = range components.Selected.Iter(world) {
		events.DamageAt.Publish(world, events.DamageEvent{
			Point:  worldPoint,
			Amount: 10,
		})
	}
}

func (a AttackAction) Valid(world donburi.World, point math.Vec2) bool {
	if components.IsOverMinimap(point, components.MinimapRect()) {
		return false
	}
	player := components.GetPlayer(world)
	if player == nil {
		return false
	}
	worldPoint := player.ScreenToWorld(point)
	return components.IsInWorld(worldPoint)
}

var Attack = AttackAction{
	Key:  ebiten.Key2,
	Name: "Attack",
}
