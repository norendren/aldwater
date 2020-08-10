package dungeonGen

import (
	"image/color"
	"sync"

	"github.com/meshiest/go-dungeon/dungeon"
)

type Floor struct {
	Area [][]*Tile
	Cols int
	Rows int
	sync.Mutex
}

func (f *Floor) InBounds(x, y int) bool {
	if x >= f.Rows || y >= f.Cols {
		return false
	}
	if x < 0 || y < 0 {
		return false
	}
	return true
}

func (f *Floor) IsOpaque(x, y int) bool {
	if f.InBounds(x, y) && f.Area[y][x].Walkable {
		return false
	}
	return true
}

func (f *Floor) Index(x, y int) (int, int) {
	return y, x
}

func New(cols, rows, fontSize int) *Floor {
	d := dungeon.NewDungeon(30, 12)

	gameMap := make([][]*Tile, cols)
	for col := range gameMap {
		gameMap[col] = make([]*Tile, rows)
	}
	x := 15
	y := 24
	for c, row := range gameMap {
		for r, _ := range row {
			switch {
			case c == (len(gameMap)-1) || r == (len(row)-1):
				gameMap[c][r] = NewTile(false, "#", x, y, color.White)
			case d.Grid[c][r] == 1:
				gameMap[c][r] = NewTile(true, " .", x, y, color.White)
			case d.Grid[c][r] == 0:
				gameMap[c][r] = NewTile(false, "#", x, y, color.White)

			}
			//displayResource.Color3 good for explored but not visible
			x += fontSize
		}
		y += fontSize
		x = 15
	}
	return &Floor{
		Area: gameMap,
		Cols: cols,
		Rows: rows,
	}
}

type Tile struct {
	Walkable bool
	Char     string
	Posx     int
	Posy     int
	Color    color.Color
	Explored bool
}

func NewTile(walkable bool, char string, posx, posy int, c color.Color) *Tile {
	return &Tile{
		Walkable: walkable,
		Char:     char,
		Posx:     posx,
		Posy:     posy,
		Color:    c,
	}
}
