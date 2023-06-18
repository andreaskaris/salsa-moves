package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Moves struct {
	List []Move
	Max  int
	Min  int
}

type Move struct {
	Name   string
	Counts int
}

type Text struct {
	Size float32
}

type Config struct {
	Moves Moves
	Text  Text
}

func Load(fileName string) (*Config, error) {
	c := &Config{}
	fContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s, err: %q", fileName, err)
	}
	if err := yaml.Unmarshal(fContent, c); err != nil {
		return nil, fmt.Errorf("error parsing YAML file %s, err: %q", fileName, err)
	}
	return c, nil
}
