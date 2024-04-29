package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// interface for things that are drawn to an image, screen
type Drawable interface {
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

func offsetPoint(p, offset Point) Point {
	return Point{p.x + offset.x, p.y + offset.y}
}
