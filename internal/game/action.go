package game

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	ActionMove string = "Move"
)

const ActionTextSize = 24

type ActionState struct {
	Name string
}

func NewAction() ActionState {
	return ActionState{
		Name: ActionMove,
	}
}

func (g *Game) HandleActionInput() {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		_, has_selected := components.Selected.First(g.World)
		if has_selected {
			g.Action.Name = ActionMove
			g.ClearDrag()
		}
	}
}

func (g *Game) HandleActions() {
	mousePoint := cursorPoint()
	switch g.Action.Name {
	case ActionMove:
		{
			if _, has_selection := components.Selected.First(g.World); has_selection {
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
					components.MoveSelectedTo(g.World, mousePoint, components.DEFAULT_STOP_DISTANCE)
				}
			}
		}
	}
}

func (g *Game) RenderActionText(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(0, float64(SCREEN_HEIGHT)-ActionTextSize)
	text.Draw(screen, g.Action.Name, &text.GoTextFace{
		Source: assets.YolkFontSource,
		Size:   ActionTextSize,
	}, op)
}
