package main

import (
	"flag"
	_ "image/png"
	"log"
	"math/rand"
	"strconv"
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
	LevelIdBricksHL = iota
	LevelIdFlappy
	LevelIdBricksHJKL
	LevelIdGemsVM
)

type Mode int

var (
	rng  *rand.Rand
	keys []ebiten.Key
)

type Game struct {
	frameCount   int
	levelHL      LevelBricksHL
	levelJK      LevelFlappy
	levelVM      LevelGemsVisualMode
	currentLevel LevelID
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
	switch g.currentLevel {
	case LevelIdBricksHL:
		g.levelHL.Draw(screen, g.frameCount)
	case LevelIdFlappy:
		g.levelJK.Draw(screen, g.frameCount)
	case LevelIdBricksHJKL:
		g.levelHL.Draw(screen, g.frameCount)
	case LevelIdGemsVM:
		g.levelVM.Draw(screen, g.frameCount)
	default:
		panic("Unknown game level " + strconv.Itoa(int(g.currentLevel)))
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

	g.levelHL.Initialize()
	g.levelJK.Initialize()
	g.levelVM.Initialize()
	g.currentLevel = LevelIdBricksHL

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
	var levelOver bool
	var err error
	switch g.currentLevel {
	case LevelIdBricksHL:
		levelOver, err = g.levelHL.Update(g.frameCount)
	case LevelIdBricksHJKL:
		levelOver, err = g.levelHL.Update(g.frameCount)
	case LevelIdFlappy:
		levelOver, err = g.levelJK.Update(g.frameCount)
	case LevelIdGemsVM:
		levelOver, err = g.levelVM.Update(g.frameCount)
	}
	if levelOver {
		// advance to next Level if current level has been won
		// bugbug: we don't handle completing the last level cleanly
		g.currentLevel += 1

		switch g.currentLevel {
		case LevelIdBricksHL:
			g.levelHL.Initialize()
		case LevelIdBricksHJKL:
			g.levelHL.level = g.currentLevel
			g.levelHL.Initialize()
		case LevelIdFlappy:
			g.levelJK.Initialize()
		case LevelIdGemsVM:
			g.levelVM.Initialize()
		}

	}
	return err
}
