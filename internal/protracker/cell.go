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
	Note       string `json:"note"`
	Instrument uint8  `json:"instrument"`
	Effect     string `json:"effect"`
}

func CellFactory() Cell {
	c, _ := CellFromString("--- -- ---")
	return c
}

func CellFromString(s string) (Cell, error) {
	c := Cell{
		Note:       "---",
		Instrument: 0,
		Effect:     "---",
	}
	if s == "---" || s == "..." || s == "-" || s == "" {
		return c, nil
	}
	c.Note = s[0:3]
	i := s[4:6]
	if i != "--" {
		u, err := strconv.ParseUint(i, 16, 16)
		if err != nil {
			return c, err
		}
		c.Instrument = uint8(u)
	}
	if len(s) >= 10 {
		c.Effect = s[7:10]
	}
	return c, nil
}

func (re Cell) StringFromCell() string {
	instr := fmt.Sprintf("%02x", re.Instrument)
	if instr == "00" {
		instr = "--"
	}
	return fmt.Sprintf("%s %s %s", re.Note, instr, re.Effect)
}
