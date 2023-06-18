package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/andreaskaris/salsa-moves/pkg/config"
	"github.com/andreaskaris/salsa-moves/pkg/display"
)

var ()

const (
	configFile = "config.yaml"
)

func print(text string, ch chan string) {
	ch <- text
}

func renderText(c *config.Config, ch chan string) {
	print("Starting ...", ch)
	for {
		sleepFor := rand.Intn(30) + 10
		for range time.Tick(1 * time.Second) {
			print(fmt.Sprintf("Get ready in %d seconds", sleepFor), ch)
			sleepFor--
			if sleepFor < 0 {
				break
			}
		}

		tasks := []string{}
		numSequences := rand.Intn(2) + 2
		for i := 0; i < numSequences; i++ {
			r := rand.Intn(len(c.Moves))
			tasks = append(tasks, c.Moves[r])
		}
		for i := 0; i < len(tasks); i++ {
			print(tasks[i], ch)
		}
		time.Sleep(10 * time.Duration(len(tasks)) * time.Second)
	}
}

func main() {
	c, err := config.Load(configFile)
	if err != nil {
		log.Fatal(err)
	}
	d := display.New()
	go renderText(c, d.Ch)
	d.Render()

}
