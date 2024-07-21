# Viple: A Game to Learn Vi
Viple is a game designed to help you learn and practice the vi text editor commands. Viple is written in Go, using the [Ebiten game engine](https://ebitengine.org/) and the [Ebiten UI widget library](https://ebitenui.github.io/). 

This project is still under development. Contributions are welcome.

## Features
- Interactive levels to teach various vi commands
- Gamified learning experience with different game modes
- Cross-platform compatibility

## MVP
- change gem, selection, cursor drawing to make all more distinguishable
- sound on pufferfish collision
- first pipes should appear sooner in flappy
- build in WebAsm
-- env GOOS=js GOARCH=wasm go build -o viple.wasm github.com/wearsunscreen/viple
- create video of playthrough

## To-Do List
- Add insert mode level, snake level where only in insert mode can you eat, change color of snake by mode
- Embed viple into a web page with instructions. 
- center dialogs?
- Animate the pufferfish character
- Implement a Snake game to practice the H, J, K, and L navigation keys
- Implement a minefield game to practice w and page up/down
- Improve the overall visual design and aesthetics
- Develop a maki game to teach word advance
- Create a main menu for selecting levels
- Script the pipe gaps in the pufferfish level
- Add support for quitting the game using commands `:q`, `:quit`, `:exit`, etc.
- set of valid and invalid key strokes per level
- Create jeopardy like game level to quiz write and exit commands
- Add a timer to levels to challenge the user
- Implement scaling, fading and rotation animations for disappearing gems and bricks
- Build and deploy the game as a web application using WebAssembly (WASM)
- Save player progress
- Test the game on different platforms for compatibility

## Refactoring
- change functions that take level pointers and change to level methods
- Isolate ebiten code from game logic
- Make Grid a package
- Add logging, perhaps with glog

## Known Defects
- mixing d and v can get the end level confused.

## Contributing
Contributions to Viple are welcome! Want to add a new level? Found a bug? Have ideas for improvements? Open an issue or submit a pull request on the project's GitLab repository.
