package main

import (
	"log"

	"github.com/andreaskaris/salsa-moves/pkg/config"
	"github.com/andreaskaris/salsa-moves/pkg/display"
)

const (
	configFile = "config.yaml"
)

func main() {
	c, err := config.Load(configFile)
	if err != nil {
		log.Fatal(err)
	}
	d := display.New(c)
	d.Render()

}
