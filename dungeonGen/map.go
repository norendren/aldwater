package dungeonGen

import (
	"fmt"
	"image/color"
	"math"
	"sync"

	"github.com/meshiest/go-dungeon/dungeon"
)

type Floor struct {
	Area    [][]*Tile
	Visible map[string]*Tile
	Cols    int
	Rows    int
	sync.Mutex
}

func (f *Floor) ComputeFOV(pX, pY, r int) {
	f.Visible = make(map[string]*Tile)
	f.Visible[fmt.Sprintf("%d%d", pX, pY)] = f.Area[pY][pX]
	f.Area[pY][pX].Explored = true
	for i := 1; i <= 8; i++ {
		f.fov(pX, pY, 1, 0, 1, i, r)
	}
}

func (f *Floor) InBounds(x, y int) bool {
	if x >= f.Cols || y >= f.Cols {
		return false
	}
	if x < 0 || y < 0 {
		return false
	}
	return true
}

func (f *Floor) fov(px, py, dist int, lowSlope, highSlope float64, oct, rad int) {
	if dist > rad {
		return
	}
	low := math.Floor(lowSlope*float64(dist) + 0.5)
	high := math.Floor(highSlope*float64(dist) + 0.5)
	inGap := false

	for height := low; height <= high; height++ {
		mapx, mapy := distHeightXY(px, py, dist, int(height), oct)
		if f.InBounds(mapx, mapy) && distTo(px, py, mapx, mapy) < rad {
			f.Visible[fmt.Sprintf("%d%d", mapx, mapy)] = f.Area[mapy][mapx]
			f.Area[mapy][mapx].Explored = true
		}

		if f.InBounds(mapx, mapy) && !f.Area[mapy][mapx].Walkable {
			if inGap {
				f.fov(px, py, dist+1, lowSlope, (height-0.5)/float64(dist), oct, rad)
			}
			lowSlope = (height + 0.5) / float64(dist)
			inGap = false
		} else {
			inGap = true
			if height == high {
				f.fov(px, py, dist+1, lowSlope, highSlope, oct, rad)
			}
		}
	}
}

func distHeightXY(px, py, d, h, oct int) (int, int) {
	if oct&0x1 > 0 {
		d = -d
	}
	if oct&0x2 > 0 {
		h = -h
	}
	if oct&0x4 > 0 {
		return px + h, py + d
	}
	return px + d, py + h
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
