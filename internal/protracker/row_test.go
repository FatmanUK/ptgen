package protracker

import (
	"fmt"
	"testing"
)

type NumTests struct {
	number     string
	isHex      bool
	expected   uint16
	shouldFail bool
}

func TestRowFactory_Must_Succeed(t *testing.T) {
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

func TestRowFactory_Must_Fail(t *testing.T) {
}

func TestRowRegexFactory_Must_Succeed(t *testing.T) {
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

func TestRowRegexFactory_Must_Fail(t *testing.T) {
}

func TestRowNumFromRowStr_Must_Succeed(t *testing.T) {
	tests := []NumTests{
		{"1e", true, 30, false},
		{"20", true, 32, false},
		{"20", false, 20, false},

		{"1e", true, 31, true},
		{"20", true, 31, true},
		{"2x", false, 20, true},
	}
	for _, test := range tests {
		out, _ := RowNumFromRowStr(test.number, test.isHex)
		if test.shouldFail {
			if out != test.expected {
				t.Logf(MSG_ROW_NUM_NOK)
			} else {
				m := fmt.Sprintf("%s %s",
					NERR_ROW_NUM, FMT_TST_NUMBER)
				t.Errorf(m, out, test.expected)
			}
		} else {
			if out == test.expected {
				t.Logf(MSG_ROW_NUM_OK)
			} else {
				m := fmt.Sprintf("%s %s",
					ERR_ROW_NUM, FMT_TST_NUMBER)
				t.Errorf(m, out, test.expected)
			}
		}
	}
}

func TestRowNumFromRowStr_Must_Fail(t *testing.T) {
}

func TestRowSave_Must_Succeed(t *testing.T) {
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

func TestRowSave_Must_Fail(t *testing.T) {
	//47:func (r *Row) Save(idx uint16) string {
}

func TestRowWrite_Must_Succeed(t *testing.T) {
	//52:func (r *Row) Write(w io.Writer, pid int, rid int) error {
}

func TestRowWrite_Must_Fail(t *testing.T) {
	//52:func (r *Row) Write(w io.Writer, pid int, rid int) error {
}
