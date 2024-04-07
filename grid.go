package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawGrid(screen *ebiten.Image, g *Game) {

	// draw background of triples
	for y, row := range g.l3.gemGrid {
		for x := range row {
			if g.l3.triplesMask[y][x] {
				g.l3.gemGrid[y][x].DrawBackground(screen, darkButter)
			}
		}
	}

	// draw cursor
	cursorColors := [2]color.Color{color.White, color.Black}
	blink := g.frameCount / blinkInverval % 2
	if g.l3.mode == InsertMode {
		// we are in swap mode, faster blink, brighter colors
		blink = g.frameCount / (blinkInverval / 2) % 2
		cursorColors = [2]color.Color{brightRed, lightButter}
	}
	g.l3.gemGrid[g.l3.cursorGem.y][g.l3.cursorGem.x].DrawBackground(screen, cursorColors[blink])

	// draw gems
	for y, row := range g.l3.gemGrid {
		for x, _ := range row {
			g.l3.gemGrid[y][x].DrawGem(screen, g.l3.gemImages[g.l3.gemGrid[y][x].color], g.frameCount)
		}
	}
}

func FindTriples(gemGrid [][]Square) (bool, [][]bool) {
	// create a local mask to mark all square that are in triples
	mask := make([][]bool, gemRows)
	for i := range mask {
		mask[i] = make([]bool, numColumns)
	}

	found := false
	// find all horizontal triples
	for y, row := range gemGrid[:len(gemGrid)] {
		for x := range gemGrid[:len(row)-2] {
			if gemGrid[y][x].color >= 0 { // if is a color
				if gemGrid[y][x].color == gemGrid[y][x+1].color && gemGrid[y][x].color == gemGrid[y][x+2].color {
					mask[y][x], mask[y][x+1], mask[y][x+2] = true, true, true
					found = true
				}
			}
		}
	}

	// find all vertical triples
	for y, row := range gemGrid[:len(gemGrid)-2] {
		for x := range gemGrid[:len(row)] {
			if gemGrid[y][x].color >= 0 { // if is a color
				if gemGrid[y][x].color == gemGrid[y+1][x].color && gemGrid[y][x].color == gemGrid[y+2][x].color {
					mask[y][x], mask[y+1][x], mask[y+2][x] = true, true, true
					found = true
				}
			}
		}
	}
	return found, mask
}

func UpdateTriples(g *Game) {
	found, mask := FindTriples(g.l3.gemGrid)

	if found {
		// now that we have completed detecting all triples we can update the game state
		for y, row := range g.l3.gemGrid {
			for x := range row {
				if mask[y][x] {
					g.l3.gemGrid[y][x].color = -1
					g.l3.triplesMask[y][x] = true
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
		for y := range gemRows {
			y = gemRows - 1 - y // work from bottom up
			if g.l3.gemGrid[y][x].color == -1 {
				above := findSquareAbove(g, Point{x, y})
				if above.y >= 0 {
					g.l3.gemGrid[y][x].color = g.l3.gemGrid[above.y][above.x].color
					g.l3.gemGrid[above.y][above.x].color = -1
					g.l3.gemGrid[y][x].AddMover(g.frameCount, dropDuration,
						g.l3.gemGrid[above.y][above.x].point,
						g.l3.gemGrid[y][x].point)
				}
			}
		}
	}

	// fill empties at the top of the gemGrid with newly generated colors
	for x := range numColumns {
		for y := range gemRows {
			if g.l3.gemGrid[y][x].color == -1 {
				g.l3.gemGrid[y][x].color = rng.Intn(g.l3.numGems)

				// there's a bit of a kludge here. The call to offsetPoint should be equal to the height
				// of the stack squares being removed, but don't calculate that height and just pass
				// cellsize * -1.
				g.l3.gemGrid[y][x].AddMover(g.frameCount, dropDuration,
					offsetPoint(g.l3.gemGrid[y][x].point, Point{0, gemCellSize * -1}),
					g.l3.gemGrid[y][x].point)
			}
		}
	}
}

func fillRandom(g *Game) {
	for y, row := range g.l3.gemGrid {
		for x := range row {
			g.l3.gemGrid[y][x].color = rng.Intn(g.l3.numGems)
		}
	}
}

func findSquareAbove(g *Game, p Point) Point {
	for y := range p.y {
		y = p.y - 1 - y
		for g.l3.gemGrid[y][p.x].color != -1 {
			return Point{p.x, y}
		}
	}
	return Point{-1, -1} // did not find a square with color
}

func gameIsWon(g *Game) bool {
	for y, row := range g.l3.gemGrid {
		for x := range row {
			if !g.l3.triplesMask[y][x] {
				return false
			}
		}
	}
	return true
}
