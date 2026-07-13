package protracker

import (
	"fmt"
)

// Start, Length and End do not refer to the Instrument but rather the
// forward loop (if present). Obviously, Start < End,
// Start + Length = End, and End < file size.
// Max file size = 65535*2 = 131070, needs a uint32 to store them.
// Start/Length = 0 for no loop.
// ID: The target slot (1-31)
// Volume: 0-64
// Finetune: 0-15 (0=0, 1=1... 8=-8... 15=-1)
// Start: Offset in BYTES
// Length: Length in BYTES (0 means no loop)
// Data: Raw 8-bit SIGNED PCM
type Instrument struct {
	ID       uint8  `json:"id"`
	Source   string `json:"source"`
	Name     string `json:"name"`
	Volume   uint8  `json:"volume"`
	Finetune uint8  `json:"finetune"`
	Start    uint32 `json:"start"`
	Length   uint32 `json:"length"`
	End      uint32 `json:"end"`
	Data     []byte `json:"data"`
}

func InstrumentFactory() Instrument {
	return Instrument{
		ID:       0,
		Source:   "",
		Name:     "",
		Volume:   64,
		Finetune: 0,
		Start:    0,
		Length:   0,
		End:      0,
		Data:     []byte{0},
	}
}

func (i *Instrument) CalculateLength() {
	if i.Start == 0 {
		i.Length = 0
		i.End = 0
	}
	if i.Start > 0 && i.Length == 0 {
		i.Length = (i.End - i.Start)
	}
	if i.Start > 0 && i.End == 0 {
		i.End = (i.Length + i.Start)
	}
}

// Save: serialise struct to a string
func (i *Instrument) Save() string {
	m := `id:{%d} source:{%s} name:{%s} volume:{%d} finetune:{%d}`
	m += ` start:{%d} length:{%d} end:{%d}`
	return fmt.Sprintf(m, i.ID, i.Source, i.Name, i.Volume,
		i.Finetune, i.Start, i.Length, i.End)
}
