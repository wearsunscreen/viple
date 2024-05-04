package main

/*
	LevelIdFlappy = iota
	LevelIdBricksHL
	LevelIdGemsDD
	LevelIdBricksHJKL
	LevelIdGemsVM
*/

var IntroText = [...]string{
	//LevelIdFlappy
	`Learn vi by playing classic games. 
In vi you use the J key to move down and the K key to move up 
Use these keys to guide the pufferfish through the gaps. 
Pass seven obstacles without fail to advance to the next level`,
	// LevelIdGemsDD
	`Now that you can move up and down
	move your cursor up and down and press the d key twice
	"dd" will delete the current line in vi.
	In this game you will delete lines to line up three 
	jewels and be rewarded by turning squares gold.
	Be careful. If you you try to delete a line that doesn't 
	match up three jewels you'll lose gold!
	Turn the whole grid gold to advance to the next level`,
	// LevelIdBricksHL
	`Move the paddle left and right by 
	pressing H and K keys.
	Clear the bricks to advance to the next level`,
	// LevelIdBricksHJKL
	`Move the horizontal paddles left and right (H, L)
	and the veritial paddle up and down (J, K)
	to defend all four edgesProtect all the edges.
	Clear the bricks to move one`,
	// LevelIdGemsVM
	`Visual Mode`,
}

var IntroTitle = [...]string{
	//LevelIdFlappy
	`Welcome to Viple`,
	// LevelIdGemsDD
	`Connect Three!`,
	// LevelIdBricksHL
	`Bricker!`,
	// LevelIdBricksHJKL
	`Bricker Hayhem!`,
	// LevelIdGemsVM
	`Visual Mode`,
}

func GetIntroText(id int) string {
	if id >= len(IntroTitle) {
		return "You have reached the end of the game!"
	}
	return IntroText[id]
}

func GetTitleText(id int) string {
	if id >= len(IntroTitle) {
		return "Kudos!"
	}
	return IntroTitle[id]
}
