package abilities

import (
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

var Attack = components.Ability{
	Key:  ebiten.Key2,
	Name: "Attack",
	Handle: func(ecs *ecs.ECS, screenPoint math.Vec2) {
		if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			return
		}

		if components.IsOverMinimap(screenPoint, components.MinimapRect()) {
			return
		}

		player := components.GetPlayer(ecs.World)
		if player == nil {
			return
		}
		worldPoint := player.ScreenToWorld(screenPoint)

		for _ = range components.Selected.Iter(ecs.World) {
			components.DamageAt(ecs.World, worldPoint, 10)
		}
	},
	Valid: func(ecs *ecs.ECS, screenPoint math.Vec2) bool {
		if components.IsOverMinimap(screenPoint, components.MinimapRect()) {
			return false
		}
		player := components.GetPlayer(ecs.World)
		if player == nil {
			return false
		}
		worldPoint := player.ScreenToWorld(screenPoint)
		return components.IsInWorld(worldPoint)
	},
}
