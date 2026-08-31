package pathing

import (
	"container/heap"
	"math"
	"rtwp_ebitengine/internal/components"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

const DefaultCellSize = 48.0

type cell struct {
	x, y   int
	g, h   float64
	parent *cell
}

type cellPriorityQueue []*cell

func (pq cellPriorityQueue) Len() int { return len(pq) }

func (pq cellPriorityQueue) Less(i, j int) bool {
	return (pq[i].g + pq[i].h) < (pq[j].g + pq[j].h)
}

func (pq cellPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *cellPriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*cell))
}

func (pq *cellPriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

func FindPath(world donburi.World, start, goal dmath.Vec2) ([]dmath.Vec2, bool) {
	start = components.ClampWorldPosition(start)
	goal = components.ClampWorldPosition(goal)
	if start.Distance(goal) <= 1.0 {
		return []dmath.Vec2{goal}, true
	}

	cellSize := DefaultCellSize
	boundsMinX := math.Min(start.X, goal.X) - cellSize
	boundsMinY := math.Min(start.Y, goal.Y) - cellSize
	boundsMaxX := math.Max(start.X, goal.X) + cellSize
	boundsMaxY := math.Max(start.Y, goal.Y) + cellSize

	for other := range components.CollisionQuery.Iter(world) {
		if !other.HasComponent(transform.Transform) {
			continue
		}

		bounds, ok := components.Rect(other)
		if !ok {
			continue
		}

		boundsMinX = math.Min(boundsMinX, float64(bounds.Min.X)-cellSize)
		boundsMinY = math.Min(boundsMinY, float64(bounds.Min.Y)-cellSize)
		boundsMaxX = math.Max(boundsMaxX, float64(bounds.Max.X)+cellSize)
		boundsMaxY = math.Max(boundsMaxY, float64(bounds.Max.Y)+cellSize)
	}

	width := int(math.Ceil((boundsMaxX-boundsMinX)/cellSize)) + 1
	height := int(math.Ceil((boundsMaxY-boundsMinY)/cellSize)) + 1
	if width <= 0 || height <= 0 || width*height > 200000 {
		return nil, false
	}

	blocked := make([][]bool, width)
	cells := make([][]*cell, width)
	for x := range width {
		blocked[x] = make([]bool, height)
		cells[x] = make([]*cell, height)
		for y := range height {
			cells[x][y] = &cell{x: x, y: y}
		}
	}

	for other := range components.CollisionQuery.Iter(world) {
		if !other.HasComponent(transform.Transform) {
			continue
		}

		bounds, ok := components.Rect(other)
		if !ok {
			continue
		}

		minCellX := int(math.Floor((float64(bounds.Min.X) - boundsMinX) / cellSize))
		maxCellX := int(math.Floor((float64(bounds.Max.X) - boundsMinX) / cellSize))
		minCellY := int(math.Floor((float64(bounds.Min.Y) - boundsMinY) / cellSize))
		maxCellY := int(math.Floor((float64(bounds.Max.Y) - boundsMinY) / cellSize))

		for x := minCellX; x <= maxCellX; x++ {
			if x < 0 || x >= width {
				continue
			}
			for y := minCellY; y <= maxCellY; y++ {
				if y < 0 || y >= height {
					continue
				}
				blocked[x][y] = true
			}
		}
	}

	startCellX := int(math.Floor((start.X - boundsMinX) / cellSize))
	startCellY := int(math.Floor((start.Y - boundsMinY) / cellSize))
	goalCellX := int(math.Floor((goal.X - boundsMinX) / cellSize))
	goalCellY := int(math.Floor((goal.Y - boundsMinY) / cellSize))

	if startCellX < 0 || startCellX >= width || startCellY < 0 || startCellY >= height || goalCellX < 0 || goalCellX >= width || goalCellY < 0 || goalCellY >= height {
		return nil, false
	}

	blocked[startCellX][startCellY] = false
	blocked[goalCellX][goalCellY] = false

	startCell := cells[startCellX][startCellY]
	goalCell := cells[goalCellX][goalCellY]
	startCell.g = 0
	startCell.h = heuristic(startCell, goalCell)

	open := &cellPriorityQueue{startCell}
	closed := make(map[[2]int]bool)

	for open.Len() > 0 {
		current := heap.Pop(open).(*cell)
		key := [2]int{current.x, current.y}
		if closed[key] {
			continue
		}
		closed[key] = true

		if current.x == goalCellX && current.y == goalCellY {
			return reconstructPath(current, goalCellX, goalCellY, boundsMinX, boundsMinY, cellSize), true
		}

		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}

				nx := current.x + dx
				ny := current.y + dy
				if nx < 0 || nx >= width || ny < 0 || ny >= height {
					continue
				}
				if blocked[nx][ny] {
					continue
				}

				neighbor := cells[nx][ny]
				neighborKey := [2]int{nx, ny}
				if closed[neighborKey] {
					continue
				}

				moveCost := 1.0
				if dx != 0 && dy != 0 {
					moveCost = 1.41421356237
				}

				tentative := current.g + moveCost
				if tentative < neighbor.g || (neighbor.parent == nil && !(neighbor.x == startCellX && neighbor.y == startCellY)) {
					neighbor.g = tentative
					neighbor.h = heuristic(neighbor, goalCell)
					neighbor.parent = current
					heap.Push(open, neighbor)
				}
			}
		}
	}

	return nil, false
}

func heuristic(a, b *cell) float64 {
	dx := math.Abs(float64(a.x - b.x))
	dy := math.Abs(float64(a.y - b.y))
	return dx + dy
}

func reconstructPath(end *cell, goalX, goalY int, minX, minY, cellSize float64) []dmath.Vec2 {
	path := []dmath.Vec2{}
	for current := end; current != nil; current = current.parent {
		point := dmath.NewVec2(
			minX+(float64(current.x)+0.5)*cellSize,
			minY+(float64(current.y)+0.5)*cellSize,
		)
		path = append(path, point)
	}

	for left, right := 0, len(path)-1; left < right; left, right = left+1, right-1 {
		path[left], path[right] = path[right], path[left]
	}

	if len(path) == 0 {
		path = []dmath.Vec2{dmath.NewVec2(
			minX+(float64(goalX)+0.5)*cellSize,
			minY+(float64(goalY)+0.5)*cellSize,
		)}
	}

	return path
}
