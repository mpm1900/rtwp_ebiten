package main

import (
	"log"
	"rtwp_ebitengine/internal/abilities"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/effects"
	"rtwp_ebitengine/internal/entities"
	"rtwp_ebitengine/internal/game"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi/features/math"
)

var actions = []*components.Ability{
	&abilities.Move,
}

var amount = 10

func main() {
	g := game.NewGame()
	player := components.GetPlayerEntity(g.ECS.World)
	for i := range amount {
		for j := range amount {
			entities.CreateActor(g.ECS, util.NewVec2(200+i*40, 200+j*40), actions, player)

		}
	}

	speed_up := effects.SpeedUp.Spawn(g.ECS, math.NewVec2(800, 800))
	effects.SpeedDown.Spawn(g.ECS, math.NewVec2(300, 250))

	components.WithDelay(g.ECS.World.Entry(speed_up), 60)
	components.WithRange(g.ECS.World.Entry(speed_up), 120)
	components.WithCollision(g.ECS.World.Entry(speed_up))

	ebiten.SetWindowSize(components.SCREEN_WIDTH, components.SCREEN_HEIGHT)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
