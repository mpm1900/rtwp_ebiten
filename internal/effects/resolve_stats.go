package effects

import (
	"rtwp_ebitengine/internal/components"
	"rtwp_ebitengine/internal/entities"
	"rtwp_ebitengine/internal/util"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

type ResolveStats struct {
	components.ModifierData
}

func (e ResolveStats) Active(world donburi.World, modifier *donburi.Entry) bool {
	return true
}
func (e ResolveStats) Apply(world donburi.World, frame *util.Frame, modifier *donburi.Entry) {
	components.EachDependent(world, modifier, func(entry *donburi.Entry) {
		if entry.HasComponent(components.Stats) {
			frame.Modify(entry, components.Stats, func(stats *components.StatsData) {
				stats.MapStages()
			})
		}
	})
}
func (e ResolveStats) Spawn(ecs *ecs.ECS, position math.Vec2) donburi.Entity {
	entity := entities.CreateEffect(ecs, e, e.ModifierConfig)
	entry := ecs.World.Entry(entity)
	entry.AddComponent(components.TargetsWhere)
	components.TargetsWhere.SetValue(entry, func(e donburi.Entity) bool {
		return true
	})

	return entity
}

var SystemResolveStats = ResolveStats{
	Priority: 1,
}
