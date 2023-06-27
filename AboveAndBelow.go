package main

import (
	//"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"fmt"
	"image/color"
	"log"
	"math/rand"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
)

func main() {
	ebiten.SetWindowTitle("Above & Below")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizable(true)
	if error := ebiten.RunGame(NewGame()); error != nil {
		log.Fatal(error)
	}
}

type Game struct {
	Pixels       []byte
	Ichi  [][]int
	Ni [][]int
	Index        int
	PixelSize    int
	//ebiten.Game
}

func NewGame() *Game {
	g := &Game{}
	g.PixelSize = 10
	g.Pixels = make([]byte, screenWidth*screenHeight*4)
	g.Ichi = make([][]int, screenHeight/g.PixelSize)
	g.Ni = make([][]int, screenHeight/g.PixelSize)
	for i := range g.Ichi {
		g.Ichi[i] = make([]int, screenWidth/g.PixelSize)
		g.Ni[i] = make([]int, screenWidth/g.PixelSize)
	}
	return g
}

func (g *Game) Update() error {
	MouseInteract(g)
	g.AliveArray()
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	for row := 0; row < len(g.Ichi); row++ {
		for col := 0; col < len(g.Ichi[row]); col++ {
			clr := ElementMap[g.Ichi[row][col]].Color
			for y := 0; y < g.PixelSize; y++ {
				for x := 0; x < g.PixelSize; x++ {
					i := ((row*g.PixelSize+y)*screenWidth + (col*g.PixelSize + x)) * 4
					g.Pixels[i+0] = clr.R
					g.Pixels[i+1] = clr.G
					g.Pixels[i+2] = clr.B
					g.Pixels[i+3] = clr.A
				}
			}
		}
	}
	screen.WritePixels(g.Pixels)
}

func MouseInteract(g *Game) {
	x, y := ebiten.CursorPosition()
	world_x, world_y := x/g.PixelSize, y/g.PixelSize
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
		g.Ichi[world_y][world_x] = g.Index
	}
	if mouse_two {
		g.Ichi[world_y][world_x] = 0
	}
}

func (g *Game) AliveArray() {
	aliveCells := make([][2]int, 0)
	for row := 0; row < len(g.Ichi); row++ {
		for col := 0; col < len(g.Ichi[row]); col++ {
			if g.Ichi[row][col] != 0 {
				aliveCells = append(aliveCells, [2]int{row, col})
			}
		}
	}
	rand.Shuffle(len(aliveCells), func(i, j int) {
		aliveCells[i], aliveCells[j] = aliveCells[j], aliveCells[i]
	})
	// Reset next buffer
	for row := range g.Ichi {
		for col := range g.Ichi[row] {
			g.Ni[row][col] = 0
		}
	}

	for _, pos := range aliveCells {
		row, col := pos[0], pos[1]
		switch g.Ichi[row][col] {
		case 6:
			g.Phys_Powder(row, col)
		case 14:
			g.Phys_Powder(row, col)
		case 22:
			g.Phys_Solid(row, col)
		default:

		}
	}
	g.Ichi, g.Ni = g.Ni, g.Ichi
}

var ElementMap = map[int]Element{
	0: {
		Color:   colornames.Black,
		Name:    "Void",
		Density: 100,
	},
	6: {
		Color:   colornames.Gray,
		Name:    "Carbon",
		Density: 22,
	},
	14: {
		Color:   colornames.Red,
		Name:    "Silicon",
		Density: 24,
	},
	22: {
		Color:   colornames.Cornflowerblue,
		Name:    "Titanium",
		Density: 45,
	},
}

type Element struct {
	Name    string
	Color   color.RGBA
	Density int
}

func (g *Game) IsFree(row, col int) bool {
	// Check if free space is available in both buffers
	if g.Ichi[row][col] == 0 && g.Ni[row][col] == 0 {
		return true
	}
	return false
}

func (g *Game) Phys_Solid(row, col int) {
	if col > 0 {
		g.Ni[row][col] = g.Ichi[row][col]
	}
}

func (g *Game) IsDenser (row, col int) bool {
	return ElementMap[g.Ichi[row][col]].Density > ElementMap[g.Ichi[row + 1][col]].Density 
}
func (g *Game) Phys_Powder(row, col int) {
	// Fall down -> -> Swap densities -> fall either side -> fall left -> fall right -> stay stationary 
	if row+1 < len(g.Ichi) && g.IsFree(row+1, col) {
		g.Ni[row+1][col] = g.Ichi[row][col]
		g.Ni[row][col] = 0
	} else if g.IsDenser(row, col) {
        Above := g.Ichi[row][col]
        Below := g.Ichi[row+1][col]
        g.Ni[row+1][col] = Above
        if g.Ni[row][col] == 0 {
            g.Ni[row][col] = Below
        }
	} else {
		
		leftFree := col-1 >= 0 && g.IsFree(row+1, col-1)
		rightFree := col+1 < len(g.Ichi[row]) && g.IsFree(row+1, col+1)

		if leftFree && rightFree {
			if rand.Float32() < 0.5 {
				g.Ni[row+1][col-1] = g.Ichi[row][col]
			} else {
				g.Ni[row+1][col+1] = g.Ichi[row][col]
			}
			g.Ni[row][col] = 0
		} else if leftFree {
			g.Ni[row+1][col-1] = g.Ichi[row][col]
			g.Ni[row][col] = 0
		} else if rightFree {
			g.Ni[row+1][col+1] = g.Ichi[row][col]
			g.Ni[row][col] = 0
		} else {
			g.Ni[row][col] = g.Ichi[row][col]
		}
	}
}