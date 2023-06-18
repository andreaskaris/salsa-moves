package display

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type Display struct {
	app      fyne.App
	window   fyne.Window
	textSize float32
	Ch       chan []string
}

func New(textSize float32) *Display {
	display := &Display{}
	display.app = app.New()
	display.window = display.app.NewWindow("Salsa moves")
	display.Ch = make(chan []string, 1)
	display.textSize = textSize
	return display
}

func (d *Display) renderVBox(text ...string) {
	var vboxObjects []fyne.CanvasObject

	for _, t := range text {
		c := canvas.NewText(t, &color.RGBA{0xff, 0x00, 0x00, 0xff})
		c.TextSize = d.textSize
		vboxObjects = append(vboxObjects, c)
	}
	vbox := container.NewVBox(vboxObjects...)
	d.window.SetContent(vbox)
}

func (d *Display) Render() {
	d.renderVBox("")

	go func() {
		for e := range d.Ch {
			d.renderVBox(e...)
		}
	}()
	d.window.ShowAndRun()
}
