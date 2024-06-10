package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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
	level LevelID
	food  Coord
	score int
	snake *Snake
}

var (
	gridWidth  = screenWidth / size
	gridHeight = screenHeight / size
	snakeColor = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	size       = 40 // size of each square in the grid
)

func (l *LevelSnake) Draw(screen *ebiten.Image, frameCount int) {
	// Draw the snake
	for _, p := range l.snake.body {
		vector.DrawFilledRect(screen, float32(p.x*size), float32(p.y*size), float32(size), float32(size), snakeColor, false)
	}

	// Draw the food
	vector.DrawFilledRect(screen, float32(l.food.x*size), float32(l.food.y*size), float32(size), float32(size), mediumScarletRed, false)

	// Draw the score
	// ebitenutil.DrawString(screen, fmt.Sprintf("Score: %d", l.score), 10, 10)
}

func (l *LevelSnake) Initialize(id LevelID) {
	l.snake = &Snake{
		body:      []Coord{{x: 5, y: gridHeight / 2}},
		direction: east,
	}
	l.level = LevelIdSnake
	l.food = l.generateFood()
	l.score = 0
}

func (l *LevelSnake) Update(frameCount int) (bool, error) {
	if isCheatKeyPressed() {
		return true, nil
	}
	// Handle keystrokes
	dir := l.snake.direction
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		if l.snake.direction != east {
			dir = west
		} else {
			PlaySound(failOgg)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		if l.snake.direction != south {
			dir = north
		} else {
			PlaySound(failOgg)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		if l.snake.direction != north {
			dir = south
		} else {
			PlaySound(failOgg)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		if l.snake.direction != west {
			dir = east
		} else {
			PlaySound(failOgg)
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
		if head == l.food {
			l.food = l.generateFood()
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

func (l *LevelSnake) generateFood() Coord {
	food := Coord{
		x: rand.Intn(gridWidth-8) + 4,
		y: rand.Intn(gridHeight-8) + 4,
	}
	for i := 0; i < len(l.snake.body); i++ {
		if l.snake.body[i] == food {
			return l.generateFood()
		}
	}
	return food
}
