package main

import(
	"math/rand"
)

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
