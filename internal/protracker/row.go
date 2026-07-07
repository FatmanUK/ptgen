package protracker

import (
	"fmt"
)

// A Row contains exactly CHANNELS_PER_ROW channels.
type Row [CHANNELS_PER_ROW]Cell

func RowFactory() Row {
	r := Row{}
	for i := range r {
		r[i] = CellFactory()
	}
	return r
}

func RowRegexFactory() string {
	rowRgx := `[0-F]{2}`
	noteRgx := `[A-G-][-#][3-5-]`
	instRgx := `[0-F-]{2}`
	efftRgx := `[0-Z-][0-Z-][0-Z-]` // ? TODO refine this one
	cellRgx := fmt.Sprintf(`%s %s %s|%s %s|\.\.\.|---|-`,
			noteRgx, instRgx, efftRgx, noteRgx, instRgx)
	sepRgx := `[:|]`
	rgFmt := `(%s) *%s *(%s) *%s *(%s) *%s *(%s) *%s *(%s) *%s?`
	return fmt.Sprintf(rgFmt,
			rowRgx, sepRgx, cellRgx, sepRgx,
			cellRgx, sepRgx, cellRgx, sepRgx,
			cellRgx, sepRgx)
}

// I could have left the row index out of this, but Rows are never
// created in isolation. Part of the definition of the Row is that it
// always has an identifying number. Paradoxically, it's not integral
// to the Row's nature or function, so I don't include the identifier
// *in* the Row type.
func (re Row) StringFromRow(idx uint16) string {
	return fmt.Sprintf("%02x | %s | %s | %s | %s |", idx,
		re[0].StringFromCell(), re[1].StringFromCell(),
		re[2].StringFromCell(), re[3].StringFromCell())
}
