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

type Level interface {
	Draw(screen *ebiten.Image, frameCount int)
	Initialize()
	// update every frame, return true if level is complete
	Update(g *Game) (bool, error)
}

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
	rng  *rand.Rand
	keys []ebiten.Key
)

type Game struct {
	frameCount int
	level1     Level1Data
	level3     Level3Data
	level      LevelID
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
		g.level1.Draw(screen, g.frameCount)
	case Level3:
		g.level3.Draw(screen, g.frameCount)
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

	g.level1.Initialize()
	g.level3.Initialize()
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
	var b bool
	var err error
	switch g.level {
	case Level1:
		b, err = g.level1.Update(g.frameCount)
	case Level3:
		b, err = g.level3.Update(g.frameCount)
	}
	if b {
		g.level += 1
	}
	return err
}
