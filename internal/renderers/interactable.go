package renderers

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi/ecs"
)

func RenderInteractableTargets(ecs *ecs.ECS, screen *ebiten.Image) {
	view := newCameraView(ecs)

	for entry := range components.InteractableQuery.Iter(ecs.World) {
		center := components.Center(entry)
		interactable := components.Interactable.Get(entry)
		target := view.Point(center.Add(interactable.TargetOffset))

		vector.FillCircle(screen, float32(target.X), float32(target.Y), 8, assets.ColorEffectTarget, true)
	}
}
