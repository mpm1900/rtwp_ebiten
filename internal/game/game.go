package game

import (
	"rtwp_ebitengine/internal/events"
	"rtwp_ebitengine/internal/renderers"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/ecs"
)

type Game struct {
	Frame *util.Frame
	ECS   *ecs.ECS
}

func (g *Game) Update() error {
	g.Frame.Restore(g.ECS.World)
	g.ECS.Update()
	events.ProcessEvents(g.ECS.World)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ECS.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return renderers.SCREEN_WIDTH, renderers.SCREEN_HEIGHT
}
