package protracker

import (
	"fmt"
	"testing"
	"github.com/FatmanUK/fatgo/utils"
)

// Test data.
const TST_OK_ROW_NUM_CELLS = `Row has %d Cells.`
const TST_NO_ROW_NUM_CELLS = `Number of Row Cells is wrong.`

const TST_OK_ROW_REGEX = `Regex is ok.`
const TST_NO_ROW_REGEX = `Regex is wrong.`

const TST_OK_ROW_NUM = `Row num is ok.`
const TST_NO_ROW_NUM = `Row num is wrong.`

const TST_OK_ROW_NUM_NO = `Row num is wrong, but that's ok.`
const TST_NO_ROW_NUM_OK = `Row num is ok, but that's wrong.`

const TST_OK_ROW_SAVE = `Save output is ok.`
const TST_NO_ROW_SAVE = `Save output is wrong.`

const TST_OK_ROW_SAVE_NO = `Save errors, but that's ok.`
const TST_NO_ROW_SAVE_OK = `Save doesn't error, but that's wrong.`

const TST_OK_ROW_WRITE = `Write is ok.`
const TST_NO_ROW_WRITE = `Write is wrong.`

const TST_OK_ROW_WRITE_NO = `Write is wrong, but that's ok.`
const TST_NO_ROW_WRITE_OK = `Write is ok, but that's wrong.`

type NumTests struct {
	number     string
	isHex      bool
	expected   uint16
}

// Only defining succeed as it's too simple for fail conditions.
func TestRowFactory_Must_Succeed(t *testing.T) {
	var errMsg [2]string
	r := RowFactory()
	l := len(r)
	errMsg = [2]string{TST_NO_ROW_NUM_CELLS, FMT_TST_NUM}
	if l == CHANNELS_PER_ROW {
		t.Logf(TST_OK_ROW_NUM_CELLS, CHANNELS_PER_ROW)
	} else {
		valPair := [2]int{CHANNELS_PER_ROW, l}
		utils.TErr(t, errMsg, valPair)
	}
	expectedCell := CellFactory()
	errMsg = [2]string{TST_NO_ROW_CELL, FMT_TST_VAR}
	for n := 0; n < l; n++ {
		if r[n] == expectedCell {
			t.Logf(TST_OK_ROW_CELL, n)
		} else {
			valPair := [2]Cell{expectedCell, r[n]}
			utils.TErr(t, errMsg, valPair)
		}
	}
}

// Only defining succeed as it's too simple for fail conditions.
func TestRowRegexFactory_Must_Succeed(t *testing.T) {
	// The regex is so big I've compressed it with zlib!
	exRgxCmp := `eJzSiDbQtXTUdUvUTYutNqrVVNCKtqqJVdDSiHbUddeNj`
	exRgxCmp += `dZVjo021jXVjVWAKtQFKUPiGNfW4FVaE6MHgjW6uro1uq`
	exRgxCmp += `PGEzbeHhAAAP__UIZrTQ`
	expectedRgx, err := utils.Decompress(exRgxCmp)
	var errMsg [2]string
	if err != nil {
		t.Fatalf("Error decompressing string: %v", err)
	}
	r := RowRegexFactory()
	if r == expectedRgx {
		t.Logf(TST_OK_ROW_REGEX)
	} else {
		errMsg = [2]string{TST_NO_ROW_REGEX, FMT_TST_STR}
		valPair := [2]string{expectedRgx, r}
		utils.TErr(t, errMsg, valPair)
	}
}

func TestRowNumFromRowStr_Must_Succeed(t *testing.T) {
	var errMsg [2]string
	tests := []NumTests{
		{"1e", true, 30}, // 1e = 30
		{"20", true, 32}, // 20 = 32
		{"20", false, 20}, // 20 = 20
	}
	errMsg = [2]string{TST_NO_ROW_NUM, FMT_TST_NUM}
	for _, test := range tests {
		out, _ := RowNumFromRowStr(test.number, test.isHex)
		valPair := [2]uint16{out, test.expected}
		if out == test.expected {
			t.Logf(TST_OK_ROW_NUM)
		} else {
			utils.TErr(t, errMsg, valPair)
		}
	}
}

func TestRowNumFromRowStr_Must_Fail(t *testing.T) {
	var errMsg [2]string
	tests := []NumTests{
		{"1e", true, 31}, // 1e = 30
		{"20", true, 31}, // 20 = 32
		{"2x", false, 20}, // 2x = nonsense
	}
	errMsg = [2]string{TST_NO_ROW_NUM_OK, FMT_TST_NUM}
	for _, test := range tests {
		out, _ := RowNumFromRowStr(test.number, test.isHex)
		valPair := [2]uint16{out, test.expected}
		if out != test.expected {
			t.Logf(TST_OK_ROW_NUM_NO)
		} else {
			utils.TErr(t, errMsg, valPair)
		}
	}
}

func TestRowSave_Must_Succeed(t *testing.T) {
	var errMsg [2]string
	r := RowFactory()
	emptyCell := CellFactory()
	blank := emptyCell.Save()
	for rn := uint8(42); rn < 64; rn += 4 { // 42 46 50 54 58 62
		s := r.Save(rn)
		expected := fmt.Sprintf(FMT_ROW_PRETTY, rn,
			blank, blank, blank, blank)
		errMsg = [2]string{TST_NO_ROW_SAVE, FMT_TST_STR}
		if s == expected {
			t.Logf(TST_OK_ROW_SAVE)
		} else {
			valPair := [2]string{expected, s}
			utils.TErr(t, errMsg, valPair)
		}
	}
}

func TestRowSave_Must_Fail(t *testing.T) {
	var errMsg [2]string
	r := RowFactory()
	expected := ""
	for rn := uint8(64); rn < 88; rn += 4 { // 64 68 72 76 80 84
		s := r.Save(rn)
		errMsg = [2]string{TST_NO_ROW_SAVE_OK, FMT_TST_STR}
		if s == expected {
			t.Logf(TST_OK_ROW_SAVE_NO)
		} else {
			valPair := [2]string{expected, s}
			utils.TErr(t, errMsg, valPair)
		}
	}
}

// Don't love this, but it creates a persistent datum in a temporary.
type BytesWriter struct {
	data *[]byte
}

func (b *BytesWriter) Init() {
	b.data = &[]byte{}
}

func (b BytesWriter) Write(p []byte) (int, error) {
	oldLen := len(*(b.data))
	newLen := oldLen + len(p)
	newData := make([]byte, newLen)
	if oldLen > 0 {
		copy(newData[0:oldLen], *(b.data))
	}
	copy(newData[oldLen:], p)
	*(b.data) = newData
	return len(p), nil
}

func (b BytesWriter) Bytes() []byte {
	if b.data == nil {
		return []byte{}
	}
	return *(b.data)
}

func TestRowWrite_Must_Succeed(t *testing.T) {
	var w BytesWriter
	w.Init()

	r := RowFactory()
	r[0] = Cell{"A#4", 1, "F06"}
	r[1] = Cell{"C#4", 2, "C40"}
	r[2] = Cell{"F#4", 3, "484"}
	r[3] = Cell{"G#4", 4, "E00"}
	expected := "[0 240 31 6 1 148 44 64 1 46 52 132 1 13 78 0]"

	patternIndex := uint8(0)
	rowIndex := uint8(0)

	err := r.Write(w, patternIndex, rowIndex) // TODO: test limits
	if err != nil {
		t.Fatalf("Error writing Row: %v", err)
	}

	s := fmt.Sprintf("%v", w.Bytes())
	if s == expected {
		t.Logf(TST_OK_ROW_WRITE)
	} else {
		errMsg := [2]string{TST_NO_ROW_WRITE, FMT_TST_STR}
		valPair := [2]string{expected, s}
		utils.TErr(t, errMsg, valPair)
	}
}

func TestRowWrite_Must_Fail(t *testing.T) {
	var w BytesWriter
	w.Init()

	r := RowFactory()
	r[0] = Cell{"A#4", 1, "F06"}
	r[1] = Cell{"C#4", 2, "C40"}
	r[2] = Cell{"F#4", 3, "484"}
	r[3] = Cell{"G#4", 4, "E00"}
	// TODO put together proper test data
	expected := "[0 240 31 6 1 148 44 64 1 46 52 132 1 13 78 1]"

	patternIndex := uint8(0)
	rowIndex := uint8(0)

	err := r.Write(w, patternIndex, rowIndex) // TODO: test limits
	if err != nil {
		t.Fatalf("Error writing Row: %v", err)
	}

	s := fmt.Sprintf("%v", w.Bytes())
	if s != expected {
		t.Logf(TST_OK_ROW_WRITE_NO)
	} else {
		// TODO
		// logic needs check... tired brain no understands...
		errMsg := [2]string{TST_NO_ROW_WRITE_OK, FMT_TST_STR}
		valPair := [2]string{expected, s}
		utils.TErr(t, errMsg, valPair)
	}
}
