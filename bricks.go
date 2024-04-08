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

func drawHBricks(screen *ebiten.Image, g *Game) {
	// Draw background
	screen.Fill(darkCoal)

	// Draw paddle
	vector.DrawFilledRect(screen, g.l1.paddleX, screenHeight-paddleHeight, paddleWidth, paddleHeight, lightAluminium, false)

	// Draw ball
	vector.DrawFilledCircle(screen, g.l1.ballX, g.l1.ballY, ballRadius, lightAluminium, false)

	// Draw bricks with borders
	for y := 0; y < len(g.l1.bricks); y++ {
		for x := 0; x < len(g.l1.bricks[y]); x++ {
			if g.l1.bricks[y][x] {
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

func initLevel1(g *Game) {
	g.l1.paddleX = screenWidth/2 - paddleWidth/2
	g.l1.paddleY = screenHeight - paddleHeight
	g.l1.ballX = screenWidth / 2
	g.l1.ballY = screenHeight / 3 * 2
	g.l1.ballDX = 2
	g.l1.ballDY = -ballSpeedY

	g.l1.bricks = make([][]bool, numBrickRows)
	for y := range g.l1.bricks {
		g.l1.bricks[y] = make([]bool, numBrickCols)
		fillSlice(g.l1.bricks[y], true)
	}
}

func updateLevel1(g *Game) error {
	// Update paddle position based on keyboard input
	heldLeft := ebiten.IsKeyPressed(ebiten.KeyH)
	heldRight := ebiten.IsKeyPressed(ebiten.KeyL)
	if heldLeft || heldRight {
		if heldLeft && !heldRight {
			g.l1.paddleX -= paddleSpeed
		} else if !heldLeft && heldRight {
			g.l1.paddleX += paddleSpeed
		}
	}

	// Clamp paddle movement within screen bounds
	if g.l1.paddleX < 0 {
		g.l1.paddleX = 0
	} else if g.l1.paddleX > screenWidth-paddleWidth {
		g.l1.paddleX = screenWidth - paddleWidth
	}

	// Check for wall collisions
	if g.l1.ballX < 0 || g.l1.ballX > screenWidth-ballRadius {
		g.l1.ballDX *= -1
	}
	if g.l1.ballY < 0 {
		g.l1.ballDY *= -1
	}

	// Check for ball off bottom of screen
	if g.l1.ballY+ballRadius > screenHeight {
		initLevel1(g)
	}

	// Check for paddle collision
	if g.l1.ballY+ballRadius > screenHeight-paddleHeight &&
		g.l1.ballX >= g.l1.paddleX && g.l1.ballX <= g.l1.paddleX+paddleWidth {
		g.l1.ballDY *= -1

		// modify angle depending on where the ball hits the paddle
		ratio := (g.l1.ballX - g.l1.paddleX) / paddleWidth
		g.l1.ballDX = ratio*4 - 2
		//log.Printf("ratio %f, dx is %f", ratio, g.l1.ballDX)
	}

	// Check for brick collision (simplified for brevity)
	for y, row := range g.l1.bricks {
		for x, brick := range row {
			if brick {
				if g.l1.ballX > float32(x*brickWidth) && g.l1.ballX < float32((x+1)*brickWidth) &&
					g.l1.ballY > float32(y*brickHeight) && g.l1.ballY < float32((y+1)*brickHeight) {
					g.l1.bricks[y][x] = false
					g.l1.ballDY *= -1
				}
			}
		}
	}

	// Update ball position
	if anyIn2DSlice(g.l1.bricks) {
		g.l1.ballX += g.l1.ballDX
		g.l1.ballY += g.l1.ballDY
	}

	return nil
}
