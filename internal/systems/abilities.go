package systems

import (
	"rtwp_ebitengine/internal/components"
	"uuid"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi/ecs"
)

type Ability struct {
	AbilityID uuid.UUID
	Key       ebiten.Key
	Name      string
	Handle    func(*ecs.ECS)
}

func HandleAbilities(registry map[uuid.UUID]Ability) ecs.System {
	return func(ecs *ecs.ECS) {
		_, has_selected := components.Selected.First(ecs.World)
		if !has_selected {
			return
		}

		player := components.GetPlayer(ecs.World)
		for _, ability := range registry {
			if inpututil.IsKeyJustPressed(ability.Key) {
				player.ActionName = ability.Name
			}

			if player.ActionName == ability.Name {
				ability.Handle(ecs)
			}
		}
	}
}
