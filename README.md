# Viple: A Game to Learn Vi
Viple is a game designed to help you learn and practice the vi text editor commands. Viple is written in Go, using the [Ebiten game engine](https://ebitengine.org/) and the [Ebiten UI widget library](https://ebitenui.github.io/). 

Find the lastest published version on https://johncrane.dev/

## Contributing
Contributions to Viple are welcome! Want to add a new level? Found a bug? Have ideas for improvements? Open an issue or submit a pull request on the project's GitLab repository.

### How to create a new level
1. Create a LevelID in viple.go
1. Add a case for the new LevelID in IntroText() and TitleText()
1. Add a case for the new LevelID in advanceLevelMode()

## Features
- Interactive levels to teach various vi commands
- Gamified learning experience with different game modes
- Cross-platform compatibility

## To-Do List
- change color of 10th puffyfish obstacle
- create video of playthrough
- loading screen for wasm
- Embed viple into a web page with instructions. - https://ebitengine.org/en/documents/webassembly.html
- sound on pufferfish collision
- Test the game on different platforms for compatibility
- Add "Next Level" and "Prev Level" buttons to intro dialog or a main menu for selecting levels
- Animate the pufferfish character
- center dialogs
- Implement a zuma game to practice w and insert
- Improve the overall visual design and aesthetics
- Develop a maki game to teach word advance
- Add support for quitting the game using commands `:q`, `:quit`, `:exit`, etc.
- Implement scaling, fading and rotation animations for disappearing gems and bricks
- Save player progress
- add in game hints (e.g. "enter insert mode to eat the fruit" when snake touches apple in normal mode)

## build in WebAsm
-- env GOOS=js GOARCH=wasm go build -o viple.wasm github.com/wearsunscreen/viple
[//]: <> (-- cp viple.wasm ../../dev-portfolio/site/viple)


