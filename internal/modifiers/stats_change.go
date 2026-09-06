package modifiers

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

type StatsChangeBehavior struct {
	Update func(stats *components.StatsData)
}

func (b StatsChangeBehavior) Active(mod *components.ModifierData, world donburi.World, modifier *donburi.Entry) bool {
	return true
}
func (b StatsChangeBehavior) Apply(mod *components.ModifierData, world donburi.World, frame *util.Frame, modifier *donburi.Entry) {
	components.EachDependent(world, modifier, func(entry *donburi.Entry) {
		if entry.HasComponent(components.Stats) {
			frame.Modify(entry, components.Stats, func(stats *components.StatsData) {
				b.Update(stats)
			})
		}
		if entry.HasComponent(components.Image) {
			frame.Modify(entry, components.Image, func(image **ebiten.Image) {
				*image = assets.BlueSquareImage
			})
		}
	})
}
func (b StatsChangeBehavior) Spawn(mod *components.ModifierData, ecs *ecs.ECS, position math.Vec2) donburi.Entity {
	entity := entities.CreateEffect(ecs, mod)
	entry := ecs.World.Entry(entity)
	components.WithImage(entry, assets.YellowSquareImage, position)
	components.WithRange(entry, 100)
	return entity
}

var SpeedUp = &components.ModifierData{
	Priority: 0,
	Behavior: StatsChangeBehavior{
		Update: func(stats *components.StatsData) {
			stats.Stages[components.StatSpeed] += 2
		},
	},
}

var SpeedDown = &components.ModifierData{
	Priority: 0,
	Behavior: StatsChangeBehavior{
		Update: func(stats *components.StatsData) {
			stats.Stages[components.StatSpeed] -= 2
		},
	},
}
