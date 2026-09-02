package renderers

import (
	"image"
	"image/color"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var renderMovementQuery = donburi.NewQuery(
	filter.Contains(components.Movement, transform.Transform, components.Selected),
)

var movementLineColor = color.RGBA{0xff, 0xff, 0xff, 0xff}
var movementPath vector.Path
var movementStrokeOptions = vector.StrokeOptions{
	Width: 1,
}
var movementDrawOptions = newMovementDrawOptions()

func newMovementDrawOptions() vector.DrawPathOptions {
	options := vector.DrawPathOptions{
		AntiAlias: true,
	}
	options.ColorScale.ScaleWithColor(movementLineColor)
	return options
}

func RenderMovement(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)
	screen_bounds := screen.Bounds()
	movementPath.Reset()
	has_segments := false

	for entry := range renderMovementQuery.Iter(ecs.World) {
		movement := components.Movement.Get(entry)
		from := components.Center(entry)
		if movement.Follow != donburi.Null {
			to, ok := components.MovementPosition(ecs.World, movement)
			if !ok {
				continue
			}

			has_segments = addMovementSegment(&movementPath, screen_bounds, view, from, to) || has_segments
			continue
		}

		for i, to := range movement.Targets {
			if i > 0 || !movement.Loop {
				has_segments = addMovementSegment(&movementPath, screen_bounds, view, from, to) || has_segments
			}
			from = to
		}

		if movement.Loop && len(movement.Targets) > 0 {
			has_segments = addMovementSegment(&movementPath, screen_bounds, view, from, movement.Targets[0]) || has_segments
		}
	}

	if !has_segments {
		return
	}

	vector.StrokePath(screen, &movementPath, &movementStrokeOptions, &movementDrawOptions)
}

func addMovementSegment(path *vector.Path, bounds image.Rectangle, view cameraView, from, to dmath.Vec2) bool {
	screen_from := view.Point(from)
	screen_to := view.Point(to)
	if !segmentIntersects(bounds, screen_from, screen_to) {
		return false
	}

	path.MoveTo(float32(screen_from.X), float32(screen_from.Y))
	path.LineTo(float32(screen_to.X), float32(screen_to.Y))
	return true
}

func segmentIntersects(bounds image.Rectangle, from, to dmath.Vec2) bool {
	min_x := min(from.X, to.X)
	max_x := max(from.X, to.X)
	min_y := min(from.Y, to.Y)
	max_y := max(from.Y, to.Y)

	return max_x >= float64(bounds.Min.X) &&
		min_x < float64(bounds.Max.X) &&
		max_y >= float64(bounds.Min.Y) &&
		min_y < float64(bounds.Max.Y)
}
