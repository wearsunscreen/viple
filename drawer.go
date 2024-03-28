package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Drawer interface {
	SetZ(z int)
	GetZ() int
	Draw(screen *ebiten.Image, frameCount int)
	AddMover(startFrame int, duration int, from Point, to Point)
}

type Modifier struct {
	startFrame  int
	endFrame    int
	startOffset Point
	endOffset   Point
}

type Point struct {
	x, y int
}

type Square struct {
	color     int
	modifiers []Modifier
	point     Point
	z         int
}

func (square *Square) AddMover(startFrame int, duration int, to Point, from Point) {
	// Bugbug: this will delete all existing modifiers
	// add animation
	mover := Modifier{
		startFrame:  startFrame,
		endFrame:    startFrame + duration,
		startOffset: Point{20, 20},
		endOffset:   Point{40, 40},
	}

	square.modifiers = []Modifier{mover}
}

func applyModifier(m *Modifier, op *ebiten.DrawImageOptions, frameCount int) {
	var completionRatio float64
	completionRatio = float64(m.endFrame-frameCount) / float64(m.endFrame-m.startFrame)
	op.GeoM.Translate(
		float64(m.startOffset.y)+(completionRatio*float64(m.endOffset.y-m.startOffset.y)),
		float64(m.startOffset.x)+(completionRatio*float64(m.endOffset.x-m.startOffset.x)))
}

func (square *Square) Draw(screen *ebiten.Image, frameCount int) {
	//vector.DrawFilledRect(screen, float32(cellSize*x+margin+2), float32(cellSize*y+margin+2), cellSize-4, cellSize-4, gameColors[color], false)
	rect := ebiten.NewImage(cellSize-4, cellSize-4)
	rect.Fill(gameColors[square.color])
	op := &ebiten.DrawImageOptions{}
	for _, m := range square.modifiers {
		applyModifier(&m, op, frameCount)
	}
	op.GeoM.Translate(cellSize*float64(square.point.x)+margin+2, cellSize*float64(square.point.y)+margin+2)
	screen.DrawImage(rect, op)
}

func (square *Square) GetZ() int {
	return square.z
}

func (square *Square) SetZ(z int) {
	square.z = z
}
