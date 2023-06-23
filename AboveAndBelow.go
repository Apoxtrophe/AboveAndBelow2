package main

import (
	//"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"fmt"
	"log"
	"math/rand"
)

const (
    screenWidth  = 1920
    screenHeight = 1080
)

func main () {
	ebiten.SetWindowTitle ("Above & Below")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizable(true)
	if error := ebiten.RunGame(NewGame()); error!= nil {
		log.Fatal(error)		
	}
}

type Game struct {
	Pixels []byte
	FirstBuffer[][] int
	SecondBuffer [][] int
	Index int
	PixelSize int
    //ebiten.Game
}

func NewGame() *Game {	
	g := &Game{}
	g.PixelSize = 10
	g.Pixels = make([]byte, screenWidth*screenHeight*4)
	g.FirstBuffer = make ([][]int, screenHeight / g.PixelSize)
	g.SecondBuffer = make ([][]int, screenHeight / g.PixelSize)
	for i := range g.FirstBuffer {
		g.FirstBuffer[i] = make([]int, screenWidth / g.PixelSize)
		g.SecondBuffer[i] = make([]int, screenWidth / g.PixelSize)
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
    for row := 0; row < len(g.FirstBuffer); row++ {
		for col := 0; col < len(g.FirstBuffer[row]); col++ {
			clr := ElementMap [g.FirstBuffer[row][col]].Color
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
	if mouse_one { g.FirstBuffer[world_y][world_x] = g.Index }
	if mouse_two { g.FirstBuffer[world_y][world_x] = 0 }
}

func (g *Game) AliveArray() {
	aliveCells := make([][2]int,0)
	for row := 0; row < len(g.FirstBuffer); row++ {
		for col := 0; col < len(g.FirstBuffer[row]); col++ {
           if g.FirstBuffer[row][col] != 0 {
			aliveCells = append(aliveCells, [2]int{row, col})
		   }
        }
	}
	rand.Shuffle(len(aliveCells), func(i,j int){
		aliveCells[i], aliveCells[j] = aliveCells[j], aliveCells[i]
	})
	// Reset next buffer
	for row := range g.FirstBuffer {
		for col := range g.FirstBuffer[row] {
			g.SecondBuffer[row][col] = 0
		}
	}

	for _, pos := range aliveCells {
		row, col := pos[0], pos[1]
		switch g.FirstBuffer[row][col] {
		case 6: 
			g.Phys_Powder(row,col)
			g.DensityCompare(row,col)
		case 14: 	
			g.Phys_Powder(row, col)
			g.DensitySwap(row,col)
		case 22:
			g.Phys_Solid(row, col)
		default: 
		
		}
	}
	g.FirstBuffer, g.SecondBuffer = g.SecondBuffer, g.FirstBuffer
}


var ElementMap = map[int] Element {
	0: {
		Color: colornames.Black,
		Name: "Void",
		Density: 0,
	},
	6: {
		Color: colornames.Gray,
        Name: "Carbon",
        Density: 22,
	},
	14: {
		Color: colornames.Red,
        Name: "Silicon",
		Density: 24,
	},
	22: {
		Color: colornames.Cornflowerblue ,
        Name: "white",
		Density: 45,
	},
}

type Element struct {
	Name string
	Color color.RGBA
	Density int
}

func (g * Game) FreeSpace (row, col int) bool {
	if g.FirstBuffer[row][col] == 0 && g.SecondBuffer[row][col] == 0{
		return true
	}
	return false
}

func (g *Game) Phys_Solid (row, col int){
	if col > 0 {
		g.SecondBuffer[row][col] = g.FirstBuffer[row][col]
	}
}

func (g *Game) DensityCompare (row, col int) bool {
	if ElementMap[g.FirstBuffer[row][col]].Density > ElementMap[g.SecondBuffer[row -1][col]].Density {
		return true
	}
	return false
}

func (g *Game) DensitySwap (row, col int){
	if g.DensityCompare(row, col){
		g.SecondBuffer[row -1][col] = g.FirstBuffer[row][col]
		g.FirstBuffer[row][col] = g.FirstBuffer[row-1][col]
	}
}

func (g *Game) Phys_Powder (row, col int) {
	if row + 1 < len(g.FirstBuffer) && g.FreeSpace(row + 1, col) {
		g.SecondBuffer[row + 1][col] = g.FirstBuffer[row][col]
		g.SecondBuffer[row][col] = 0
		} else {
			// If space below is not free, then check if diagonal spaces are free
			// If both diagonal spaces are free, choose a random direction
			// If one is free, move in that direction
			// If none are free, stay in place
			leftFree := col - 1 >= 0 && g.FreeSpace(row + 1, col - 1)
			rightFree := col + 1 < len(g.FirstBuffer[row]) && g.FreeSpace(row + 1, col + 1)
	
			if leftFree && rightFree {
				if rand.Float32() < 0.5 {
					g.SecondBuffer[row+1][col-1] = g.FirstBuffer[row][col]
				} else {
					g.SecondBuffer[row+1][col+1] = g.FirstBuffer[row][col]
				}
				g.SecondBuffer[row][col] = 0
			} else if leftFree {
				g.SecondBuffer[row+1][col-1] = g.FirstBuffer[row][col]
				g.SecondBuffer[row][col] = 0
			} else if rightFree {
				g.SecondBuffer[row+1][col+1] = g.FirstBuffer[row][col]
				g.SecondBuffer[row][col] = 0
			} else {
				g.SecondBuffer[row][col] = g.FirstBuffer[row][col]
			}
		}
}

