package effects

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/util"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type StatsChange struct {
	components.ModifierData
	Update func(stats *components.StatsData)
}

func (e StatsChange) Active(world donburi.World, modifier *donburi.Entry) bool {
	return true
}
func (e StatsChange) Apply(world donburi.World, frame *util.Frame, modifier *donburi.Entry) {
	components.EachDependent(world, modifier, func(entry *donburi.Entry) {
		if entry.HasComponent(components.Stats) {
			frame.Modify(entry, components.Stats, func(stats *components.StatsData) {
				e.Update(stats)
			})
		}
		if entry.HasComponent(components.Image) {
			frame.Modify(entry, components.Image, func(image **ebiten.Image) {
				*image = assets.BlueSquareImage
			})
		}
	})
}

var AttackUp StatsChange = StatsChange{
	EffectID: uuid.MustParse("01a03c4e-06de-7469-9bb8-efc31688ee16"),
	Priority: 0,
	Update: func(stats *components.StatsData) {
		stats.Base[components.StatMelee] = stats.Base[components.StatMelee] * 2
	},
}
