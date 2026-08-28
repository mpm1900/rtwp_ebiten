package main

import (
	"log"
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/data/effects"
	"rtwp_ebitengine/internal/game"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

func main() {
	assets.MustLoadAssets()
	world := donburi.NewWorld()
	frame := util.NewFrame()

	ebiten.SetWindowSize(game.SCREEN_WIDTH, game.SCREEN_HEIGHT)

	game := game.Game{
		Frame:     frame,
		World:     world,
		Action:    game.NewAction(),
		Selection: game.NewSelection(),
	}

	components.CreateActor(world, math.NewVec2(100, 100))
	components.CreateActor(world, math.NewVec2(50, 200))

	attack_up := components.CreateEffect(world, effects.SpeedUp)
	attack_down := components.CreateEffect(world, effects.SpeedDown)
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
