package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ballRadius   = 10
	ballSpeedY   = 3.5
	brickWidth   = screenWidth / numBrickCols
	brickHeight  = 50
	numBrickRows = 3
	numBrickCols = 5
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
	brickWidth   int
	brickHeight  int
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

func (level *LevelBricksHL) Draw(screen *ebiten.Image, frameCount int) {
	// Draw background
	screen.Fill(darkCoal)

	// Draw paddle
	vector.DrawFilledRect(screen, level.paddleX, screenHeight-paddleHeight, paddleWidth, paddleHeight, lightAluminium, false)

	// Draw ball
	vector.DrawFilledCircle(screen, level.ballX, level.ballY, ballRadius, lightAluminium, false)

	// Draw bricks with borders
	for y := 0; y < len(level.bricks); y++ {
		for x := 0; x < len(level.bricks[y]); x++ {
			if level.bricks[y][x] {
				// Draw brick
				vector.DrawFilledRect(screen, float32(x*brickWidth), float32(y*brickHeight),
					brickWidth, brickHeight, brightRed, false)
				// Draw border
				vector.StrokeRect(screen, float32(x*brickWidth), float32(y*brickHeight),
					brickWidth, brickHeight, outlineWidth, mediumCoal, false)
			}
		}
	}
}

func (level *LevelBricksHL) Initialize() {
	level.brickWidth = screenWidth / numBrickCols
	level.brickHeight = 50
	level.numBrickRows = 3
	level.numBrickCols = 5

	level.paddleX = screenWidth/2 - paddleWidth/2
	level.paddleY = screenHeight - paddleHeight
	level.ballX = screenWidth / 2
	level.ballY = screenHeight / 3 * 2
	level.ballDX = 2
	level.ballDY = -ballSpeedY

	level.bricks = make([][]bool, numBrickRows)
	for y := range level.bricks {
		level.bricks[y] = make([]bool, numBrickCols)
		fillSlice(level.bricks[y], true)
	}
}

func (level *LevelBricksHL) Update(frameCount int) (bool, error) {
	// cheat to complete level
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		for y := range level.bricks {
			level.bricks[y] = make([]bool, numBrickCols)
			fillSlice(level.bricks[y], false)
		}
	}

	// Update paddle position based on keyboard input
	heldLeft := ebiten.IsKeyPressed(ebiten.KeyH)
	heldRight := ebiten.IsKeyPressed(ebiten.KeyL)
	if heldLeft || heldRight {
		if heldLeft && !heldRight {
			level.paddleX -= paddleSpeed
		} else if !heldLeft && heldRight {
			level.paddleX += paddleSpeed
		}
	}

	// Clamp paddle movement within screen bounds
	if level.paddleX < 0 {
		level.paddleX = 0
	} else if level.paddleX > screenWidth-paddleWidth {
		level.paddleX = screenWidth - paddleWidth
	}

	// Check for wall collisions
	if level.ballX < 0 || level.ballX > screenWidth-ballRadius {
		level.ballDX *= -1
	}
	if level.ballY < 0 {
		level.ballDY *= -1
	}

	// Check for ball off bottom of screen
	if level.ballY+ballRadius > screenHeight {
		level.Initialize()
	}

	// Check for paddle collision
	if level.ballY+ballRadius > screenHeight-paddleHeight &&
		level.ballX >= level.paddleX && level.ballX <= level.paddleX+paddleWidth {
		level.ballDY *= -1

		// modify angle depending on where the ball hits the paddle
		ratio := (level.ballX - level.paddleX) / paddleWidth
		level.ballDX = ratio*4 - 2
		//log.Printf("ratio %f, dx is %f", ratio, level.ballDX)

		PlaySound(paddleOgg)
	}

	// Check for brick collision
	for y, row := range level.bricks {
		for x, brick := range row {
			if brick {
				if level.ballX > float32(x*brickWidth) && level.ballX < float32((x+1)*brickWidth) &&
					level.ballY > float32(y*brickHeight) && level.ballY < float32((y+1)*brickHeight) {
					level.bricks[y][x] = false
					level.ballDY *= -1
					PlaySound(brickOgg)

				}
			}
		}
	}

	// Update ball position
	if anyIn2DSlice(level.bricks) {
		level.ballX += level.ballDX
		level.ballY += level.ballDY
	} else {
		return true, nil
	}

	return false, nil
}
