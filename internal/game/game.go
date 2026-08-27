package game

import (
	"fmt"
	"image/color"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/data/effects"
	"rtwp_ebitengine/internal/util"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
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

	components.DecrementDelays(g.World)
	components.DecrementDurations(g.World)
	components.RemoveCompleted(g.World)
	components.RemoveCompletedDelays(g.World)

	components.ResolveModifiers(g.World, g.Frame, effects.EffectRegistry)
	components.MoveEntities(g.World)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var text strings.Builder
	g.drawRanges(g.World, screen)
	effect_entry, ok := components.Duration.First(g.World)
	if ok {
		if effect_entry.HasComponent(components.Duration) {
			duration := components.Duration.Get(effect_entry)
			text.WriteString(fmt.Sprintf("Duration = %d \n", *duration))
		}
	}

	for selected := range components.Selected.Iter(g.World) {
		if selected.HasComponent(components.Stats) {
			stats := components.Stats.Get(selected)
			text.WriteString(fmt.Sprintf("Speed = %f, Entity = %d \n", stats.Base[components.StatSpeed], selected.Id()))
			continue
		}
		text.WriteString(fmt.Sprintf("Selected = %s \n", selected.Entity()))
	}

	components.RenderEntries(screen, g.World)
	g.drawDragRect(screen)

	ebitenutil.DebugPrint(screen, text.String())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 680, 480
}

func cursorPoint() dmath.Vec2 {
	x, y := ebiten.CursorPosition()
	return dmath.Vec2{
		X: float64(x),
		Y: float64(y),
	}
}

func (g *Game) drawDragRect(screen *ebiten.Image) {
	rect, ok := g.DragRect()
	if !ok {
		return
	}

	borderColor := color.RGBA{0, 0xff, 0, 0xff}
	vx := float32(rect.Min.X)
	vy := float32(rect.Min.Y)
	vheight := float32(rect.Dy())
	vwidth := float32(rect.Dx())
	vector.FillRect(screen, vx, vy, vwidth, 1, borderColor, false)
	vector.FillRect(screen, vx, vy+vheight-1, vwidth, 1, borderColor, false)
	vector.FillRect(screen, vx, vy, 1, vheight, borderColor, false)
	vector.FillRect(screen, vx+vwidth-1, vy, 1, vheight, borderColor, false)
}

func (g *Game) drawRanges(world donburi.World, screen *ebiten.Image) {
	for entry := range components.RangeQuery.Iter(world) {
		t := transform.GetTransform(entry)
		r := components.Range.Get(entry)
		vector.StrokeCircle(
			screen,
			float32(t.LocalPosition.X),
			float32(t.LocalPosition.Y),
			float32(*r),
			2,
			color.RGBA{0xff, 0xff, 0, 0xff},
			false,
		)
	}
}
