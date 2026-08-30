package renderers

import (
	"rtwp_ebitengine/internal/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi/ecs"
)

func RenderPausedText(ecs *ecs.ECS, screen *ebiten.Image) {
	if !ecs.IsPaused() {
		return
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(SCREEN_WIDTH)-120, float64(SCREEN_HEIGHT)-24)
	text.Draw(screen, "Paused", &text.GoTextFace{
		Source: assets.YolkFontSource,
		Size:   24,
	}, op)
}
