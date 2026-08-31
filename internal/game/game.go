package game

import (
	"image/color"
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/effects"
	"rtwp_ebitengine/internal/entities"
	"rtwp_ebitengine/internal/events"
	"rtwp_ebitengine/internal/renderers"
	"rtwp_ebitengine/internal/systems"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type Game struct {
	Frame *util.Frame
	ECS   *ecs.ECS
}

func NewGame() Game {
	g := Game{
		Frame: util.NewFrame(),
		ECS:   ecs.NewECS(donburi.NewWorld()),
	}

	assets.MustLoadAssets()
	events.Load(g.ECS.World)
	systems.Load(g.ECS, g.Frame)
	renderers.Load(g.ECS)

	effects.LoadSystemModifiers(g.ECS)
	entities.CreatePlayer(g.ECS)

	return g
}

func (g *Game) Update() error {
	g.Frame.Restore(g.ECS.World)
	g.ECS.Update()
	events.ProcessEvents(g.ECS.World)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(assets.ColorBackground)
	player := components.GetPlayer(g.ECS.World)
	camera_surface := player.Camera.Surface
	camera_surface.Fill(color.Black)

	g.ECS.DrawLayer(renderers.RenderLayerBackground, camera_surface)
	g.ECS.DrawLayer(renderers.RenderLayerEffects, camera_surface)
	g.ECS.DrawLayer(renderers.RenderLayerActors, camera_surface)
	player.Camera.Blit(screen)

	g.ECS.DrawLayer(renderers.RenderLayerSelection, screen)
	g.ECS.DrawLayer(renderers.RenderLayerUI, screen)
	g.ECS.DrawLayer(renderers.RenderLayerCursor, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return components.SCREEN_WIDTH, components.SCREEN_HEIGHT
}
