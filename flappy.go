package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	birdHeight   = 30
	birdSpeed    = 3.0
	birdWidth    = 30
	birdX        = 150
	gapHeight    = 100
	pipeWidth    = 60
	pipeInterval = 5 * 60
)

type LevelFlappy struct {
	birdColor     color.RGBA
	birdY         float32
	pipeColor     color.RGBA
	pipes         []*Pipe
	startingFrame int
}

type Pipe struct {
	gapY          float32
	startingFrame int
	x             float32
}

func (l *LevelFlappy) Draw(screen *ebiten.Image, frameCount int) {
	// Draw background
	screen.Fill(mediumSkyBlue)

	// Draw bird
	vector.DrawFilledRect(screen, birdX, l.birdY, birdHeight, birdWidth, l.birdColor, false)

	// top pipe
	for _, p := range l.pipes {
		vector.DrawFilledRect(screen, p.x, 0, pipeWidth, p.gapY, l.pipeColor, false)
		vector.DrawFilledRect(screen, p.x, p.gapY+gapHeight, pipeWidth, screenHeight-p.gapY+gapHeight, l.pipeColor, false)
	}
}

func (l *LevelFlappy) Initialize() {
	l.birdY = screenHeight / 2
	l.birdColor = mediumButter
	l.pipeColor = darkAluminium
	l.startingFrame = 0
}

func (l *LevelFlappy) Update(frameCount int) (bool, error) {
	if l.startingFrame == 0 {
		l.startingFrame = frameCount
	}
	if (frameCount+pipeInterval)%pipeInterval == 0 {
		p := new(Pipe)
		p.startingFrame = frameCount
		p.gapY = float32(rng.Intn(screenHeight - (gapHeight / 2) + gapHeight/4))
		p.x = screenWidth
		l.pipes = append(l.pipes, p)
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
	for _, p := range l.pipes {
		p.x -= 1
	}

	// remove pipe that are off screen
	var newSlice []*Pipe
	for _, p := range l.pipes {
		if p.x > -pipeWidth {
			newSlice = append(newSlice, p)
		}
	}
	l.pipes = newSlice

	return false, nil
}
