package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

type Square struct {
	color int
	mover *Mover
	point Point
	z     int
}

func (square *Square) AddMover(startFrame int, duration int, from Point, to Point) {
	// add animation
	mover := new(Mover)

	mover.startFrame = startFrame
	mover.endFrame = startFrame + duration
	mover.startPoint = from
	mover.endPoint = to

	square.mover = mover
}

func drawGem(screen *ebiten.Image, image *ebiten.Image, p Point) {
}

// convert the x,y of the square into screen coordinates
func SquareToScreenPoint(squareXY Point) Point {
	return Point{
		cellSize*squareXY.x + margin + 2,
		cellSize*squareXY.y + margin + 2,
	}
}

func applyMover(mover *Mover, op *ebiten.DrawImageOptions, frameCount int) {
	completionRatio := 1 - float64(mover.endFrame-frameCount)/float64(mover.endFrame-mover.startFrame)
	startPosition := SquareToScreenPoint(mover.startPoint)
	endPosition := SquareToScreenPoint(mover.endPoint)
	op.GeoM.Translate(
		float64(startPosition.x)+(completionRatio*float64(endPosition.x-startPosition.x)),
		float64(startPosition.y)+(completionRatio*float64(endPosition.y-startPosition.y)))
}

func (square *Square) DrawBackground(screen *ebiten.Image, color color.Color) {
	pos := SquareToScreenPoint(square.point)
	vector.DrawFilledRect(screen, float32(pos.x), float32(pos.y), cellSize-4, cellSize-4, color, false)
	/*
		rect := ebiten.NewImage(cellSize-4, cellSize-4)
		rect.Fill(gameColors[square.color])
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(rect, op)
	*/
}

func (square *Square) DrawGem(screen *ebiten.Image, gemImage *ebiten.Image, frameCount int) {
	if square.color >= 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(gemScale, gemScale)
		if square.mover != nil {
			applyMover(square.mover, op, frameCount)
		} else {
			pos := SquareToScreenPoint(square.point)
			op.GeoM.Translate(float64(pos.x), float64(pos.y))
		}
		screen.DrawImage(gemImage, op)
		/*
			//vector.DrawFilledRect(screen, float32(cellSize*x+margin+2), float32(cellSize*y+margin+2), cellSize-4, cellSize-4, gameColors[color], false)
			rect := ebiten.NewImage(cellSize-4, cellSize-4)
			rect.Fill(gameColors[square.color])
			op := &ebiten.DrawImageOptions{}
			screen.DrawImage(rect, op)
		*/
	}
}

func (square *Square) GetZ() int {
	return square.z
}

func (square *Square) SetZ(z int) {
	square.z = z
}
