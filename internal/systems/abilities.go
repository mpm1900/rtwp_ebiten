package systems

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi/ecs"
)

func HandleActions(ecs *ecs.ECS) {
	player := components.GetPlayer(ecs.World)
	if player == nil || player.Action == nil {
		return
	}

	actions := map[components.Action]int{}
	count := 0

	for selected := range components.SelectedActorsQuery.Iter(ecs.World) {
		count++
		actor := components.Actor.Get(selected)
		for _, action := range actor.Actions {
			actions[action]++
		}
	}

	for action := range actions {
		if inpututil.IsKeyJustPressed(action.Data().Key) && actions[action] == count {
			player.Action = action
		}
	}

	screenPoint := util.CursorPoint()
	if !player.Action.Valid(ecs.World, screenPoint) {
		return
	}

	player.Action.Handle(ecs.World, screenPoint)
}
