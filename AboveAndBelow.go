package main

//ANCHOR Imports
import (
	//"image/color"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
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

// ANCHOR Button Struct
type Button struct {
	x, y, size int
	img *ebiten.Image
	onClick    func()
}


//ANCHOR New Button Function
func NewButton(x, y, size int, imgPath string, action func()) *Button {
    img, err := NewSizedImageFromFile(imgPath, size)
    if err != nil {
        log.Fatal(err)
    }
    return &Button{
        x:       x,
        y:       y,
		size: size,
        img:     img,
        onClick: action,
    }
}

//ANCHOR Image Resizer
func NewSizedImageFromFile(imgPath string, size int) (*ebiten.Image, error){
	img, _, err := ebitenutil.NewImageFromFile(imgPath)
	if err != nil{
		return nil, err
	}
	resizedImg := ebiten.NewImage(size, size)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(size)/float64(img.Bounds().Dx()), float64(size)/float64(img.Bounds().Dy()))
	resizedImg.DrawImage(img,op)
	return resizedImg, nil
}

// ANCHOR Update Button
func (b *Button) Update() {
	x, y := ebiten.CursorPosition()
	if x >= b.x && y >= b.y && x < b.x+b.size && y < b.y+b.size {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.onClick()
		}
	}
}

// ANCHOR Draw Button
func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(b.img, op)
}

// ANCHOR Game Struct
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

// ANCHOR Game Constructor
func NewGame() *Game {
	g := &Game{}
	g.buttons = []*Button{
		NewButton(100, 50, ButtonSize, "Elements/Carbon.png", func() { g.Index = 6 }),
		NewButton(220, 50, ButtonSize, "Elements/Oxygen.png", func() { g.Index = 8 }),
		NewButton(340, 50, ButtonSize, "Elements/Silicon.png", func() { g.Index = 14 }),
		NewButton(460, 50, ButtonSize, "Elements/Titanium.png", func() { g.Index = 22 }),
		NewButton(580, 50, ButtonSize, "Elements/Mercury.png", func() { g.Index = 80 }),
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

// ANCHOR Update
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

// ANCHOR Layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// ANCHOR DRAW
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
	//Print debug information
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.2f\nNumber of Particles: %d\nElement: %s\nBrush Size: %d", g.FPS, g.ParticleCount, g.SelectedElement, g.BrushSize))

	// Draw brush Size
	g.DrawBrushGhost(screen)
}

// ANCHOR Update Brush Image
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


//ANCHOR Draw Brush Ghost
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



//ANCHOR Draw UI
func (g *Game) DrawUI(screen *ebiten.Image) {
	for _, button := range g.buttons {
		button.Draw(screen)
	}
}

// ANCHOR Mouse Work


//TODO - Fix this ugly ass function
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
		g.BrushSize ++
	} else if wheelY < 0 {
		g.BrushSize --
	}

	if g.BrushSize < 1 {
		g.BrushSize = 1
	}
	if g.BrushSize > 100 {
		g.BrushSize = 100
	}
	radius := float64(g.BrushSize) / 2.0
	// Clicking detection
	if mouse_one || mouse_two {
		dx := float64(world_x - prevMouseX)
		dy := float64(world_y - prevMouseY)
		length := math.Sqrt(float64(dx*dx + dy*dy))
		if length > 0{
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
						if mouse_one {
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
		if row <= 1 || col <= 1 || row >= len(g.Ichi)-2 || col >= len(g.Ichi[0])-2 {
			g.Ichi[row][col] = 0
			continue
		}
		switch g.Ichi[row][col] {
		case 6:
			g.Phys_Powder(row, col)
		case 8:
			g.Phys_Gas(row, col)
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
	//Icky Icky debug stuff
	g.SelectedElement = ElementMap[g.Index].Name
	g.FPS = ebiten.ActualTPS()
	g.ParticleCount = len(aliveCells)
}

// ANCHOR Element Map
var ElementMap = map[int]Element{
	0: {
		Color:   colornames.Black,
		Name:    "Void",
		Density: 0,
		isFluid: true,
	},
	6: {
		Color:   colornames.Gray,
		Name:    "Carbon",
		Density: 22,
		isFluid: false,
	},
	8: {
		Color:   colornames.Aqua,
		Name:    "Oxygen",
		Density: 1,
		isFluid: true,
	},
	14: {
		Color:   colornames.Red,
		Name:    "Silicon",
		Density: 24,
		isFluid: false,
	},
	22: {
		Color:   colornames.Cornflowerblue,
		Name:    "Titanium",
		Density: 45,
		isFluid: false,
	},
	80: {
		Color:   colornames.White,
		Name:    "Mercury",
		Density: 13,
		isFluid: true,
	},
}

// ANCHOR Element Struct
type Element struct {
	Name    string
	Color   color.RGBA
	Density int
	isFluid bool
}

// ANCHOR SolidPhysics
func (g *Game) Phys_Solid(row, col int) {
	if col > 0 {
		g.Ni[row][col] = g.Ichi[row][col]
	}
}

// ANCHOR PowderPhysics
func (g *Game) Phys_Powder(row, col int) {
	// Fall down -> fall either side -> fall left -> fall right -> stay stationary
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

// ANCHOR GasPhysics
func (g *Game) Phys_Gas(row, col int) {
	newRow, newCol := g.randomPosition(row, col)
	if g.canSwapTo(row, col, newRow, newCol) {
		g.swapParticle(row, col, newRow, newCol)
	} else {
		g.swapParticle(row, col, row, col)
	}
}

// ANCHOR Helper Functions
func (g *Game) randomPosition(row, col int) (int, int) {
	positions := [8][2]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}
	randValue := globalRand.Intn(8)
	return row + positions[randValue][0], col + positions[randValue][1]
}
func (g *Game) canSwapTo(sourceRow, sourceCol, targetRow, targetCol int) bool {
	return targetRow < len(g.Ichi) && g.isMoreDense(sourceRow, sourceCol, targetRow, targetCol) && g.NiFree(sourceRow, sourceCol, targetRow, targetCol) && ElementMap[g.Ichi[targetRow][targetCol]].isFluid
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

func (g *Game) NiFree(sourceRow, sourceCol, targetRow, targetCol int) bool {
	return g.Ni[sourceRow][sourceCol] == 0 && g.Ni[targetRow][targetCol] == 0
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



// Used in determining brush size
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
