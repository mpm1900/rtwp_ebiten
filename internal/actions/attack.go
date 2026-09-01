package actions

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/events"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type AttackAction struct {
	components.ActionData
}

func (a AttackAction) Data() components.ActionData {
	return a.ActionData
}
func (a AttackAction) Publish(world donburi.World, point math.Vec2) {
	for selected := range components.Selected.Iter(world) {
		events.Actions.Publish(world, components.ActionEvent{
			Action: a,
			Source: selected.Entity(),
			Point:  point,
		})
	}
}
func (a AttackAction) Handle(world donburi.World, source donburi.Entity, point math.Vec2) {
	events.DamageAt.Publish(world, events.DamageEvent{
		Point:  point,
		Amount: 10,
	})
}

func (a AttackAction) Valid(world donburi.World, point math.Vec2) bool {
	return components.IsInWorld(point)
}

var Attack = AttackAction{
	Key:  ebiten.Key2,
	Name: "Attack",
}
