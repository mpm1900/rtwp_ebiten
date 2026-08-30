package main

import (
	"log"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/effects"
	"rtwp_ebitengine/internal/entities"
	"rtwp_ebitengine/internal/game"
	"rtwp_ebitengine/internal/renderers"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

func main() {
	g := game.NewGame()
	entities.CreateActor(g.ECS, math.NewVec2(100, 100))
	entities.CreateActor(g.ECS, math.NewVec2(50, 200))
	speed_up := effects.SpeedUp.Spawn(g.ECS, math.NewVec2(200, 200))
	effects.SpeedDown.Spawn(g.ECS, math.NewVec2(300, 250))

	components.WithDelay(g.ECS.World.Entry(speed_up), 60)
	components.WithRange(g.ECS.World.Entry(speed_up), 120)

	g.ECS.World.Entry(speed_up).RemoveComponent(components.Range)
	components.WithTargetsWhere(g.ECS.World.Entry(speed_up), func(e donburi.Entity) bool {
		return true
	})

	ebiten.SetWindowSize(renderers.SCREEN_WIDTH, renderers.SCREEN_HEIGHT)
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
