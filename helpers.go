package main

// Used in determining brush size
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

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