package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"image/color"
)

var prevMouseX, prevMouseY int

func MouseInteract(g *Game) {
	x, y := ebiten.CursorPosition()

	// Clamp to world bounds
	world_x := clamp(x/PixelSize, 0, len(g.Ichi[0])-1)
	world_y := clamp(y/PixelSize, 0, len(g.Ichi)-1)
	mouse_one := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	mouse_two := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	// Brush Sizing
	_, wheelY := ebiten.Wheel()
	if wheelY > 0 {
		g.BrushSize = g.BrushSize +2
	} else if wheelY < 0 {
		g.BrushSize = g.BrushSize -2
	}

	if g.BrushSize < 1 {
		g.BrushSize = 1
	}
	if g.BrushSize > 101 {
		g.BrushSize = 101
	}
	radius := float64(g.BrushSize) / 2.0
	// Clicking detection
	if mouse_one || mouse_two {
		dx := float64(world_x - prevMouseX)
		dy := float64(world_y - prevMouseY)
		length := math.Sqrt(float64(dx*dx + dy*dy))
		if length > 0 {
			dx /= (length)
			dy /= (length)
		}

		for i := 0; i <= int(length); i++ {
			x := prevMouseX + int(float64(i)*dx)
			y := prevMouseY + int(float64(i)*dy)
			for row := -radius; row <= radius; row++ {
				for col := -radius; col <= radius; col++ {
					dist := math.Hypot(float64(row), float64(col))
					if dist <= radius {
						ix := clamp(x+int(col), 0, len(g.Ichi[0])-1)
						iy := clamp(y+int(row), 0, len(g.Ichi)-1)
						if mouse_one  {
							g.Ichi[iy][ix] = g.Index
						} else {
							g.Ichi[iy][ix] = 0
						}
					}
				}
			}
		}
	}

	prevMouseX = world_x
	prevMouseY = world_y
}

// Used in Mouse selection currently
func clamp(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

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