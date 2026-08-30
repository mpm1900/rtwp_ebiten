package systems

import (
	"rtwp_ebitengine/internal/components"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi/ecs"
)

func HandleAbilities(registry map[uuid.UUID]*components.Ability) ecs.System {
	return func(ecs *ecs.ECS) {
		player := components.GetPlayer(ecs.World)
		if player == nil {
			return
		}

		for _, ability := range registry {
			if inpututil.IsKeyJustPressed(ability.Key) {
				player.Ability = ability
			}
		}

		if player.Ability != nil {
			player.Ability.Handle(ecs)
		}
	}
}
