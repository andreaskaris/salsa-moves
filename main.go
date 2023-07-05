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
		log.Fatalf("error loading config file %s, %q", configFile, err)
	}
	d := display.New(c)
	d.Render()
}
