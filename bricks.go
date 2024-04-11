package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ballRadius   = 10
	ballSpeedY   = 3.5
	outlineWidth = 2
	paddleWidth  = 100
	paddleHeight = 20
	paddleSpeed  = 5
)

type LevelBricksHL struct {
	ballDX       float32
	ballDY       float32
	ballX        float32
	ballY        float32
	bricks       [][]bool
	brickHeight  int
	brickWidth   int
	brickLeft    int
	brickTop     int
	level        LevelID
	numBrickRows int
	numBrickCols int
	paddleX      float32
	paddleY      float32
}

// return true is any value in the 2D slice is true
func anyIn2DSlice(bools [][]bool) bool {
	for _, row := range bools {
		if b := anyInSlice(row); b {
			return true
		}
	}
	return false
}

// return true is any value in the slice is true
func anyInSlice(bools []bool) bool {
	for _, b := range bools {
		if b {
			return true
		}
	}
	return false
}

func (l *LevelBricksHL) Draw(screen *ebiten.Image, frameCount int) {
	// Draw background
	screen.Fill(darkCoal)

	// Draw paddle
	vector.DrawFilledRect(screen, l.paddleX, screenHeight-paddleHeight, paddleWidth, paddleHeight, lightAluminium, false)

	// Draw ball
	vector.DrawFilledCircle(screen, l.ballX, l.ballY, ballRadius, lightAluminium, false)

	// Draw bricks with borders
	for y := 0; y < len(l.bricks); y++ {
		for x := 0; x < len(l.bricks[y]); x++ {
			if l.bricks[y][x] {
				// Draw brick
				vector.DrawFilledRect(screen, float32(x*l.brickWidth+l.brickLeft), float32(y*l.brickHeight+l.brickTop),
					float32(l.brickWidth), float32(l.brickHeight), brightRed, false)
				// Draw border
				vector.StrokeRect(screen, float32(x*l.brickWidth+l.brickLeft), float32(y*l.brickHeight+l.brickTop),
					float32(l.brickWidth), float32(l.brickHeight), outlineWidth, mediumCoal, false)
			}
		}
	}
}

func (l *LevelBricksHL) Initialize() {
	if l.level == LevelIdBricksHL {
		l.numBrickRows = 3
		l.numBrickCols = 5
		l.brickWidth = screenWidth / l.numBrickCols
		l.brickHeight = 50
		l.brickWidth = screenWidth / l.numBrickCols
		l.brickHeight = 50

		l.paddleX = screenWidth/2 - paddleWidth/2
		l.paddleY = screenHeight - paddleHeight
		l.ballX = screenWidth / 2
		l.ballY = screenHeight / 3 * 2
		l.ballDX = 2
		l.ballDY = -ballSpeedY
	} else {
		l.brickWidth = 50
		l.brickHeight = 50
		l.numBrickRows = 3
		l.numBrickCols = 5
		l.brickLeft = (screenWidth - l.brickWidth*l.numBrickCols) / 2
		l.brickTop = (screenHeight - l.brickHeight*l.numBrickRows) / 2

		l.paddleX = screenWidth/2 - paddleWidth/2
		l.paddleY = screenHeight - paddleHeight
		l.ballX = screenWidth / 2
		l.ballY = screenHeight / 3 * 2
		l.ballDX = 2
		l.ballDY = -ballSpeedY
	}

	l.bricks = make([][]bool, l.numBrickRows)
	for y := range l.bricks {
		l.bricks[y] = make([]bool, l.numBrickCols)
		fillSlice(l.bricks[y], true)
	}
}

func (l *LevelBricksHL) Update(frameCount int) (bool, error) {
	// cheat to complete level
	cheatKey := ebiten.KeyA
	if l.level == LevelIdBricksHJKL {
		cheatKey = ebiten.KeyB
	}
	if ebiten.IsKeyPressed(cheatKey) {
		for y := range l.bricks {
			l.bricks[y] = make([]bool, l.numBrickCols)
			fillSlice(l.bricks[y], false)
		}
	}

	// Update paddle position based on keyboard input
	heldLeft := ebiten.IsKeyPressed(ebiten.KeyH)
	heldRight := ebiten.IsKeyPressed(ebiten.KeyL)
	if heldLeft || heldRight {
		if heldLeft && !heldRight {
			l.paddleX -= paddleSpeed
		} else if !heldLeft && heldRight {
			l.paddleX += paddleSpeed
		}
	}

	// Clamp paddle movement within screen bounds
	if l.paddleX < 0 {
		l.paddleX = 0
	} else if l.paddleX > screenWidth-paddleWidth {
		l.paddleX = screenWidth - paddleWidth
	}

	// Check for wall collisions
	if l.ballX < 0 || l.ballX > screenWidth-ballRadius {
		l.ballDX *= -1
	}
	if l.ballY < 0 {
		l.ballDY *= -1
	}

	// Check for ball off bottom of screen
	if l.ballY+ballRadius > screenHeight {
		l.Initialize()
	}

	// Check for paddle collision
	if l.ballY+ballRadius > screenHeight-paddleHeight &&
		l.ballX >= l.paddleX && l.ballX <= l.paddleX+paddleWidth {
		l.ballDY *= -1

		// modify angle depending on where the ball hits the paddle
		ratio := (l.ballX - l.paddleX) / paddleWidth
		l.ballDX = ratio*4 - 2
		//log.Printf("ratio %f, dx is %f", ratio, l.ballDX)

		PlaySound(paddleOgg)
	}

	// Check for brick collision
	for y, row := range l.bricks {
		for x, brick := range row {
			if brick {
				if l.ballX > float32(x*l.brickWidth+l.brickLeft) && l.ballX < float32((x+1)*l.brickWidth+l.brickLeft) &&
					l.ballY > float32(y*l.brickHeight+l.brickTop) && l.ballY < float32((y+1)*l.brickHeight+l.brickTop) {
					l.bricks[y][x] = false
					l.ballDY *= -1
					PlaySound(brickOgg)

				}
			}
		}
	}

	// Update ball position
	if anyIn2DSlice(l.bricks) {
		l.ballX += l.ballDX
		l.ballY += l.ballDY
	} else {
		return true, nil
	}

	return false, nil
}
