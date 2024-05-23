package main

/*
	LevelIdFlappy = iota
	LevelIdBricksHL
	LevelIdBricksHJKL
	LevelIdGemsDD
	LevelIdGemsVM
*/

var IntroText = [...]string{
	//LevelIdFlappy
	`Learn vi by playing classic games. 

In vi you use the J key to move down and the K key to move up
Use these keys to guide the pufferfish through the gaps. 

Pass seven obstacles without fail to advance to the next level`,
	// LevelIdBricksHL
	`Move the paddle left and right by 
pressing H and K keys.

Clear the bricks to advance to the next level`,
	// LevelIdBricksHJKL
	`Move the horizontal paddles left and right (H, L)
and the veritial paddle up and down (J, K) to defend 
all four edges

Clear all bricks to advance to the next level.`,
	// LevelIdGemsDD
	`Delete lines to connect 3 matching jewels. 
Pressing d twice will delete the current line in vi.
Pressing d a number and return will multiple lines.

Delete lines to match 3 identical jewels and turn
the squares gold. Turn all squares gold to advance 
to the next level. 

Be careful. If you you try to delete a line that doesn't 
match up three jewels you'll lose gold!`,
	// LevelIdGemsVM
	`Visual Mode in VI lets you make a text selection.

Press V to enter visual and the navigation keys (H,J,K,L)
to select jewels. Press D to delete the selection.
Escape to exit visual mode.

Make sure deleting connects three identical jewels!`,
	// LevelIdGemsEnd
	`Congratulations you have completed all the learning levels.
Use all the skills you've learned toto complete this level!`,
}

var IntroTitle = [...]string{
	//LevelIdFlappy
	`Welcome to Viple`,
	// LevelIdBricksHL
	`Bricker!`,
	// LevelIdBricksHJKL
	`Bricker Hayhem!`,
	// LevelIdGemsDD
	`Connect Three!`,
	// LevelIdGemsVM
	`Visual Mode`,
	// LevelIdGemsEnd
	`Challenge Level!`,
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
