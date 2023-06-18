package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Moves []string
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
