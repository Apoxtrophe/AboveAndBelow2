package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

type Game struct {
	buttons         []*Button
	Pixels          []byte
	Ichi            [][]int
	Ni              [][]int
	Index           int
	ParticleCount   int
	FPS             float64
	SelectedElement string
	BrushSize       int
	BrushImage      *ebiten.Image
}

func NewGame() *Game {
	g := &Game{}
	g.buttons = []*Button{
		NewButton(1*ButtonSize*1.2, 50, ButtonSize, "Elements/Hydrogen.png", func() { g.Index = 1 }),
		NewButton(2*ButtonSize*1.2, 50, ButtonSize, "Elements/Carbon.png", func() { g.Index = 6 }),
		NewButton(3*ButtonSize*1.2, 50, ButtonSize, "Elements/Oxygen.png", func() { g.Index = 8 }),
		NewButton(4*ButtonSize*1.2, 50, ButtonSize, "Elements/Silicon.png", func() { g.Index = 14 }),
		NewButton(5*ButtonSize*1.2, 50, ButtonSize, "Elements/Titanium.png", func() { g.Index = 22 }),
		NewButton(6*ButtonSize*1.2, 50, ButtonSize, "Elements/Mercury.png", func() { g.Index = 80 }),
	}
	g.BrushSize = 2
	g.Pixels = make([]byte, screenWidth*screenHeight*4)
	g.Ichi = make([][]int, screenHeight/PixelSize)
	g.Ni = make([][]int, screenHeight/PixelSize)
	for i := range g.Ichi {
		g.Ichi[i] = make([]int, screenWidth/PixelSize)
		g.Ni[i] = make([]int, screenWidth/PixelSize)
	}
	return g
}

func (g *Game) Update() error {
	g.UpdateBrushImage()
	g.UpdateUI()
	MouseInteract(g)
	g.AliveArray()
	return nil
}

func (g *Game) UpdateUI() {
	for _, button := range g.buttons {
		button.Update()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	for row := 0; row < len(g.Ichi); row++ {
		for col := 0; col < len(g.Ichi[row]); col++ {
			clr := ElementMap[g.Ichi[row][col]].Color
			for y := 0; y < PixelSize; y++ {
				for x := 0; x < PixelSize; x++ {
					i := ((row*PixelSize+y)*screenWidth + (col*PixelSize + x)) * 4
					g.Pixels[i+0] = clr.R
					g.Pixels[i+1] = clr.G
					g.Pixels[i+2] = clr.B
					g.Pixels[i+3] = clr.A
				}
			}
		}
	}
	screen.WritePixels(g.Pixels)
	g.DrawUI(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.2f\nNumber of Particles: %d\nElement: %s\nBrush Size: %d", g.FPS, g.ParticleCount, g.SelectedElement, g.BrushSize))
	g.DrawBrushGhost(screen)
}

func (g *Game) DrawUI(screen *ebiten.Image) {
	for _, button := range g.buttons {
		button.Draw(screen)
	}
}