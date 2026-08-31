package systems

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi/ecs"
)

func HandleAbilities(ecs *ecs.ECS) {
	player := components.GetPlayer(ecs.World)
	actions := map[*components.Ability]struct{}{}

	for selected := range components.SelectedActorsQuery.Iter(ecs.World) {
		actor := components.Actor.Get(selected)
		for _, action := range actor.Abilities {
			actions[action] = struct{}{}
		}
	}

	for ability := range actions {
		if inpututil.IsKeyJustPressed(ability.Key) {
			player.Ability = ability
		}
	}

	if player.Ability != nil {
		screenPoint := util.CursorPoint()
		if !player.Ability.Valid(ecs, screenPoint) {
			return
		}

		player.Ability.Handle(ecs, screenPoint)
	}
}
