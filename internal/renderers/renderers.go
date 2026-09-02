package renderers

import "github.com/yohamta/donburi/ecs"

func Load(ecs *ecs.ECS) {
	ecs.AddRenderer(RenderLayerBackground, RenderBackground)
	ecs.AddRenderer(RenderLayerActors, RenderActors)
	ecs.AddRenderer(RenderLayerEffects, RenderMovement)
	ecs.AddRenderer(RenderLayerEffects, RenderTargets)
	ecs.AddRenderer(RenderLayerEffects, RenderRanges)
	ecs.AddRenderer(RenderLayerEffects, RenderEffect)
	ecs.AddRenderer(RenderLayerSelection, RenderDragRect)
	ecs.AddRenderer(RenderLayerActors, RenderHealthbars)
	ecs.AddRenderer(RenderLayerUI, RenderMinimap)
	ecs.AddRenderer(RenderLayerUI, RenderActionText)
	ecs.AddRenderer(RenderLayerUI, RenderPausedText)
	ecs.AddRenderer(RenderLayerCursor, RenderCursor)
}
