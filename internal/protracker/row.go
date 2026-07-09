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
	cellRgx := CellRegexFactory()
	return fmt.Sprintf(RGX_ROW, RGX_ROWNUM, RGX_CELL_SEPARATOR,
		cellRgx, RGX_CELL_SEPARATOR,
		cellRgx, RGX_CELL_SEPARATOR,
		cellRgx, RGX_CELL_SEPARATOR,
		cellRgx, RGX_CELL_SEPARATOR)
}

// I could have left the row index out of this, but Rows are never
// created in isolation. Part of the definition of the Row is that it
// always has an identifying number. Paradoxically, it's not integral
// to the Row's nature or function, so I don't include the identifier
// *in* the Row type.
// Save: serialise struct to a string
func (re *Row) Save(idx uint16) string {
	return fmt.Sprintf(FMT_ROW_PRETTY, idx,
		re[0].Save(), re[1].Save(),
		re[2].Save(), re[3].Save())
}
