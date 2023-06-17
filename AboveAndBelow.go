package main

import (

	"github.com/hajimehoshi/ebiten/v2"
	//"image/color"
	"log"
	
)

const (win_width, win_height, pix_size = 640, 480, 10)

type Game struct {
	Pixels []byte //Pixel buffer
	X int
	Y int
	Arr [][] int
}

// Constructor for game
func NewGame() *Game {
    g := &Game{}
    g.Pixels = make([]byte, win_width * win_height * 4)
	//Init 2D array with side length width / pixel size
	g.Arr = make([][]int, win_height/pix_size)
	for i := range g.Arr {
		g.Arr[i] = make([]int, win_width/pix_size)
	}
    // initialize g.Pixels here...
    return g
}

func (g *Game) Update() error {

	return nil
}

func main(){
	ebiten.SetWindowSize(win_width, win_height)
	ebiten.SetWindowTitle("Above & Below")
	
	if err := ebiten.RunGame(NewGame()); err!= nil {
        log.Fatal(err)
    }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return win_width, win_height
}


func (g *Game) Draw(screen *ebiten.Image) {
	for row := 0; row < len(g.Arr); row++ {
		for col := 0; col < len(g.Arr[row]); col++ {
            i := (row * win_width + col) * 4
			if g.Arr[row][col] == 1 {
				g.Pixels[i+0] = 255 
				g.Pixels[i+1] = 255
				g.Pixels[i+2] = 255
				g.Pixels[i+3] = 255	
			}else {
				g.Pixels[i+0] = 0
				g.Pixels[i+1] = 0
				g.Pixels[i+2] = 0
				g.Pixels[i+3] = 0	
			}
        }
	}
	screen.WritePixels(g.Pixels)
}


func Scramble(array2D [][]int) [][]int{
	g.Arr = 




	return array2D
}