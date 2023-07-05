package display

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/andreaskaris/salsa-moves/pkg/config"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

const (
	sleepForRand  = 15
	sleepForConst = 5
	bpm           = 60
)

type Display struct {
	app             fyne.App
	window          fyne.Window
	content         *fyne.Container
	options         *fyne.Container
	movesFunction   func() (*config.Config, chan []string, chan bool)
	movesContext    context.Context
	movesCancelFunc context.CancelFunc
	Ch              chan []string
	config          *config.Config
}

func New(c *config.Config) *Display {
	display := &Display{}
	display.app = app.New()
	display.window = display.app.NewWindow("Salsa moves")
	display.config = c
	return display
}

func (d *Display) Render() {
	content := container.NewMax()
	title := widget.NewLabel("title")
	screen := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)
	setScreen := func(t Screen) {
		title.SetText(t.Title)
		content.Objects = []fyne.CanvasObject{t.View()}
		content.Refresh()
	}

	split := container.NewHSplit(d.makeNav(setScreen, false), screen)
	split.Offset = 0.2
	d.window.SetContent(split)

	d.window.Resize(fyne.NewSize(1024, 800))
	d.window.ShowAndRun()
}

func (d *Display) optionsScreen() fyne.CanvasObject {
	if d.movesCancelFunc != nil {
		d.movesCancelFunc()
	}

	bpmFloat := float64(d.config.GetBPM())
	bpmData := binding.BindFloat(&bpmFloat)
	bpmData.AddListener(binding.NewDataListener(func() {
		bpm, err := bpmData.Get()
		if err != nil {
			log.Fatalf("error getting bpmData %f, err: %q", bpmData, err)
		}
		d.config.SetBPM(int(bpm))
	}))
	bpmLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(bpmData, "BPM: %0.2f"))
	bpmEntry := widget.NewEntryWithData(binding.FloatToString(bpmData))
	bpmFloats := container.NewGridWithColumns(2, bpmLabel, bpmEntry)
	bpmSlide := widget.NewSliderWithData(60, 240, bpmData)
	bpmSlide.Step = 1.0

	sleepForRandFloat := float64(d.config.GetSleepForRand())
	sleepForRandData := binding.BindFloat(&sleepForRandFloat)
	sleepForRandData.AddListener(binding.NewDataListener(func() {
		data, err := sleepForRandData.Get()
		if err != nil {
			log.Fatalf("error getting sleepForRandData %f, err: %q", sleepForRandData, err)
		}
		d.config.SetSleepForRand(int(data))
	}))
	sleepForRandLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(sleepForRandData, "Sleep for rand: %0.2f"))
	sleepForRandEntry := widget.NewEntryWithData(binding.FloatToString(sleepForRandData))
	sleepForRandFloats := container.NewGridWithColumns(2, sleepForRandLabel, sleepForRandEntry)
	sleepForRandSlide := widget.NewSliderWithData(0, 60, sleepForRandData)
	sleepForRandSlide.Step = 1.0

	sleepForConstFloat := float64(d.config.GetSleepForConst())
	sleepForConstData := binding.BindFloat(&sleepForConstFloat)
	sleepForConstData.AddListener(binding.NewDataListener(func() {
		data, err := sleepForConstData.Get()
		if err != nil {
			log.Fatalf("error getting sleepForConstData %f, err: %q", sleepForConstData, err)
		}
		d.config.SetSleepForConst(int(data))
	}))
	sleepForConstLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(sleepForConstData, "Sleep for const: %0.2f"))
	sleepForConstEntry := widget.NewEntryWithData(binding.FloatToString(sleepForConstData))
	sleepForConstFloats := container.NewGridWithColumns(2, sleepForConstLabel, sleepForConstEntry)
	sleepForConstSlide := widget.NewSliderWithData(0, 60, sleepForConstData)
	sleepForConstSlide.Step = 1.0

	minMovesFloat := float64(d.config.GetMinMoves())
	minMovesData := binding.BindFloat(&minMovesFloat)
	minMovesData.AddListener(binding.NewDataListener(func() {
		data, err := minMovesData.Get()
		if err != nil {
			log.Fatalf("error getting minMovesData %f, err: %q", minMovesData, err)
		}
		d.config.SetMinMoves(int(data))
	}))
	minMovesLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(minMovesData, "Min moves: %0.2f"))
	minMovesEntry := widget.NewEntryWithData(binding.FloatToString(minMovesData))
	minMovesFloats := container.NewGridWithColumns(2, minMovesLabel, minMovesEntry)
	minMovesSlide := widget.NewSliderWithData(0, 10, minMovesData)
	minMovesSlide.Step = 1.0

	maxMovesFloat := float64(d.config.GetMaxMoves())
	maxMovesData := binding.BindFloat(&maxMovesFloat)
	maxMovesData.AddListener(binding.NewDataListener(func() {
		data, err := maxMovesData.Get()
		if err != nil {
			log.Fatalf("error getting maxMovesData %f, err: %q", maxMovesData, err)
		}
		d.config.SetMaxMoves(int(data))
	}))
	maxMovesLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(maxMovesData, "Max moves: %0.2f"))
	maxMovesEntry := widget.NewEntryWithData(binding.FloatToString(maxMovesData))
	maxMovesFloats := container.NewGridWithColumns(2, maxMovesLabel, maxMovesEntry)
	maxMovesSlide := widget.NewSliderWithData(1, 10, maxMovesData)
	maxMovesSlide.Step = 1.0

	textSizeFloat := float64(d.config.GetTextSize())
	textSizeData := binding.BindFloat(&textSizeFloat)
	textSizeData.AddListener(binding.NewDataListener(func() {
		data, err := textSizeData.Get()
		if err != nil {
			log.Fatalf("error getting textSizeData %f, err: %q", textSizeData, err)
		}
		d.config.SetTextSize(int(data))
	}))
	textSizeLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(textSizeData, "Text size: %0.2f"))
	textSizeEntry := widget.NewEntryWithData(binding.FloatToString(textSizeData))
	textSizeFloats := container.NewGridWithColumns(2, textSizeLabel, textSizeEntry)
	textSizeSlide := widget.NewSliderWithData(10, 40, textSizeData)
	textSizeSlide.Step = 1.0

	moves := d.config.GetMoveStringList(config.DefaultMoveList)
	moveList := binding.BindStringList(&moves)
	list := widget.NewListWithData(
		moveList,
		func() fyne.CanvasObject {
			return container.NewBorder(
				nil,
				nil,
				nil,
				widget.NewButton("-", nil),
				widget.NewLabel(""),
			)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			f := item.(binding.String)
			text := obj.(*fyne.Container).Objects[0].(*widget.Label)
			text.Bind(f)

			btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {
				val, _ := f.Get()
				move, err := config.ParseMove(val)
				if err != nil {
					log.Fatalf("error parsing move value on delete, err: %q", err)
				}
				d.config.DeleteMove(config.DefaultMoveList, move.Name)
				moveList.Set(d.config.GetMoveStringList(config.DefaultMoveList))
			}
		})

	moveName := "Move Name"
	moveNameData := binding.BindString(&moveName)
	moveNameWidget := widget.NewEntryWithData(moveNameData)
	moveCounts := 8
	moveCountsData := binding.BindInt(&moveCounts)
	moveCountsWidget := widget.NewEntryWithData(binding.IntToString(moveCountsData))
	addMoveLine := container.NewGridWithColumns(2, moveNameWidget, moveCountsWidget)
	addMoveButton := widget.NewButton("Add move", func() {
		moveName, err := moveNameData.Get()
		if err != nil {
			log.Fatalf("error getting move name, err: %q", err)
		}
		moveCounts, err := moveCountsData.Get()
		if err != nil {
			log.Fatalf("error getting move counts, err: %q", err)
		}
		move := config.Move{Name: moveName, Counts: moveCounts}
		d.config.AddMove(config.DefaultMoveList, move)
		if err != nil {
			log.Fatal(err)
		}
		moveList.Set(d.config.GetMoveStringList(config.DefaultMoveList))
	})

	saveButton := widget.NewButton("Save configuration", func() {
		err := d.config.Save()
		if err != nil {
			log.Fatalf("error saving data to file, err: %q", err)
		}
	})

	item := container.NewVBox(
		saveButton, widget.NewSeparator(),
		bpmFloats, bpmSlide, widget.NewSeparator(),
		sleepForRandFloats, sleepForRandSlide, widget.NewSeparator(),
		sleepForConstFloats, sleepForConstSlide, widget.NewSeparator(),
		minMovesFloats, minMovesSlide, widget.NewSeparator(),
		maxMovesFloats, maxMovesSlide, widget.NewSeparator(),
		textSizeFloats, textSizeSlide, widget.NewSeparator(),
		addMoveLine, widget.NewSeparator(),
		addMoveButton, widget.NewSeparator(),
	)
	return container.NewBorder(item, nil, nil, nil, list)
}

func (d *Display) movesScreen() fyne.CanvasObject {
	d.movesContext, d.movesCancelFunc = context.WithCancel(context.TODO())

	container := container.NewVBox()
	refresh := func(co color.Color, text ...string) {
		var vboxObjects []fyne.CanvasObject

		for _, t := range text {
			c := canvas.NewText(t, co)
			c.TextSize = float32(d.config.Text.Size)
			c.Alignment = fyne.TextAlignCenter
			vboxObjects = append(vboxObjects, c, widget.NewSeparator())
		}

		container.Objects = vboxObjects
		container.Refresh()
	}

	go func() {
		for {
			sleepFor := rand.Intn(d.config.GetSleepForRand()) + d.config.GetSleepForConst()
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()
		inner:
			for {
				refresh(&color.RGBA{0x00, 0x00, 0xff, 0xff}, fmt.Sprintf("Get ready in %d seconds", sleepFor))
				sleepFor--
				if sleepFor < 0 {
					ticker.Stop()
					break inner
				}

				select {
				case <-d.movesContext.Done():
					return
				case <-ticker.C:
					continue
				}
			}

			tasks := []string{""}
			numSequences := d.config.GetMinMoves() + rand.Intn(d.config.GetMaxMoves()-d.config.GetMinMoves())
			movesCounts := 0
			moveList := d.config.GetMoveList(config.DefaultMoveList)
			for i := 0; i < numSequences; i++ {
				r := rand.Intn(len(moveList))
				tasks = append(tasks, moveList[r].Name)
				movesCounts += moveList[r].Counts
			}
			refresh(&color.RGBA{0xff, 0x00, 0x00, 0xff}, tasks...)

			waitForMoves := float32(movesCounts) / float32(d.config.Song.BPM) * 60
			timer := time.NewTimer(time.Duration(waitForMoves) * time.Second)
			defer timer.Stop()
			select {
			case <-d.movesContext.Done():
				return
			case <-timer.C:
				timer.Stop()
				break
			}
		}
	}()
	return container
}

// Screen defines the data structure for a tutorial
type Screen struct {
	Title, Intro string
	View         func() fyne.CanvasObject
	SupportWeb   bool
}

// directly copied from fyne demo app
func (d *Display) makeNav(setTutorial func(tutorial Screen), loadPrevious bool) fyne.CanvasObject {
	Screens := map[string]Screen{
		"moves": {
			"Moves",
			"",
			d.movesScreen,
			true},
		"options": {"Options",
			"",
			d.optionsScreen,
			true,
		},
	}
	// ScreenIndex defines how screens should be laid out in the index tree
	ScreenIndex := map[string][]string{
		"": {"moves", "options"},
	}

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return ScreenIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := ScreenIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := Screens[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			if unsupportedTutorial(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := Screens[uid]; ok {
				if unsupportedTutorial(t) {
					return
				}
				d.app.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	tree.Select("moves")

	return container.NewBorder(nil, nil, nil, nil, tree)
}

func unsupportedTutorial(t Screen) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}

const preferenceCurrentTutorial = "moves"
