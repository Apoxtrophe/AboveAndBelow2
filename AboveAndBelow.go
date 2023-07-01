package main

//ANCHOR Imports
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

// ANCHOR Global Constants
const (
	screenWidth  = 1920
	screenHeight = 1080
)

// ANCHOR Main
func main() {
	ebiten.SetWindowTitle("Above & Below")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizable(true)
	if error := ebiten.RunGame(NewGame()); error != nil {
		log.Fatal(error)
	}
}

// ANCHOR Button Struct
type Button struct {
	x, y, w, h int
	color      color.Color
	onClick    func()
}

// ANCHOR Update Button
func (b *Button) Update() {
	x, y := ebiten.CursorPosition()
	if x >= b.x && y >= b.y && x < b.x+b.w && y < b.y+b.h {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.onClick()
		}
	}
}

// ANCHOR Draw Button
func (b *Button) Draw(screen *ebiten.Image) {
	button := ebiten.NewImage(b.w, b.h)
	button.Fill(b.color)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(button, op)
}

// ANCHOR Game Struct
type Game struct {
	buttons   []*Button
	Pixels    []byte
	Ichi      [][]int
	Ni        [][]int
	Index     int
	PixelSize int
}

// ANCHOR Game Constructor
func NewGame() *Game {
	g := &Game{}
	g.buttons = []*Button{
		NewButton(50, 50, 100, 50, ElementMap[6], func() { g.Index = 6 }),
		NewButton(200, 50, 100, 50, ElementMap[14], func() { g.Index = 14 }),
	}
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

// ANCHOR Update
func (g *Game) Update() error {
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

// ANCHOR Layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// ANCHOR DRAW
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
	g.DrawUI(screen)
}

func (g *Game) DrawUI(screen *ebiten.Image) {
	for _, button := range g.buttons {
		button.Draw(screen)
	}
}

func NewButton(x, y, w, h int, element Element, action func()) *Button {
	return &Button{
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		color:   element.Color,
		onClick: action,
	}
}

// ANCHOR Mouse Work
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

// ANCHOR Alive Array
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
		case 80:
			g.Phys_Liquid(row, col)
		default:
			
		}
	}
	g.Ichi, g.Ni = g.Ni, g.Ichi
	length := len(aliveCells)
	fmt.Println("          ", length)
}

// ANCHOR Element Map
var ElementMap = map[int]Element{
	0: {
		Color:   colornames.Black,
		Name:    "Void",
		Density: 0,
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
	80: {
		Color:   colornames.White,
		Name:    "Mercury",
		Density: 135,
	},
}

// ANCHOR Element Struct
type Element struct {
	Name    string
	Color   color.RGBA
	Density int
}

// ANCHOR SolidPhysics
func (g *Game) Phys_Solid(row, col int) {
	if col > 0 {
		g.Ni[row][col] = g.Ichi[row][col]
	}
}

// ANCHOR PowderPhysics
func (g *Game) Phys_Powder(row, col int) {
	// Fall down -> ?Swap densities -> fall either side -> fall left -> fall right -> stay stationary
	if g.canSwapTo(row, col, row+1, col) {
		g.swapParticle(row, col, row+1, col)
		return
	}
	leftFree := g.canSwapTo(row, col, row+1, col-1)
	rightFree := g.canSwapTo(row, col, row+1, col+1)
	if leftFree && rightFree {
		if rand.Float32() < 0.5 {
			g.swapParticle(row, col, row+1, col-1)
		} else {
			g.swapParticle(row, col, row+1, col+1)
		}
	} else if leftFree {
		g.swapParticle(row, col, row+1, col-1)
	} else if rightFree {
		g.swapParticle(row, col, row+1, col+1)
	} else {
		g.swapParticle(row, col, row, col)
	}
}

// ANCHOR LiquidPhysics
func (g *Game) Phys_Liquid(row, col int) {
	if g.canSwapTo(row, col, row+1, col) {
		g.swapParticle(row, col, row+1, col)
		return
	}
	leftFree := g.canSwapTo(row, col, row, col-1)
	rightFree := g.canSwapTo(row, col, row, col+1)
	if leftFree && rightFree {
		if rand.Float32() < 0.5 {
			g.swapParticle(row, col, row, col-1)
		} else {
			g.swapParticle(row, col, row, col+1)
		}
	} else if leftFree {
		g.swapParticle(row, col, row, col-1)
	} else if rightFree {
		g.swapParticle(row, col, row, col+1)
	} else {
		g.swapParticle(row, col, row, col)
	}

}

// ANCHOR Helper Functions
func (g *Game) IsFree(row, col int) bool {
	// Check if free space is available in both buffers
	if g.Ichi[row][col] == 0 && g.Ni[row][col] == 0 {
		return true
	}
	return false
}

func (g *Game) canSwapTo(sourceRow, sourceCol, targetRow, targetCol int) bool {
	return targetRow < len(g.Ichi) && g.isMoreDense(sourceRow, sourceCol, targetRow, targetCol) && g.NiFreeDebug(sourceRow, sourceCol, targetRow, targetCol)
}

func (g *Game) swapParticle(sourceRow, sourceCol, targetRow, targetCol int) {
	Source := g.Ichi[sourceRow][sourceCol]
	Target := g.Ichi[targetRow][targetCol]
    g.Ichi[sourceRow][sourceCol] = Target
    g.Ichi[targetRow][targetCol] = Source

    g.Ni[sourceRow][sourceCol] = g.Ichi[sourceRow][sourceCol]
    g.Ni[targetRow][targetCol] = g.Ichi[targetRow][targetCol]
}

func (g *Game) isMoreDense(sourceRow, sourceCol, targetRow, targetCol int) bool {
	return ElementMap[g.Ichi[sourceRow][sourceCol]].Density > ElementMap[g.Ichi[targetRow][targetCol]].Density
}

func (g *Game) NiFreeDebug(sourceRow, sourceCol, targetRow, targetCol int) bool{
	return g.Ni[sourceRow][sourceCol] == 0 && g.Ni[targetRow][targetCol] == 0
}