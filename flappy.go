package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	fishHeight   = 60
	fishScale    = 1.0
	fishSpeed    = 3.0
	fishWidth    = 60
	fishX        = 150
	gapHeight    = 100
	pipeWidth    = 60
	pipeInterval = 5 * 60
)

type LevelFlappy struct {
	fishColor     color.RGBA
	fishImage     *ebiten.Image
	fishY         float32
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

	// Draw fish
	vector.DrawFilledRect(screen, fishX-fishWidth/2, l.fishY-fishHeight/2, fishHeight, fishWidth, l.fishColor, false)

	// top pipe
	for _, p := range l.pipes {
		vector.DrawFilledRect(screen, p.x, 0, pipeWidth, p.gapY, l.pipeColor, false)
		vector.DrawFilledRect(screen, p.x, p.gapY+gapHeight, pipeWidth, screenHeight-p.gapY+gapHeight, l.pipeColor, false)
	}
	// draw fish
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(fishScale, fishScale)
	op.GeoM.Translate(fishX-fishWidth/2, float64(l.fishY)-fishHeight/2)
	op.GeoM.Translate(0, 0)
	screen.DrawImage(l.fishImage, op)

}

func (l *LevelFlappy) Initialize() {
	l.fishY = screenHeight / 2
	l.fishColor = mediumButter
	l.pipeColor = darkAluminium
	l.startingFrame = 0
	if l.fishImage == nil {
		l.fishImage = loadImage("resources/pufferfish80.png")

	}
}

func (l *LevelFlappy) Update(frameCount int) (bool, error) {
	if IsCheatKeyPressed() {
		return true, nil
	}

	if l.startingFrame == 0 {
		l.startingFrame = frameCount
	}
	if (frameCount+pipeInterval)%pipeInterval == 0 {
		p := new(Pipe)
		p.startingFrame = frameCount
		p.gapY = float32(rng.Intn(screenHeight-(fishHeight*2)) + fishHeight/2)
		p.x = screenWidth
		l.pipes = append(l.pipes, p)
	}

	// Update paddle vertical position based on keyboard input
	heldDown := ebiten.IsKeyPressed(ebiten.KeyJ)
	heldUp := ebiten.IsKeyPressed(ebiten.KeyK)
	if heldDown || heldUp {
		if heldDown && !heldUp {
			l.fishY += fishSpeed
		} else if !heldDown && heldUp {
			l.fishY -= fishSpeed
		}
		l.fishY = limitToRange(l.fishY, screenHeight-fishHeight/2)
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
