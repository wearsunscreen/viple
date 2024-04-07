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
	gemCellSize   = 50
	dropDuration  = 60
	gemScale      = float64(gemCellSize-4) / float64(gemWidth)
	gemWidth      = 100
	gemsMargin    = 20
	gemRows       = 11
	numColumns    = 5
	screenWidth   = 800
	screenHeight  = 600
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

type Level3 struct {
	cursorGem   Point
	gemGrid     [][]Square
	gemImages   []*ebiten.Image
	mode        Mode
	numGems     int
	swapGem     Point
	triplesMask [][]bool
}

type Game struct {
	frameCount int
	keyInput   string
	keys       []ebiten.Key
	player     *AudioPlayer
	l3         Level3
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
	return screenWidth, screenHeight
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
	g.l3.gemImages = make([]*ebiten.Image, g.l3.numGems)
	for i := range g.l3.numGems {
		image := loadImage("resources/Gem " + strconv.Itoa(i+1) + ".png")
		g.l3.gemImages[i] = image
	}
}

func newGame() *Game {
	g := Game{}

	g.l3.numGems = 5
	g.l3.cursorGem = Point{numColumns / 2, gemRows / 2}
	g.l3.swapGem = Point{-1, -1}
	g.l3.gemGrid = make([][]Square, gemRows)
	for y := range g.l3.gemGrid {
		g.l3.gemGrid[y] = make([]Square, numColumns)
	}

	for y, row := range g.l3.gemGrid {
		for x, _ := range row {
			g.l3.gemGrid[y][x].point = Point{x, y}
		}
	}

	g.l3.triplesMask = make([][]bool, gemRows)
	for i := range g.l3.triplesMask {
		g.l3.triplesMask[i] = make([]bool, numColumns)
	}

	g.l3.numGems = 5
	g.l3.mode = CommandMode
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
// player attempts to move off the gemGrid). The exchange fails if both points have the
// same value.
func SwapSquares(g *Game) bool {
	if g.l3.swapGem == g.l3.cursorGem {
		return false
	}
	if g.l3.gemGrid[g.l3.swapGem.y][g.l3.swapGem.x].point == g.l3.gemGrid[g.l3.cursorGem.y][g.l3.cursorGem.x].point {
		return false
	}

	// swap colors
	fromSquare := &g.l3.gemGrid[g.l3.swapGem.y][g.l3.swapGem.x]
	toSquare := &g.l3.gemGrid[g.l3.cursorGem.y][g.l3.cursorGem.x]
	temp := fromSquare.color
	fromSquare.color = toSquare.color
	toSquare.color = temp

	// check if the swap will create a triple
	makesATriple, _ := FindTriples(g.l3.gemGrid)

	if makesATriple {
		toSquare.AddMover(g.frameCount, 60, fromSquare.point, toSquare.point)
		fromSquare.AddMover(g.frameCount, 60, toSquare.point, fromSquare.point)
	} else {
		// restore original colors
		temp := fromSquare.color
		fromSquare.color = toSquare.color
		toSquare.color = temp
		g.l3.swapGem = g.l3.cursorGem // return the cursor to the original location

		// tell the user he made an invalid move
		return false
	}
	return true
}

func handleKeyCommand(g *Game, key ebiten.Key) {
	switch key {
	case ebiten.KeyH:
		g.l3.cursorGem.x = max(g.l3.cursorGem.x-1, 0)
	case ebiten.KeyL:
		g.l3.cursorGem.x = min(g.l3.cursorGem.x+1, numColumns-1)
	case ebiten.KeyK:
		g.l3.cursorGem.y = max(g.l3.cursorGem.y-1, 0)
	case ebiten.KeyJ:
		g.l3.cursorGem.y = min(g.l3.cursorGem.y+1, gemRows-1)
	case ebiten.KeyI:
		// entering InsertMode (where we do swaps)
		g.l3.swapGem = g.l3.cursorGem
		g.l3.mode = InsertMode
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
		g.l3.swapGem.x = max(g.l3.swapGem.x-1, 0)
	case ebiten.KeyL:
		g.l3.swapGem.x = min(g.l3.swapGem.x+1, numColumns-1)
	case ebiten.KeyK:
		g.l3.swapGem.y = max(g.l3.swapGem.y-1, 0)
	case ebiten.KeyJ:
		g.l3.swapGem.y = min(g.l3.swapGem.y+1, gemRows-1)
	case ebiten.KeyI:
		PlaySound(failOgg)
	case ebiten.KeyEscape:
		g.l3.mode = CommandMode
		g.l3.swapGem = Point{-1, -1}
	}
	if g.l3.swapGem.x != -1 && g.l3.swapGem != g.l3.cursorGem {
		if result := SwapSquares(g); !result {
			PlaySound(failOgg)
		}
		g.l3.cursorGem = g.l3.swapGem
	}

}

func (g *Game) Update() error {
	g.frameCount++
	// clear movers if expired
	for y, row := range g.l3.gemGrid {
		for x, _ := range row {
			if g.l3.gemGrid[y][x].mover != nil {
				if g.l3.gemGrid[y][x].mover.endFrame < g.frameCount {
					g.l3.gemGrid[y][x].mover = nil
					UpdateTriples(g)
				}
			}
		}
	}

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	for _, key := range g.keys {
		if inpututil.IsKeyJustPressed(key) {
			switch g.l3.mode {
			case CommandMode:
				handleKeyCommand(g, key)
			case InsertMode:
				handleKeyInsert(g, key)
			}
		}
	}
	return nil
}
