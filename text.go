package main

import (
	"log"
)

func IntroText(level LevelID) string {
	switch level {
	case LevelIdFlappy:
		return `Learn vi by playing classic games. 

In vi you use the J key to move down and the K key to move up
Use these keys to guide the pufferfish through the gaps. 

Pass seven obstacles without fail to advance to the next level.

Q to quit`
	case LevelIdBricksHL:
		return `Move the paddle left and right by 
pressing H and K keys.

Clear the bricks to advance to the next level

Q to quit`
	case LevelIdSnake:
		return `Guide the snake using the H, J, K, L keys.
Eat the apples to grow the snake longer.

Q to quit`
	case LevelIdBricksHJKL:
		return `Move the horizontal paddles left and right (H, L)
and the veritial paddle up and down (J, K) to defend 
all four edges

Clear all bricks to advance to the next level.

Q to quit`
	case LevelIdGemsDD:
		return `Delete lines to connect 3 matching jewels. 
Pressing d twice will delete the current line in vi.
Pressing d a number and return will multiple lines.

Delete lines to match 3 identical jewels and turn
the squares gold. Turn all squares gold to advance 
to the next level. 

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
Use all the skills you've learned toto complete this level!

Q to quit`
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
