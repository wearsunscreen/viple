package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	version       = "Viple 0.1"
	margin        = 20
	cellSize      = 30
	numRows       = 11
	numColumns    = 5
	blinkInverval = 60 / 2
)

var (
	brightRed        = color.RGBA{0xfc, 0x10, 0x10, 0xff}
	lightButter      = color.RGBA{0xfc, 0xe9, 0x4f, 0xff}
	lightOrange      = color.RGBA{0xfc, 0xaf, 0x3e, 0xff}
	lightChocolate   = color.RGBA{0xe9, 0xb9, 0x6e, 0xff}
	lightGreen       = color.RGBA{0x8a, 0xe2, 0x34, 0xff}
	lightSkyBlue     = color.RGBA{0x72, 0x9f, 0xcf, 0xff}
	lightPlum        = color.RGBA{0xad, 0x7f, 0xa8, 0xff}
	lightScarletRed  = color.RGBA{0xef, 0x29, 0x29, 0xff}
	lightAluminium   = color.RGBA{0xee, 0xee, 0xec, 0xff}
	lightCoal        = color.RGBA{0x88, 0x8a, 0x85, 0xff}
	mediumButter     = color.RGBA{0xed, 0xd4, 0x00, 0xff}
	mediumOrange     = color.RGBA{0xf5, 0x79, 0x00, 0xff}
	mediumChocolate  = color.RGBA{0xc1, 0x7d, 0x11, 0xff}
	mediumGreen      = color.RGBA{0x73, 0xd2, 0x16, 0xff}
	mediumSkyBlue    = color.RGBA{0x34, 0x65, 0xa4, 0xff}
	mediumPlum       = color.RGBA{0x75, 0x50, 0x7b, 0xff}
	mediumScarletRed = color.RGBA{0xcc, 0x00, 0x00, 0xff}
	mediumAluminium  = color.RGBA{0xd3, 0xd7, 0xcf, 0xff}
	mediumCoal       = color.RGBA{0x55, 0x57, 0x53, 0xff}
	darkButter       = color.RGBA{0xc4, 0xa0, 0x00, 0xff}
	darkOrange       = color.RGBA{0xce, 0x5c, 0x00, 0xff}
	darkChocolate    = color.RGBA{0x8f, 0x59, 0x02, 0xff}
	darkGreen        = color.RGBA{0x4e, 0x9a, 0x06, 0xff}
	darkSkyBlue      = color.RGBA{0x20, 0x4a, 0x87, 0xff}
	darkPlum         = color.RGBA{0x5c, 0x35, 0x66, 0xff}
	darkScarletRed   = color.RGBA{0xa4, 0x00, 0x00, 0xff}
	darkAluminium    = color.RGBA{0xba, 0xbd, 0xb6, 0xff}
	darkCoal         = color.RGBA{0x2e, 0x34, 0x36, 0xff}

	gameColors = [6]color.Color{darkButter, mediumGreen, darkChocolate, lightPlum, mediumSkyBlue, darkScarletRed}
	rng        *rand.Rand
)

type Game struct {
	cursorSquare Point
	swapSquare   Point
	frameCount   int
	grid         [][]Square
	keys         []ebiten.Key
	maxColors    int
	triplesMask  [][]bool
}

type Modifier struct {
	startFrame  int
	endFrame    int
	startOffset [2]int
	endOffset   [2]int
}

type Square struct {
	color     int
	modifiers []Modifier
	point     Point
}

type Point struct {
	x, y int
}

/* todo
- quit with ":q", ":x", ":exit"
- animation
*/

func main() {
	ebiten.SetWindowSize(gameDimensions())
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}

func init() {
	source := rand.NewSource(3) // seeding the random number generator can be useful in debugging
	//source := rand.NewSource(time.Now().UnixNano())
	rng = rand.New(source)
}

func detectTriples(g *Game) {
	// find all horizontal triples
	for y, row := range g.grid[:len(g.grid)] {
		for x := range g.grid[:len(row)-2] {
			if g.grid[y][x].color == g.grid[y][x+1].color && g.grid[y][x].color == g.grid[y][x+2].color {
				g.triplesMask[y][x], g.triplesMask[y][x+1], g.triplesMask[y][x+2] = true, true, true
			}
		}
	}

	// find all vertical triples
	for y, row := range g.grid[:len(g.grid)-2] {
		for x := range g.grid[:len(row)] {
			if g.grid[y][x].color == g.grid[y+1][x].color && g.grid[y][x].color == g.grid[y+2][x].color {
				g.triplesMask[y][x], g.triplesMask[y+1][x], g.triplesMask[y+2][x] = true, true, true
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	screen.Fill(mediumCoal)
	ebitenutil.DebugPrint(screen, version)

	// find triples
	detectTriples(g)

	// draw cells
	for y, row := range g.grid {
		for x, _ := range row {
			DrawSquare(screen, &g.grid[y][x], g.frameCount)
		}
	}

	// draw outlines of triples
	for y, row := range g.grid {
		for x := range row {
			if g.triplesMask[y][x] {
				vector.StrokeRect(screen, float32(cellSize*x+margin), float32(cellSize*y+margin), cellSize, cellSize, 4, lightGreen, false)
			}
		}
	}

	// draw cursor
	cursorColors := [2]color.Color{color.White, color.Black}
	blink := g.frameCount / blinkInverval % 2
	var cursorWidth float32 = 4
	if g.swapSquare.x != -1 {
		// we are in swap mode, faster blink, brighter colors
		blink = g.frameCount / (blinkInverval / 2) % 2
		cursorColors = [2]color.Color{brightRed, lightButter}
		cursorWidth = 6
	}
	vector.StrokeRect(screen, float32(cellSize*g.cursorSquare.x+margin), float32(cellSize*g.cursorSquare.y+margin),
		cellSize, cellSize, cursorWidth, cursorColors[blink], false)
}

func applyModifier(m *Modifier, op *ebiten.DrawImageOptions, frameCount int) {
	var completionRatio float64
	completionRatio = float64(m.endFrame-frameCount) / float64(m.endFrame-m.startFrame)
	op.GeoM.Translate(
		float64(m.startOffset[0])+(completionRatio*float64(m.endOffset[0]-m.startOffset[0])),
		float64(m.startOffset[1])+(completionRatio*float64(m.endOffset[0]-m.startOffset[1])))
}

func DrawSquare(screen *ebiten.Image, square *Square, frameCount int) {
	//vector.DrawFilledRect(screen, float32(cellSize*x+margin+2), float32(cellSize*y+margin+2), cellSize-4, cellSize-4, gameColors[color], false)
	rect := ebiten.NewImage(cellSize-4, cellSize-4)
	rect.Fill(gameColors[square.color])
	op := &ebiten.DrawImageOptions{}
	for _, m := range square.modifiers {
		applyModifier(&m, op, frameCount)
	}
	op.GeoM.Translate(cellSize*float64(square.point.x)+margin+2, cellSize*float64(square.point.y)+margin+2)
	screen.DrawImage(rect, op)
}

func fillRandom(g *Game, upTo int) {
	for y, row := range g.grid {
		for x := range row {
			g.grid[y][x].color = rng.Intn(upTo)
		}
	}
}

func gameDimensions() (width int, height int) {
	return (margin * 2) + (numColumns * cellSize), (margin * 2) + (numRows * cellSize)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w, h := gameDimensions()
	return w, h
}

func newGame() *Game {
	g := Game{
		maxColors:    5,
		cursorSquare: Point{numColumns / 2, numRows / 2},
		swapSquare:   Point{-1, -1},
	}

	g.grid = make([][]Square, numRows)
	for y := range g.grid {
		g.grid[y] = make([]Square, numColumns)
	}

	for y, row := range g.grid {
		for x, _ := range row {
			g.grid[y][x].point = Point{x, y}
		}
	}

	g.triplesMask = make([][]bool, numRows)
	for i := range g.triplesMask {
		g.triplesMask[i] = make([]bool, numColumns)
	}

	fillRandom(&g, 6)
	return &g
}

// Exchange positions of two neighboring squares, return false if unable to exchange.
// The exchange fails if swap point and cursor point are the same (this can happen when
// player attempts to move off the grid). The exchange fails if both points have the
// same value.
func SwapSquares(g *Game) bool {
	if g.swapSquare == g.cursorSquare {
		g.swapSquare = Point{-1, -1} // indicates we are no longer attempting to swap
		return false
	}
	if g.grid[g.swapSquare.y][g.swapSquare.x].point == g.grid[g.cursorSquare.y][g.cursorSquare.x].point {
		g.swapSquare = Point{-1, -1} // indicates we are no longer attempting to swap
		return false
	}

	temp := g.grid[g.swapSquare.y][g.swapSquare.x].color
	g.grid[g.swapSquare.y][g.swapSquare.x].color = g.grid[g.cursorSquare.y][g.cursorSquare.x].color
	g.grid[g.cursorSquare.y][g.cursorSquare.x].color = temp

	// Bugbug: this will delete all existing modifiers
	// add animation
	mover := Modifier{
		startFrame:  g.frameCount,
		endFrame:    g.frameCount + 120,
		startOffset: [2]int{20, 20},
		endOffset:   [2]int{40, 40}}

	newModifiers := []Modifier{mover}

	g.grid[g.swapSquare.y][g.swapSquare.x].modifiers = newModifiers
	g.grid[g.swapSquare.y][g.swapSquare.x].modifiers = newModifiers

	g.swapSquare = Point{-1, -1} // indicates we are no longer attempting to swap
	return true
}

func (g *Game) Update() error {
	g.frameCount++
	// clear all expired Modifiers
	for y, row := range g.grid {
		for x, _ := range row {
			if len(g.grid[y][x].modifiers) > 0 {
				newModifiers := []Modifier{}
				for _, m := range g.grid[y][x].modifiers {
					if g.frameCount < m.endFrame {
						newModifiers = append(newModifiers, m)
					}
				}
				g.grid[y][x].modifiers = newModifiers
			}
		}
	}

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	for _, k := range g.keys {
		if inpututil.IsKeyJustPressed(k) {
			switch k {
			case ebiten.KeyJ:
				g.cursorSquare.x = max(g.cursorSquare.x-1, 0)
			case ebiten.KeyL:
				g.cursorSquare.x = min(g.cursorSquare.x+1, numColumns-1)
			case ebiten.KeyI:
				g.cursorSquare.y = max(g.cursorSquare.y-1, 0)
			case ebiten.KeyK:
				g.cursorSquare.y = min(g.cursorSquare.y+1, numRows-1)
			case ebiten.KeyV:
				if g.swapSquare.x == -1 {
					// initiating a swap
					g.swapSquare = g.cursorSquare
				}
			}
			if g.swapSquare.x != -1 && g.swapSquare != g.cursorSquare {
				if result := SwapSquares(g); !result {
					// TODO: play a buzzer noise to indicate failure
					log.Println("BUZZZ")
				}
			}
		}
	}
	return nil
}
