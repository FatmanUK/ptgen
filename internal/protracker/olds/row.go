package protracker

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Regexes.
const RGX_ROWNUM = `[0-9A-Fa-f]{2}`
const RGX_CELL_SEP = `[:|]`

// Formatting patterns.
const FMT_LOCATION = `Pattern %d, Row %d, Channel %d: %w`

type Row struct {
	Cells []Cell
}

func RowFactory(numChannels uint8) Row {
	r := Row{
		Cells: make([]Cell, numChannels),
	}
	for i := range r.Cells {
		r.Cells[i] = CellFactory()
	}
	return r
}

func RowRegexFactory(numChannels uint8) string {
	cellRgx := CellRegexFactory()
	var rv []string
	rv[0] = fmt.Sprintf("(%s)", RGX_ROWNUM)
	for ; numChannels > 0; numChannels-- {
		rv = append(rv, fmt.Sprintf("*%s", RGX_CELL_SEP))
		rv = append(rv, fmt.Sprintf("*(%s)", cellRgx))
	}
	rv = append(rv, fmt.Sprintf("*%s?", RGX_CELL_SEP))
	return strings.Join(rv, " ")
}

func RowNumFromRowStr(rowStr string, isHex bool) (uint8, error) {
	rowBase := 10
	if isHex {
		rowBase = 16
	}
	u, err := strconv.ParseUint(rowStr, rowBase, 8)
	if err != nil {
		return 0, err
	}
	return uint8(u), nil
}

// I could have left the row index out of this, but Rows are never
// created in isolation. Part of the definition of the Row is that it
// always has an identifying number. Paradoxically, it's not integral
// to the Row's nature or function, so I don't include the identifier
// *in* the Row type.
// Save: serialise struct to a string
func (r *Row) Save(idx uint8, chans uint8) string {
	rowArr := []string{}
	rowArr = append(rowArr, fmt.Sprintf("(%02x)", idx))
	for n := uint8(0); n < chans; n++ {
		c := r.Cells[n].Save()
		rowArr = append(rowArr, fmt.Sprintf(" %s |", c))
	}
	return strings.Join(rowArr, " ")
}

func (r *Row) Write(w io.Writer, pi uint8, ri uint8) error {
	for ci, cell := range r.Cells {
		err := cell.Write(w)
		if err != nil {
			m := fmt.Errorf(FMT_LOCATION, pi, ri, ci, err)
			return m
		}
	}
	return nil
}
