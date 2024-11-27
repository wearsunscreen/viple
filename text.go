package main

import (
	"log"
)

func IntroText(level LevelID) string {
	switch level {
	// intro messages
	case LevelIdMsgWelcome:
		return `Learn vi by playing classic games. Your first challenge 
is to navigate the pufferfish through the obstacles. 
Pass seven obstacles without fail to advance to the next level.

J -- Up
K -- Down
Q -- Quit


Game by John Crane, https://github.com/wearsunscreen/viple

Music and SFX by Gianni Bellucci, all rights reserved`
	case LevelIdMsgBricksIntro:
		return `Clear the bricks to advance to the next level

H to move left
K to move right`
	case LevelIdMsgSnakeIntro:
		return `Guide the snake using the H, J, K, L keys.
Eat the apples to grow the snake longer.`
	case LevelIdMsgInsertIntro:
		return `Enter Insert Mode to eat the apple.
Exit Insert Mode to move the snake.

I - enter insert mode
Esc - exit insert mode`
	case LevelIdMsgDeleteIntro:
		return `Delete lines to connect 3 matching jewels.
D, D -- Delete line
D, [2, 3, 4, ...], Enter -- Delete multiple lines

Delete lines to line up 3 identical jewels in a vertical column. 
Matching gems will turn the squares gold. 
Turn all squares gold to advance to the next level. 

Be careful. If you you try to delete a line that doesn't 
match up three jewels you'll lose gold!`
	case LevelIdMsgVMIntro:
		return `Visual Mode in VI lets you make a text selection.

Press V to enter visual and the navigation keys (H,J,K,L)
to select jewels. Press D to delete the selection.
Escape to exit visual mode.

Make sure deleting connects three identical jewels!`
	case LevelIdMsgChallengeIntro:
		return `Congratulations you have completed all the learning levels.
Use all the skills you've learned toto complete this level!`

	// done messages
	case LevelIdMsgFlappyDone:
		return `You have completed the first level.`
	case LevelIdMsgBricksDone:
		fallthrough
	case LevelIdMsgSnakeDone:
		fallthrough
	case LevelIdMsgInsertDone:
		fallthrough
	case LevelIdMsgDeleteDone:
		fallthrough
	case LevelIdMsgVMDone:
		fallthrough
	case LevelIdMsgAllLevelsDone:
		return `You have completed the level.`
	case LevelIdMsgChallengeDone:
		return `You have completed all levels.`

	case LevelIdBricksHJKL: // deprecated
		return `Move the horizontal paddles left and right (H, L)
and the veritial paddle up and down (J, K) to defend 
all four edges

Clear all bricks to advance to the next level.`

	default:
		log.Println("Unknown Level ", level)
		return "Unknown Level!"
	}
}

func TitleText(level LevelID) string {
	switch level {
	// intro messages
	case LevelIdMsgWelcome:
		return `Welcome to Viple!`
	case LevelIdMsgBricksIntro:
		return `Bricker!`
	case LevelIdMsgSnakeIntro:
		return `Snake!`
	case LevelIdMsgInsertIntro:
		return `Insert Mode!`
	case LevelIdMsgDeleteIntro:
		return `Connect Three!`
	case LevelIdMsgVMIntro:
		return `Visual Mode!`
	case LevelIdMsgChallengeIntro:
		return `Challenge Level!`

	// done messages
	case LevelIdMsgFlappyDone:
		fallthrough
	case LevelIdMsgBricksDone:
		fallthrough
	case LevelIdMsgSnakeDone:
		fallthrough
	case LevelIdMsgInsertDone:
		fallthrough
	case LevelIdMsgDeleteDone:
		fallthrough
	case LevelIdMsgVMDone:
		fallthrough
	case LevelIdMsgChallengeDone:
		fallthrough
	case LevelIdMsgAllLevelsDone:
		return `Congratulations!`

	// play messages
	case LevelIdBricksHJKL: // deprecated
		return `Bricker Hayhem!`
	default:
		log.Println("Unknown Level ", level)
		return "Unknown Level!"
	}
}
