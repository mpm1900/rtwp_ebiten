package components

import (
	"fmt"
	"maps"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type Stat int

const (
	StatHealth  Stat = iota
	StatMelee   Stat = iota
	StatDefense Stat = iota
	StatSpeed   Stat = iota

	StatAccuracy Stat = iota
	StatEvasion  Stat = iota
)

func (s Stat) String() string {
	switch s {
	case StatHealth:
		return "Health"
	case StatMelee:
		return "Melee"
	case StatDefense:
		return "Defense"
	case StatSpeed:
		return "Speed"
	default:
		return fmt.Sprintf("Stat(%d)", s)
	}
}

type StatsData struct {
	Base   map[Stat]float64
	Stages map[Stat]int
	Stats  map[Stat]float64
}

func NewStatsData(base map[Stat]float64) *StatsData {
	stages := map[Stat]int{}
	for stat := range base {
		stages[stat] = 0
	}

	return &StatsData{
		Base:   base,
		Stages: stages,
		Stats:  maps.Clone(base),
	}
}

func getStageMult(stage int, factor float64) float64 {
	n := factor
	d := factor
	if stage > 0 {
		n += float64(stage)
	}
	if stage < 0 {
		d -= float64(stage)
	}

	return n / d
}

func resolveBaseStat(value float64, level float64) float64 {
	return (((2 * value) * level) / 100) + level + 10
}

func (s *StatsData) ResolveBaseStats() {
	for stat, value := range s.Stats {
		if stat == StatAccuracy || stat == StatEvasion {
			continue
		}

		s.Stats[stat] = resolveBaseStat(value, 25)
	}
}

func (s *StatsData) MapStages() {
	for stat, value := range s.Stats {
		if stat == StatAccuracy || stat == StatEvasion {
			value *= getStageMult(s.Stages[stat], 3)
		} else {
			value *= getStageMult(s.Stages[stat], 2)
		}

		s.Stats[stat] = value
	}
}

func (stats StatsData) Clone() StatsData {
	stats.Base = maps.Clone(stats.Base)
	stats.Stages = maps.Clone(stats.Stages)
	stats.Stats = maps.Clone(stats.Stats)
	return stats
}

var Stats = donburi.NewComponentType[StatsData]()
var StatsQuery = donburi.NewQuery(filter.Contains(Stats))
