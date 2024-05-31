package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Snake struct {
	body      []Coord
	direction ebiten.Key
	size      int
}

type LevelSnake struct {
	snake *Snake
	food  Coord
	score int
}

func (l *LevelSnake) Initialize(id LevelID) {
	l.snake = &Snake{
		body:      []Coord{{x: 100, y: screenHeight / 2}},
		direction: ebiten.KeyRight,
		size:      20,
	}
	l.food = l.generateFood()
	l.score = 0
}

func (l *LevelSnake) Draw(screen *ebiten.Image, frameCount int) {
	// Draw the snake
	for _, p := range l.snake.body {
		vector.DrawFilledRect(screen, float32(p.x), float32(p.y), float32(l.snake.size), float32(l.snake.size), color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}, false)
	}

	// Draw the food
	vector.DrawFilledRect(screen, float32(l.food.x), float32(l.food.y), float32(l.snake.size), float32(l.snake.size), color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, false)

	// Draw the score
	// ebitenutil.DrawString(screen, fmt.Sprintf("Score: %d", l.score), 10, 10)
}

func (l *LevelSnake) Update(frameCount int) (bool, error) {
	// Update the snake's position
	head := l.snake.body[0]

	// Handle keystrokes
	dir := l.snake.direction
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		dir = ebiten.KeyLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		dir = ebiten.KeyUp
	}
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		dir = ebiten.KeyDown
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		dir = ebiten.KeyRight
	}
	l.snake.direction = dir

	switch l.snake.direction {
	case ebiten.KeyUp:
		head.y -= l.snake.size
	case ebiten.KeyDown:
		head.y += l.snake.size
	case ebiten.KeyLeft:
		head.x -= l.snake.size
	case ebiten.KeyRight:
		head.x += l.snake.size
	}

	// Check if the snake has collided with the food
	if head == l.food {
		l.snake.body = append(l.snake.body, head)
		l.food = l.generateFood()
		l.score++
		log.Println("Score: ", l.score)
	} else {
		// Remove the tail if the snake has grown
		if len(l.snake.body) > l.score+1 {
			l.snake.body = l.snake.body[1:]
		}
	}

	// Check if the snake has collided with the boundaries or itself
	if head.x < 0 || head.x >= screenWidth || head.y < 0 || head.y >= screenHeight {
		return true, nil
	}
	for i := 1; i < len(l.snake.body); i++ {
		if head == l.snake.body[i] {
			return true, nil
		}
	}

	// Update the snake's body
	l.snake.body[0] = head

	return false, nil
}

func (l *LevelSnake) gameIsWon() bool {
	return len(l.snake.body) >= 10
}

func (l *LevelSnake) generateFood() Coord {
	food := Coord{
		x: rand.Intn(screenWidth),
		y: rand.Intn(screenHeight),
	}
	for i := 0; i < len(l.snake.body); i++ {
		if l.snake.body[i] == food {
			return l.generateFood()
		}
	}
	return food
}
