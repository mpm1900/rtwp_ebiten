package main

import (
	"fmt"
	"log"
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/effects"
	"rtwp_ebitengine/internal/entities"
	"rtwp_ebitengine/internal/game"
	"rtwp_ebitengine/internal/systems"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi"
	ecslib "github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

func System(ecs *ecslib.ECS) {
	fmt.Println("test")
}

var effectLayer ecslib.LayerID = 0
var actorLayer ecslib.LayerID = 1

func main() {
	assets.MustLoadAssets()
	world := donburi.NewWorld()
	ecs := ecslib.NewECS(world)
	//ecs.AddSystem(System)
	frame := util.NewFrame()

	ebiten.SetWindowSize(game.SCREEN_WIDTH, game.SCREEN_HEIGHT)

	game := game.Game{
		Frame:     frame,
		ECS:       ecs,
		Action:    game.NewAction(),
		Selection: game.NewSelection(),
	}

	systems.LoadSystems(ecs, game.Frame)

	entities.CreateActor(ecs, actorLayer, math.NewVec2(100, 100))
	entities.CreateActor(ecs, actorLayer, math.NewVec2(50, 200))

	attack_up := entities.CreateEffect(ecs, effects.SpeedUp, effectLayer)
	attack_down := entities.CreateEffect(ecs, effects.SpeedDown, effectLayer)
	// components.WithDuration(attack_up, 60)
	components.WithDelay(attack_up, 60)
	// components.WithTargets(attack_up, one.Entity())
	components.WithRange(attack_up, 120)
	components.WithRange(attack_down, 100)
	components.WithImage(attack_up, assets.YellowSquareImage, math.NewVec2(200, 200))
	components.WithImage(attack_down, assets.YellowSquareImage, math.NewVec2(300, 250))

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
