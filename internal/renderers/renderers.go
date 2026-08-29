package renderers

import "github.com/yohamta/donburi/ecs"

const (
	SCREEN_HEIGHT int = 480
	SCREEN_WIDTH  int = 640
)

func Load(ecs *ecs.ECS) {
	ecs.AddRenderer(RenderLayerActors, RenderActors)
	ecs.AddRenderer(RenderLayerEffects, RenderMovement)
	ecs.AddRenderer(RenderLayerEffects, RenderRanges)
	ecs.AddRenderer(RenderLayerEffects, RenderEffect)
	ecs.AddRenderer(RenderLayerSelection, RenderDragRect)
	ecs.AddRenderer(RenderLayerUI, RenderActionText)
}
