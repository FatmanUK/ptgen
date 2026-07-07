package protracker

type ModInfo struct {
	Title    string
	Speed    uint8
	BPM      uint8
	SequenceLen int
	PatternsLen int
}

// ModProject represents your input JSON structure.
// Speed: Optional: 1-31 (0: default, always 6)
// BPM: Optional: 32-255 (0: default, always 125)
// Patterns: a slice of Patterns
type ModProject struct {
	Title    string    `json:"title"`
	Speed    uint8     `json:"speed"`
	BPM      uint8     `json:"bpm"`
	Sequence []uint8   `json:"orderList"`
	Patterns []Pattern `json:"patterns"`
}

func ModProjectFactory() ModProject {
	return ModProject{
		Title:    "A Song With No Name",
		Speed:    DEFAULT_SPEED,
		BPM:      DEFAULT_BPM,
		Sequence: []uint8{0},
		Patterns: []Pattern{PatternFactory()},
	}
}

func (re ModProject) ModInfoFactory() ModInfo {
	return ModInfo{
		Title: re.Title,
		Speed: re.Speed,
		BPM: re.BPM,
		SequenceLen: len(re.Sequence),
		PatternsLen: len(re.Patterns),
	}
}

func (re ModProject) IsOrderListValid() bool {
	// TODO: check that each entry in the order list refers to
	// a pattern.
	return true
}
