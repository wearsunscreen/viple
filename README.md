# viple
A game to help you learn vi.


## License
Copyright John Crane, 2024

## Defects
* The first swap will allow act like a swap if there are triples on the board even if the swap does not create a triple.

## To Do
* snake level
* Level description
* split into modules
* use a generic list - https://gobyexample.com/generics
* use composition to define levels https://www.tutorialspoint.com/composition-in-golang
** each level would contain a Drawer, Updater, Initializer, StartInformer, StartEnder
* level to teach delete lines in gems
* cut and paste lines in gems
* maki game and word advance
* brickout where you switch side with paddle using page up, page down commands
* adventure game with command line
* bricks 1 and 2 using function pointers for polymorphism
** to introduce a level
** to congratulate completion of a level and game, options to repeat, go to next, go to main menu
** Main menu, lets you choose level
* move handling of : commands to shared location
* pufferfish animation
* pufferfish have its own rng to repeat level
* Scaler to animate gems
* Fader to animate disappearing gems, bricks
* Flush keystrokes between levels
* Deploy as Open Source?
* deploy as web page
* Snake