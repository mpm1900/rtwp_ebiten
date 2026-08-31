package systems

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi/ecs"
)

func HandleAbilities(ecs *ecs.ECS) {
	player := components.GetPlayer(ecs.World)
	if player == nil {
		return
	}

	actions := map[*components.Ability]struct{}{}

	for selected := range components.SelectedActorsQuery.Iter(ecs.World) {
		actor := components.Actor.Get(selected)
		for _, action := range actor.Abilities {
			if action != nil {
				actions[action] = struct{}{}
			}
		}
	}

	for ability := range actions {
		if inpututil.IsKeyJustPressed(ability.Key) {
			player.Ability = ability
		}
	}

	if player.Ability != nil {
		screenPoint := util.CursorPoint()
		if player.Ability.Valid != nil && !player.Ability.Valid(ecs, screenPoint) {
			return
		}

		if player.Ability.Handle != nil {
			player.Ability.Handle(ecs, screenPoint)
		}
	}
}
