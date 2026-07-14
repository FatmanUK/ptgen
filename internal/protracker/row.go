package protracker

import (
	"fmt"
	"io"
	"strconv"
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

func RowNumFromRowStr(rowStr string, isHex bool) (uint16, error) {
	rowBase := 10
	if isHex {
		rowBase = 16
	}
	u, err := strconv.ParseUint(rowStr, rowBase, 16)
	if err != nil {
		return 0, err
	}
	return uint16(u), nil
}

// I could have left the row index out of this, but Rows are never
// created in isolation. Part of the definition of the Row is that it
// always has an identifying number. Paradoxically, it's not integral
// to the Row's nature or function, so I don't include the identifier
// *in* the Row type.
// Save: serialise struct to a string
func (r *Row) Save(idx uint16) string {
	return fmt.Sprintf(FMT_ROW_PRETTY, idx, r[0].Save(),
		r[1].Save(), r[2].Save(), r[3].Save())
}

func (r *Row) Write(w io.Writer, pid int, rid int) error {
	for cid, c := range *r {
		err := c.Write(w)
		if err != nil {
			m := fmt.Errorf(ERR_LOCATION,
				pid, rid, cid, err)
			return m
		}
	}
	return nil
}
