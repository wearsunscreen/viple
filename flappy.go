package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	fishHeight   = 60
	fishRadius   = (fishHeight / 2) - 5
	fishScale    = 1.0
	fishSpeed    = 3.0
	fishWidth    = 60
	fishX        = 150
	gapHeight    = 100
	lastPipe     = 7
	pipeWidth    = 60
	pipeInterval = 5 * 60
)

var (
	colorPipe           = darkAluminium
	colorPastPipe       = mediumGreen
	colorHitPipe        = mediumScarletRed
	colorInivisiblePipe = mediumSkyBlue
)

type LevelFlappy struct {
	numPipesPast  int
	fishImage     *ebiten.Image
	fishY         float32
	pipes         []*Pipe
	startingFrame int
}

type Pipe struct {
	color         color.RGBA
	completed     bool
	gapY          float32
	startingFrame int
	x             float32
}

func (l *LevelFlappy) addPipe(frameCount int) {
	if l.numPipesPast <= lastPipe {
		p := new(Pipe)
		p.startingFrame = frameCount
		p.gapY = float32(rng.Intn(screenHeight-(fishHeight*2)) + fishHeight/2)
		p.x = screenWidth
		l.pipes = append(l.pipes, p)
		p.color = colorPipe
		p.completed = false
	}
}

func (l *LevelFlappy) gameIsWon() bool {
	return l.numPipesPast > lastPipe
}

func (l *LevelFlappy) CheckPipeCollisions() {
	for _, p := range l.pipes {
		if fishX+fishRadius > p.x && fishX-fishRadius < p.x+pipeWidth {
			if l.fishY-fishRadius < p.gapY || l.fishY+fishRadius > p.gapY+gapHeight {
				p.color = darkScarletRed
				l.numPipesPast = 0
			}
		}
	}
}

func (l *LevelFlappy) Draw(screen *ebiten.Image, frameCount int) {
	// Draw background
	screen.Fill(mediumSkyBlue)

	// top pipe
	for _, p := range l.pipes {
		vector.DrawFilledRect(screen, p.x, 0, pipeWidth, p.gapY, p.color, false)
		vector.DrawFilledRect(screen, p.x, p.gapY+gapHeight, pipeWidth, screenHeight-p.gapY+gapHeight, p.color, false)
	}

	// Draw fish
	//vector.DrawFilledRect(screen, fishX-fishWidth/2, l.fishY-fishHeight/2, fishHeight, fishWidth, l.fishColor, false)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(fishScale, fishScale)
	op.GeoM.Translate(fishX-fishWidth/2, float64(l.fishY)-fishHeight/2)
	screen.DrawImage(l.fishImage, op)
}

func (l *LevelFlappy) Initialize(id LevelID) {
	l.fishY = screenHeight / 2
	l.startingFrame = 0
	l.numPipesPast = 0
	if l.fishImage == nil {
		l.fishImage = loadImage("resources/pufferfish80.png")
	}
}

func (l *LevelFlappy) updateFish() {
	// Update vertical position based on keyboard input
	heldDown := ebiten.IsKeyPressed(ebiten.KeyJ)
	heldUp := ebiten.IsKeyPressed(ebiten.KeyK)
	if heldDown || heldUp {
		if heldDown && !heldUp {
			l.fishY += fishSpeed
		} else if !heldDown && heldUp {
			l.fishY -= fishSpeed
		}
		l.fishY = limitToRange(l.fishY, fishHeight/2, screenHeight-fishHeight/2)
	}
}

func (l *LevelFlappy) updatePipes(frameCount int) {
	// if l.startingFrame == 0 {
	// 	l.startingFrame = frameCount
	// 	l.addPipe(frameCount)
	// }
	if (frameCount+pipeInterval)%pipeInterval == 0 {
		l.addPipe(frameCount)
	}

	// move pipes forward
	for _, p := range l.pipes {
		p.x -= 1
		if p.x < fishX && p.color == colorPipe {
			l.numPipesPast += 1
			p.color = colorPastPipe
		}
	}

	// remove pipe that are off screen
	var newSlice []*Pipe
	for _, p := range l.pipes {
		if p.x > -pipeWidth {
			newSlice = append(newSlice, p)
		}
	}
	l.pipes = newSlice
}

func (l *LevelFlappy) Update(frameCount int) (bool, error) {
	if isCheatKeyPressed() {
		return true, nil
	}

	if l.gameIsWon() {
		return true, nil
	}

	l.updateFish()
	l.updatePipes(frameCount)
	l.CheckPipeCollisions()

	return false, nil
}
