package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
)

type Button struct {
	x, y, size int
	img        *ebiten.Image
	onClick    func()
}

func (b *Button) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(b.img, op)
	ebitenutil.DebugPrintAt(screen, "Button", b.x, b.y)
}

func (b *Button) update() {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x > b.x && x < b.x+b.size && y > b.y && y < b.y+b.size {
			b.onClick()
		}
	}
}

func NewButton(x, y, size int, onClick func()) *Button {
	img := ebiten.NewImage(size, size)
	img.Fill(color.White)
	return &Button{
		x:       x,
		y:       y,
		size:    size,
		img:     img,
		onClick: onClick,
	}
}
