package renderers

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi/ecs"
)

func RenderActionText(ecs *ecs.ECS, screen *ebiten.Image) {
	player := components.GetPlayer(ecs.World)
	if player.Action == nil {
		return
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(0, float64(components.SCREEN_HEIGHT)-24)
	text.Draw(screen, player.Action.Data().Name, &text.GoTextFace{
		Source: assets.YolkFontSource,
		Size:   24,
	}, op)
}
