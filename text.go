package main

import (
	"log"
)

func IntroText(level LevelID) string {
	switch level {
	case LevelIdFlappy:
		return `Learn vi by playing classic games. Your first challenge 
is to navigate the pufferfish through the obstacles. 
Pass seven obstacles without fail to advance to the next level.

J -- Up
K -- Down
Q -- Quit`
	case LevelIdBricksHL:
		return `Clear the bricks to advance to the next level

H to move left
K to move right`
	case LevelIdSnake:
		return `Guide the snake using the H, J, K, L keys.
Eat the apples to grow the snake longer.`

	case LevelIdInsertMode:
		return `Enter Insert Mode to eat the apple.
Exit Insert Mode to move the snake.

I - enter insert mode
Esc - exit insert mode`

	case LevelIdBricksHJKL:
		return `Move the horizontal paddles left and right (H, L)
and the veritial paddle up and down (J, K) to defend 
all four edges

Clear all bricks to advance to the next level.`

	case LevelIdGemsDD:
		return `Delete lines to connect 3 matching jewels.
D, D -- Delete line
D, [2, 3, 4, ...], Enter -- Delete multiple lines

Delete lines to line up 3 identical jewels in a vertical column. 
Matching gems will turn the squares gold. 
Turn all squares gold to advance to the next level. 

Be careful. If you you try to delete a line that doesn't 
match up three jewels you'll lose gold!`
	case LevelIdGemsVM:
		return `Visual Mode in VI lets you make a text selection.

Press V to enter visual and the navigation keys (H,J,K,L)
to select jewels. Press D to delete the selection.
Escape to exit visual mode.

Make sure deleting connects three identical jewels!`
	case LevelIdGemsEnd:
		return `Congratulations you have completed all the learning levels.
Use all the skills you've learned toto complete this level!`

	default:
		log.Println("Unknown Level ", level)
		return "Unknown Level!"
	}
}

func TitleText(level LevelID) string {
	switch level {
	case LevelIdFlappy:
		return `Welcome to Viple`
	case LevelIdBricksHL:
		return `Bricker!`
	case LevelIdSnake:
		return `Snake!`
	case LevelIdInsertMode:
		return `Insert Mode!`
	case LevelIdBricksHJKL:
		return `Bricker Hayhem!`
	case LevelIdGemsDD:
		return `Connect Three!`
	case LevelIdGemsVM:
		return `Visual Mode`
	case LevelIdGemsEnd:
		return `Challenge Level!`
	default:
		log.Println("Unknown Level ", level)
		return "Unknown Level!"
	}
}
