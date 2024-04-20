package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ballRadius     = 10
	ballSpeedY     = 3.5
	outlineWidth   = 2
	paddlesXWidth  = 100
	paddlesXHeight = 20
	paddlesYWidth  = paddlesXHeight
	paddlesYHeight = paddlesXWidth
	paddleSpeed    = 5
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
	minimumSpeed float64
	numBrickRows int
	numBrickCols int
	paddlesX     float32
	paddlesY     float32
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

func (l *LevelBricksHL) CheckBrickCollisions() {
	// Check for brick collision
	for y, row := range l.bricks {
		for x, brick := range row {
			if brick {
				if l.ballX > float32(x*l.brickWidth+l.brickLeft) && l.ballX < float32((x+1)*l.brickWidth+l.brickLeft) &&
					l.ballY > float32(y*l.brickHeight+l.brickTop) && l.ballY < float32((y+1)*l.brickHeight+l.brickTop) {

					// we have a collision,clear the brick
					l.bricks[y][x] = false

					//let's find the closest edge to determine which direction to change
					leftD := math.Abs(float64(l.ballX - float32(x*l.brickWidth+l.brickLeft)))
					rightD := math.Abs(float64(l.ballX - float32((x+1)*l.brickWidth+l.brickLeft)))
					topD := math.Abs(float64(l.ballY - float32(y*l.brickHeight+l.brickTop)))
					bottomD := math.Abs(float64(l.ballY - float32((y+1)*l.brickHeight+l.brickTop)))
					minD := math.Min(leftD, rightD)
					minD = math.Min(minD, topD)
					minD = math.Min(minD, bottomD)
					switch {
					case minD == leftD:
						l.ballDX *= -1
					case minD == rightD:
						l.ballDX *= -1
					case minD == topD:
						l.ballDY *= -1
					case minD == bottomD:
						l.ballDY *= -1
					default:
						panic("Bad collision calculation")
					}

					PlaySound(brickOgg)
				}
			}
		}
	}
}

func (l *LevelBricksHL) CheckPaddleCollisions() {
	// Check for bottom paddle collision
	if l.ballY+ballRadius > screenHeight-paddlesXHeight &&
		l.ballX >= l.paddlesX && l.ballX <= l.paddlesX+paddlesXWidth {
		l.ballDY *= -1

		// modify angle depending on where the ball hits the paddle
		ratio := (l.ballX - l.paddlesX) / paddlesXWidth
		l.ballDX = (ratio*4 - 2) * 1.5
		//log.Printf("ratio %f, dx is %f, dy is %f", ratio, l.ballDX, l.ballDY)

		PlaySound(paddleOgg)
	}

	if l.level == LevelIdBricksHJKL {
		// check top paddle collision
		if l.ballY-ballRadius < paddlesXHeight &&
			l.ballX >= l.paddlesX && l.ballX <= l.paddlesX+paddlesXWidth {
			l.ballDY *= -1

			// modify angle depending on where the ball hits the paddle
			ratio := (l.ballX - l.paddlesX) / paddlesXWidth
			l.ballDX = ratio*4 - 2
			//log.Printf("ratio %f, dx is %f, dy is %f", ratio, l.ballDX, l.ballDY)

			PlaySound(paddleOgg)
		}

		// check left paddle
		if l.ballX-ballRadius < paddlesYWidth &&
			l.ballY >= l.paddlesY && l.ballY <= l.paddlesY+paddlesYHeight {
			l.ballDX *= -1

			// modify angle depending on where the ball hits the paddle
			ratio := (l.ballY - l.paddlesY) / paddlesYHeight
			l.ballDY = ratio*4 - 2
			//log.Printf("ratio %f, dx is %f, dy is %f", ratio, l.ballDX, l.ballDY)

			PlaySound(paddleOgg)
		}

		// check right paddle
		if l.ballX+ballRadius > screenWidth-paddlesYWidth &&
			l.ballY >= l.paddlesY && l.ballY <= l.paddlesY+paddlesYHeight {
			l.ballDX *= -1

			// modify angle depending on where the ball hits the paddle
			ratio := (l.ballY - l.paddlesY) / paddlesYHeight
			l.ballDY = ratio*4 - 2
			//log.Printf("ratio %f, dx is %f, dy is %f", ratio, l.ballDX, l.ballDY)

			PlaySound(paddleOgg)
		}
	}

	if l.ballDX != 0.0 {
		// ensure the ball speed is not too slow
		for math.Abs(float64(l.ballDX))+math.Abs(float64(l.ballDY)) < l.minimumSpeed {
			l.ballDX *= 1.1
			l.ballDY *= 1.1
		}
	}
}

func (l *LevelBricksHL) CheckWallCollisions() {
	// Check for wall collisions
	if l.level == LevelIdBricksHL {
		if l.ballX < 0 || l.ballX > screenWidth-ballRadius {
			l.ballDX *= -1
		}
		if l.ballY < 0 {
			l.ballDY *= -1
		}
	}

	// Check for ball off bottom of screen
	if l.ballY+ballRadius > screenHeight {
		l.Initialize(l.level)
	}

	if l.level == LevelIdBricksHJKL {
		// Check for ball off top of screen
		if l.ballY+ballRadius < 0 {
			l.Initialize(l.level)
		}
		// Check for ball off left of screen
		if l.ballX+ballRadius < 0 {
			l.Initialize(l.level)
		}
		// Check for ball off right of screen
		if l.ballX+ballRadius > screenWidth {
			l.Initialize(l.level)
		}
	}
}

func (l *LevelBricksHL) Draw(screen *ebiten.Image, frameCount int) {
	// Draw background
	screen.Fill(darkCoal)

	// Draw paddle
	vector.DrawFilledRect(screen, l.paddlesX, screenHeight-paddlesXHeight, paddlesXWidth, paddlesXHeight, darkAluminium, false)
	if l.level == LevelIdBricksHJKL {
		vector.DrawFilledRect(screen, l.paddlesX, 0, paddlesXWidth, paddlesXHeight, darkAluminium, false)
		vector.DrawFilledRect(screen, 0, l.paddlesY, paddlesYWidth, paddlesYHeight, darkAluminium, false)
		vector.DrawFilledRect(screen, screenWidth-paddlesYWidth, l.paddlesY, paddlesYWidth, paddlesYHeight, darkAluminium, false)
	}
	// Draw ball
	vector.DrawFilledCircle(screen, l.ballX, l.ballY, ballRadius, darkAluminium, false)

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

func (l *LevelBricksHL) Initialize(id LevelID) {
	l.level = id
	if l.level == LevelIdBricksHL {
		l.numBrickRows = 3
		l.numBrickCols = 5
		l.brickWidth = screenWidth / l.numBrickCols
		l.brickHeight = 50
		l.brickWidth = screenWidth / l.numBrickCols
		l.brickHeight = 50

		l.paddlesX = screenWidth/2 - paddlesXWidth/2
		l.ballX = screenWidth / 2
		l.ballY = screenHeight / 3 * 2
	} else {
		l.brickWidth = 50
		l.brickHeight = 50
		l.numBrickRows = 3
		l.numBrickCols = 5
		l.brickLeft = (screenWidth - l.brickWidth*l.numBrickCols) / 2
		l.brickTop = (screenHeight - l.brickHeight*l.numBrickRows) / 2

		l.paddlesX = screenWidth/2 - paddlesXWidth/2
		l.paddlesY = screenHeight/2 - paddlesYHeight/2
		l.ballY = float32(l.brickTop) - ballRadius
		l.ballX = float32(l.brickLeft + (l.brickHeight * (l.numBrickRows / 2)))
	}

	l.ballDX = 0
	l.ballDY = 0
	l.minimumSpeed = 3.0
	l.bricks = make([][]bool, l.numBrickRows)
	for y := range l.bricks {
		l.bricks[y] = make([]bool, l.numBrickCols)
		fillSlice(l.bricks[y], true)
	}
}

func (l *LevelBricksHL) IntroText() string {
	return `In the first level you will 
	learn to move left and right 
	by pressing H and K keys.`
}

func (l *LevelBricksHL) TitleText() string {
	return `Welcome to Viple
	VI Play to Learn.`
}

func (l *LevelBricksHL) Update(frameCount int) (bool, error) {
	if isCheatKeyPressed() {
		return true, nil
	}
	l.UpdateBallPosition()
	l.UpdatePaddlePositions()
	l.CheckWallCollisions()
	l.CheckBrickCollisions()
	l.CheckPaddleCollisions()

	// check for end of level
	if !anyIn2DSlice(l.bricks) {
		return true, nil
	}

	return false, nil
}

func (l *LevelBricksHL) initBallMovement() {
	if l.ballDX == 0 {
		if l.level == LevelIdBricksHL {
			l.ballDX = 0.1
			l.ballDY = -ballSpeedY
		} else {
			l.ballDX = -2.0
			l.ballDY = 0.1

		}
	}
}

func (l *LevelBricksHL) UpdateBallPosition() {
	if anyIn2DSlice(l.bricks) {
		l.ballX += l.ballDX
		l.ballY += l.ballDY
	}
}

func (l *LevelBricksHL) UpdatePaddlePositions() {
	// Update paddle horizontal position based on keyboard input
	heldLeft := ebiten.IsKeyPressed(ebiten.KeyH)
	heldRight := ebiten.IsKeyPressed(ebiten.KeyL)
	if heldLeft || heldRight {
		if heldLeft && !heldRight {
			l.paddlesX -= paddleSpeed
		} else if !heldLeft && heldRight {
			l.paddlesX += paddleSpeed
		}

		// the level waits for the first paddle move before starting the ball
		l.initBallMovement()
	}

	if l.level == LevelIdBricksHJKL {
		// Update paddle vertical position based on keyboard input
		heldDown := ebiten.IsKeyPressed(ebiten.KeyJ)
		heldUp := ebiten.IsKeyPressed(ebiten.KeyK)
		if heldDown || heldUp {
			if heldDown && !heldUp {
				l.paddlesY += paddleSpeed
			} else if !heldDown && heldUp {
				l.paddlesY -= paddleSpeed
			}
			l.initBallMovement()
		}
	}
	// limit paddle movement within screen bounds
	l.paddlesX = limitToRange(l.paddlesX, 0, screenWidth-paddlesXWidth)
	if l.level == LevelIdBricksHJKL {
		l.paddlesY = limitToRange(l.paddlesY, 0, screenHeight-paddlesYHeight)
	}

}
