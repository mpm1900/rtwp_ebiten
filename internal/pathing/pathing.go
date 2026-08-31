package pathing

import (
	"image"
	"math"
	"sync"

	"rtwp_ebitengine/internal/components"

	qpathing "github.com/quasilyte/pathing"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

const (
	DefaultCellSize  = 12.0
	MaxPointDistance = 48.0
)

var (
	gridMu sync.Mutex
	grid   *qpathing.Grid
	astar  *qpathing.AStar
	layer  = qpathing.MakeGridLayer([8]uint8{0: 1})
)

func init() {
	initGridAndAStar(DefaultCellSize)
}

func initGridAndAStar(cellSize float64) {
	_, _, worldMaxX, worldMaxY := components.WorldRect()
	maxX := worldMaxX + components.WORLD_BORDER
	maxY := worldMaxY + components.WORLD_BORDER

	cols := uint(math.Ceil(maxX / cellSize))
	rows := uint(math.Ceil(maxY / cellSize))

	grid = qpathing.NewGrid(qpathing.GridConfig{
		WorldWidth:  cols * uint(cellSize),
		WorldHeight: rows * uint(cellSize),
		CellWidth:   uint(cellSize),
		CellHeight:  uint(cellSize),
	})

	astar = qpathing.NewAStar(qpathing.AStarConfig{
		NumCols: uint(grid.NumCols()),
		NumRows: uint(grid.NumRows()),
	})
}

// FindPath finds a path between start and goal world coordinates around collision obstacles,
// with diagonal line-of-sight path smoothing and multi-point interpolation.
func FindPath(world donburi.World, start, goal dmath.Vec2) ([]dmath.Vec2, bool) {
	start = components.ClampWorldPosition(start)
	goal = components.ClampWorldPosition(goal)
	if start.Distance(goal) <= 1.0 {
		return []dmath.Vec2{goal}, true
	}

	gridMu.Lock()
	defer gridMu.Unlock()

	blocked := blockObstacles(world, start, goal)
	defer func() {
		for _, c := range blocked {
			grid.SetCellIsBlocked(c, false)
		}
	}()

	startCoord := grid.PosToCoord(start.X, start.Y)
	goalCoord := grid.PosToCoord(goal.X, goal.Y)

	if !inBounds(startCoord) || !inBounds(goalCoord) {
		return nil, false
	}
	if startCoord == goalCoord {
		return []dmath.Vec2{goal}, true
	}

	grid.SetCellIsBlocked(startCoord, false)
	grid.SetCellIsBlocked(goalCoord, false)

	if lineOfSight(startCoord, goalCoord) {
		return subdividePath([]dmath.Vec2{start, goal}, MaxPointDistance), true
	}

	coords := findAStarPath(startCoord, goalCoord)
	if len(coords) <= 1 {
		return nil, false
	}

	smoothed := smoothPathCoords(coords)
	rawPoints := make([]dmath.Vec2, len(smoothed))
	rawPoints[0] = start
	for i, c := range smoothed[1:] {
		x, y := grid.CoordToPos(c)
		rawPoints[i+1] = dmath.NewVec2(x, y)
	}
	rawPoints[len(rawPoints)-1] = goal

	return subdividePath(rawPoints, MaxPointDistance), true
}

func blockObstacles(world donburi.World, start, goal dmath.Vec2) []qpathing.GridCoord {
	var blocked []qpathing.GridCoord

	for other := range components.CollisionQuery.Iter(world) {
		if !other.HasComponent(transform.Transform) {
			continue
		}

		bounds, ok := components.Rect(other)
		if !ok || rectContains(bounds, start) || rectContains(bounds, goal) {
			continue
		}

		minCoord := grid.PosToCoord(float64(bounds.Min.X), float64(bounds.Min.Y))
		maxCoord := grid.PosToCoord(float64(bounds.Max.X-1), float64(bounds.Max.Y-1))

		for x := minCoord.X; x <= maxCoord.X; x++ {
			for y := minCoord.Y; y <= maxCoord.Y; y++ {
				c := qpathing.GridCoord{X: x, Y: y}
				if inBounds(c) && !grid.GetCellIsBlocked(c) {
					grid.SetCellIsBlocked(c, true)
					blocked = append(blocked, c)
				}
			}
		}
	}

	return blocked
}

func findAStarPath(from, to qpathing.GridCoord) []qpathing.GridCoord {
	curr := from
	coords := []qpathing.GridCoord{from}
	visited := map[qpathing.GridCoord]bool{from: true}

	for _ = range 10 {
		result := astar.BuildPath(grid, curr, to, layer)
		if result.Steps.Len() == 0 {
			break
		}

		added := 0
		for result.Steps.HasNext() {
			curr = curr.Move(result.Steps.Next())
			coords = append(coords, curr)
			added++
		}

		if added == 0 || !result.Partial || result.Finish == to || visited[result.Finish] {
			break
		}
		visited[result.Finish] = true
		curr = result.Finish
	}

	return coords
}

func smoothPathCoords(coords []qpathing.GridCoord) []qpathing.GridCoord {
	if len(coords) <= 2 {
		return coords
	}

	smoothed := []qpathing.GridCoord{coords[0]}
	for curr := 0; curr < len(coords)-1; {
		furthest := curr + 1
		for next := len(coords) - 1; next > curr+1; next-- {
			if lineOfSight(coords[curr], coords[next]) {
				furthest = next
				break
			}
		}
		smoothed = append(smoothed, coords[furthest])
		curr = furthest
	}

	return smoothed
}

func lineOfSight(from, to qpathing.GridCoord) bool {
	dx := max(from.X, to.X) - min(from.X, to.X)
	dy := max(from.Y, to.Y) - min(from.Y, to.Y)

	sx, sy := 1, 1
	if from.X > to.X {
		sx = -1
	}
	if from.Y > to.Y {
		sy = -1
	}

	x, y := from.X, from.Y
	err := dx - dy

	for {
		c := qpathing.GridCoord{X: x, Y: y}
		if !inBounds(c) || grid.GetCellIsBlocked(c) {
			return false
		}
		if x == to.X && y == to.Y {
			return true
		}

		e2 := 2 * err
		if e2 > -dy {
			if e2 < dx {
				if grid.GetCellIsBlocked(qpathing.GridCoord{X: x + sx, Y: y}) ||
					grid.GetCellIsBlocked(qpathing.GridCoord{X: x, Y: y + sy}) {
					return false
				}
			}
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

func subdividePath(points []dmath.Vec2, maxSpacing float64) []dmath.Vec2 {
	if len(points) <= 1 {
		return points
	}

	path := make([]dmath.Vec2, 0, len(points)*2)
	for i := 0; i < len(points)-1; i++ {
		p1, p2 := points[i], points[i+1]
		dist := p1.Distance(p2)
		if dist <= 0.001 {
			continue
		}

		steps := int(math.Ceil(dist / maxSpacing))
		for s := 1; s <= steps; s++ {
			t := float64(s) / float64(steps)
			path = append(path, p1.Add(p2.Sub(p1).MulScalar(t)))
		}
	}

	if len(path) == 0 {
		return []dmath.Vec2{points[len(points)-1]}
	}

	return path
}

func inBounds(c qpathing.GridCoord) bool {
	return c.X >= 0 && c.X < grid.NumCols() && c.Y >= 0 && c.Y < grid.NumRows()
}

func rectContains(r image.Rectangle, p dmath.Vec2) bool {
	return float64(r.Min.X) <= p.X && p.X <= float64(r.Max.X) &&
		float64(r.Min.Y) <= p.Y && p.Y <= float64(r.Max.Y)
}
