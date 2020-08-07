package dungeonGen

import (
	"aldwater/player"
	"image/color"
	"math"
	"sync"

	"github.com/meshiest/go-dungeon/dungeon"
)

type Floor struct {
	Area     [][]*Tile
	Visible  []*Tile
	Explored []*Tile
	Cols     int
	Rows     int
	sync.Mutex
}

func (f *Floor) ComputeFOV(p *player.Player, r int) {
	f.Visible = nil
	for i := 1; i <= 8; i++ {
		f.fov(p, p.X, p.Y, 1, 0, 1, i, r)
	}
}

func (f *Floor) fov(p *player.Player, x, y, dist int, startSlope, endSlope float64, oct, rad int) {
	if dist > rad {
		return
	}
	for startSlope < endSlope {

	}

}

func mapXY(x, y, dist, oct int) (int, int) {
	switch oct {
	case 1:
		return x + dist, y + dist
	}
	return 0, 0
}

func distTo(x1, y1, x2, y2 int) int {
	vx := math.Pow(float64(x1-x2), 2)
	vy := math.Pow(float64(y1-y2), 2)
	return int(math.Sqrt(vx + vy))
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
