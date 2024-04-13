package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	birdHeight = 30
	birdSpeed  = 3.0
	birdWidth  = 30
	birdX      = 150
	gapHeight  = 100
	pipeWidth  = 60
)

type LevelFlappy struct {
	birdColor     color.RGBA
	birdY         float32
	gapY          float32
	pipeColor     color.RGBA
	pipeX         float32
	startingFrame int
}

func (l *LevelFlappy) Draw(screen *ebiten.Image, frameCount int) {
	// Draw background
	screen.Fill(mediumSkyBlue)

	// Draw bird
	vector.DrawFilledRect(screen, birdX, l.birdY, birdHeight, birdWidth, l.birdColor, false)

	// top pipe
	vector.DrawFilledRect(screen, l.pipeX, 0, pipeWidth, l.gapY, l.pipeColor, false)
	vector.DrawFilledRect(screen, l.pipeX, l.gapY+gapHeight, pipeWidth, screenHeight-l.gapY+gapHeight, l.pipeColor, false)
}

func (l *LevelFlappy) Initialize() {
	l.birdY = screenHeight / 2
	l.pipeX = screenWidth / 4 * 3
	l.gapY = 100
	l.birdColor = mediumButter
	l.pipeColor = darkAluminium
	l.startingFrame = 0
}

func (l *LevelFlappy) Update(frameCount int) (bool, error) {
	if l.startingFrame == 0 {
		l.startingFrame = frameCount
	}

	// Update paddle vertical position based on keyboard input
	heldDown := ebiten.IsKeyPressed(ebiten.KeyJ)
	heldUp := ebiten.IsKeyPressed(ebiten.KeyK)
	if heldDown || heldUp {
		if heldDown && !heldUp {
			l.birdY += birdSpeed
		} else if !heldDown && heldUp {
			l.birdY -= birdSpeed
		}
	}

	// move pipes forward
	l.pipeX = float32(screenWidth - (frameCount - l.startingFrame))

	return false, nil
}
