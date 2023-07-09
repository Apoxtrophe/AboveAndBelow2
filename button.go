package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
)

type Button struct {
	x, y, size int
	img        *ebiten.Image
	onClick    func()
}

func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(b.img, op)
	ebitenutil.DebugPrintAt(screen, "Button", b.x, b.y)
}

func (b *Button) Update() {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x > b.x && x < b.x+b.size && y > b.y && y < b.y+b.size {
			b.onClick()
		}
	}
}

func NewButton(x, y, size int, imgPath string, action func()) *Button {
	img, err := NewSizedImageFromFile(imgPath, size)
	if err != nil {
		log.Fatal(err)
	}
	return &Button{
		x:       x,
		y:       y,
		size:    size,
		img:     img,
		onClick: action,
	}
}

func NewSizedImageFromFile(imgPath string, size int) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(imgPath)
	if err != nil {
		return nil, err
	}
	resizedImg := ebiten.NewImage(size, size)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(size)/float64(img.Bounds().Dx()), float64(size)/float64(img.Bounds().Dy()))
	resizedImg.DrawImage(img, op)
	return resizedImg, nil
}

