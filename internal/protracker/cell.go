package protracker

import (
	"fmt"
	"strconv"
)

// A Cell represents one channel on one row.
// Note: e.g., "C-4" or "---" or ""
// Instrument: 1-31 (0 means no instrument)
// Effect: e.g., "C40", "047", "F03", or ""
type Cell struct {
	Note   string `json:"note"`
	Instr  uint8  `json:"instr"`
	Effect string `json:"effect"`
}

func CellFactory() Cell {
	return Cell{
		Note:   "---",
		Instr:  0,
		Effect: "---",
	}
}

func CellRegexFactory() string {
	return fmt.Sprintf(RGX_CELL,
		RGX_NOTE, RGX_INSTR, RGX_EFFECT,
		RGX_NOTE, RGX_INSTR)
}

// Load: deserialise struct from a string
func (c *Cell) Load(s string) error {
	if s == "---" || s == "..." || s == "-" || s == "" {
		return nil
	}
	l := len(s)
	if l >= 3 {
		c.Note = s[0:3]
	}
	if l >= 6 {
		i := s[4:6]
		if i != "--" {
			u, err := strconv.ParseUint(i, 16, 16)
			if err != nil {
				return err
			}
			c.Instr = uint8(u)
		}
	}
	if l >= 10 {
		c.Effect = s[7:10]
	}
	return nil
}

// Save: serialise struct to a string
func (c *Cell) Save() string {
	i := fmt.Sprintf("%02x", c.Instr)
	if i == "00" {
		i = "--"
	}
	return fmt.Sprintf("%s %s %s", c.Note, i, c.Effect)
}
