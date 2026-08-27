package components

import (
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
)

type Effect interface {
	Modifier() ModifierData
	Active(world donburi.World, modifier *donburi.Entry) bool
	Apply(world donburi.World, frame *util.Frame, modifier *donburi.Entry)
}

func CreateEffect(world donburi.World, effect Effect) *donburi.Entry {
	return CreateModifier(world, effect.Modifier())
}

func WithEffect(parent *donburi.Entry, effect *donburi.Entry) {
	WithTargets(effect, parent.Entity())
}
