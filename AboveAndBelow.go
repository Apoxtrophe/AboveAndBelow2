package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	//"image/color"
	"fmt"
	"log"
	"math/rand"
)

const (
	win_width, win_height, pix_size = 640, 480, 10
)


// ANCHOR Game struct
type Game struct {
	Pixels   []byte //Pixel buffer
	X        int
	Y        int
	Arr      [][]*Element
	Arr_Next [][]*Element
	Index    int
}


// ANCHOR Game Constructor
func NewGame() *Game {
	g := &Game{}
	g.Pixels = make([]byte, win_width*win_height*4)
	//Init 2D array with side length width / pixel size
	g.Arr = make([][]*Element, win_height/pix_size)
	g.Arr_Next = make([][]*Element, win_height/pix_size)
	for i := range g.Arr {
		g.Arr[i] = make([]*Element, win_width/pix_size)
		g.Arr_Next[i] = make([]*Element, win_width/pix_size)
	}
	g.Index = 0
	// initialize g.Pixels here...
	return g
}

// ANCHOR Update Function
func (g *Game) Update() error {
	MouseInteract(g)
	g.AliveArray()
	return nil
}


// ANCHOR Cell Update
func (g *Game) AliveArray() {
	aliveCells := make([][2]int, 0)
	for row := 0; row < len(g.Arr); row++ {
		for col := 0; col < len(g.Arr[row]); col++ {
			if g.Arr[row][col] != nil && g.Arr[row][col].Type != 0{
				aliveCells = append(aliveCells, [2]int{row, col})
			}
		}
	}
	// shuffle update order of alive cells
	rand.Shuffle(len(aliveCells), func(i, j int) {
		aliveCells[i], aliveCells[j] = aliveCells[j], aliveCells[i]
	})

	// reset Arr_Next
	for row := range g.Arr_Next {
		for col := range g.Arr_Next[row] {
			g.Arr_Next[row][col] = Elements[0]
		}
	}

	// Swap Arr and Arr_Next
	g.Arr, g.Arr_Next = g.Arr_Next, g.Arr
}

// ANCHOR Mouse Interaction
func MouseInteract(g *Game) {
	x, y := ebiten.CursorPosition()
	world_x, world_y := x/pix_size, y/pix_size
	mouse_one := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	mouse_two := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	//Scroll detection
	_, deltaY := ebiten.Wheel()
	if deltaY > 0 {
		g.Index++
	} else if deltaY < 0 {
		g.Index--
	}
	fmt.Println(g.Index)
	//Clicking detection
	if mouse_one {
		g.Arr[world_y][world_x] = Elements[g.Index]
	}
	if mouse_two {
		g.Arr[world_y][world_x] = Elements[0]
	}
}


// ANCHOR MAIN
func main() {
	ebiten.SetWindowSize(win_width, win_height)
	ebiten.SetWindowTitle("Above & Below")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}


// ANCHOR Layout Function
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return win_width, win_height
}

// ANCHOR Draw Function
func (g *Game) Draw(screen *ebiten.Image) {
	for row := 0; row < len(g.Arr); row++ {
		for col := 0; col < len(g.Arr[row]); col++ {
			// Calculate the color based on g.Arr[row][col]
			clr := g.Arr[row][col].Color
			// Set the color of a 10x10 block of pixels in g.Pixels
			for y := 0; y < pix_size; y++ {
				for x := 0; x < pix_size; x++ {
					i := ((row*pix_size+y)*win_width + (col*pix_size + x)) * 4
					g.Pixels[i+0] = clr.R
					g.Pixels[i+1] = clr.G
					g.Pixels[i+2] = clr.B
					g.Pixels[i+3] = clr.A
				}
			}
		}
	}

	//fmt.Println(ebiten.ActualFPS())
	screen.WritePixels(g.Pixels)
}
// ANCHOR Element Struct
type Element struct {
	Type int
	Color color.RGBA
	Weight int
	Gravity int
}


// ANCHOR Element Properties
var Elements = map[int]*Element{
	0: &Element{
		Type: 0,
		Color: colornames.Black,
		Weight: 0,
        Gravity: 0,
	},
	14: &Element{
		Type: 14,
		Color: colornames.Yellow,
		Weight: 10,
        Gravity: 0,
	},
	22: &Element{
		Type: 22,
		Color: colornames.Darkgray,
		Weight: 100,
        Gravity: 0,
	},
}

type Solid interface {
	IsSolid() bool
}
type Powder interface {
	IsPowder() bool
}
type Liquid interface {
	IsLiquid() bool
}
type Gas interface {
	IsGas() bool
}

type Titanium struct {
	Element
}
func (t Titanium) IsSolid() bool {
    return true
}

type Silicon struct {
	Element
}
func (s Silicon) IsPowder() bool {
	return true
}

