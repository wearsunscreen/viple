# Viple: A Game to Learn Vi

Viple is a game designed to help you learn and practice the vi text editor commands. Viple is written in Go, using the [Ebiten game engine](https://ebitengine.org/) and the [Ebiten UI widget library](https://ebitenui.github.io/). This project is still under development. Contributions are welcome.

## Features

- Interactive levels to teach various vi commands
- Gamified learning experience with different game modes
- Cross-platform compatibility

## To-Do List

- DD level
- - don't allow other key strokes
- - make wider
- - allow delete multiple lines, d2, d3 etc
- - require "dd" instead of "d"
- on exit without write, restart level
- on exit with write, advance to next level
- Implement a Snake game to practice the H, J, K, and L navigation keys
- Improve the overall visual design and aesthetics
- Develop a maki game to teach word advance 
- Create a main menu for selecting levels 
- Add support for quitting the game using commands `:q`, `:quit`, `:exit`, etc.
- Introduce a quiz level for learning write and exit commands
- Add a timer to levels to challenge the user
- Animate the pufferfish character
- Script the pipe gaps in the pufferfish level
- Implement scaling, fading and rotation animations for disappearing gems and bricks
- Build and deploy the game as a web application using WebAssembly (WASM)
- Save player progress
- Test the game on different platforms for compatibility

## Known Defects

- The first and second pipes in the Flappy Bird level are not evenly spaced
- LevelGemsDD can have unplayable grid, no possible triples to make

## Contributing

Contributions to Viple are welcome! If you find any issues or have ideas for improvements, please feel free to open an issue or submit a pull request on the project's GitLab repository.
