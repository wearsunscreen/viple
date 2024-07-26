package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

/*
 * LevelSnake implements two levels. The first (LevelIdSnake) provides the play practice of the
 * navigation keys, h,j,k,l. The second (LevelIDInsertMode) provides the play practice of the navigation
 * and entering and exiting insert mode.
 */
type Direction int

const (
	north Direction = iota
	east
	south
	west
)

const (
	lengthForWin = 25
)

type Snake struct {
	body      []Coord
	direction Direction
}

type LevelSnake struct {
	level  LevelID
	food   Coord
	score  int
	snake  *Snake
	viMode VIMode
}

var (
	gridWidth  = screenWidth / size
	gridHeight = screenHeight / size
	snakeColor = color.RGBA{R: 0x20, G: 0xFF, B: 0x20, A: 0xFF}
	foodColor  = color.RGBA{0xcc, 0x00, 0x00, 0xa0} // mediumScarletRed
	size       = 40                                 // size of each square in the grid
)

func (l *LevelSnake) Draw(screen *ebiten.Image, frameCount int) {
	sc := snakeColor
	if l.level == LevelIdInsertMode && l.viMode == InsertMode {
		// show different color if able to eat food in insert level
		sc = color.RGBA{R: 0xBB, G: 0x80, B: 0x80, A: 0xFF}
	}

	// Draw the snake
	green := sc.G - (0x08 * uint8(len(l.snake.body)))
	for _, p := range l.snake.body {
		sc = color.RGBA{sc.R, green, sc.B, sc.A}
		vector.DrawFilledRect(screen, float32(p.x*size), float32(p.y*size), float32(size), float32(size), sc, false)
		green += 0x08
	}

	// Draw the food
	vector.DrawFilledRect(screen, float32(l.food.x*size), float32(l.food.y*size), float32(size), float32(size), foodColor, false)

	// Draw the score
	// ebitenutil.DrawString(screen, fmt.Sprintf("Score: %d", l.score), 10, 10)
}

func (l *LevelSnake) Initialize(id LevelID) {
	l.viMode = NormalMode
	l.snake = &Snake{
		body:      []Coord{{x: 5, y: gridHeight / 2}},
		direction: east,
	}
	l.level = id
	l.food = l.generateFood(Coord{x: 0, y: 0})
	l.score = 0
}

func (l *LevelSnake) Update(frameCount int) (bool, error) {
	if isCheatKeyPressed() {
		return true, nil
	}
	// Handle keystrokes
	canTurn := l.level == LevelIdSnake || (l.level == LevelIdInsertMode && l.viMode == NormalMode)
	dir := l.snake.direction
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		if l.snake.direction != east && canTurn {
			dir = west
		} else {
			PlaySound(failOgg)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		if l.snake.direction != south && canTurn {
			dir = north
		} else {
			PlaySound(failOgg)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		if l.snake.direction != north && canTurn {
			dir = south
		} else {
			PlaySound(failOgg)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		if l.snake.direction != west && canTurn {
			dir = east
		} else {
			PlaySound(failOgg)
		}
	}
	if l.level == LevelIdInsertMode {
		if ebiten.IsKeyPressed(ebiten.KeyI) {
			l.viMode = InsertMode
		}
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			l.viMode = NormalMode
		}
	}
	l.snake.direction = dir

	if frameCount%30 == 0 {
		// Only move the snake every 30 frames
		// Update the snake's position
		head := l.snake.body[len(l.snake.body)-1]

		switch l.snake.direction {
		case north:
			head.y -= 1
		case south:
			head.y += 1
		case west:
			head.x -= 1
		case east:
			head.x += 1
		}

		// Check if the snake has collided with the food
		canEat := l.level == LevelIdSnake || (l.level == LevelIdInsertMode && l.viMode == InsertMode)
		if head == l.food && canEat {
			l.food = l.generateFood(head)
			l.score++
			if l.score == lengthForWin {
				PlaySound(winOgg)
			}
			// log.Println("Score: ", l.score)
		} else {
			// Remove the tail
			l.snake.body = l.snake.body[1:]

			// Check if the snake has collided with the boundaries or itself
			if head.x < 0 || head.x >= gridWidth || head.y < 0 || head.y >= gridHeight {
				l.Initialize(l.level)
				return false, nil
			}
			for i := 1; i < len(l.snake.body); i++ {
				if head == l.snake.body[i] {
					l.Initialize(l.level)
					return false, nil
				}
			}
		}

		// Update the snake's body
		l.snake.body = append(l.snake.body, head)
	}

	// continue level until snake dies
	return l.gameIsWon(), nil
}

func (l *LevelSnake) gameIsWon() bool {
	win := len(l.snake.body) >= lengthForWin
	if win {
		PlaySound(winOgg)
	}
	return win
}

func (l *LevelSnake) generateFood(head Coord) Coord {
	// don't put food on the edges
	food := Coord{
		x: rand.Intn(gridWidth-4) + 2,
		y: rand.Intn(gridHeight-4) + 2,
	}
	// don't put food on the snake
	if food == head {
		return l.generateFood(head)
	} else {
		for i := 0; i < len(l.snake.body); i++ {
			if l.snake.body[i] == food {
				return l.generateFood(head)
			}
		}
	}
	return food
}
