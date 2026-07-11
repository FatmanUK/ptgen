package protracker

type ModInfo struct {
	Title       string
	Speed       uint8
	BPM         uint8
	OrderLen    uint8
	Patterns    uint8
	Instruments uint8
	Samples     uint8
}

// TODO: enforce these limits.
// ModProject represents your input JSON structure.
// Speed: Optional: 1-31 (0: default, always 6)
// BPM: Optional: 32-255 (0: default, always 125)
// Patterns: a slice of Patterns
// Up to 31 instruments
type ModProject struct {
	Title       string       `json:"title"`
	Speed       uint8        `json:"speed"` 
	BPM         uint8        `json:"bpm"`   
	OrderList   []uint8      `json:"orderList" yaml:"orderList"`
	Patterns    []Pattern    `json:"patterns"`
	Instruments []Instrument `json:"instruments"`
	Samples     []Sample     `json:"samples"`
}

func ModProjectFactory() ModProject {
	return ModProject{
		Title:       DEFAULT_TITLE,
		Speed:       DEFAULT_SPEED,
		BPM:         DEFAULT_BPM,
		OrderList:   []uint8{0},
		Patterns:    []Pattern{PatternFactory()},
		Instruments: []Instrument{},
		Samples:     []Sample{},
	}
}

func (re *ModProject) ModInfoFactory() ModInfo {
	return ModInfo{
		Title:       re.Title,
		Speed:       re.Speed,
		BPM:         re.BPM,
		OrderLen:    uint8(len(re.OrderList)),
		Patterns:    uint8(len(re.Patterns)),
		Instruments: uint8(len(re.Instruments)),
		Samples:     uint8(len(re.Samples)),
	}
}

func (re *ModProject) IsOrderListValid() bool {
	l := len(re.OrderList)
	if l == 0 || l > MAX_PATTERNS {
		return false
	}
	for _, j := range re.OrderList {
		if int(j) >= len(re.Patterns) {
			return false
		}
	}
	return true
}
