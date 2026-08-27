package components

import (
	"maps"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

const (
	StatMelee = iota
	StatSpeed
)

type StatsValue map[int]float64

type StatsData struct {
	Base StatsValue
}

func NewStatsData(base StatsValue) *StatsData {
	return &StatsData{
		Base: base,
	}
}

func (stats StatsData) Clone() StatsData {
	stats.Base = maps.Clone(stats.Base)
	return stats
}

var Stats = donburi.NewComponentType[StatsData]()
var StatsQuery = donburi.NewQuery(filter.Contains(Stats))
