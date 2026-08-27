package effects

import (
	"rtwp_ebitengine/internal/ecs"
	"rtwp_ebitengine/internal/util"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type StatsChange struct {
	ecs.ModifierData
	Update func(stats *ecs.StatsData)
}

func (e StatsChange) Active(world donburi.World, modifier *donburi.Entry) bool {
	return true
}
func (e StatsChange) Apply(world donburi.World, frame *util.Frame, modifier *donburi.Entry) {
	ecs.EachDependent(world, modifier, func(entry *donburi.Entry) {
		if entry.HasComponent(ecs.Stats) {
			frame.Modify(entry, ecs.Stats, func(stats *ecs.StatsData) {
				e.Update(stats)
			})
		}
		if entry.HasComponent(ecs.Image) {
			frame.Modify(entry, ecs.Image, func(image **ebiten.Image) {
				*image = ecs.BlueSquareImage
			})
		}
	})
}

var AttackUp StatsChange = StatsChange{
	EffectID: uuid.MustParse("01a03c4e-06de-7469-9bb8-efc31688ee16"),
	Priority: 0,
	Update: func(stats *ecs.StatsData) {
		stats.Base[ecs.StatMelee] = stats.Base[ecs.StatMelee] * 2
	},
}
