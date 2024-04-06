package main

import (
	"flag"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	blinkInverval = 60 / 2
	cellSize      = 50
	dropDuration  = 60
	gemScale      = float64(cellSize-4) / float64(gemWidth)
	gemWidth      = 100
	margin        = 20
	numRows       = 11
	numColumns    = 5
	swapDuration  = 40
	version       = "Viple 0.1"
)

type Mode int

const (
	CommandMode = iota
	LineMode
	InsertMode
)

var (
	rng *rand.Rand
)

type Game struct {
	cursorSquare Point
	swapSquare   Point
	frameCount   int
	grid         [][]Square
	gemImages    []*ebiten.Image
	keyInput     string
	keys         []ebiten.Key
	maxColors    int
	mode         Mode
	numColors    int
	player       *AudioPlayer
	triplesMask  [][]bool
}

func main() {
	var seed int
	flag.IntVar(&seed, "seed", 0, "Seed for random number generation")
	flag.Parse()
	seedRNG(int64(seed))

	ebiten.SetWindowSize(gameDimensions())
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	screen.Fill(mediumCoal)

	DrawGrid(screen, g)

}

func gameDimensions() (width int, height int) {
	return (margin * 2) + (numColumns * cellSize), (margin * 2) + (numRows * cellSize)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w, h := gameDimensions()
	return w, h
}

func loadImage(path string) *ebiten.Image {
	image, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Error loading image: %v", err)
	}
	return image
}

func loadImages(g *Game) {
	g.gemImages = make([]*ebiten.Image, g.numColors)
	for i := range g.numColors {
		image := loadImage("resources/Gem " + strconv.Itoa(i+1) + ".png")
		g.gemImages[i] = image
	}
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

	g.numColors = 5
	g.mode = CommandMode
	fillRandom(&g)

	loadImages(&g)

	return &g
}

func offsetPoint(p, offset Point) Point {
	return Point{p.x + offset.x, p.y + offset.y}
}

func seedRNG(seed int64) {
	if seed == 0 {
		seed = time.Now().UnixNano() % 10000
	}
	log.Println("Random seed is ", seed)
	rng = rand.New(rand.NewSource(seed))
}

// Exchange positions of two neighboring squares, return false if unable to exchange.
// The exchange fails if swap point and cursor point are the same (this can happen when
// player attempts to move off the grid). The exchange fails if both points have the
// same value.
func SwapSquares(g *Game) bool {
	if g.swapSquare == g.cursorSquare {
		return false
	}
	if g.grid[g.swapSquare.y][g.swapSquare.x].point == g.grid[g.cursorSquare.y][g.cursorSquare.x].point {
		return false
	}

	// swap colors
	fromSquare := &g.grid[g.swapSquare.y][g.swapSquare.x]
	toSquare := &g.grid[g.cursorSquare.y][g.cursorSquare.x]
	temp := fromSquare.color
	fromSquare.color = toSquare.color
	toSquare.color = temp

	// check if the swap will create a triple
	makesATriple, _ := FindTriples(g.grid)

	if makesATriple {
		toSquare.AddMover(g.frameCount, 60, fromSquare.point, toSquare.point)
		fromSquare.AddMover(g.frameCount, 60, toSquare.point, fromSquare.point)
	} else {
		// restore original colors
		temp := fromSquare.color
		fromSquare.color = toSquare.color
		toSquare.color = temp
		g.swapSquare = g.cursorSquare // return the cursor to the original location

		// tell the user he made an invalid move
		return false
	}
	return true
}

func handleKeyCommand(g *Game, key ebiten.Key) {
	switch key {
	case ebiten.KeyH:
		g.cursorSquare.x = max(g.cursorSquare.x-1, 0)
	case ebiten.KeyL:
		g.cursorSquare.x = min(g.cursorSquare.x+1, numColumns-1)
	case ebiten.KeyK:
		g.cursorSquare.y = max(g.cursorSquare.y-1, 0)
	case ebiten.KeyJ:
		g.cursorSquare.y = min(g.cursorSquare.y+1, numRows-1)
	case ebiten.KeyI:
		// entering InsertMode (where we do swaps)
		g.swapSquare = g.cursorSquare
		g.mode = InsertMode
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
		g.swapSquare.x = max(g.swapSquare.x-1, 0)
	case ebiten.KeyL:
		g.swapSquare.x = min(g.swapSquare.x+1, numColumns-1)
	case ebiten.KeyK:
		g.swapSquare.y = max(g.swapSquare.y-1, 0)
	case ebiten.KeyJ:
		g.swapSquare.y = min(g.swapSquare.y+1, numRows-1)
	case ebiten.KeyI:
		PlaySound(failOgg)
	case ebiten.KeyEscape:
		g.mode = CommandMode
		g.swapSquare = Point{-1, -1}
	}
	if g.swapSquare.x != -1 && g.swapSquare != g.cursorSquare {
		if result := SwapSquares(g); !result {
			PlaySound(failOgg)
		}
		g.cursorSquare = g.swapSquare
	}

}

func (g *Game) Update() error {
	g.frameCount++
	// clear movers if expired
	for y, row := range g.grid {
		for x, _ := range row {
			if g.grid[y][x].mover != nil {
				if g.grid[y][x].mover.endFrame < g.frameCount {
					g.grid[y][x].mover = nil
					UpdateTriples(g)
				}
			}
		}
	}

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	for _, key := range g.keys {
		if inpututil.IsKeyJustPressed(key) {
			switch g.mode {
			case CommandMode:
				handleKeyCommand(g, key)
			case InsertMode:
				handleKeyInsert(g, key)
			}
		}
	}
	return nil
}
