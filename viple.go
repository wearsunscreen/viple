package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	version    = "Viple 0.1"
	margin     = 20
	cellSize   = 30
	numRows    = 20
	numColumns = 10
)

var (
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

type Point struct {
	x, y int
}

type Game struct {
	grid        [][]int
	maxColors   int
	triplesMask [][]bool
	cursorPoint Point
}

/* todo
- detect triples
- create randdom number generator with seed so I can replay specific games
*/

func main() {
	ebiten.SetWindowSize(gameDimensions())
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}

func init() {
	//source := rand.NewSource(2) // seeding the random number generator can be useful in debugging
	source := rand.NewSource(time.Now().UnixNano())
	rng = rand.New(source)
}

func any(bools []bool) bool {
	for _, value := range bools {
		if value {
			return true
		}
	}
	return false
}

func detectTriples(g *Game) {
	// find all horizontal triples
	for y, row := range g.grid[:len(g.grid)] {
		for x := range g.grid[:len(row)-2] {
			if g.grid[y][x] == g.grid[y][x+1] && g.grid[y][x] == g.grid[y][x+2] {
				g.triplesMask[y][x], g.triplesMask[y][x+1], g.triplesMask[y][x+2] = true, true, true
			}
		}
	}

	// find all vertical triples
	for y, row := range g.grid[:len(g.grid)-2] {
		for x := range g.grid[:len(row)] {
			if g.grid[y][x] == g.grid[y+1][x] && g.grid[y][x] == g.grid[y+2][x] {
				g.triplesMask[y][x], g.triplesMask[y+1][x], g.triplesMask[y+2][x] = true, true, true
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, version)

	// draw background
	screen.Fill(mediumCoal)

	detectTriples(g)

	// draw outlines of triples
	for y, row := range g.grid {
		for x := range row {
			if g.triplesMask[y][x] {
				rect := ebiten.NewImage(cellSize, cellSize)
				rect.Fill(lightGreen)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(cellSize*x+margin), float64(cellSize*y+margin))
				screen.DrawImage(rect, op)
			}
		}
	}

	// draw cells
	for y, row := range g.grid {
		for x, col := range row {
			rect := ebiten.NewImage(cellSize-4, cellSize-4)
			rect.Fill(gameColors[col])
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(cellSize*x+margin+2), float64(cellSize*y+margin+2))
			screen.DrawImage(rect, op)
		}
	}

	// draw cursor
	vector.StrokeRect(screen, 50, 50, 100, 100, 5, color.White, false)

}

func fillRandom(arr [][]int, upTo int) {
	for i := range arr {
		for j := range arr[i] {
			arr[i][j] = rng.Intn(upTo) // Generate random number between 0 (inclusive) and upTo (exclusive)
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
		maxColors:   5,
		cursorPoint: Point{numRows / 2, numColumns / 2},
	}

	g.grid = make([][]int, numRows)
	for i := range g.grid {
		g.grid[i] = make([]int, numColumns)
	}

	g.triplesMask = make([][]bool, numRows)
	for i := range g.triplesMask {
		g.triplesMask[i] = make([]bool, numColumns)
	}

	fillRandom(g.grid, 6)
	return &g
}

func (g *Game) Update() error {
	return nil
}
