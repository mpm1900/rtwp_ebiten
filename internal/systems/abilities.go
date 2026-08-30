package systems

import (
	"rtwp_ebitengine/internal/components"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi/ecs"
)

func HandleAbilities(registry map[uuid.UUID]components.Ability) ecs.System {
	return func(ecs *ecs.ECS) {
		_, has_selected := components.Selected.First(ecs.World)
		if !has_selected {
			return
		}

		player := components.GetPlayer(ecs.World)
		for _, ability := range registry {
			if inpututil.IsKeyJustPressed(ability.Key) {
				player.Ability.AbilityID = ability.AbilityID
			}

			if player.Ability.AbilityID == ability.AbilityID {
				ability.Handle(ecs)
			}
		}
	}
}
