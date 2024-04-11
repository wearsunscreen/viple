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
	NormalMode = iota
	VisualMode
	InsertMode
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

type LevelGemsVisualMode struct {
	cursorGem   Point
	gemGrid     [][]Square
	gemImages   []*ebiten.Image
	keyInput    string
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

// return the highest and lowest points in that order where highest is closest to the top left of the screen
func highLow(p1, p2 Point) (Point, Point) {
	if p1.y == p2.y {
		if p1.x >= p2.x {
			return p2, p1
		} else {
			return p1, p2
		}
	} else {
		if p1.y > p2.y {
			return p2, p1
		} else {
			return p1, p2
		}
	}
}

func (level *LevelGemsVisualMode) Draw(screen *ebiten.Image, frameCount int) {

	screen.Fill(mediumCoal)

	// draw background of triples
	for y, row := range level.gemGrid {
		for x := range row {
			if level.triplesMask[y][x] {
				level.gemGrid[y][x].drawBackground(screen, darkButter)
			}
		}
	}

	// draw cursor
	cursorColors := [2]color.Color{color.White, color.Black}
	blink := frameCount / blinkInverval % 2
	if level.mode == VisualMode {
		// we are in swap mode, faster blink, brighter colors
		blink = frameCount / (blinkInverval / 2) % 2
		cursorColors = [2]color.Color{brightRed, lightButter}

		// draw visualmode cursor
		cursorStart, cursorEnd := highLow(level.cursorGem, level.swapGem)
		startX := cursorStart.x
		for y := cursorStart.y; y <= cursorEnd.y; y++ {
			for x := startX; x < numBrickCols; x++ {
				level.gemGrid[y][x].drawBackground(screen, darkGreen)
				if x == cursorEnd.x && y == cursorEnd.y {
					level.gemGrid[y][x].drawBackground(screen, darkGreen)
					break
				}
				// start next line at left edge
				startX = 0
			}
		}
		level.gemGrid[level.swapGem.y][level.swapGem.x].drawBackground(screen, cursorColors[blink])
	} else {
		// draw cursor
		level.gemGrid[level.cursorGem.y][level.cursorGem.x].drawBackground(screen, cursorColors[blink])
	}

	// draw gems
	for y, row := range level.gemGrid {
		for x := range row {
			level.gemGrid[y][x].drawGem(screen, level.gemImages[level.gemGrid[y][x].color], frameCount)
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

func handleKeyCommand(level *LevelGemsVisualMode, key ebiten.Key) {
	switch key {
	case ebiten.KeyH:
		level.cursorGem.x = max(level.cursorGem.x-1, 0)
	case ebiten.KeyL:
		level.cursorGem.x = min(level.cursorGem.x+1, gemColumns-1)
	case ebiten.KeyK:
		level.cursorGem.y = max(level.cursorGem.y-1, 0)
	case ebiten.KeyJ:
		level.cursorGem.y = min(level.cursorGem.y+1, gemRows-1)
	case ebiten.KeyV:
		// entering VisualMode (where we do swaps)
		level.swapGem = level.cursorGem
		level.mode = VisualMode
	case ebiten.KeySemicolon:
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			level.keyInput = level.keyInput + ":"
		}
	case ebiten.KeyQ:
		fallthrough
	case ebiten.KeyX:
		// quit on ":q", ":x"
		level.keyInput = level.keyInput + key.String()
		if len(level.keyInput) > 1 {
			if level.keyInput[len(level.keyInput)-2:] == ":Q" ||
				level.keyInput[len(level.keyInput)-2:] == ":X" {
				os.Exit(0)
			}
		}
	}
}

func handleKeyVisual(level *LevelGemsVisualMode, key ebiten.Key, frameCount int) {
	switch key {
	case ebiten.KeyH:
		level.swapGem.x = max(level.swapGem.x-1, 0)
	case ebiten.KeyL:
		level.swapGem.x = min(level.swapGem.x+1, gemColumns-1)
	case ebiten.KeyK:
		level.swapGem.y = max(level.swapGem.y-1, 0)
	case ebiten.KeyJ:
		level.swapGem.y = min(level.swapGem.y+1, gemRows-1)
	case ebiten.KeyV:
		PlaySound(failOgg)
	case ebiten.KeyEscape:
		// exit visual mode without swapping
		level.mode = NormalMode
		level.cursorGem = level.swapGem
		level.swapGem = Point{-1, -1}
	case ebiten.KeyY:
		// attempt swap
		if level.swapGem.x != -1 && level.swapGem != level.cursorGem {
			if result := swapSquares(level, frameCount); result {
				// swap successful
				// exiting visual mode
				level.mode = NormalMode
				level.cursorGem = level.swapGem
				level.swapGem = Point{-1, -1}
			} else {
				PlaySound(failOgg)
			}
		}
	}
}

func (level *LevelGemsVisualMode) Initialize() {
	level.numGems = 5
	level.cursorGem = Point{gemColumns / 2, gemRows / 2}
	level.swapGem = Point{-1, -1}
	level.gemGrid = make([][]Square, gemRows)
	for y := range level.gemGrid {
		level.gemGrid[y] = make([]Square, gemColumns)
	}

	for y, row := range level.gemGrid {
		for x := range row {
			level.gemGrid[y][x].point = Point{x, y}
		}
	}

	level.triplesMask = make([][]bool, gemRows)
	for i := range level.triplesMask {
		level.triplesMask[i] = make([]bool, gemColumns)
	}

	level.numGems = 5
	level.mode = NormalMode
	fillRandom(level)

	loadGems(level)
}

func loadGems(level *LevelGemsVisualMode) {
	level.gemImages = make([]*ebiten.Image, level.numGems)
	for i := range level.numGems {
		image := loadImage("resources/Gem " + strconv.Itoa(i+1) + ".png")
		level.gemImages[i] = image
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

func (level *LevelGemsVisualMode) Update(frameCount int) (bool, error) {
	// clear movers if expired
	for y, row := range level.gemGrid {
		for x := range row {
			if level.gemGrid[y][x].mover != nil {
				if level.gemGrid[y][x].mover.endFrame < frameCount {
					level.gemGrid[y][x].mover = nil
					updateTriples(level, frameCount)
				}
			}
		}
	}

	keys = inpututil.AppendPressedKeys(keys[:0])
	for _, key := range keys {
		if inpututil.IsKeyJustPressed(key) {
			switch level.mode {
			case NormalMode:
				handleKeyCommand(level, key)
			case VisualMode:
				handleKeyVisual(level, key, frameCount)
			}
		}
		// cheat code to fill
		if ebiten.IsKeyPressed(ebiten.KeyZ) {
			for y := range level.triplesMask {
				level.triplesMask[y] = make([]bool, numBrickCols)
				fillSlice(level.triplesMask[y], true)
			}
		}
	}

	if gameIsWon(level) {
		return true, nil
	}
	return false, nil
}

func updateTriples(level *LevelGemsVisualMode, frameCount int) {
	found, mask := findTriples(level.gemGrid)

	if found {
		// now that we have completed detecting all triples we can update the game state
		for y, row := range level.gemGrid {
			for x := range row {
				if mask[y][x] {
					level.gemGrid[y][x].color = -1
					level.triplesMask[y][x] = true
				}
			}
		}
		if gameIsWon(level) {
			PlaySound(winOgg)
		} else {
			PlaySound(tripleOgg)
		}
	}

	fillEmpties(level, frameCount)
}

func fillEmpties(level *LevelGemsVisualMode, frameCount int) {
	// find empty square and move squares from above down to fill
	for x := range gemColumns {
		for y := range gemRows {
			y = gemRows - 1 - y // work from bottom up
			if level.gemGrid[y][x].color == -1 {
				above := findSquareAbove(level, Point{x, y})
				if above.y >= 0 {
					level.gemGrid[y][x].color = level.gemGrid[above.y][above.x].color
					level.gemGrid[above.y][above.x].color = -1
					level.gemGrid[y][x].AddMover(frameCount, dropDuration,
						level.gemGrid[above.y][above.x].point,
						level.gemGrid[y][x].point)
				}
			}
		}
	}

	// fill empties at the top of the gemGrid with newly generated colors
	for x := range gemColumns {
		for y := range gemRows {
			if level.gemGrid[y][x].color == -1 {
				level.gemGrid[y][x].color = rng.Intn(level.numGems)

				// there's a bit of a kludge here. The call to offsetPoint should be equal to the height
				// of the stack squares being removed, but don't calculate that height and just pass
				// cellsize * -1.
				level.gemGrid[y][x].AddMover(frameCount, dropDuration,
					offsetPoint(level.gemGrid[y][x].point, Point{0, gemCellSize * -1}),
					level.gemGrid[y][x].point)
			}
		}
	}
}

func fillRandom(level *LevelGemsVisualMode) {
	for y, row := range level.gemGrid {
		for x := range row {
			level.gemGrid[y][x].color = rng.Intn(level.numGems)
		}
	}
}

func findSquareAbove(level *LevelGemsVisualMode, p Point) Point {
	for y := range p.y {
		y = p.y - 1 - y
		for level.gemGrid[y][p.x].color != -1 {
			return Point{p.x, y}
		}
	}
	return Point{-1, -1} // did not find a square with color
}

func gameIsWon(level *LevelGemsVisualMode) bool {
	for y, row := range level.gemGrid {
		for x := range row {
			if !level.triplesMask[y][x] {
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
func swapSquares(level *LevelGemsVisualMode, frameCount int) bool {
	if level.swapGem == level.cursorGem {
		return false
	}
	/*
		if level.gemGrid[level.swapGem.y][level.swapGem.x].point == level.gemGrid[level.cursorGem.y][level.cursorGem.x].point {
			return false
		}
	*/

	// swap colors
	fromSquare := &level.gemGrid[level.swapGem.y][level.swapGem.x]
	toSquare := &level.gemGrid[level.cursorGem.y][level.cursorGem.x]
	temp := fromSquare.color
	fromSquare.color = toSquare.color
	toSquare.color = temp

	// check if the swap will create a triple
	makesATriple, _ := findTriples(level.gemGrid)

	if makesATriple {
		toSquare.AddMover(frameCount, 60, fromSquare.point, toSquare.point)
		fromSquare.AddMover(frameCount, 60, toSquare.point, fromSquare.point)
	} else {
		// restore original colors
		temp := fromSquare.color
		fromSquare.color = toSquare.color
		toSquare.color = temp
		level.swapGem = level.cursorGem // return the cursor to the original location

		// tell the user he made an invalid move
		return false
	}
	return true
}
