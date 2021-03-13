package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"os"
	"time"
)

var (
	screenWidth  = 100
	screenHeight = 30
)

func main() {
	floor := NewTile(true, '.', tcell.StyleDefault)
	wall := NewTile(false, '#', tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorBlack))

	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Println(err)
	}

	if err := s.Init(); err != nil {
		fmt.Println(err)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.Clear()

	player := Entity{
		x:     10,
		y:     10,
		char:  '@',
		style: tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(tcell.ColorBlack),
	}

	gameMap := make([][]*Tile, screenWidth)
	for i := range gameMap {
		gameMap[i] = make([]*Tile, screenHeight)
	}

	for i := range gameMap {
		for j := range gameMap[i] {
			gameMap[i][j] = floor
		}
	}

	gameMap[22][8] = wall
	gameMap[22][9] = wall
	gameMap[22][10] = wall
	gameMap[22][11] = wall

	gameMap[27][8] = wall
	gameMap[27][9] = wall
	gameMap[27][10] = wall
	gameMap[27][11] = wall

	gameMap[24][8] = wall
	gameMap[25][8] = wall
	gameMap[26][8] = wall

	gameMap[23][11] = wall
	gameMap[24][11] = wall
	gameMap[25][11] = wall
	gameMap[26][11] = wall
	gameMap[27][11] = wall

	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()

			switch ev := ev.(type) {
			case *tcell.EventKey:
				//fmt.Println(ev.Name())
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyUp:
					player.move(gameMap, 0, -1)
				case tcell.KeyDown:
					player.move(gameMap, 0, 1)
				case tcell.KeyLeft:
					player.move(gameMap, -1, 0)
				case tcell.KeyRight:
					player.move(gameMap, 1, 0)
				case tcell.KeyRune:
					switch ev.Rune() {
					case '4':
						player.move(gameMap, -1, 0)
					case 'h':
						player.move(gameMap, -1, 0)
					case 'j':
						player.move(gameMap, 0, 1)
					case 'k':
						player.move(gameMap, 0, -1)
					case 'l':
						player.move(gameMap, 1, 0)
					case 'y':
						player.move(gameMap, -1, -1)
					case 'u':
						player.move(gameMap, 1, -1)
					case 'b':
						player.move(gameMap, -1, 1)
					case 'n':
						player.move(gameMap, 1, 1)
					}
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	for {
		select {
		case <-quit:
			s.Fini()
			os.Exit(0)
		case <-time.After(40 * time.Millisecond):

		}

		for i := range gameMap {
			for j := range gameMap[i] {
				s.SetContent(i, j, gameMap[i][j].char, nil, gameMap[i][j].style)
			}
		}
		s.SetContent(player.x, player.y, player.char, nil, player.style)
		s.Show()

	}
}

type Entity struct {
	x     int
	y     int
	char  rune
	style tcell.Style
}

func (e *Entity) move(gameMap [][]*Tile, x, y int) {
	destX := e.x + x
	destY := e.y + y

	if (destX < 0 || destY < 0) || (destX >= screenWidth || destY >= screenHeight) {
		return
	}
	if !gameMap[destX][destY].walkable {
		return
	}
	e.x = destX
	e.y = destY
	return
}

type Tile struct {
	walkable bool
	char     rune
	style    tcell.Style
}

func NewTile(walkable bool, char rune, style tcell.Style) *Tile {
	return &Tile{
		walkable: walkable,
		char:     char,
		style:    style,
	}
}
