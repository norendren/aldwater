package dungeonGen

import (
	"image/color"
	"math"
	"sync"

	"github.com/meshiest/go-dungeon/dungeon"
)

type Floor struct {
	Area    [][]*Tile
	Visible []*Tile
	Cols    int
	Rows    int
	sync.Mutex
}

func (f *Floor) ComputeFOV(pX, pY, r int) {
	f.Visible = nil
	f.Visible = append(f.Visible, f.Area[pY][pX])
	f.Area[pY][pX].Explored = true
	for i := 1; i <= 2; i++ {
		f.fov(pX, pY, pX, pY, 1, 1, 0, i, r)
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

func (f *Floor) fov(px, py, x, y, dist int, startSlope, endSlope float64, oct, rad int) {
	if dist > rad {
		return
	}
	initX, initY := initCoords(px, py, dist, oct)

	for startSlope != endSlope {
		if f.InBounds(initX, initY) && distTo(px, py, initX, initY) < rad {
			f.Visible = append(f.Visible, f.Area[initY][initX])
			f.Area[initY][initX].Explored = true
		}
		initX, initY = progress(initX, initY, oct)
		if oct == 2 {
			startSlope = invSlope(px, py, initX, initY)
		} else {
			startSlope = invSlope(px, py, initX, initY)
		}

		if startSlope == endSlope {
			f.fov(px, py, px, py, dist+1, 1, 0, oct, rad)
		}
	}

}
func progress(x, y, oct int) (int, int) {
	switch oct {
	case 1:
		return x + 1, y
	case 2:
		return x - 1, y
	}
	return 0, 0
}

func invSlope(x1, y1, x2, y2 int) float64 {
	fx1 := float64(x1)
	fx2 := float64(x2)
	fy1 := float64(y1)
	fy2 := float64(y2)
	return (fx2 - fx1) / (fy2 - fy1)
}

func slope(x1, y1, x2, y2 int) float64 {
	fx1 := float64(x1)
	fx2 := float64(x2)
	fy1 := float64(y1)
	fy2 := float64(y2)
	return (fy2 - fy1) / (fx2 - fx1)
}

func initCoords(x, y, dist, oct int) (int, int) {
	switch oct {
	case 1:
		return x - dist, y - dist
	case 2:
		return x + dist, y - dist
	}
	return 0, 0
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
