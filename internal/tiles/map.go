package tiles

import (
	"os"
	"strings"
)

// Map is a 2D grid of tile ids. A space or '.' cell is empty (grass shows
// through). A multi-cell tile (W x H) is written as its single character at the
// top-left cell; the cells it covers to the right and below are consumed and
// treated as part of that tile.
type Map struct {
	Width  int
	Height int
	Cells  []byte // char at the tile's origin cell; EmptyCell elsewhere
	// OriginX/OriginY give, for each cell, the origin cell of the tile that
	// occupies it (or -1 if the cell is empty).
	OriginX []int
	OriginY []int
}

// EmptyCell is the character used for "no tile here".
const EmptyCell = ' '

// IsEmpty reports whether a cell character means "no tile".
func IsEmpty(c byte) bool {
	return c == EmptyCell || c == '.'
}

// LoadMap reads a plain-text map file. Each line is a row; each character is a
// cell. Rows of differing lengths are right-padded with spaces. A trailing
// newline is ignored. Multi-cell tiles consume the cells they cover.
func LoadMap(path string) (*Map, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	height := len(lines)
	m := &Map{Width: width, Height: height}
	if width == 0 || height == 0 {
		return m, nil
	}

	cells := make([]byte, width*height)
	for i := range cells {
		cells[i] = EmptyCell
	}
	for y, l := range lines {
		for x := 0; x < width && x < len(l); x++ {
			cells[y*width+x] = l[x]
		}
	}

	// Resolve multi-cell tiles: for each origin char, consume the w x h block
	// it covers and record origin coordinates for every covered cell.
	originX := make([]int, width*height)
	originY := make([]int, width*height)
	for i := range originX {
		originX[i] = -1
		originY[i] = -1
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := cells[y*width+x]
			if IsEmpty(c) {
				continue
			}
			tile, ok := Get(c)
			if !ok {
				// Unknown char: leave as-is, treat as 1x1 origin.
				originX[y*width+x] = x
				originY[y*width+x] = y
				continue
			}
			w, h := tile.W, tile.H
			if w < 1 {
				w = 1
			}
			if h < 1 {
				h = 1
			}
			for dy := 0; dy < h && y+dy < height; dy++ {
				for dx := 0; dx < w && x+dx < width; dx++ {
					ox := x + dx
					oy := y + dy
					originX[oy*width+ox] = x
					originY[oy*width+ox] = y
					if dx != 0 || dy != 0 {
						cells[oy*width+ox] = EmptyCell // covered, not an origin
					}
				}
			}
		}
	}

	m.Cells = cells
	m.OriginX = originX
	m.OriginY = originY
	return m, nil
}

// Cell returns the origin character at (x, y). Out-of-range cells are empty.
func (m *Map) Cell(x, y int) byte {
	if x < 0 || y < 0 || x >= m.Width || y >= m.Height {
		return EmptyCell
	}
	return m.Cells[y*m.Width+x]
}

// TileAt returns the Tile occupying grid cell (x, y), resolving multi-cell
// tiles back to their origin. ok is false if the cell is empty or undefined.
func (m *Map) TileAt(x, y int) (Tile, bool) {
	if x < 0 || y < 0 || x >= m.Width || y >= m.Height {
		return Tile{}, false
	}
	idx := y*m.Width + x
	ox := m.OriginX[idx]
	oy := m.OriginY[idx]
	if ox < 0 {
		return Tile{}, false
	}
	c := m.Cells[oy*m.Width+ox]
	if IsEmpty(c) {
		return Tile{}, false
	}
	return Get(c)
}
