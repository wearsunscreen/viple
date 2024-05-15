# Viple: A Game to Learn Vi
Viple is a game designed to help you learn and practice the vi text editor commands. Viple is written in Go, using the [Ebiten game engine](https://ebitengine.org/) and the [Ebiten UI widget library](https://ebitenui.github.io/). This project is still under development. Contributions are welcome.

## Features
- Interactive levels to teach various vi commands
- Gamified learning experience with different game modes
- Cross-platform compatibility

## MVP
- clean end of game
- build in WebAsm

## To-Do List
- Animate the pufferfish character
- Implement a Snake game to practice the H, J, K, and L navigation keys
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
- Isolate ebiten code from game logic
- Make Grid a package

## Known Defects
- "dd" causes fail.ogg sound in VM game level

## Contributing

Contributions to Viple are welcome! Want to add a new level? Found a but? Have ideas for improvements? Please feel free to open an issue or submit a pull request on the project's GitLab repository.
