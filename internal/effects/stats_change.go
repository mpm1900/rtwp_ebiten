package effects

import (
	"rtwp_ebitengine/internal/assets"
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/entities"
	"rtwp_ebitengine/internal/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
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

func (e StatsChange) Spawn(ecs *ecs.ECS, position math.Vec2) donburi.Entity {
	entity := entities.CreateEffect(ecs, e, e.ModifierConfig)
	entry := ecs.World.Entry(entity)
	components.WithImage(entry, assets.YellowSquareImage, position)
	components.WithRange(entry, 100)
	return entity
}

var SpeedUp StatsChange = StatsChange{
	Priority: 0,
	Update: func(stats *components.StatsData) {
		stats.Base[components.StatSpeed] = stats.Base[components.StatSpeed] * 2
	},
}

var SpeedDown StatsChange = StatsChange{
	Priority: 0,
	Update: func(stats *components.StatsData) {
		stats.Base[components.StatSpeed] = stats.Base[components.StatSpeed] / 2
	},
}
