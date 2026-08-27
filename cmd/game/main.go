package main

import (
	"image/color"
	"log"
	"rtwp_ebitengine/internal/data/effects"
	"rtwp_ebitengine/internal/ecs"
	"rtwp_ebitengine/internal/game"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi"
)

func main() {
	world := donburi.NewWorld()
	frame := util.NewFrame()

	ebiten.SetWindowSize(680, 480)
	ecs.RedSquareImage = ebiten.NewImage(24, 24)
	ecs.RedSquareImage.Fill(color.RGBA{0xff, 0, 0, 0xff})
	ecs.BlueSquareImage = ebiten.NewImage(24, 24)
	ecs.BlueSquareImage.Fill(color.RGBA{0, 0, 0xff, 0xff})
	ecs.GreenSquareImage = ebiten.NewImage(24, 24)
	ecs.GreenSquareImage.Fill(color.RGBA{0, 0xff, 0, 0xff})

	game := game.Game{
		Frame: frame,
		World: world,
		State: game.NewState(),
	}

	one := ecs.MakeActor(world, ecs.Point{X: 100, Y: 100})
	ecs.MakeActor(world, ecs.Point{X: 50, Y: 200})

	attack_up := ecs.CreateEffect(world, effects.AttackUp)
	ecs.WithDuration(attack_up, 60)
	ecs.WithDelay(attack_up, 60)
	ecs.WithEffect(one, attack_up)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
