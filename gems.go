package main

import (
	"image/color"
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
	blinkInverval = 60 / 3
	emptyGem      = -1
	dropDuration  = 60
	gemCellSize   = 50
	gemScale      = float64(gemCellSize-4) / float64(gemWidth)
	gemWidth      = 100
	numGemRows    = 11
	swapDuration  = 40
)

var (
	lightGold   = color.RGBA{0xff, 0xff, 0x80, 0xff}
	redCursor   = color.RGBA{0xfc, 0x10, 0x10, 0x80}
	whiteCursor = color.RGBA{0xfc, 0xfc, 0xfc, 0x80}
)

var numGemColumns int

type GemMover struct {
	startFrame int
	endFrame   int
	startGrid  Coord // grid coordinates
	endCoord   Coord // grid coordinates
}
type LevelGems struct {
	cursorGem   Coord
	gemGrid     Grid[Square]
	gemImages   []*ebiten.Image
	level       LevelID
	viMode      int
	numGems     int
	swapGem     Coord
	triplesMask Grid[bool]
}

type Square struct {
	gem    int
	mover  *GemMover
	coords Coord // position in the grid
}

func applyMover(mover *GemMover, op *ebiten.DrawImageOptions, frameCount int) {
	completionRatio := 1 - float64(mover.endFrame-frameCount)/float64(mover.endFrame-mover.startFrame)
	startPosition := squareToScreenPoint(mover.startGrid)
	endPosition := squareToScreenPoint(mover.endCoord)
	op.GeoM.Translate(
		float64(startPosition.x)+(completionRatio*float64(endPosition.x-startPosition.x)),
		float64(startPosition.y)+(completionRatio*float64(endPosition.y-startPosition.y)))
}

// fillFromBelow() generates new gems to fill in empty squares at the bottom of the gem grid
func fillFromBelow(gemGrid *Grid[Square], numGems, frameCount int, addMover bool) {
	// fill empties at the bottom of the gemGrid with newly generated gems
	for x := range numGemColumns {
		for y := range numGemRows {
			if gemGrid.Get(Coord{x, y}).gem == emptyGem {
				sqPtr := gemGrid.GetPtr(Coord{x, y})
				sqPtr.gem = rng.Intn(numGems)
				if addMover {
					sqPtr.addMover(frameCount, dropDuration,
						Coord{sqPtr.coords.x, numGemRows + 1},
						Coord{x, y})
				}
			}
		}
	}
}

// highLow() returns the highest and lowest points in that order where highest is closest to the top left of the screen
// disregard a point if it has negative values
func highLow(p1, p2 Coord) (Coord, Coord) {
	if p2.x < 0 || p2.y < 0 {
		return p1, p2
	}
	if p1.x < 0 || p1.y < 0 {
		return p2, p1
	}
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

// moveUpFromBelow() moves gems from below to fill in empty squares
func moveUpFromBelow(gemGrid *Grid[Square], frameCount int, addMover bool) {
	// for each column
	for x := range numGemColumns {
		// find empty square and move squares from below to fill
		for y := 0; y < numGemRows; y++ {
			pt := Coord{x, y}
			if gemGrid.Get(pt).gem == emptyGem {
				below := findSquareBelow(gemGrid, pt)
				if below.y >= 0 {
					sqPtr := gemGrid.GetPtr(pt)
					sqPtr.gem = gemGrid.Get(below).gem
					sqPtr = gemGrid.GetPtr(below)
					sqPtr.gem = emptyGem
					if addMover {
						sqPtr = gemGrid.GetPtr(pt)
						sqPtr.addMover(frameCount, dropDuration, below, pt)
					}
				}
			}
		}
	}
}

// newGridOfSquares() creates a grid of squares
func newGridOfSquares(width, height int) Grid[Square] {
	r := make([][]Square, height)
	for i := range r {
		r[i] = make([]Square, width)
		for j := range r[i] {
			r[i][j] = Square{coords: Coord{j, i}}
		}
	}
	return r
}

func setGem(grid *Grid[Square], pt Coord, gem int) {
	sq := grid.Get(pt)
	sq.gem = gem
	grid.Set(pt, sq)
}

// Set all gems in a selection to a value
func setSelection(gemGrid *Grid[Square], p1, p2 Coord, gem int) {
	cursorStart, cursorEnd := highLow(p1, p2)
	startX := cursorStart.x
	for y := cursorStart.y; y <= cursorEnd.y; y++ {
		for x := startX; x < numGemColumns; x++ {
			setGem(gemGrid, Coord{x, y}, gem)
			if x == cursorEnd.x && y == cursorEnd.y {
				break
			}
			// start next line at left edge
			startX = 0
		}
	}
}

// convert the x,y of the square into screen coordinates
func squareToScreenPoint(squareXY Coord) Coord {
	// get leftmost x
	widthOfGrid := gemCellSize * numGemColumns
	xMargin := (screenWidth - widthOfGrid) / 2
	// get top y
	heightOfGrid := gemCellSize * numGemRows
	yMargin := (screenHeight - heightOfGrid) / 2
	return Coord{
		gemCellSize*squareXY.x + xMargin,
		gemCellSize*squareXY.y + yMargin,
	}
}

func (l *LevelGems) Draw(screen *ebiten.Image, frameCount int) {

	screen.Fill(mediumCoal)

	// draw background of triples
	l.gemGrid.ForEach(func(p Coord, s Square) {
		if l.triplesMask.Get(p) {
			s.drawBackground(screen, lightGold)
		}
	})

	// draw cursor
	cursorColors := [2]color.Color{redCursor, whiteCursor}
	blink := frameCount / blinkInverval % 2

	switch l.level {
	case LevelIdGemsEnd:
		fallthrough
	case LevelIdGemsVM:
		// we are in swap mode, faster blink, brighter colors
		blink = frameCount / (blinkInverval / 2) % 2

		// draw visualmode cursor
		if l.viMode == VisualMode {
			cursorStart, cursorEnd := highLow(l.cursorGem, l.swapGem)
			startX := cursorStart.x
			for y := cursorStart.y; y <= cursorEnd.y; y++ {
				for x := startX; x < numGemColumns; x++ {
					s := l.gemGrid.Get(Coord{x, y})
					s.drawBackground(screen, darkGreen)
					if x == cursorEnd.x && y == cursorEnd.y {
						break
					}
					// start next line at left edge
					startX = 0
				}
			}
			s := l.gemGrid.Get(l.swapGem)
			s.drawBackground(screen, cursorColors[blink])
		} else {
			s := l.gemGrid.Get(l.cursorGem)
			s.drawBackground(screen, cursorColors[blink])
		}
	case LevelIdGemsDD:
		l.gemGrid.ForEach(func(p Coord, s Square) {
			if p.y == l.cursorGem.y {
				s.drawBackground(screen, cursorColors[blink])
			}
		})
	}
	// draw gems
	l.gemGrid.ForEach(func(p Coord, s Square) {
		if s.gem >= 0 {
			s.drawGem(screen, l.gemImages[s.gem], frameCount)
		} else {
			// shouldn't get here
			s.drawBackground(screen, darkGreen)
		}
	})

}

// Delete all gems in a row. If it does nor result in a triple the delete will fail and restore to original state.
func (l *LevelGems) deleteRows(numRows, frameCount int) bool {
	// copy the gem grid
	newGrid := l.gemGrid.Copy()

	// remove the rows
	for i := 0; i < numRows; i++ {
		row := l.cursorGem.y
		if row >= newGrid.NumRows() {
			break
		}
		newGrid = newGrid.DeleteRow(row)
	}

	// check if the swap will create a triple
	makesATriple, _ := findTriples(newGrid)

	if makesATriple {
		// mark squares in deleted rows as EMPTY_GEM
		for i := 0; i < numRows; i++ {
			row := l.cursorGem.y + i
			if row >= l.gemGrid.NumRows() {
				break
			}
			for x := range l.gemGrid.NumColumns() {
				p := Coord{x, row}
				setGem(&l.gemGrid, p, emptyGem)
			}
		}
	} else {
		PlaySound(failOgg)
		// penalize player for invalid move
		// mark row to be deleted as EMPTY_GEM
		for x := range l.gemGrid.NumColumns() {
			l.triplesMask.Set(Coord{x, l.cursorGem.y}, false)
		}
		return false
	}
	l.fillEmpties(frameCount, true)
	return true
}

// Delete all gems selected in visual mode.
// If it does nor result in a triple the delete will fail and restore to original state.
func (l *LevelGems) deleteSelectionReplaceFromBelow(frameCount int) bool {

	// copy the gem grid to test if removing the selected squares will result in a triple
	newGrid := l.gemGrid.Copy()

	// set all selected squares to EMPTY_GEM
	setSelection(&newGrid, l.cursorGem, l.swapGem, emptyGem)
	moveUpFromBelow(&newGrid, frameCount, false)

	// check if the swap will create a triple
	makesATriple, _ := findTriples(newGrid)

	if makesATriple {
		// set all selected squares to EMPTY_GEM
		setSelection(&l.gemGrid, l.cursorGem, l.swapGem, emptyGem)
	} else {
		PlaySound(failOgg)
		// penalize player for invalid move
		// mark row to be deleted as EMPTY_GEM
		for x := range l.gemGrid.NumColumns() {
			l.triplesMask.Set(Coord{x, l.cursorGem.y}, false)
		}
		return false
	}
	l.fillEmpties(frameCount, true)
	return true
}

// Not used.
// Delete all gems selected in visual mode.
// If it does nor result in a triple the delete will fail and restore to original state.
// This moves gems left to fill the empty space and wraps gems from the next row.
func (l *LevelGems) deleteSelectionReplaceFromRight(frameCount int) bool {
	// copy the gem grid to test if removing the selected squares will result in a triple
	newGrid := l.gemGrid.Copy()

	{ // using bracket notation to scope local variables

		// set all selected squares to EMPTY_GEM
		selectionStart, selectionEnd := highLow(l.cursorGem, l.swapGem)
		dest := newGrid.IndexOf(selectionStart)
		src := newGrid.IndexOf(selectionEnd) + 1

		var gem int
		// move squares past selection into the selection
		for i := dest; i < newGrid.LastIndex(); i++ {
			if src >= newGrid.LastIndex() {
				gem = emptyGem
			} else {
				gem = newGrid.GetAtIndex(src).gem
			}
			setGem(&newGrid, newGrid.IndexToCoord(i), gem)
			src++
		}
	}

	// check if the swap will create a triple
	makesATriple, _ := findTriples(newGrid)

	if makesATriple {
		// set all selected squares to EMPTY_GEM
		selectionStart, selectionEnd := highLow(l.cursorGem, l.swapGem)
		dest := l.gemGrid.IndexOf(selectionStart)
		src := l.gemGrid.IndexOf(selectionEnd) + 1

		// move squares past selection into the selection
		for i := dest; i < l.gemGrid.LastIndex(); i++ {
			if src >= l.gemGrid.LastIndex() {
				setGem(&l.gemGrid, l.gemGrid.IndexToCoord(i), emptyGem)
			} else {
				setGem(&l.gemGrid, l.gemGrid.IndexToCoord(i), l.gemGrid.GetAtIndex(src).gem)

				sqPtr := l.gemGrid.GetPtr(l.gemGrid.IndexToCoord(i))
				sqPtr.addMover(frameCount, dropDuration,
					l.gemGrid.IndexToCoord(src),
					l.gemGrid.IndexToCoord(i))
			}
			src++
		}
	} else {
		PlaySound(failOgg)
		// penalize player for invalid move
		// mark row to be deleted as EMPTY_GEM
		for x := range l.gemGrid.NumColumns() {
			l.triplesMask.Set(Coord{x, l.cursorGem.y}, false)
		}
		return false
	}

	// note this is identical to bottom of fillEmpties()
	// fill empties at the bottom of the gemGrid with newly generated gems
	for x := range numGemColumns {
		for y := range numGemRows {
			if l.gemGrid.Get(Coord{x, y}).gem == emptyGem {
				sqPtr := l.gemGrid.GetPtr(Coord{x, y})
				sqPtr.gem = rng.Intn(l.numGems)
				sqPtr.addMover(frameCount, dropDuration,
					Coord{sqPtr.coords.x, numGemRows + 1},
					Coord{x, y})
			}
		}
	}
	return true
}

func (l *LevelGems) fillEmpties(frameCount int, addMover bool) {
	moveUpFromBelow(&l.gemGrid, frameCount, addMover)
	fillFromBelow(&l.gemGrid, l.numGems, frameCount, addMover)
}

// fills the entire gemGrid with random gems
func (l *LevelGems) fillRandom() {
	for y := range l.gemGrid.NumRows() {
		for x := range l.gemGrid.NumColumns() {
			setGem(&l.gemGrid, Coord{x, y}, rng.Intn(l.numGems))
		}
	}
}

func findSquareBelow(gemGrid *Grid[Square], p Coord) Coord {
	for y := p.y; y < numGemRows; y++ {
		if gemGrid.Get(Coord{p.x, y}).gem != emptyGem {
			return Coord{p.x, y}
		}
	}
	return Coord{-1, -1} // did not find a square with color
}

func findTriples(gemGrid Grid[Square]) (bool, Grid[bool]) {
	// create a local mask to mark all square that are in triples
	mask := NewGridOfBools(numGemColumns, numGemRows)

	found := false
	// find all horizontal triples
	for y := range gemGrid.NumRows() {
		for x := range gemGrid.NumColumns() - 2 {
			p := Coord{x, y}
			if gemGrid.Get(p).gem >= 0 { // if is a gem
				if gemGrid.Get(p).gem == gemGrid.Get(Coord{x + 1, y}).gem && gemGrid.Get(p).gem == gemGrid.Get(Coord{x + 2, y}).gem {
					mask.Set(p, true)
					mask.Set(Coord{x + 1, y}, true)
					mask.Set(Coord{x + 2, y}, true)
					found = true
				}
			}
		}
	}

	// find all vertical triples
	for y := range gemGrid.NumRows() - 2 {
		for x := range gemGrid.NumColumns() {
			p := Coord{x, y}
			if gemGrid.Get(p).gem >= 0 { // if is a gem
				if gemGrid.Get(p).gem == gemGrid.Get(Coord{x, y + 1}).gem && gemGrid.Get(p).gem == gemGrid.Get(Coord{x, y + 2}).gem {
					mask.Set(p, true)
					mask.Set(Coord{x, y + 1}, true)
					mask.Set(Coord{x, y + 2}, true)
					found = true
				}
			}
		}
	}
	return found, mask
}

func (l *LevelGems) gameIsWon() bool {
	for y := range l.gemGrid.NumRows() {
		for x := range l.gemGrid.NumColumns() {
			if !l.triplesMask.Get(Coord{x, y}) {
				return false
			}
		}
	}
	return true
}

func (l *LevelGems) handleKeyNormalMode(key ebiten.Key, frameCount int) {
	switch key {
	case ebiten.KeyD:
		if equals(tail(globalKeys, 2), []ebiten.Key{ebiten.KeyD, ebiten.KeyD}) {
			if l.level != LevelIdGemsVM {
				l.deleteRows(1, frameCount)
			}
			clearKeystrokes()
		}
	case ebiten.KeyEnter:
		if l.level != LevelIdGemsVM {
			t := tail(globalKeys, 3)
			if len(t) > 1 {
				if t[0] == ebiten.KeyD {
					n := t[1] - ebiten.Key0
					if n > 0 && n < 10 {
						l.deleteRows(int(n), frameCount)
						clearKeystrokes()
					}
				}
			}
		}
	case ebiten.KeyH:
		l.cursorGem.x = max(l.cursorGem.x-1, 0)
		clearKeystrokes()
	case ebiten.KeyL:
		l.cursorGem.x = min(l.cursorGem.x+1, numGemColumns-1)
		clearKeystrokes()
	case ebiten.KeyK:
		l.cursorGem.y = max(l.cursorGem.y-1, 0)
		clearKeystrokes()
	case ebiten.KeyJ:
		l.cursorGem.y = min(l.cursorGem.y+1, numGemRows-1)
		clearKeystrokes()
	case ebiten.KeyV:
		// entering VisualMode (where we do swaps)
		if l.level != LevelIdGemsDD {
			l.swapGem = l.cursorGem
			l.viMode = VisualMode
		}
		clearKeystrokes()
	}
}

func (l *LevelGems) handleKeyVisualMode(key ebiten.Key, frameCount int) {
	switch key {
	case ebiten.KeyH:
		l.swapGem.x = max(l.swapGem.x-1, 0)
		clearKeystrokes()
	case ebiten.KeyL:
		l.swapGem.x = min(l.swapGem.x+1, numGemColumns-1)
		clearKeystrokes()
	case ebiten.KeyK:
		l.swapGem.y = max(l.swapGem.y-1, 0)
		clearKeystrokes()
	case ebiten.KeyJ:
		l.swapGem.y = min(l.swapGem.y+1, numGemRows-1)
		clearKeystrokes()
	case ebiten.KeyV:
		PlaySound(failOgg)
		clearKeystrokes()
	case ebiten.KeyEscape:
		// exit visual mode without swapping
		l.viMode = NormalMode
		l.cursorGem = l.swapGem
		l.swapGem = Coord{-1, -1}
		clearKeystrokes()
	case ebiten.KeyD:
		// attempt swap
		if result := l.deleteSelectionReplaceFromBelow(frameCount); result {
			// swap successful
			// exiting visual mode
			l.viMode = NormalMode
			l.cursorGem = l.swapGem
			l.swapGem = Coord{-1, -1}
		} else {
			PlaySound(failOgg)
		}
		clearKeystrokes()
	}
}

func (l *LevelGems) Initialize(id LevelID) {
	l.level = id
	switch id {
	case LevelIdGemsEnd:
		l.numGems = 6
		numGemColumns = 10
	case LevelIdGemsVM:
		l.numGems = 5
		numGemColumns = 5
	case LevelIdGemsDD:
		l.numGems = 4
		numGemColumns = 8
	}
	l.cursorGem = Coord{numGemColumns / 2, numGemRows / 2}
	l.swapGem = Coord{-1, -1}
	l.gemGrid = newGridOfSquares(numGemColumns, numGemRows)
	l.triplesMask = NewGridOfBools(numGemColumns, numGemRows)

	l.viMode = NormalMode
	l.fillRandom()
	l.loadGems()
}

func (l *LevelGems) loadGems() {
	if len(l.gemImages) == 0 {
		l.gemImages = make([]*ebiten.Image, l.numGems)
		for i := range l.numGems {
			image := loadImage("resources/Gem " + strconv.Itoa(i+1) + ".png")
			l.gemImages[i] = image
		}
	}
}

func (l *LevelGems) Update(frameCount int) (bool, error) {
	if isCheatKeyPressed() {
		return true, nil
	}
	// clear movers if expired
	l.gemGrid.ForEach(func(p Coord, s Square) {
		if s.mover != nil {
			if s.mover.endFrame < frameCount {
				sqPtr := l.gemGrid.GetPtr(p)
				sqPtr.mover = nil
				l.updateTriples(frameCount)
			}
		}
	})

	for _, key := range globalKeys {
		if inpututil.IsKeyJustPressed(key) {
			switch l.viMode {
			case NormalMode:
				l.handleKeyNormalMode(key, frameCount)
			case VisualMode:
				l.handleKeyVisualMode(key, frameCount)
			}
		}
	}

	if l.gameIsWon() {
		return true, nil
	}
	return false, nil
}

func (l *LevelGems) updateTriples(frameCount int) {
	found, mask := findTriples(l.gemGrid)

	if found {
		// now that we have completed detecting all triples we can update the game state
		for y := range l.gemGrid.NumRows() {
			for x := range l.gemGrid.NumColumns() {
				if mask.Get(Coord{x, y}) {
					sq := l.gemGrid.GetPtr(Coord{x, y})
					sq.gem = emptyGem
					l.triplesMask.Set(Coord{x, y}, true)
				}
			}
		}
		// play triple sound unless the level is complete
		if !l.gameIsWon() {
			PlaySound(tripleOgg)
			l.fillEmpties(frameCount, true)
		} else {
			// fill empties so Draw() continutes to work
			l.fillEmpties(frameCount, false)
			// remove all movers
			l.gemGrid.ForEach(func(p Coord, s Square) {
				if s.mover != nil {
					sqPtr := l.gemGrid.GetPtr(p)
					sqPtr.mover = nil
					l.updateTriples(frameCount)
				}
			})

		}
	}
}

func (square *Square) addMover(startFrame int, duration int, from Coord, to Coord) {
	// add animation
	mover := new(GemMover)

	mover.startFrame = startFrame
	mover.endFrame = startFrame + duration
	mover.startGrid = from
	mover.endCoord = to

	square.mover = mover
}

func (square *Square) drawBackground(screen *ebiten.Image, color color.Color) {
	pos := squareToScreenPoint(square.coords)
	vector.DrawFilledRect(screen, float32(pos.x), float32(pos.y), gemCellSize-4, gemCellSize-4, color, false)
}

func (square *Square) drawGem(screen *ebiten.Image, gemImage *ebiten.Image, frameCount int) {
	if square.gem >= 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(gemScale, gemScale)

		// Bugbug: applyMove should be called from Update(), not Draw()
		if square.mover != nil {
			applyMover(square.mover, op, frameCount)
		} else {
			pos := squareToScreenPoint(square.coords)
			op.GeoM.Translate(float64(pos.x), float64(pos.y))
		}
		screen.DrawImage(gemImage, op)
	}
}
