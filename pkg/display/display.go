package display

import (
	"log"

	"github.com/andreaskaris/salsa-moves/pkg/config"
	"gopkg.in/yaml.v3"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Display struct {
	app            fyne.App
	window         fyne.Window
	content        *fyne.Container
	movesContainer *fyne.Container
	menu           *widget.Menu
	options        *fyne.Container
	Ch             chan []string
	config         *config.Config
}

func New(c *config.Config) *Display {
	display := &Display{}
	display.app = app.New()
	display.window = display.app.NewWindow("Salsa moves")
	display.Ch = make(chan []string, 1)
	display.config = c
	return display
}

func (d *Display) refreshMovesPanel(text ...string) {
	var vboxObjects []fyne.CanvasObject

	for _, t := range text {
		c := canvas.NewText(t, &color.RGBA{0xff, 0x00, 0x00, 0xff})
		c.TextSize = d.config.Text.Size
		vboxObjects = append(vboxObjects, c)
	}

	d.movesContainer.Objects = vboxObjects
	d.movesContainer.Refresh()
}

func (d *Display) Render() {
	d.menu = widget.NewMenu(fyne.NewMenu(
		"test",
		fyne.NewMenuItem("Main", func() {
			d.options.Hide()
			d.movesContainer.Show()
		}),
		fyne.NewMenuItem("Options", func() {
			d.options.Show()
			d.movesContainer.Hide()
		}),
	))

	d.movesContainer = container.NewVBox()

	d.options = drawOptionsPanel(d.config)
	d.options.Hide()
	d.content = container.NewHBox(d.menu, d.movesContainer, d.options)
	d.window.SetContent(d.content)

	go func() {
		for e := range d.Ch {
			d.refreshMovesPanel(e...)
		}
	}()

	d.window.ShowAndRun()
}

func drawOptionsPanel(c *config.Config) *fyne.Container {
	text, err := yaml.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	return container.NewVBox(widget.NewLabel(string(text)))
}
