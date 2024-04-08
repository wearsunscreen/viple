package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// interface for things that are drawn to an image, screen
type Drawer interface {
	SetZ(z int)
	GetZ() int
	Draw(screen *ebiten.Image, frameCount int)
	AddMover(startFrame int, duration int, from Point, to Point)
}

type Mover struct {
	startFrame int
	endFrame   int
	startPoint Point
	endPoint   Point
}

type Point struct {
	x, y int
}
