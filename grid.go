package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawGrid(screen *ebiten.Image, g *Game) {

	// draw background of triples
	for y, row := range g.grid {
		for x := range row {
			if g.triplesMask[y][x] {
				g.grid[y][x].DrawBackground(screen, darkButter)
			}
		}
	}

	// draw cursor
	cursorColors := [2]color.Color{color.White, color.Black}
	blink := g.frameCount / blinkInverval % 2
	if g.swapSquare.x != -1 {
		// we are in swap mode, faster blink, brighter colors
		blink = g.frameCount / (blinkInverval / 2) % 2
		cursorColors = [2]color.Color{brightRed, lightButter}
	}
	g.grid[g.cursorSquare.y][g.cursorSquare.x].DrawBackground(screen, cursorColors[blink])

	// draw gems
	for y, row := range g.grid {
		for x, _ := range row {
			g.grid[y][x].DrawGem(screen, g.gemImages[g.grid[y][x].color], g.frameCount)
		}
	}
}

func FindTriples(grid [][]Square) (bool, [][]bool) {
	// create a local mask to mark all square that are in triples
	mask := make([][]bool, numRows)
	for i := range mask {
		mask[i] = make([]bool, numColumns)
	}

	found := false
	// find all horizontal triples
	for y, row := range grid[:len(grid)] {
		for x := range grid[:len(row)-2] {
			if grid[y][x].color >= 0 { // if is a color
				if grid[y][x].color == grid[y][x+1].color && grid[y][x].color == grid[y][x+2].color {
					mask[y][x], mask[y][x+1], mask[y][x+2] = true, true, true
					found = true
				}
			}
		}
	}

	// find all vertical triples
	for y, row := range grid[:len(grid)-2] {
		for x := range grid[:len(row)] {
			if grid[y][x].color >= 0 { // if is a color
				if grid[y][x].color == grid[y+1][x].color && grid[y][x].color == grid[y+2][x].color {
					mask[y][x], mask[y+1][x], mask[y+2][x] = true, true, true
					found = true
				}
			}
		}
	}
	return found, mask
}

func UpdateTriples(g *Game) {
	found, mask := FindTriples(g.grid)

	if found {
		// now that we have completed detecting all triples we can update the game state
		for y, row := range g.grid {
			for x := range row {
				if mask[y][x] {
					g.grid[y][x].color = -1
					g.triplesMask[y][x] = true
				}
			}
		}
		if gameIsWon(g) {
			g.player, _ = PlaySound(winOgg)
		} else {
			g.player, _ = PlaySound(tripleOgg)
		}
	}

	FillEmpties(g)
}

func FillEmpties(g *Game) {
	// find empty square and move squares from above down to fill
	for x := range numColumns {
		for y := range numRows {
			y = numRows - 1 - y // work from bottom up
			if g.grid[y][x].color == -1 {
				above := findSquareAbove(g, Point{x, y})
				if above.y >= 0 {
					g.grid[y][x].color = g.grid[above.y][above.x].color
					g.grid[above.y][above.x].color = -1
					g.grid[y][x].AddMover(g.frameCount, dropDuration,
						g.grid[above.y][above.x].point,
						g.grid[y][x].point)
				}
			}
		}
	}

	// fill empties at the top of the grid with newly generated colors
	for x := range numColumns {
		for y := range numRows {
			if g.grid[y][x].color == -1 {
				g.grid[y][x].color = rng.Intn(g.numColors)

				// there's a bit of a kludge here. The call to offsetPoint should be equal to the height
				// of the stack squares being removed, but don't calculate that height and just pass
				// cellsize * -1.
				g.grid[y][x].AddMover(g.frameCount, dropDuration,
					offsetPoint(g.grid[y][x].point, Point{0, cellSize * -1}),
					g.grid[y][x].point)
			}
		}
	}
}

func fillRandom(g *Game) {
	for y, row := range g.grid {
		for x := range row {
			g.grid[y][x].color = rng.Intn(g.numColors)
		}
	}
}

func findSquareAbove(g *Game, p Point) Point {
	for y := range p.y {
		y = p.y - 1 - y
		for g.grid[y][p.x].color != -1 {
			return Point{p.x, y}
		}
	}
	return Point{-1, -1} // did not find a square with color
}

func gameIsWon(g *Game) bool {
	for y, row := range g.grid {
		for x := range row {
			if !g.triplesMask[y][x] {
				return false
			}
		}
	}
	return true
}
