package main

import (
	"image/color"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	blinkInverval = 60 / 2
	gemCellSize   = 50
	dropDuration  = 60
	gemScale      = float64(gemCellSize-4) / float64(gemWidth)
	gemWidth      = 100
	gemRows       = 11
	gemColumns    = 5
	swapDuration  = 40
)

type Level3Data struct {
	cursorGem   Point
	gemGrid     [][]Square
	gemImages   []*ebiten.Image
	mode        Mode
	numGems     int
	swapGem     Point
	triplesMask [][]bool
}

type Square struct {
	color int
	mover *Mover
	point Point
	z     int
}

func (square *Square) AddMover(startFrame int, duration int, from Point, to Point) {
	// add animation
	mover := new(Mover)

	mover.startFrame = startFrame
	mover.endFrame = startFrame + duration
	mover.startPoint = from
	mover.endPoint = to

	square.mover = mover
}

func applyMover(mover *Mover, op *ebiten.DrawImageOptions, frameCount int) {
	completionRatio := 1 - float64(mover.endFrame-frameCount)/float64(mover.endFrame-mover.startFrame)
	startPosition := squareToScreenPoint(mover.startPoint)
	endPosition := squareToScreenPoint(mover.endPoint)
	op.GeoM.Translate(
		float64(startPosition.x)+(completionRatio*float64(endPosition.x-startPosition.x)),
		float64(startPosition.y)+(completionRatio*float64(endPosition.y-startPosition.y)))
}

func (square *Square) drawBackground(screen *ebiten.Image, color color.Color) {
	pos := squareToScreenPoint(square.point)
	vector.DrawFilledRect(screen, float32(pos.x), float32(pos.y), gemCellSize-4, gemCellSize-4, color, false)
}

func (square *Square) drawGem(screen *ebiten.Image, gemImage *ebiten.Image, frameCount int) {
	if square.color >= 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(gemScale, gemScale)
		if square.mover != nil {
			applyMover(square.mover, op, frameCount)
		} else {
			pos := squareToScreenPoint(square.point)
			op.GeoM.Translate(float64(pos.x), float64(pos.y))
		}
		screen.DrawImage(gemImage, op)
	}
}

func (level3 *Level3Data) Draw(screen *ebiten.Image, g *Game) {

	// draw background of triples
	for y, row := range level3.gemGrid {
		for x := range row {
			if level3.triplesMask[y][x] {
				level3.gemGrid[y][x].drawBackground(screen, darkButter)
			}
		}
	}

	// draw cursor
	cursorColors := [2]color.Color{color.White, color.Black}
	blink := g.frameCount / blinkInverval % 2
	if g.level3.mode == InsertMode {
		// we are in swap mode, faster blink, brighter colors
		blink = g.frameCount / (blinkInverval / 2) % 2
		cursorColors = [2]color.Color{brightRed, lightButter}
	}
	g.level3.gemGrid[g.level3.cursorGem.y][g.level3.cursorGem.x].drawBackground(screen, cursorColors[blink])

	// draw gems
	for y, row := range g.level3.gemGrid {
		for x := range row {
			g.level3.gemGrid[y][x].drawGem(screen, g.level3.gemImages[g.level3.gemGrid[y][x].color], g.frameCount)
		}
	}
}

func findTriples(gemGrid [][]Square) (bool, [][]bool) {
	// create a local mask to mark all square that are in triples
	mask := make([][]bool, gemRows)
	for i := range mask {
		mask[i] = make([]bool, gemColumns)
	}

	found := false
	// find all horizontal triples
	for y, row := range gemGrid[:] {
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

func (square *Square) GetZ() int {
	return square.z
}

func handleKeyCommand(g *Game, key ebiten.Key) {
	switch key {
	case ebiten.KeyH:
		g.level3.cursorGem.x = max(g.level3.cursorGem.x-1, 0)
	case ebiten.KeyL:
		g.level3.cursorGem.x = min(g.level3.cursorGem.x+1, gemColumns-1)
	case ebiten.KeyK:
		g.level3.cursorGem.y = max(g.level3.cursorGem.y-1, 0)
	case ebiten.KeyJ:
		g.level3.cursorGem.y = min(g.level3.cursorGem.y+1, gemRows-1)
	case ebiten.KeyI:
		// entering InsertMode (where we do swaps)
		g.level3.swapGem = g.level3.cursorGem
		g.level3.mode = InsertMode
	case ebiten.KeySemicolon:
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			g.keyInput = g.keyInput + ":"
		}
	case ebiten.KeyQ:
		fallthrough
	case ebiten.KeyX:
		// quit on ":q", ":x"
		g.keyInput = g.keyInput + key.String()
		if len(g.keyInput) > 1 {
			if g.keyInput[len(g.keyInput)-2:] == ":Q" ||
				g.keyInput[len(g.keyInput)-2:] == ":X" {
				os.Exit(0)
			}
		}
	}
}

func handleKeyInsert(g *Game, key ebiten.Key) {
	switch key {
	case ebiten.KeyH:
		g.level3.swapGem.x = max(g.level3.swapGem.x-1, 0)
	case ebiten.KeyL:
		g.level3.swapGem.x = min(g.level3.swapGem.x+1, gemColumns-1)
	case ebiten.KeyK:
		g.level3.swapGem.y = max(g.level3.swapGem.y-1, 0)
	case ebiten.KeyJ:
		g.level3.swapGem.y = min(g.level3.swapGem.y+1, gemRows-1)
	case ebiten.KeyI:
		PlaySound(failOgg)
	case ebiten.KeyEscape:
		g.level3.mode = NormalMode
		g.level3.swapGem = Point{-1, -1}
	}
	if g.level3.swapGem.x != -1 && g.level3.swapGem != g.level3.cursorGem {
		if result := swapSquares(g); !result {
			PlaySound(failOgg)
		}
		g.level3.cursorGem = g.level3.swapGem
	}

}

func initLevel3(g *Game) {
	g.level3.numGems = 5
	g.level3.cursorGem = Point{gemColumns / 2, gemRows / 2}
	g.level3.swapGem = Point{-1, -1}
	g.level3.gemGrid = make([][]Square, gemRows)
	for y := range g.level3.gemGrid {
		g.level3.gemGrid[y] = make([]Square, gemColumns)
	}

	for y, row := range g.level3.gemGrid {
		for x := range row {
			g.level3.gemGrid[y][x].point = Point{x, y}
		}
	}

	g.level3.triplesMask = make([][]bool, gemRows)
	for i := range g.level3.triplesMask {
		g.level3.triplesMask[i] = make([]bool, gemColumns)
	}

	g.level3.numGems = 5
	g.level3.mode = NormalMode
	fillRandom(g)

	loadGems(g)
}

func loadGems(g *Game) {
	g.level3.gemImages = make([]*ebiten.Image, g.level3.numGems)
	for i := range g.level3.numGems {
		image := loadImage("resources/Gem " + strconv.Itoa(i+1) + ".png")
		g.level3.gemImages[i] = image
	}
}

func (square *Square) SetZ(z int) {
	square.z = z
}

// convert the x,y of the square into screen coordinates
func squareToScreenPoint(squareXY Point) Point {
	// get leftmost x
	widthOfGrid := gemCellSize * gemColumns
	xMargin := (screenWidth - widthOfGrid) / 2
	// get top y
	heightOfGrid := gemCellSize * gemRows
	yMargin := (screenHeight - heightOfGrid) / 2
	return Point{
		gemCellSize*squareXY.x + xMargin,
		gemCellSize*squareXY.y + yMargin,
	}
}

func (level3 *Level3Data) Update(g *Game) (bool, error) {
	// clear movers if expired
	for y, row := range g.level3.gemGrid {
		for x := range row {
			if g.level3.gemGrid[y][x].mover != nil {
				if g.level3.gemGrid[y][x].mover.endFrame < g.frameCount {
					g.level3.gemGrid[y][x].mover = nil
					updateTriples(g)
				}
			}
		}
	}

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	for _, key := range g.keys {
		if inpututil.IsKeyJustPressed(key) {
			switch g.level3.mode {
			case NormalMode:
				handleKeyCommand(g, key)
			case InsertMode:
				handleKeyInsert(g, key)
			}
		}
	}

	if gameIsWon(g) {
		return true, nil
	}
	return false, nil
}

func updateTriples(g *Game) {
	found, mask := findTriples(g.level3.gemGrid)

	if found {
		// now that we have completed detecting all triples we can update the game state
		for y, row := range g.level3.gemGrid {
			for x := range row {
				if mask[y][x] {
					g.level3.gemGrid[y][x].color = -1
					g.level3.triplesMask[y][x] = true
				}
			}
		}
		if gameIsWon(g) {
			g.player, _ = PlaySound(winOgg)
		} else {
			g.player, _ = PlaySound(tripleOgg)
		}
	}

	fillEmpties(g)
}

func fillEmpties(g *Game) {
	// find empty square and move squares from above down to fill
	for x := range gemColumns {
		for y := range gemRows {
			y = gemRows - 1 - y // work from bottom up
			if g.level3.gemGrid[y][x].color == -1 {
				above := findSquareAbove(g, Point{x, y})
				if above.y >= 0 {
					g.level3.gemGrid[y][x].color = g.level3.gemGrid[above.y][above.x].color
					g.level3.gemGrid[above.y][above.x].color = -1
					g.level3.gemGrid[y][x].AddMover(g.frameCount, dropDuration,
						g.level3.gemGrid[above.y][above.x].point,
						g.level3.gemGrid[y][x].point)
				}
			}
		}
	}

	// fill empties at the top of the gemGrid with newly generated colors
	for x := range gemColumns {
		for y := range gemRows {
			if g.level3.gemGrid[y][x].color == -1 {
				g.level3.gemGrid[y][x].color = rng.Intn(g.level3.numGems)

				// there's a bit of a kludge here. The call to offsetPoint should be equal to the height
				// of the stack squares being removed, but don't calculate that height and just pass
				// cellsize * -1.
				g.level3.gemGrid[y][x].AddMover(g.frameCount, dropDuration,
					offsetPoint(g.level3.gemGrid[y][x].point, Point{0, gemCellSize * -1}),
					g.level3.gemGrid[y][x].point)
			}
		}
	}
}

func fillRandom(g *Game) {
	for y, row := range g.level3.gemGrid {
		for x := range row {
			g.level3.gemGrid[y][x].color = rng.Intn(g.level3.numGems)
		}
	}
}

func findSquareAbove(g *Game, p Point) Point {
	for y := range p.y {
		y = p.y - 1 - y
		for g.level3.gemGrid[y][p.x].color != -1 {
			return Point{p.x, y}
		}
	}
	return Point{-1, -1} // did not find a square with color
}

func gameIsWon(g *Game) bool {
	for y, row := range g.level3.gemGrid {
		for x := range row {
			if !g.level3.triplesMask[y][x] {
				return false
			}
		}
	}
	return true
}

// Exchange positions of two neighboring squares, return false if unable to exchange.
// The exchange fails if swap point and cursor point are the same (this can happen when
// player attempts to move off the gemGrid). The exchange fails if both points have the
// same value.
func swapSquares(g *Game) bool {
	if g.level3.swapGem == g.level3.cursorGem {
		return false
	}
	if g.level3.gemGrid[g.level3.swapGem.y][g.level3.swapGem.x].point == g.level3.gemGrid[g.level3.cursorGem.y][g.level3.cursorGem.x].point {
		return false
	}

	// swap colors
	fromSquare := &g.level3.gemGrid[g.level3.swapGem.y][g.level3.swapGem.x]
	toSquare := &g.level3.gemGrid[g.level3.cursorGem.y][g.level3.cursorGem.x]
	temp := fromSquare.color
	fromSquare.color = toSquare.color
	toSquare.color = temp

	// check if the swap will create a triple
	makesATriple, _ := findTriples(g.level3.gemGrid)

	if makesATriple {
		toSquare.AddMover(g.frameCount, 60, fromSquare.point, toSquare.point)
		fromSquare.AddMover(g.frameCount, 60, toSquare.point, fromSquare.point)
	} else {
		// restore original colors
		temp := fromSquare.color
		fromSquare.color = toSquare.color
		toSquare.color = temp
		g.level3.swapGem = g.level3.cursorGem // return the cursor to the original location

		// tell the user he made an invalid move
		return false
	}
	return true
}
