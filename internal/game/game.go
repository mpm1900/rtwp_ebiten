package game

import (
	"fmt"
	"image/color"
	"rtwp_ebitengine/internal/data/effects"
	"rtwp_ebitengine/internal/ecs"
	"rtwp_ebitengine/internal/util"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

type Game struct {
	Frame *util.Frame
	World donburi.World
	State State
}

func (g *Game) Update() error {
	mousePoint := cursorPoint()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.BeginDrag(mousePoint)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.UpdateDrag(mousePoint)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.EndDrag(mousePoint)
	}

	g.Frame.Restore(g.World)
	ecs.DecrementDelays(g.World)
	ecs.DecrementDurations(g.World)
	ecs.RemoveCompleted(g.World)
	ecs.RemoveCompletedDelays(g.World)
	ecs.ResolveModifiers(g.World, g.Frame, effects.EffectRegistry)
	ecs.MoveEntities(g.World)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var text strings.Builder
	effect_entry, ok := ecs.Duration.First(g.World)
	if ok {
		if effect_entry.HasComponent(ecs.Duration) {
			duration := ecs.Duration.Get(effect_entry)
			text.WriteString(fmt.Sprintf("Duration = %d \n", *duration))
		}
	}

	for selected := range ecs.Selected.Iter(g.World) {
		if selected.HasComponent(ecs.Stats) {
			stats := ecs.Stats.Get(selected)
			text.WriteString(fmt.Sprintf("Attack Power = %f, Entity = %d \n", stats.Base[ecs.StatMelee], selected.Id()))
			continue
		}
		text.WriteString(fmt.Sprintf("Selected = %s \n", selected.Entity()))
	}

	ecs.RenderEntries(screen, g.World)
	g.drawDragRect(screen)
	ebitenutil.DebugPrint(screen, text.String())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 680, 480
}

func cursorPoint() ecs.Point {
	x, y := ebiten.CursorPosition()
	return ecs.Point{
		X: float64(x),
		Y: float64(y),
	}
}

func (g *Game) drawDragRect(screen *ebiten.Image) {
	x, y, width, height, ok := g.DragRect()
	if !ok {
		return
	}

	borderColor := color.RGBA{0, 0xff, 0, 0xff}
	ebitenutil.DrawRect(screen, x, y, width, 1, borderColor)
	ebitenutil.DrawRect(screen, x, y+height-1, width, 1, borderColor)
	ebitenutil.DrawRect(screen, x, y, 1, height, borderColor)
	ebitenutil.DrawRect(screen, x+width-1, y, 1, height, borderColor)
}
