package main

import (
	"log"
	"rtwp_ebitengine/internal/abilities"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/effects"
	"rtwp_ebitengine/internal/entities"
	"rtwp_ebitengine/internal/game"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi/features/math"
)

var actions = []*components.Ability{
	&abilities.Move,
}

func main() {
	g := game.NewGame()
	player := components.GetPlayerEntity(g.ECS.World)
	entities.CreateActor(g.ECS, math.NewVec2(100, 100), actions, player)
	entities.CreateActor(g.ECS, math.NewVec2(130, 130), actions, player)
	entities.CreateActor(g.ECS, math.NewVec2(130, 190), actions, player)
	entities.CreateActor(g.ECS, math.NewVec2(160, 100), actions, player)
	entities.CreateActor(g.ECS, math.NewVec2(160, 160), actions, player)
	entities.CreateActor(g.ECS, math.NewVec2(160, 210), actions, player)
	entities.CreateActor(g.ECS, math.NewVec2(190, 130), actions, player)
	entities.CreateActor(g.ECS, math.NewVec2(190, 210), actions, player)

	speed_up := effects.SpeedUp.Spawn(g.ECS, math.NewVec2(200, 200))
	effects.SpeedDown.Spawn(g.ECS, math.NewVec2(300, 250))

	components.WithDelay(g.ECS.World.Entry(speed_up), 60)
	components.WithRange(g.ECS.World.Entry(speed_up), 120)

	ebiten.SetWindowSize(components.SCREEN_WIDTH, components.SCREEN_HEIGHT)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
