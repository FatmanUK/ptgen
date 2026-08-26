package protracker

import (
)

type Instrument struct {
	ID     uint8  `yaml:"id"`
	Volume uint8  `yaml:"volume"`
	Source string `yaml:"source"`
	Name   string `yaml:"name"`
}
