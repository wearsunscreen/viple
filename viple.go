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
	player     *AudioPlayer
	l3         Level3
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

	drawGrid(screen, g)
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

func newGame() *Game {
	g := Game{}

	initLevel3(&g)

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

func (g *Game) Update() error {
	g.frameCount++
	return updateLevel3(g)
}
