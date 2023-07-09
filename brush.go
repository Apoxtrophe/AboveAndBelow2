package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BrushAlpha = 100
)

func (game *Game) UpdateBrushImage() {
	radius := float64(game.BrushSize) / 2.0

	game.BrushImage = ebiten.NewImage(game.BrushSize+1, game.BrushSize+1)

	factor := float64(BrushAlpha) / 255
	red, green, blue := uint8(factor*255), uint8(factor*255), uint8(factor*255)

	for row := -radius; row <= radius; row++ {
		for col := -radius; col <= radius; col++ {
			dist := math.Hypot(float64(row), float64(col))
			if dist <= radius {
				ix := int(math.Round(radius + col))
				iy := int(math.Round(radius + row))
				game.BrushImage.Set(ix, iy, color.RGBA{red, green, blue, uint8(BrushAlpha)})
			}
		}
	}
}

func (game *Game) DrawBrushGhost(screen *ebiten.Image) {
	radius := float64(game.BrushSize) / 2.0
	offsetX := radius * float64(PixelSize)
	offsetY := radius * float64(PixelSize)

	mouseX, mouseY := ebiten.CursorPosition()
	mouseX = (mouseX / PixelSize) * PixelSize
	mouseY = (mouseY / PixelSize) * PixelSize

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(PixelSize), float64(PixelSize))
	op.GeoM.Translate(float64(mouseX)-offsetX, float64(mouseY)-offsetY)
	screen.DrawImage(game.BrushImage, op)
}