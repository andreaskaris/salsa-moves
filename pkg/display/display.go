package display

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Display struct {
	app    fyne.App
	window fyne.Window
	Ch     chan string
}

func New() *Display {
	display := &Display{}
	display.app = app.New()
	display.window = display.app.NewWindow("Hello")
	display.Ch = make(chan string, 1)
	return display
}

func (d *Display) renderVBox(text string) {
	c := canvas.NewText(text, &color.RGBA{0xff, 0x00, 0x00, 0xff})
	c.TextSize = 20
	vbox := container.NewVBox(
		c,
		widget.NewButton("Options", func() {
		}))
	d.window.SetContent(vbox)
}

func (d *Display) Render() {
	d.renderVBox("")

	go func() {
		for e := range d.Ch {
			d.renderVBox(e)
		}
	}()
	d.window.ShowAndRun()
}
