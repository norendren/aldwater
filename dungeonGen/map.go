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
		f.fov(pX, pY, 1, 1, 0, i, r)
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

func (f *Floor) fov(px, py, dist int, startSlope, endSlope float64, oct, rad int) {
	if dist > rad {
		return
	}
	initX, initY := initCoords(px, py, dist, oct, startSlope)
	//consecutiveBlocked := false

	//for i := 0; i < 200; i++ {
	for startSlope != endSlope {
		if f.InBounds(initX, initY) && distTo(px, py, initX, initY) < rad {
			f.Visible[fmt.Sprintf("%d%d", initX, initY)] = f.Area[initY][initX]
			f.Area[initY][initX].Explored = true

			//if oct == 1 {
			//	if !f.Area[initY][initX].Walkable && !consecutiveBlocked {
			//		f.fov(px, py, dist+1, startSlope, mSlope(px, py, initX, initY, oct), oct, rad)
			//		consecutiveBlocked = true
			//	} else if f.Area[initY][initX].Walkable && consecutiveBlocked {
			//		startSlope = mSlope(px, py, initX, initY, oct)
			//		consecutiveBlocked = false
			//	}
			//}
		}
		initX, initY = progress(initX, initY, px, py, oct)

		if mSlope(px, py, initX, initY, oct) == endSlope {
			if f.InBounds(initX, initY) && distTo(px, py, initX, initY) < rad {
				f.Visible[fmt.Sprintf("%d%d", initX, initY)] = f.Area[initY][initX]
				f.Area[initY][initX].Explored = true
				if !f.Area[initY][initX].Walkable {
					return
				}
			}

			f.fov(px, py, dist+1, startSlope, endSlope, oct, rad)
			return
		}
	}

}
func progress(x, y, px, py, oct int) (int, int) {
	switch oct {
	case 1, 6:
		x += 1
		return x, y
	case 2, 5:
		x -= 1
		return x, y
	case 3, 8:
		y += 1
		return x, y
	case 4, 7:
		y -= 1
		return x, y
	}
	return 0, 0
}

func initCoords(x, y, dist, oct int, s float64) (int, int) {
	//switch oct {
	//case 1:
	//	return x - int(s*float64(dist)), y - dist
	//case 8:
	//	return x - dist, y - int(s*float64(dist))
	//case 2, 3:
	//	return x + dist, y - dist
	//case 4, 5:
	//	return x + dist, y + dist
	//case 6, 7:
	//	return x - dist, y + dist
	//}
	switch oct {
	case 1:
		return x - int(s*float64(dist)), y - dist
	case 2:
		return x + int(s*float64(dist)), y - dist
	case 3:
		return x + dist, y - int(s*float64(dist))
	case 4:
		return x + dist, y + int(s*float64(dist))
	case 5:
		return x + int(s*float64(dist)), y + dist
	case 6:
		return x - int(s*float64(dist)), y + dist
	case 7:
		return x - dist, y + int(s*float64(dist))
	case 8:
		return x - dist, y - int(s*float64(dist))
	}
	return 0, 0
}

func mSlope(x1, y1, x2, y2, oct int) float64 {
	switch oct {
	case 1, 2, 5, 6:
		return invSlope(x1, y1, x2, y2)
	default:
		return slope(x1, y1, x2, y2)
	}
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
