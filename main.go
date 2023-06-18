package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/andreaskaris/salsa-moves/pkg/config"
	"github.com/andreaskaris/salsa-moves/pkg/display"
)

const (
	configFile    = "config.yaml"
	sleepForRand  = 15
	sleepForConst = 5
	bpm           = 120
)

func print(c *config.Config, ch chan []string, text ...string) {
	buffer := make([]string, c.Moves.Max)
	for i, t := range text {
		buffer[i] = t
	}
	ch <- buffer
}

func renderText(c *config.Config, ch chan []string) {
	print(c, ch, "Starting ...")
	for {
		sleepFor := rand.Intn(sleepForRand) + sleepForConst
		for range time.Tick(1 * time.Second) {
			print(c, ch, fmt.Sprintf("Get ready in %d seconds", sleepFor))
			sleepFor--
			if sleepFor < 0 {
				break
			}
		}

		tasks := make([]string, c.Moves.Max)
		numSequences := c.Moves.Min + rand.Intn(c.Moves.Max-c.Moves.Min)
		movesCounts := 0
		for i := 0; i < numSequences; i++ {
			r := rand.Intn(len(c.Moves.List))
			tasks[i] = c.Moves.List[r].Name
			movesCounts += c.Moves.List[r].Counts
		}
		print(c, ch, tasks...)
		time.Sleep(time.Duration(float32(movesCounts)/bpm*60) * time.Second)
	}
}

func main() {
	c, err := config.Load(configFile)
	if err != nil {
		log.Fatal(err)
	}
	d := display.New(c.Text.Size)
	go renderText(c, d.Ch)
	d.Render()

}
