package protracker

import (
	"testing"
)

/*
// A Cell represents one channel on one row.
// Note: e.g., "C-4" or "---" or ""
// Instrument: 1-31 (0 means no instrument)
// Effect: e.g., "C40", "047", "F03", or ""
type Cell struct {
	Note       string `json:"note"`
	Instrument uint8  `json:"instrument"`
	Effect     string `json:"effect"`
}
*/

func TestCellFactory(t *testing.T) {
	// func CellFactory() Cell
}

func TestCellFromString(t *testing.T) {
	// func CellFromString(s string) (Cell, error)
}

func TestStringFromCell(t *testing.T) {
	// func (re Cell) StringFromCell() string
}
