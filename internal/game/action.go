package game

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
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
					var found *donburi.Entry

					components.EachActorAtPoint(g.World, util.ToPoint(mousePoint), func(e *donburi.Entry) {
						found = e
					})

					if found != nil {
						components.MoveSelectedFollow(g.World, found.Entity(), components.DEFAULT_STOP_DISTANCE)
					} else {
						components.MoveSelectedTo(g.World, mousePoint, components.DEFAULT_STOP_DISTANCE)
					}
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
