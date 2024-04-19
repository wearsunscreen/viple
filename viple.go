package main

import (
	"flag"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/exp/constraints"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	version      = "Viple 0.1"
)

type DialogText struct {
	title string
	intro string
	outro string
}

var Level1Dialog = DialogText{
	`Welcome to Viple
VI Play to Learn. `,
	`In the first level you will 
learn to move left and right 
by pressing H and K keys.`,
	`Congrats, you're ready to 
move onto the next level`,
}

type Level interface {
	Draw(screen *ebiten.Image, frameCount int)
	Initialize()
	// update every frame, return true if level is complete
	Update(g *Game) (bool, error)
}

type LevelID int

const (
	LevelIdBricksHL = iota
	LevelIdFlappy
	LevelIdBricksHJKL
	LevelIdGemsVM
)

type Mode int

var (
	rng  *rand.Rand
	keys []ebiten.Key
)

type Game struct {
	currentLevel LevelID
	frameCount   int
	levelHL      LevelBricksHL
	levelJK      LevelFlappy
	levelVM      LevelGemsVisualMode
	ui           *ebitenui.UI
}

type Number interface {
	constraints.Integer | constraints.Float
}

func main() {
	var seed int
	flag.IntVar(&seed, "seed", 0, "Seed for random number generation")
	flag.Parse()
	seedRNG(int64(seed))

	ebiten.SetWindowSize(gameDimensions())
	ebiten.SetWindowTitle(version)

	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}

// function to fill slice of any type
func fillSlice[T any](s []T, value T) []T {
	if s == nil {
		panic("slice cannot be nil")
	}

	for i := range s {
		s[i] = value
	}
	return s
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	switch g.currentLevel {
	case LevelIdBricksHL:
		g.levelHL.Draw(screen, g.frameCount)
	case LevelIdFlappy:
		g.levelJK.Draw(screen, g.frameCount)
	case LevelIdBricksHJKL:
		g.levelHL.Draw(screen, g.frameCount)
	case LevelIdGemsVM:
		g.levelVM.Draw(screen, g.frameCount)
	default:
		panic("Unknown game level " + strconv.Itoa(int(g.currentLevel)))
	}

	// the UI
	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w, h := gameDimensions()
	return w, h
}

func (g *Game) Update() error {
	g.ui.Update()

	var levelOver bool
	var err error
	g.frameCount++
	switch g.currentLevel {
	case LevelIdBricksHL:
		levelOver, err = g.levelHL.Update(g.frameCount)
	case LevelIdBricksHJKL:
		levelOver, err = g.levelHL.Update(g.frameCount)
	case LevelIdFlappy:
		levelOver, err = g.levelJK.Update(g.frameCount)
	case LevelIdGemsVM:
		levelOver, err = g.levelVM.Update(g.frameCount)
	}
	if levelOver {
		// advance to next Level if current level has been won
		// bugbug: we don't handle completing the last level cleanly
		g.currentLevel += 1

		switch g.currentLevel {
		case LevelIdBricksHL:
			g.levelHL.Initialize()
		case LevelIdBricksHJKL:
			g.levelHL.level = g.currentLevel
			g.levelHL.Initialize()
		case LevelIdFlappy:
			g.levelJK.Initialize()
		case LevelIdGemsVM:
			g.levelVM.Initialize()
		}
	}
	return err
}

func gameDimensions() (width int, height int) {
	return screenWidth, screenHeight
}

func isCheatKeyPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		return true
	}
	return false
}

func limitToRange[T Number](input, low, high T) (output T) {
	output = input
	if input < low {
		output = low
	} else if input > high {
		output = high
	}
	return output
}

func loadImage(path string) *ebiten.Image {
	image, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Error loading image: %v", err)
	}
	return image
}

func closeUI(res *uiResources) {
	res.close()
}
func newSeparator(res *uiResources, ld interface{}) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}))),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(ld)))

	c.AddChild(widget.NewGraphic(
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch:   true,
			MaxHeight: 2,
		})),
		widget.GraphicOpts.ImageNineSlice(image.NewNineSliceColor(res.separatorColor)),
	))

	return c
}

func newGame() *Game {
	res, err := newUIResources()
	if err != nil {
		return nil
	}

	// This loads a font and creates a font face.
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal("Error Parsing Font", err)
	}
	fontFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 32,
	})

	//This creates the root container for this UI.
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0x80})),
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			//Set how much padding before displaying content
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(100)),
		)),
	)

	// This adds the root container to the UI, so that it will be rendered.
	ui := &ebitenui.UI{
		Container: rootContainer,
	}

	defer closeUI(res)

	innerContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(mediumButter)),
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			//Define number of columns in the grid
			widget.GridLayoutOpts.Columns(1),
			//Define how much padding to inset the child content
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			//Define how far apart the rows and columns should be
			widget.GridLayoutOpts.Spacing(20, 10),
			//Define how to stretch the rows and columns. Note it is required to
			//specify the Stretch for each row and column.
			widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{false, true}),
		)),
	)
	rootContainer.AddChild(innerContainer)

	titleText := widget.NewText(
		widget.TextOpts.Text(Level1Dialog.title, fontFace, color.White),
	)
	innerContainer.AddChild(titleText)

	level1IntroText := widget.NewText(
		widget.TextOpts.Text(Level1Dialog.intro, fontFace, color.White),
	)
	innerContainer.AddChild(level1IntroText)

	innerContainer.AddChild(newSeparator(res, widget.RowLayoutData{
		Stretch: true,
	}))

	b := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.Text("Ok", res.button.face, res.button.text),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) { fmt.Println("Cursor Entered: " + args.Button.Text().Label) }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { fmt.Println("Cursor Exited: " + args.Button.Text().Label) }),
	)
	innerContainer.AddChild(b)

	g := Game{
		ui: ui,
	}

	g.levelHL.Initialize()
	g.levelJK.Initialize()
	g.levelVM.Initialize()
	g.currentLevel = LevelIdBricksHL

	return &g
}

func seedRNG(seed int64) {
	if seed == 0 {
		seed = time.Now().UnixNano() % 10000
	}
	log.Println("Random seed is ", seed)
	rng = rand.New(rand.NewSource(seed))
}
