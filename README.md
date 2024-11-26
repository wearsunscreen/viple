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
- Implement a zuma game to practice w, yank and put
- better colors for pufferfish level
- add message that says you have completed all levels, 
- animate the pufferfish character
- center dialogs
- add next level, current level, previous level props to all levels to allow forward backward and repeat
- Add "Next Level" and "Prev Level" buttons to intro dialog or a main menu for selecting levels
- toggle music on and off
- toggle sound on and off
- add difficulty level
- pause key
- Implement scaling, fading and rotation animations for disappearing gems and bricks
- Improve the overall visual design and aesthetics
- Save player progress
- environment variable for build type, wasm or app

## build in WebAsm
-- env GOOS=js GOARCH=wasm go build -o viple.wasm github.com/wearsunscreen/viple
[//]: <> (-- cp viple.wasm ../../dev-portfolio/site/viple)


