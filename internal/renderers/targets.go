package renderers

import (
	"image/color"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var renderTargetsQuery = donburi.NewQuery(
	filter.Contains(components.Targets, transform.Transform, components.Selected),
)

var targetsLineColor = color.RGBA{0xff, 0, 0, 0xff}
var targetsPath vector.Path
var targetsStrokeOptions = vector.StrokeOptions{
	Width: 1,
}
var targetsDrawOptions = newTargetsDrawOptions()

func newTargetsDrawOptions() vector.DrawPathOptions {
	options := vector.DrawPathOptions{
		AntiAlias: true,
	}
	options.ColorScale.ScaleWithColor(targetsLineColor)
	return options
}

func RenderTargets(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)
	screen_bounds := screen.Bounds()
	targetsPath.Reset()
	has_segments := false

	for entry := range renderTargetsQuery.Iter(ecs.World) {
		from := components.Center(entry)
		components.EachTarget(ecs.World, entry, func(target_entry *donburi.Entry) {
			to := components.Center(target_entry)
			has_segments = addMovementSegment(&targetsPath, screen_bounds, view, from, to) || has_segments
		})
	}

	if !has_segments {
		return
	}

	vector.StrokePath(screen, &targetsPath, &targetsStrokeOptions, &targetsDrawOptions)
}
