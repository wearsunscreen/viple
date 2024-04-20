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
	numGemColumns = 5
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

/* ==================================
/* ==================================
	LevelGemsVisualMode methods
/* ==================================
/* ==================================
*/

func (l *LevelGemsVisualMode) Draw(screen *ebiten.Image, frameCount int) {

	screen.Fill(mediumCoal)

	// draw background of triples
	for y, row := range l.gemGrid {
		for x := range row {
			if l.triplesMask[y][x] {
				l.gemGrid[y][x].drawBackground(screen, darkButter)
			}
		}
	}

	// draw cursor
	cursorColors := [2]color.Color{color.White, color.Black}
	blink := frameCount / blinkInverval % 2
	if l.mode == VisualMode {
		// we are in swap mode, faster blink, brighter colors
		blink = frameCount / (blinkInverval / 2) % 2
		cursorColors = [2]color.Color{brightRed, lightButter}

		// draw visualmode cursor
		cursorStart, cursorEnd := highLow(l.cursorGem, l.swapGem)
		startX := cursorStart.x
		for y := cursorStart.y; y <= cursorEnd.y; y++ {
			for x := startX; x < numGemColumns; x++ {
				l.gemGrid[y][x].drawBackground(screen, darkGreen)
				if x == cursorEnd.x && y == cursorEnd.y {
					l.gemGrid[y][x].drawBackground(screen, darkGreen)
					break
				}
				// start next line at left edge
				startX = 0
			}
		}
		l.gemGrid[l.swapGem.y][l.swapGem.x].drawBackground(screen, cursorColors[blink])
	} else {
		// draw cursor
		l.gemGrid[l.cursorGem.y][l.cursorGem.x].drawBackground(screen, cursorColors[blink])
	}

	// draw gems
	for y, row := range l.gemGrid {
		for x := range row {
			l.gemGrid[y][x].drawGem(screen, l.gemImages[l.gemGrid[y][x].color], frameCount)
		}
	}
}

func (l *LevelGemsVisualMode) Initialize(id LevelID) {
	l.numGems = 5
	l.cursorGem = Point{numGemColumns / 2, gemRows / 2}
	l.swapGem = Point{-1, -1}
	l.gemGrid = make([][]Square, gemRows)
	for y := range l.gemGrid {
		l.gemGrid[y] = make([]Square, numGemColumns)
	}

	for y, row := range l.gemGrid {
		for x := range row {
			l.gemGrid[y][x].point = Point{x, y}
		}
	}

	l.triplesMask = make([][]bool, gemRows)
	for i := range l.triplesMask {
		l.triplesMask[i] = make([]bool, numGemColumns)
	}

	l.numGems = 5
	l.mode = NormalMode
	fillRandom(l)

	l.loadGems()
}

func (l *LevelGemsVisualMode) Update(frameCount int) (bool, error) {
	// clear movers if expired
	for y, row := range l.gemGrid {
		for x := range row {
			if l.gemGrid[y][x].mover != nil {
				if l.gemGrid[y][x].mover.endFrame < frameCount {
					l.gemGrid[y][x].mover = nil
					updateTriples(l, frameCount)
				}
			}
		}
	}

	keys = inpututil.AppendPressedKeys(keys[:0])
	for _, key := range keys {
		if inpututil.IsKeyJustPressed(key) {
			switch l.mode {
			case NormalMode:
				handleKeyCommand(l, key)
			case VisualMode:
				handleKeyVisual(l, key, frameCount)
			}
		}
		// cheat code to fill
		if isCheatKeyPressed() {
			return true, nil
		}
	}

	if l.gameIsWon() {
		return true, nil
	}
	return false, nil
}

/* ==================================
/* ==================================
	Square methods
/* ==================================
/* ==================================
*/

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

func (square *Square) GetZ() int {
	return square.z
}

func (square *Square) SetZ(z int) {
	square.z = z
}

/* ==================================
/* ==================================
	Other functions
/* ==================================
/* ==================================
*/

func applyMover(mover *Mover, op *ebiten.DrawImageOptions, frameCount int) {
	completionRatio := 1 - float64(mover.endFrame-frameCount)/float64(mover.endFrame-mover.startFrame)
	startPosition := squareToScreenPoint(mover.startPoint)
	endPosition := squareToScreenPoint(mover.endPoint)
	op.GeoM.Translate(
		float64(startPosition.x)+(completionRatio*float64(endPosition.x-startPosition.x)),
		float64(startPosition.y)+(completionRatio*float64(endPosition.y-startPosition.y)))
}

func fillEmpties(l *LevelGemsVisualMode, frameCount int) {
	// find empty square and move squares from above down to fill
	for x := range numGemColumns {
		for y := range gemRows {
			y = gemRows - 1 - y // work from bottom up
			if l.gemGrid[y][x].color == -1 {
				above := findSquareAbove(l, Point{x, y})
				if above.y >= 0 {
					l.gemGrid[y][x].color = l.gemGrid[above.y][above.x].color
					l.gemGrid[above.y][above.x].color = -1
					l.gemGrid[y][x].AddMover(frameCount, dropDuration,
						l.gemGrid[above.y][above.x].point,
						l.gemGrid[y][x].point)
				}
			}
		}
	}

	// fill empties at the top of the gemGrid with newly generated colors
	for x := range numGemColumns {
		for y := range gemRows {
			if l.gemGrid[y][x].color == -1 {
				l.gemGrid[y][x].color = rng.Intn(l.numGems)

				// there's a bit of a kludge here. The call to offsetPoint should be equal to the height
				// of the stack squares being removed, but don't calculate that height and just pass
				// cellsize * -1.
				l.gemGrid[y][x].AddMover(frameCount, dropDuration,
					offsetPoint(l.gemGrid[y][x].point, Point{0, gemCellSize * -1}),
					l.gemGrid[y][x].point)
			}
		}
	}
}

func fillRandom(l *LevelGemsVisualMode) {
	for y, row := range l.gemGrid {
		for x := range row {
			l.gemGrid[y][x].color = rng.Intn(l.numGems)
		}
	}
}

func findSquareAbove(l *LevelGemsVisualMode, p Point) Point {
	for y := range p.y {
		y = p.y - 1 - y
		for l.gemGrid[y][p.x].color != -1 {
			return Point{p.x, y}
		}
	}
	return Point{-1, -1} // did not find a square with color
}

func findTriples(gemGrid [][]Square) (bool, [][]bool) {
	// create a local mask to mark all square that are in triples
	mask := make([][]bool, gemRows)
	for i := range mask {
		mask[i] = make([]bool, numGemColumns)
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

func handleKeyCommand(l *LevelGemsVisualMode, key ebiten.Key) {
	switch key {
	case ebiten.KeyH:
		l.cursorGem.x = max(l.cursorGem.x-1, 0)
	case ebiten.KeyL:
		l.cursorGem.x = min(l.cursorGem.x+1, numGemColumns-1)
	case ebiten.KeyK:
		l.cursorGem.y = max(l.cursorGem.y-1, 0)
	case ebiten.KeyJ:
		l.cursorGem.y = min(l.cursorGem.y+1, gemRows-1)
	case ebiten.KeyV:
		// entering VisualMode (where we do swaps)
		l.swapGem = l.cursorGem
		l.mode = VisualMode
	case ebiten.KeySemicolon:
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			l.keyInput = l.keyInput + ":"
		}
	case ebiten.KeyQ:
		fallthrough
	case ebiten.KeyX:
		// quit on ":q", ":x"
		l.keyInput = l.keyInput + key.String()
		if len(l.keyInput) > 1 {
			if l.keyInput[len(l.keyInput)-2:] == ":Q" ||
				l.keyInput[len(l.keyInput)-2:] == ":X" {
				os.Exit(0)
			}
		}
	}
}

func handleKeyVisual(l *LevelGemsVisualMode, key ebiten.Key, frameCount int) {
	switch key {
	case ebiten.KeyH:
		l.swapGem.x = max(l.swapGem.x-1, 0)
	case ebiten.KeyL:
		l.swapGem.x = min(l.swapGem.x+1, numGemColumns-1)
	case ebiten.KeyK:
		l.swapGem.y = max(l.swapGem.y-1, 0)
	case ebiten.KeyJ:
		l.swapGem.y = min(l.swapGem.y+1, gemRows-1)
	case ebiten.KeyV:
		PlaySound(failOgg)
	case ebiten.KeyEscape:
		// exit visual mode without swapping
		l.mode = NormalMode
		l.cursorGem = l.swapGem
		l.swapGem = Point{-1, -1}
	case ebiten.KeyY:
		// attempt swap
		if l.swapGem.x != -1 && l.swapGem != l.cursorGem {
			if result := swapSquares(l, frameCount); result {
				// swap successful
				// exiting visual mode
				l.mode = NormalMode
				l.cursorGem = l.swapGem
				l.swapGem = Point{-1, -1}
			} else {
				PlaySound(failOgg)
			}
		}
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

func (l *LevelGemsVisualMode) gameIsWon() bool {
	for y, row := range l.gemGrid {
		for x := range row {
			if !l.triplesMask[y][x] {
				return false
			}
		}
	}
	return true
}

func (l *LevelGemsVisualMode) IntroText() string {
	return `In the first level you will 
	learn to move left and right 
	by pressing H and K keys.`
}

func (l *LevelGemsVisualMode) TitleText() string {
	return `Welcome to Viple
	VI Play to Learn.`
}

func (l *LevelGemsVisualMode) loadGems() {
	if len(l.gemImages) == 0 {
		l.gemImages = make([]*ebiten.Image, l.numGems)
		for i := range l.numGems {
			image := loadImage("resources/Gem " + strconv.Itoa(i+1) + ".png")
			l.gemImages[i] = image
		}
	}
}

// convert the x,y of the square into screen coordinates
func squareToScreenPoint(squareXY Point) Point {
	// get leftmost x
	widthOfGrid := gemCellSize * numGemColumns
	xMargin := (screenWidth - widthOfGrid) / 2
	// get top y
	heightOfGrid := gemCellSize * gemRows
	yMargin := (screenHeight - heightOfGrid) / 2
	return Point{
		gemCellSize*squareXY.x + xMargin,
		gemCellSize*squareXY.y + yMargin,
	}
}

// Exchange positions of two neighboring squares, return false if unable to exchange.
// The exchange fails if swap point and cursor point are the same (this can happen when
// player attempts to move off the gemGrid). The exchange fails if both points have the
// same value.
func swapSquares(l *LevelGemsVisualMode, frameCount int) bool {
	if l.swapGem == l.cursorGem {
		return false
	}
	/*
		if l.gemGrid[l.swapGem.y][l.swapGem.x].point == l.gemGrid[l.cursorGem.y][l.cursorGem.x].point {
			return false
		}
	*/

	// swap colors
	fromSquare := &l.gemGrid[l.swapGem.y][l.swapGem.x]
	toSquare := &l.gemGrid[l.cursorGem.y][l.cursorGem.x]
	temp := fromSquare.color
	fromSquare.color = toSquare.color
	toSquare.color = temp

	// check if the swap will create a triple
	makesATriple, _ := findTriples(l.gemGrid)

	if makesATriple {
		toSquare.AddMover(frameCount, 60, fromSquare.point, toSquare.point)
		fromSquare.AddMover(frameCount, 60, toSquare.point, fromSquare.point)
	} else {
		// restore original colors
		temp := fromSquare.color
		fromSquare.color = toSquare.color
		toSquare.color = temp
		l.swapGem = l.cursorGem // return the cursor to the original location

		// tell the user he made an invalid move
		return false
	}
	return true
}

func updateTriples(l *LevelGemsVisualMode, frameCount int) {
	found, mask := findTriples(l.gemGrid)

	if found {
		// now that we have completed detecting all triples we can update the game state
		for y, row := range l.gemGrid {
			for x := range row {
				if mask[y][x] {
					l.gemGrid[y][x].color = -1
					l.triplesMask[y][x] = true
				}
			}
		}
		if l.gameIsWon() {
			PlaySound(winOgg)
		} else {
			PlaySound(tripleOgg)
		}
	}

	fillEmpties(l, frameCount)
}
