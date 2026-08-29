package game

import (
	"fmt"
	"image/color"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

const (
	SCREEN_HEIGHT int = 480
	SCREEN_WIDTH  int = 640
)

type Game struct {
	Frame     *util.Frame
	ECS       *ecs.ECS
	Selection SelectionState
	Action    ActionState
}

func (g *Game) Update() error {
	// pre resolve, mutate things that are modified
	g.HandleSelection()
	g.HandleActions()
	g.HandleActionInput()
	g.Frame.Restore(g.ECS.World)

	// resolve pipeline
	g.ECS.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var text strings.Builder
	g.drawRanges(g.ECS.World, screen)
	components.DrawMovement(screen, g.ECS.World)
	effect_entry, ok := components.Duration.First(g.ECS.World)
	if ok {
		if effect_entry.HasComponent(components.Duration) {
			duration := components.Duration.Get(effect_entry)
			text.WriteString(fmt.Sprintf("Duration = %d \n", *duration))
		}
	}

	for selected := range components.Selected.Iter(g.ECS.World) {
		if selected.HasComponent(components.Stats) {
			stats := components.Stats.Get(selected)
			text.WriteString(fmt.Sprintf("Speed = %f, Entity = %d \n", stats.Base[components.StatSpeed], selected.Id()))
			continue
		}
		text.WriteString(fmt.Sprintf("Selected = %s \n", selected.Entity()))
	}

	components.RenderEntries(screen, g.ECS.World)
	g.drawDragRect(screen)

	ebitenutil.DebugPrint(screen, text.String())

	g.RenderActionText(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
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
	vector.StrokeRect(screen, vx, vy, vwidth, vheight, 1, borderColor, false)
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
