package main

import (
	"flag"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	version      = "Viple 0.1"
)

type LevelID int

const (
	Level1 = iota
	Level3
)

type Mode int

const (
	NormalMode = iota
	LineMode
	InsertMode
)

var (
	rng *rand.Rand
)

type Game struct {
	frameCount int
	keyInput   string
	keys       []ebiten.Key
	l1         Level1Data
	l3         Level3Data
	level      LevelID
	player     *AudioPlayer
}

func main() {
	var seed int
	flag.IntVar(&seed, "seed", 0, "Seed for random number generation")
	flag.Parse()
	seedRNG(int64(seed))

	ebiten.SetWindowSize(gameDimensions())
	ebiten.SetWindowTitle(version)
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	screen.Fill(lightCoal)
	switch g.level {
	case Level1:
		drawHBricks(screen, g)
	case Level3:
		drawGrid(screen, g)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w, h := gameDimensions()
	return w, h
}

// function to fill slice of any type
func fillSlice[T any](s []T, value T) []T {
	if s == nil {
		panic("slice cannot be nil")
	}

	for i := range s {
		s[i] = value
	}
	return s
}

func gameDimensions() (width int, height int) {
	return screenWidth, screenHeight
}

func loadImage(path string) *ebiten.Image {
	image, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Error loading image: %v", err)
	}
	return image
}

func newGame() *Game {
	g := Game{}

	initLevel1(&g)
	initLevel3(&g)
	g.level = Level1

	return &g
}

func seedRNG(seed int64) {
	if seed == 0 {
		seed = time.Now().UnixNano() % 10000
	}
	log.Println("Random seed is ", seed)
	rng = rand.New(rand.NewSource(seed))
}

func (g *Game) Update() error {
	g.frameCount++
	switch g.level {
	case Level1:
		return updateLevel1(g)
	case Level3:
		return updateLevel3(g)
	}
	return nil
}
