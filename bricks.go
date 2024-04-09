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

type Level1Data struct {
	bricks  [][]bool
	ballDX  float32
	ballDY  float32
	ballX   float32
	ballY   float32
	paddleX float32
	paddleY float32
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

func (level1 *Level1Data) Draw(screen *ebiten.Image, frameCount int) {
	// Draw background
	screen.Fill(darkCoal)

	// Draw paddle
	vector.DrawFilledRect(screen, level1.paddleX, screenHeight-paddleHeight, paddleWidth, paddleHeight, lightAluminium, false)

	// Draw ball
	vector.DrawFilledCircle(screen, level1.ballX, level1.ballY, ballRadius, lightAluminium, false)

	// Draw bricks with borders
	for y := 0; y < len(level1.bricks); y++ {
		for x := 0; x < len(level1.bricks[y]); x++ {
			if level1.bricks[y][x] {
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

func (level1 *Level1Data) Initialize() {
	level1.paddleX = screenWidth/2 - paddleWidth/2
	level1.paddleY = screenHeight - paddleHeight
	level1.ballX = screenWidth / 2
	level1.ballY = screenHeight / 3 * 2
	level1.ballDX = 2
	level1.ballDY = -ballSpeedY

	level1.bricks = make([][]bool, numBrickRows)
	for y := range level1.bricks {
		level1.bricks[y] = make([]bool, numBrickCols)
		fillSlice(level1.bricks[y], true)
	}
}

func (level1 *Level1Data) Update(frameCount int) (bool, error) {
	if cheat := ebiten.IsKeyPressed(ebiten.KeyC); cheat {
		for y := range level1.bricks {
			level1.bricks[y] = make([]bool, numBrickCols)
			fillSlice(level1.bricks[y], false)
		}
	}

	// Update paddle position based on keyboard input
	heldLeft := ebiten.IsKeyPressed(ebiten.KeyH)
	heldRight := ebiten.IsKeyPressed(ebiten.KeyL)
	if heldLeft || heldRight {
		if heldLeft && !heldRight {
			level1.paddleX -= paddleSpeed
		} else if !heldLeft && heldRight {
			level1.paddleX += paddleSpeed
		}
	}

	// Clamp paddle movement within screen bounds
	if level1.paddleX < 0 {
		level1.paddleX = 0
	} else if level1.paddleX > screenWidth-paddleWidth {
		level1.paddleX = screenWidth - paddleWidth
	}

	// Check for wall collisions
	if level1.ballX < 0 || level1.ballX > screenWidth-ballRadius {
		level1.ballDX *= -1
	}
	if level1.ballY < 0 {
		level1.ballDY *= -1
	}

	// Check for ball off bottom of screen
	if level1.ballY+ballRadius > screenHeight {
		level1.Initialize()
	}

	// Check for paddle collision
	if level1.ballY+ballRadius > screenHeight-paddleHeight &&
		level1.ballX >= level1.paddleX && level1.ballX <= level1.paddleX+paddleWidth {
		level1.ballDY *= -1

		// modify angle depending on where the ball hits the paddle
		ratio := (level1.ballX - level1.paddleX) / paddleWidth
		level1.ballDX = ratio*4 - 2
		//log.Printf("ratio %f, dx is %f", ratio, level1.ballDX)
	}

	// Check for brick collision (simplified for brevity)
	for y, row := range level1.bricks {
		for x, brick := range row {
			if brick {
				if level1.ballX > float32(x*brickWidth) && level1.ballX < float32((x+1)*brickWidth) &&
					level1.ballY > float32(y*brickHeight) && level1.ballY < float32((y+1)*brickHeight) {
					level1.bricks[y][x] = false
					level1.ballDY *= -1
				}
			}
		}
	}

	// Update ball position
	if anyIn2DSlice(level1.bricks) {
		level1.ballX += level1.ballDX
		level1.ballY += level1.ballDY
	} else {
		return true, nil
	}

	return false, nil
}
