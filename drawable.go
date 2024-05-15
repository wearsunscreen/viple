package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// interface for things that are drawn to an image, screen
type Drawable interface {
	Draw(screen *ebiten.Image, frameCount int)
	AddMover(startFrame int, duration int, from Coord, to Coord)
}
