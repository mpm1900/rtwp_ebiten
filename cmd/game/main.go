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

	ebiten.SetWindowSize(680, 480)

	game := game.Game{
		Frame: frame,
		World: world,
		State: game.NewState(),
	}

	one := components.MakeActor(world, math.Vec2{X: 100, Y: 100})
	components.MakeActor(world, math.Vec2{X: 50, Y: 200})

	attack_up := components.CreateEffect(world, effects.AttackUp)
	components.WithDuration(attack_up, 60)
	components.WithDelay(attack_up, 60)
	components.WithEffect(one, attack_up)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
