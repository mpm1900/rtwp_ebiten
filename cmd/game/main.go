package main

import (
	"log"
	"rtwp_ebitengine/internal/abilities"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/effects"
	"rtwp_ebitengine/internal/entities"
	"rtwp_ebitengine/internal/game"
	"rtwp_ebitengine/internal/renderers"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi/features/math"
)

var actions = []*components.Ability{
	&abilities.Move,
}

func main() {
	g := game.NewGame()
	entities.CreateActor(g.ECS, math.NewVec2(100, 100), actions)
	entities.CreateActor(g.ECS, math.NewVec2(50, 200), actions)
	speed_up := effects.SpeedUp.Spawn(g.ECS, math.NewVec2(200, 200))
	effects.SpeedDown.Spawn(g.ECS, math.NewVec2(300, 250))

	components.WithDelay(g.ECS.World.Entry(speed_up), 60)
	components.WithRange(g.ECS.World.Entry(speed_up), 120)

	ebiten.SetWindowSize(renderers.SCREEN_WIDTH, renderers.SCREEN_HEIGHT)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
