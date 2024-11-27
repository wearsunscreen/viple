package main

import (
	"log"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type LevelMessage struct {
	game          *Game
	startingFrame int
}

func showIntroDialog(g *Game) {
	// This loads a font and creates a font face.
	ttfFont, err := truetype.Parse(goregular.TTF)

	if err != nil {
		log.Fatal("Error Parsing Font", err)
	}

	// release resources
	g.ui.Container.RemoveChildren()

	textFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 20,
	})

	titleFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 32,
	})

	innerContainer := widget.NewContainer(

		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(dlgBackground)),
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

		widget.TextOpts.Text(TitleText(g.currentLevel), titleFace, dlgText),
	)
	innerContainer.AddChild(titleText)

	level1IntroText := widget.NewText(

		widget.TextOpts.Text(IntroText(g.currentLevel), textFace, dlgText),
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
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			advanceLevelMode(g)
		}),
	)
	innerContainer.AddChild(b)
}

func (l *LevelMessage) Draw(screen *ebiten.Image, frameCount int) {
	// Draw background
	// screen.Fill(seaColor)

	l.game.ui.Draw(screen)
}

func (l *LevelMessage) Initialize(id LevelID, g *Game) {
	l.game = g
	showIntroDialog(l.game)
}

func (l *LevelMessage) Update(frameCount int) (bool, error) {
	l.game.ui.Update()
	checkForKeystroke(ebiten.KeyEnter, func() { advanceLevelMode(l.game) })

	return false, nil
}
