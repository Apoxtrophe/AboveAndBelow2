package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// ANCHOR Global Constants
const (
	screenWidth  = 1920
	screenHeight = 1080
	PixelSize    = 5
	Tickrate     = 60
	BrushAlpha   = 100
	ButtonSize   = 100
)

// ANCHOR Main
func main() {
	ebiten.SetTPS(Tickrate)
	ebiten.SetWindowTitle("Above & Below")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizable(true)
	if error := ebiten.RunGame(NewGame()); error != nil {
		log.Fatal(error)
	}
}
