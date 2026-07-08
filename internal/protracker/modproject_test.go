package protracker

import (
	"testing"
)

/*
type ModInfo struct {
	Title       string
	Speed       uint8
	BPM         uint8
	SequenceLen uint8
	PatternsLen uint8
}

// ModProject represents your input JSON structure.
// Speed: Optional: 1-31 (0: default, always 6)
// BPM: Optional: 32-255 (0: default, always 125)
// Patterns: a slice of Patterns
type ModProject struct {
	Title    string    `json:"title"`
	Speed    uint8     `json:"speed"`
	BPM      uint8     `json:"bpm"`
	Sequence []uint8   `json:"orderList" yaml:"orderList"`
	Patterns []Pattern `json:"patterns"`
}
*/

func TestModProjectFactory(t *testing.T) {
	// func ModProjectFactory() ModProject
}

func TestModInfoFactory(t *testing.T) {
	// func (re ModProject) ModInfoFactory() ModInfo
}

func TestIsOrderListValid(t *testing.T) {
	// func (re ModProject) IsOrderListValid() bool
}
