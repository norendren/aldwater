package main

import (
	"aldwater/displayResource"
	"aldwater/dungeonGen"
	"aldwater/player"
	"errors"
	"github.com/norendren/go-fov/fov"
	"image/color"
	"log"
	"strconv"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

//30x30
var (
	cols               = 31
	rows               = 31
	fontSize   float64 = 24
	normalFont font.Face
	width      int
	height     int
)

var p = player.Player{
	X:    4,
	Y:    4,
	Char: "@",
}

func init() {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72

	normalFont = truetype.NewFace(tt, &truetype.Options{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	width = cols*int(fontSize) + 20
	height = rows*int(fontSize) + 20
}

type Game struct {
	Pressed []ebiten.Key
	Level   *dungeonGen.Floor
	FOVCalc *fov.View
}

func (g *Game) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("ended by player")
	}
	p.HandleMovement(g.Level)
	//g.FOVCalc.Compute(g.Level, p.X, p.Y, 6)
	g.Level.ComputeFOV(p.X, p.Y, 6)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y, row := range g.Level.Area {
		for x, tile := range row {
			if g.Level.IsVisible(x, y) {
				text.Draw(screen, tile.Char, normalFont, tile.Posx, tile.Posy, tile.Color)
			}
			if tile.Explored && !g.Level.IsVisible(x, y) {
				text.Draw(screen, tile.Char, normalFont, tile.Posx, tile.Posy, displayResource.Color3)
			}
		}
	}

	text.Draw(screen,
		p.Char,
		normalFont,
		g.Level.Area[p.Y][p.X].Posx,
		g.Level.Area[p.Y][p.X].Posy,
		color.White)

	text.Draw(screen, strconv.Itoa(int(ebiten.CurrentTPS())), normalFont, 24, 700, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	g := &Game{
		Level:   dungeonGen.New(rows, cols, int(fontSize)),
		FOVCalc: fov.New(),
	}

	p.StartingPosition(g.Level)

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Aldwater")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
