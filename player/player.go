package player

import (
	"aldwater/floor"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type Player struct {
	X    int
	Y    int
	Char string
}

func (p *Player) move(f *floor.Floor, addX, addY int) {
	destX := p.X + addX
	destY := p.Y + addY

	if (destX < 0 || destY < 0) || (destX >= f.Cols || destY >= f.Rows) {
		return
	}
	if !f.Area[destY][destX].Walkable {
		return
	}
	p.X = destX
	p.Y = destY
	return
}

func (p *Player) StartingPosition(f *floor.Floor) {
	for c, row := range f.Area {
		for r, tile := range row {
			if tile.Walkable {
				p.X = r
				p.Y = c
				return
			}
		}
	}
	return
}

func (p *Player) HandleMovement(f *floor.Floor) {
	if inpututil.IsKeyJustPressed(ebiten.KeyKP1) || inpututil.IsKeyJustPressed(ebiten.KeyB) {
		p.move(f, -1, 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP2) || inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		p.move(f, 0, 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP3) || inpututil.IsKeyJustPressed(ebiten.KeyN) {
		p.move(f, 1, 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP4) || inpututil.IsKeyJustPressed(ebiten.KeyH) {
		p.move(f, -1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP6) || inpututil.IsKeyJustPressed(ebiten.KeyL) {
		p.move(f, 1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP7) || inpututil.IsKeyJustPressed(ebiten.KeyY) {
		p.move(f, -1, -1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP8) || inpututil.IsKeyJustPressed(ebiten.KeyK) {
		p.move(f, 0, -1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyKP9) || inpututil.IsKeyJustPressed(ebiten.KeyU) {
		p.move(f, 1, -1)
	}

}
