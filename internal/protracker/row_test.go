package protracker

import (
	"fmt"
	"testing"
)

func TestRowFactory(t *testing.T) {
	r := RowFactory()
	l := len(r)
	if l == CHANNELS_PER_ROW {
		t.Logf(MSG_ROW_CELLS_OK, CHANNELS_PER_ROW)
	} else {
		m := fmt.Sprintf("%s %s", ERR_ROW_CELLS,
			FMT_TST_STRING)
		t.Errorf(m, CHANNELS_PER_ROW, l)
	}
	expectedCell := CellFactory()
	for n := 0; n < l; n++ {
		if r[n] == expectedCell {
			t.Logf(MSG_ROW_CELLX_OK, n)
		} else {
			m := fmt.Sprintf("%s %s", ERR_ROW_CELLX,
				FMT_TST_STRING)
			t.Errorf(m, n, expectedCell, r[n])
		}
	}
}

func TestRowRegexFactory(t *testing.T) {
	cellRgx := CellRegexFactory()
	expectedRgx := fmt.Sprintf(RGX_ROW, RGX_ROWNUM,
		RGX_CELL_SEPARATOR, cellRgx,
		RGX_CELL_SEPARATOR, cellRgx,
		RGX_CELL_SEPARATOR, cellRgx,
		RGX_CELL_SEPARATOR, cellRgx,
		RGX_CELL_SEPARATOR)
	r := RowRegexFactory()
	if r == expectedRgx {
		t.Logf(MSG_ROW_REGEX_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_ROW_REGEX,
			FMT_TST_STRING)
		t.Errorf(m, expectedRgx, r)
	}
}

func TestRowSave(t *testing.T) {
	r := RowFactory()
	s := r.Save(TEST_ROW_ROWNUM)
	emptyCell := CellFactory()
	blank := emptyCell.Save()
	expected := fmt.Sprintf(FMT_ROW_PRETTY, TEST_ROW_ROWNUM,
		blank, blank, blank, blank)
	if s == expected {
		t.Logf(MSG_ROW_SAVE_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_ROW_SAVE,
			FMT_TST_STRING)
		t.Errorf(m, expected, s)
	}
}
