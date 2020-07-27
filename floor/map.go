package floor

import (
	"github.com/meshiest/go-dungeon/dungeon"
	"image/color"
)

type Floor struct {
	Area [][]*Tile
	Cols int
	Rows int
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
				gameMap[c][r] = NewTile(false, "#", x, y)
			case d.Grid[c][r] == 1:
				gameMap[c][r] = NewTile(true, " .", x, y)
			case d.Grid[c][r] == 0:
				gameMap[c][r] = NewTile(false, "#", x, y)

			}
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
}

func NewTile(walkable bool, char string, posx, posy int) *Tile {
	return &Tile{
		Walkable: walkable,
		Char:     char,
		Posx:     posx,
		Posy:     posy,
		Color:    color.White,
	}
}
