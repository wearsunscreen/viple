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
	EMPTY_GEM     = -1
	gemCellSize   = 50
	dropDuration  = 60
	gemScale      = float64(gemCellSize-4) / float64(gemWidth)
	gemWidth      = 100
	numGemRows    = 11
	swapDuration  = 40
)

type LevelGemsVisualMode struct {
	cursorGem   Point
	gemGrid     Grid[Square]
	gemImages   []*ebiten.Image
	keyInput    string
	level       LevelID
	mode        int
	numGems     int
	swapGem     Point
	triplesMask Grid[bool]
}

type GemMover struct {
	startFrame int
	endFrame   int
	startPoint Point // grid coordinates
	endPoint   Point // grid coordinates
}

var numGemColumns int

/* ==================================
/* ==================================
	LevelGemsVisualMode methods
/* ==================================
/* ==================================
*/

func (l *LevelGemsVisualMode) Draw(screen *ebiten.Image, frameCount int) {

	screen.Fill(mediumCoal)

	// draw background of triples
	l.gemGrid.ForEach(func(p Point, s Square) {
		if l.triplesMask.Get(p) {
			s.drawBackground(screen, darkButter)
		}
	})

	// draw cursor
	cursorColors := [2]color.Color{color.White, color.RGBA{0, 0, 0, 0}}
	blink := frameCount / blinkInverval % 2

	switch l.level {
	case LevelIdGemsVM:
		// we are in swap mode, faster blink, brighter colors
		blink = frameCount / (blinkInverval / 2) % 2
		cursorColors = [2]color.Color{brightRed, lightButter}

		// draw visualmode cursor
		if l.mode == VisualMode {
			cursorStart, cursorEnd := highLow(l.cursorGem, l.swapGem)
			startX := cursorStart.x
			for y := cursorStart.y; y <= cursorEnd.y; y++ {
				for x := startX; x < numGemColumns; x++ {
					s := l.gemGrid.Get(Point{x, y})
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
		l.gemGrid.ForEach(func(p Point, s Square) {
			if p.y == l.cursorGem.y {
				s.drawBackground(screen, cursorColors[blink])
			}
		})
	}
	// draw gems
	l.gemGrid.ForEach(func(p Point, s Square) {
		s.drawGem(screen, l.gemImages[s.gem], frameCount)
	})

}

func (l *LevelGemsVisualMode) Initialize(id LevelID) {
	l.level = id
	switch id {
	case LevelIdGemsVM:
		l.numGems = 5
		numGemColumns = 5
	case LevelIdGemsDD:
		l.numGems = 4
		numGemColumns = 8
	}
	l.cursorGem = Point{numGemColumns / 2, numGemRows / 2}
	l.swapGem = Point{-1, -1}
	l.gemGrid = NewGridOfSquares(numGemColumns, numGemRows)
	l.triplesMask = NewGridOfBools(numGemColumns, numGemRows)

	l.mode = NormalMode
	fillRandom(l)

	l.loadGems()
}

func (l *LevelGemsVisualMode) Update(frameCount int) (bool, error) {
	// clear movers if expired
	l.gemGrid.ForEach(func(p Point, s Square) {
		if s.mover != nil {
			if s.mover.endFrame < frameCount {
				sqPtr := l.gemGrid.GetPtr(p)
				sqPtr.mover = nil
				updateTriples(l, frameCount)
			}
		}
	})

	for _, key := range globalKeys {
		if inpututil.IsKeyJustPressed(key) {
			switch l.mode {
			case NormalMode:
				handleKeyDeleteRows(l, key, frameCount)
			case VisualMode:
				handleKeyVisualMode(l, key, frameCount)
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
	gem    int
	mover  *GemMover
	coords Point // position in the grid
}

func NewGridOfSquares(width, height int) Grid[Square] {
	r := make([][]Square, height)
	for i := range r {
		r[i] = make([]Square, width)
		for j := range r[i] {
			r[i][j] = Square{coords: Point{j, i}}
		}
	}
	return Grid[Square]{rows: r}
}

func (square *Square) AddMover(startFrame int, duration int, from Point, to Point) {
	// add animation
	mover := new(GemMover)

	mover.startFrame = startFrame
	mover.endFrame = startFrame + duration
	mover.startPoint = from
	mover.endPoint = to

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

/* ==================================
/* ==================================
	Other functions
/* ==================================
/* ==================================
*/

func applyMover(mover *GemMover, op *ebiten.DrawImageOptions, frameCount int) {
	completionRatio := 1 - float64(mover.endFrame-frameCount)/float64(mover.endFrame-mover.startFrame)
	startPosition := squareToScreenPoint(mover.startPoint)
	endPosition := squareToScreenPoint(mover.endPoint)
	op.GeoM.Translate(
		float64(startPosition.x)+(completionRatio*float64(endPosition.x-startPosition.x)),
		float64(startPosition.y)+(completionRatio*float64(endPosition.y-startPosition.y)))
}

// TODO: rewrite this as a generic function for any type
func copyGrid(grid [][]Square) [][]Square {
	clone := make([][]Square, len(grid))
	for i := range grid {
		clone[i] = make([]Square, len(grid[i]))
		copy(clone[i], grid[i])
	}
	return clone
}

// Delete all gems in a row. If it does nor result in a triple the delete will fail and restore to original state.
func deleteRows(l *LevelGemsVisualMode, numRows, frameCount int) bool {
	// copy the gem grid
	newGrid := l.gemGrid.Copy()

	// remove the rows
	for i := 0; i < numRows; i++ {
		row := l.cursorGem.y
		if row >= newGrid.NumRows() {
			break
		}
		newGrid.deleteRow(row)
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
				p := Point{x, row}
				SetGem(&l.gemGrid, p, EMPTY_GEM)
			}
		}
	} else {
		PlaySound(failOgg)
		// penalize player for invalid move
		// mark row to be deleted as EMPTY_GEM
		for x := range l.gemGrid.NumColumns() {
			l.triplesMask.Set(Point{x, l.cursorGem.y}, false)
		}
		return false
	}
	fillEmpties(l, frameCount)
	return true
}

func SetGem(grid *Grid[Square], pt Point, gem int) {
	sq := grid.Get(pt)
	sq.gem = gem
	grid.Set(pt, sq)
}

// Delete all gems selected in visual mode.
// If it does nor result in a triple the delete will fail and restore to original state.
func deleteSelection(l *LevelGemsVisualMode, frameCount int) bool {

	// copy the gem grid to test if removing the selected squares will result in a triple
	newGrid := l.gemGrid.Copy()

	{
		// set all selected squares to EMPTY_GEM
		selectionStart, selectionEnd := highLow(l.cursorGem, l.swapGem)
		dest := newGrid.IndexOf(selectionStart)
		src := newGrid.IndexOf(selectionEnd) + 1

		var gem int
		// move squares past selection into the selection
		for i := dest; i < newGrid.LastIndex(); i++ {
			if src >= newGrid.LastIndex() {
				gem = EMPTY_GEM
			} else {
				gem = newGrid.GetAtIndex(src).gem
			}
			SetGem(&newGrid, newGrid.IndexToPoint(i), gem)
			src++
		}
	}

	// check if the swap will create a triple
	makesATriple, _ := findTriples(newGrid)

	if makesATriple {
		l.gemGrid = newGrid
	} else {
		PlaySound(failOgg)
		// penalize player for invalid move
		// mark row to be deleted as EMPTY_GEM
		for x := range l.gemGrid.NumColumns() {
			l.triplesMask.Set(Point{x, l.cursorGem.y}, false)
		}
		return false
	}

	// note this is identical to bottom of fillEmpties()
	// fill empties at the bottom of the gemGrid with newly generated gems
	for x := range numGemColumns {
		for y := range numGemRows {
			if l.gemGrid.Get(Point{x, y}).gem == EMPTY_GEM {
				sqPtr := l.gemGrid.GetPtr(Point{x, y})
				sqPtr.gem = rng.Intn(l.numGems)
				sqPtr.AddMover(frameCount, dropDuration,
					Point{sqPtr.coords.x, numGemRows + 1},
					Point{x, y})
			}
		}
	}
	return true
}

func fillEmpties(l *LevelGemsVisualMode, frameCount int) {
	// for each column
	for x := range numGemColumns {
		// find empty square and move squares from below to fill
		for y := 0; y < numGemRows; y++ {
			pt := Point{x, y}
			if l.gemGrid.Get(pt).gem == EMPTY_GEM {
				below := findSquareBelow(&l.gemGrid, pt)
				if below.y >= 0 {
					sqPtr := l.gemGrid.GetPtr(pt)
					sqPtr.gem = l.gemGrid.Get(below).gem
					sqPtr = l.gemGrid.GetPtr(below)
					sqPtr.gem = EMPTY_GEM
					sqPtr = l.gemGrid.GetPtr(pt)
					sqPtr.AddMover(frameCount, dropDuration, below, pt)
				}
			}
		}
	}

	// fill empties at the bottom of the gemGrid with newly generated gems
	for x := range numGemColumns {
		for y := range numGemRows {
			if l.gemGrid.Get(Point{x, y}).gem == EMPTY_GEM {
				sqPtr := l.gemGrid.GetPtr(Point{x, y})
				sqPtr.gem = rng.Intn(l.numGems)
				sqPtr.AddMover(frameCount, dropDuration,
					Point{sqPtr.coords.x, numGemRows + 1},
					Point{x, y})
			}
		}
	}
}

// fills the entire gemGrid with random gems
func fillRandom(l *LevelGemsVisualMode) {
	for y := range l.gemGrid.NumRows() {
		for x := range l.gemGrid.NumColumns() {
			SetGem(&l.gemGrid, Point{x, y}, rng.Intn(l.numGems))
		}
	}
}

func findSquareAbove(gemGrid *Grid[Square], p Point) Point {
	for y := range p.y {
		y = p.y - 1 - y
		for gemGrid.Get(p).gem != EMPTY_GEM {
			return Point{p.x, y}
		}
	}
	return Point{-1, -1} // did not find a square with color
}

func findSquareBelow(gemGrid *Grid[Square], p Point) Point {
	for y := p.y; y < numGemRows; y++ {
		if gemGrid.Get(Point{p.x, y}).gem != EMPTY_GEM {
			return Point{p.x, y}
		}
	}
	return Point{-1, -1} // did not find a square with color
}
func findTriples(gemGrid Grid[Square]) (bool, Grid[bool]) {
	// create a local mask to mark all square that are in triples
	mask := NewGridOfBools(numGemColumns, numGemRows)

	found := false
	// find all horizontal triples
	for y := range gemGrid.NumRows() {
		for x := range gemGrid.NumColumns() - 2 {
			p := Point{x, y}
			if gemGrid.Get(p).gem >= 0 { // if is a gem
				if gemGrid.Get(p).gem == gemGrid.Get(Point{x + 1, y}).gem && gemGrid.Get(p).gem == gemGrid.Get(Point{x + 2, y}).gem {
					mask.Set(p, true)
					mask.Set(Point{x + 1, y}, true)
					mask.Set(Point{x + 2, y}, true)
					found = true
				}
			}
		}
	}

	// find all vertical triples
	for y := range gemGrid.NumRows() - 2 {
		for x := range gemGrid.NumColumns() {
			p := Point{x, y}
			if gemGrid.Get(p).gem >= 0 { // if is a gem
				if gemGrid.Get(p).gem == gemGrid.Get(Point{x, y + 1}).gem && gemGrid.Get(p).gem == gemGrid.Get(Point{x, y + 2}).gem {
					mask.Set(p, true)
					mask.Set(Point{x, y + 1}, true)
					mask.Set(Point{x, y + 2}, true)
					found = true
				}
			}
		}
	}
	return found, mask
}

func handleKeyDeleteRows(l *LevelGemsVisualMode, key ebiten.Key, frameCount int) {
	switch key {
	case ebiten.KeyD:
		if equals(tail(globalKeys, 2), []ebiten.Key{ebiten.KeyD, ebiten.KeyD}) {
			deleteRows(l, 1, frameCount)
			clearKeystrokes()
		}
	case ebiten.KeyEnter:
		t := tail(globalKeys, 3)
		if len(t) > 1 {
			if t[0] == ebiten.KeyD {
				n := t[1] - ebiten.Key0
				if n > 0 && n < 10 {
					deleteRows(l, int(n), frameCount)
					clearKeystrokes()
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
		l.swapGem = l.cursorGem
		l.mode = VisualMode
		clearKeystrokes()
	}
}

func handleKeyVisualMode(l *LevelGemsVisualMode, key ebiten.Key, frameCount int) {
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
		l.mode = NormalMode
		l.cursorGem = l.swapGem
		l.swapGem = Point{-1, -1}
		clearKeystrokes()
	case ebiten.KeyY:
		// attempt swap
		if l.swapGem.x != -1 && l.swapGem != l.cursorGem {
			if result := deleteSelection(l, frameCount); result {
				// swap successful
				// exiting visual mode
				l.mode = NormalMode
				l.cursorGem = l.swapGem
				l.swapGem = Point{-1, -1}
			} else {
				PlaySound(failOgg)
			}
		}
		clearKeystrokes()
	}
}

// return the highest and lowest points in that order where highest is closest to the top left of the screen
// disregard a point if it has negative values
func highLow(p1, p2 Point) (Point, Point) {
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

func (l *LevelGemsVisualMode) gameIsWon() bool {
	for y := range l.gemGrid.NumRows() {
		for x := range l.gemGrid.NumColumns() {
			if !l.triplesMask.Get(Point{x, y}) {
				return false
			}
		}
	}
	return true
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
	heightOfGrid := gemCellSize * numGemRows
	yMargin := (screenHeight - heightOfGrid) / 2
	return Point{
		gemCellSize*squareXY.x + xMargin,
		gemCellSize*squareXY.y + yMargin,
	}
}

func updateTriples(l *LevelGemsVisualMode, frameCount int) {
	found, mask := findTriples(l.gemGrid)

	if found {
		// now that we have completed detecting all triples we can update the game state
		for y := range l.gemGrid.NumRows() {
			for x := range l.gemGrid.NumColumns() {
				if mask.Get(Point{x, y}) {
					sq := l.gemGrid.GetPtr(Point{x, y})
					sq.gem = EMPTY_GEM
					l.triplesMask.Set(Point{x, y}, true)
				}
			}
		}
		// play triple sound unless the level is complete
		if !l.gameIsWon() {
			PlaySound(tripleOgg)
			fillEmpties(l, frameCount)
		}
	}
}
