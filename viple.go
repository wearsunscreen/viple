package main

import (
	"flag"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
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
}

type Level interface {
	Draw(screen *ebiten.Image, frameCount int)
	Initialize(id LevelID)
	IntroText() string
	TitleText() string
	Update(frameCount int) (bool, error)
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
	curLevel     Level
	frameCount   int
	showUI       bool
	ui           *ebitenui.UI
	uiRes        *uiResources
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

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	g.curLevel.Draw(screen, g.frameCount)

	// the UI
	if g.showUI {
		g.ui.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w, h := gameDimensions()
	return w, h
}

func (g *Game) Update() error {

	var levelOver bool
	var err error
	g.frameCount++

	if g.showUI {
		g.ui.Update()
	} else {
		levelOver, err = g.curLevel.Update(g.frameCount)
		if levelOver {
			// advance to next Level if current level has been won
			// bugbug: we don't handle completing the last level cleanly
			g.currentLevel += 1
			g.showUI = true
			switch g.currentLevel {
			case LevelIdBricksHL:
				g.curLevel = Level(&LevelBricksHL{})
			case LevelIdBricksHJKL:
				g.curLevel = Level(&LevelBricksHL{})
			case LevelIdFlappy:
				g.curLevel = Level(&LevelFlappy{})
			case LevelIdGemsVM:
				g.curLevel = Level(&LevelGemsVisualMode{})
			}
			g.curLevel.Initialize(g.currentLevel)
		}
	}
	return err
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

func gameDimensions() (width int, height int) {
	return screenWidth, screenHeight
}

func hideUI(g *Game) {
	g.showUI = false
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
	g := Game{}

	g.showUI = true
	g.curLevel = Level(&LevelBricksHL{})
	g.curLevel.Initialize(LevelIdBricksHL)

	res, err := newUIResources()
	if err != nil {
		return nil
	}
	g.uiRes = res

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

	g.ui = ui

	showDialog(&g)

	return &g
}

func seedRNG(seed int64) {
	if seed == 0 {
		seed = time.Now().UnixNano() % 10000
	}
	log.Println("Random seed is ", seed)
	rng = rand.New(rand.NewSource(seed))
}

func showDialog(g *Game) {
	// This loads a font and creates a font face.
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal("Error Parsing Font", err)
	}
	fontFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 32,
	})

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
	g.ui.Container.AddChild(innerContainer)

	titleText := widget.NewText(
		widget.TextOpts.Text(g.curLevel.TitleText(), fontFace, color.White),
	)
	innerContainer.AddChild(titleText)

	level1IntroText := widget.NewText(
		widget.TextOpts.Text(g.curLevel.IntroText(), fontFace, color.White),
	)
	innerContainer.AddChild(level1IntroText)

	innerContainer.AddChild(newSeparator(g.uiRes, widget.RowLayoutData{
		Stretch: true,
	}))

	b := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(g.uiRes.button.image),
		widget.ButtonOpts.Text("Ok", g.uiRes.button.face, g.uiRes.button.text),
		widget.ButtonOpts.TextPadding(g.uiRes.button.padding),
		// widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) { fmt.Println("Cursor Entered: " + args.Button.Text().Label) }),
		// widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { fmt.Println("Cursor Exited: " + args.Button.Text().Label) }),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			hideUI(g)
		}),
	)
	innerContainer.AddChild(b)

}
