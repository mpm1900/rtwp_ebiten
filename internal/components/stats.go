package components

import (
	"maps"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type Stat int

const (
	StatHealth Stat = iota
	StatMelee  Stat = iota
	StatSpeed  Stat = iota
)

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

func (s *StatsData) MapStages() {
	stats := maps.Clone(s.Base)
	for stat := range stats {
		value := stats[stat]
		stage := s.Stages[stat]
		value *= getStageMult(stage, 2)

		stats[stat] = value
	}
	s.Stats = stats
}

func (stats StatsData) Clone() StatsData {
	stats.Base = maps.Clone(stats.Base)
	stats.Stages = maps.Clone(stats.Stages)
	stats.Stats = maps.Clone(stats.Stats)
	return stats
}

var Stats = donburi.NewComponentType[StatsData]()
var StatsQuery = donburi.NewQuery(filter.Contains(Stats))
