package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var Duration = donburi.NewComponentType[int]()
var DurationQuery = donburi.NewQuery(filter.And(
	filter.Contains(Duration),
	filter.Not(filter.Contains(Delay)),
))

func WithDuration(entry *donburi.Entry, duration int) {
	entry.AddComponent(Duration)
	Duration.SetValue(entry, duration)
}
