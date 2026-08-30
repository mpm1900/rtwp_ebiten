package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var Selected = donburi.NewTag("Selected")

var SelectedActorsQuery = donburi.NewQuery(filter.Contains(
	Actor, Selected,
))
