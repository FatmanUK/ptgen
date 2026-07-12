package protracker

import (
	"bytes"
	"strings"
	"testing"
)

// Mostly written by Google Gemini Pro. Tweaked extensively by me.

type CellTestData struct {
	name     string
	input    Cell
	expected [4]byte
}

type CellTestError struct {
	name        string
	input       Cell
	expectedErr string
}

type ModProjectTestError struct {
	name        string
	project     ModProject
	shouldFail  bool
	expectedLen int
}

// TestEncodeCell_Valid Matrix validates big-endian bitpacking output
func TestEncodeCell_Valid(t *testing.T) {
	t.Parallel()
	var tests []CellTestData
	cellBlank := Cell{Note: "", Instr: 0, Effect: ""}
	cellEmpty := Cell{Note: "---", Instr: 0, Effect: "---"}
	cellLowBound := Cell{Note: "C-3", Instr: 0, Effect: ""}
	cellMidNote := Cell{Note: "C-4", Instr: 0, Effect: ""}
	cellHighNote := Cell{Note: "B-5", Instr: 0, Effect: ""}
	cellCase := Cell{Note: "g#4", Instr: 0, Effect: ""}
	cellEffectOnly := Cell{Note: "", Instr: 0, Effect: "C40"}
	cellEffectCase := Cell{Note: "", Instr: 0, Effect: "f03"}
	cellInstrBound := Cell{Note: "", Instr: 31, Effect: ""}
	cellPacked := Cell{Note: "C-4", Instr: 17, Effect: "C40"}
	tests = append(tests, CellTestData{
		name:     "Completely blank cell",
		input:    cellBlank,
		expected: [4]byte{0x00, 0x00, 0x00, 0x00},
	})
	tests = append(tests, CellTestData{
		name:     "Tracker-style empty placeholders",
		input:    cellEmpty,
		expected: [4]byte{0x00, 0x00, 0x00, 0x00},
	})
	tests = append(tests, CellTestData{
		name:     "Low-bound note (C-3)",
		input:    cellLowBound,
		expected: [4]byte{0x03, 0x58, 0x00, 0x00},
	})
	tests = append(tests, CellTestData{
		name:     "Mid-range note (C-4)",
		input:    cellMidNote,
		expected: [4]byte{0x01, 0xAC, 0x00, 0x00},
	})
	tests = append(tests, CellTestData{
		name:     "High-bound note (B-5)",
		input:    cellHighNote,
		expected: [4]byte{0x00, 0x71, 0x00, 0x00},
	})
	tests = append(tests, CellTestData{
		name:     "Note formatting case insensitivity",
		input:    cellCase,
		expected: [4]byte{0x01, 0x0D, 0x00, 0x00},
	})
	tests = append(tests, CellTestData{
		name:     "Effect Only (C40 - Max Volume)",
		input:    cellEffectOnly,
		expected: [4]byte{0x00, 0x00, 0x0C, 0x40},
	})
	tests = append(tests, CellTestData{
		name:     "Effect Only lower-case hex parsing",
		input:    cellEffectCase,
		expected: [4]byte{0x00, 0x00, 0x0F, 0x03},
	})
	tests = append(tests, CellTestData{
		name:     "Instrument boundary upper limit (31)",
		input:    cellInstrBound,
		expected: [4]byte{0x10, 0x00, 0xF0, 0x00},
		// 0x1F split: 0x10 in byte 0, 0xF0 in byte 2
	})
	tests = append(tests, CellTestData{
		name:     "Full Packed State (C-4/Ins 17/Eff C40)",
		input:    cellPacked,
		expected: [4]byte{0x11, 0xAC, 0x1C, 0x40},
		// Structural Verification:
		// Ins=0x11, Period=0x01AC, Cmd=0x0C, Param=0x40
		// B0: Ins High (0x10) | Period High (0x01) = 0x11
		// B1: Period Low = 0xAC
		// B2: Ins Low (0x01<<4 = 0x10) | Cmd (0x0C) = 0x1C
		// B3: Param = 0x40
	})

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := encodeCell(tt.input)
			if err != nil {
				t.Fatalf(ERR_CELL_FAILED, err)
			}
			if result != tt.expected {
				t.Errorf(FMT_TST_VAR_X, tt.expected,
					result)
			}
		})
	}
}

// TestEncodeCell_Errors ensures validation blocks spec-breaking input
func TestEncodeCell_Errors(t *testing.T) {
	t.Parallel()
	var tests []CellTestError
	cellTooLow := Cell{Note: "B-2"}
	cellTooHigh := Cell{Note: "C-6"}
	cellRedDwarfError := Cell{Note: "H-4"}
	cellInstrTooHigh := Cell{Instr: 32}
	cellEffectTooShort := Cell{Effect: "C4"}
	cellEffectTooLong := Cell{Effect: "C400"}
	cellEffectWtf := Cell{Effect: "G00"}
	cellEffectMangled := Cell{Effect: "C0X"}
	tests = append(tests, CellTestError{
		name:        "Octave below limit (B-2)",
		input:       cellTooLow,
		expectedErr: "Must be a valid note in octaves 3-5",
	})
	tests = append(tests, CellTestError{
		name:        "Octave above limit (C-6)",
		input:       cellTooHigh,
		expectedErr: "Must be a valid note in octaves 3-5",
	})
	tests = append(tests, CellTestError{
		name:        "Invalid note format name",
		input:       cellRedDwarfError,
		expectedErr: "Must be a valid note in octaves 3-5",
	})
	tests = append(tests, CellTestError{
		name:        "Instrument range overflow (>31)",
		input:       cellInstrTooHigh,
		expectedErr: "Instrument 32 is wrong",
	})
	tests = append(tests, CellTestError{
		name:        "Effect string length underflow",
		input:       cellEffectTooShort,
		expectedErr: "Must be 3 hex characters",
	})
	tests = append(tests, CellTestError{
		name:        "Effect string length overflow",
		input:       cellEffectTooLong,
		expectedErr: "Must be 3 hex characters",
	})
	tests = append(tests, CellTestError{
		name:        "Non-hex effect command byte",
		input:       cellEffectWtf,
		expectedErr: "Invalid effect command",
	})
	tests = append(tests, CellTestError{
		name:        "Non-hex effect parameter bytes",
		input:       cellEffectMangled,
		expectedErr: "Invalid effect parameter",
	})

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := encodeCell(tt.input)
			if err == nil {
				t.Fatalf(ERR_CELL_SUCCESS)
			}
			hasSubstr := strings.Contains(err.Error(),
				tt.expectedErr)
			if !hasSubstr {
				t.Errorf(ERR_CELL_WRONG_ERROR,
					tt.expectedErr, err)
			}
		})
	}
}

// TestWriteMod_Integration ensures structural math and pattern
// layouts align perfectly
func TestWriteMod_Integration(t *testing.T) {
	// Generate pattern mock payload (64 rows, 4 channels)
	mockPattern := Pattern{}
	cellMock := Cell{Note: "C-4", Instr: 1, Effect: "000"}
	for i := 0; i < 64; i++ {
		mockPattern[i] = Row{
			cellMock, Cell{}, Cell{}, Cell{},
		}
	}

	var tests []ModProjectTestError
	mpCompliant := ModProject{
		Title:     "Pro Validation",
		OrderList: []uint8{0, 0, 1},
		Patterns:  []Pattern{mockPattern, mockPattern},
	}
	mpEmptyOrderList := ModProject{
		Title:     "Order List Too Short (Empty)",
		OrderList: []uint8{},
		Patterns:  []Pattern{mockPattern},
	}
	mpOrderListTooLong := ModProject{
		Title:     "Order List Too Long",
		OrderList: make([]uint8, 129), // Cap is 128
		Patterns:  []Pattern{mockPattern},
	}
	tests = append(tests, ModProjectTestError{
		name:        "Standard compliant export compilation",
		project:     mpCompliant,
		shouldFail:  false,
		expectedLen: 1084 + (2 * 1024), // Header + 2 patterns
	})
	tests = append(tests, ModProjectTestError{
		name:       "Reject order tracking length underflow",
		project:    mpEmptyOrderList,
		shouldFail: true,
	})
	tests = append(tests, ModProjectTestError{
		name:       "Reject order tracking length overflow",
		project:    mpOrderListTooLong,
		shouldFail: true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteMod(&buf, &tt.project)

			if tt.shouldFail {
				if err == nil {
					t.Errorf(ERR_MOD_CORRUPT)
				}
				return
			}

			if err != nil {
				t.Fatalf(ERR_MOD_FAILED, err)
			}

			// Size Assertion
			if buf.Len() != tt.expectedLen {
				t.Errorf(ERR_MOD_SIZE, tt.expectedLen,
					buf.Len(),
				)
			}

			outBytes := buf.Bytes()
			// Check title padding execution
			title := []byte("Pro Validation")
			expTitle := append(title, make([]byte, 6)...)
			if !bytes.Equal(outBytes[0:20], expTitle) {
				t.Errorf(ERR_MOD_TITLE_PAD)
			}

			// Validate Magic Marker Offset placement
			// 20 (Title) + 930 (Samples) + 1 (Len)
			// + 1 (Restart) + 128 (Seq Array)
			// = Offset 1080
			magicMarker := string(outBytes[1080:1084])
			if magicMarker != MAGIC_BYTES {
				t.Errorf(ERR_MOD_MAGIC, magicMarker)
			}

			// Sample header integrity loop check
			// Ensure empty sample arrays maintain default
			// loop size lengths (Byte 29 must be 0x01)
			for i := 0; i < 31; i++ {
				offset := 20 + (i * 30) + 29
				if outBytes[offset] != 0x01 {
					t.Errorf(ERR_MOD_INSTR, i+1)
				}
			}
		})
	}
}

func TestPackBytes(t *testing.T) {
	// in:  b0101 0101 bXXXX 0001  0101 0101 bXXXX 0101 b01010101
	//    0hi^^^^ ^^^^2hi 0lo^^^^  ^^^^ ^^^^1   2lo^^^^ 3^^^^^^^^
	// out: b0101 0001 b0101 0101 b0101 0101 b0101 0101
	//expectedBytes := []byte{81, 85, 85, 85}
	expectedBytes := []byte{0x51, 0x55, 0x55, 0x55}
	//output := packBytes(85, 341, 5, 85)
	output := packBytes(0x55, 0x155, 0x5, 0x55)
	if compareSlices(output[:], expectedBytes) {
		t.Logf(`Packing is ok.`)
	} else {
		t.Errorf(`Failed. Expected %v, got %v.`,
			expectedBytes, output)
	}
}

func TestWriteMagic(t *testing.T) {
	expectedBytes := []byte{'M', '.', 'K', '.'}
	buf := bytes.Buffer{}
	err := writeMagic(&buf)
	if err != nil {
		t.Errorf(`Failed. %v`, err)
	}
	b := buf.Bytes()
	if compareSlices(b, expectedBytes) {
		t.Logf(`Magic is ok.`)
	} else {
		t.Errorf(`Failed. Expected %v, got %v.`,
			expectedBytes, b)
	}
}

func TestParseNote(t *testing.T) {
	//func parseNote(n string) (uint16, bool) {
}

func TestParseEffect(t *testing.T) {
	//func parseEffect(e string) (uint8, uint8, error) {
}

func TestPrepareCommands(t *testing.T) {
	//func prepareCommands(s uint8, b uint8) ([]string, error) {
}

func TestInjectCommands(t *testing.T) {
	//func injectCommands(p *ModProject, i uint8, commands []string) uint8 {
}

func TestInjectInitialTempo(t *testing.T) {
	//func injectInitialTempo(proj *ModProject) error {
}

func TestEncodeInstrumentHeader(t *testing.T) {
	//func encodeInstrumentHeader(inst Instrument) ([30]byte, error) {
}

func TestWriteCell(t *testing.T) {
	//func writeCell(w io.Writer, row Row, p int, r int) error {
}

func TestInjectAndCheck(t *testing.T) {
	//func injectAndCheck(proj *ModProject) error {
}

func TestIsSlotPopulated(t *testing.T) {
	//func isSlotPopulated(slot *Instrument) bool {
}

func TestPreProcessInstruments(t *testing.T) {
	//func preProcessInstruments(proj *ModProject) ([31]Instrument, error) {
}

func TestWriteTitle(t *testing.T) {
	//func writeTitle(w io.Writer, title string) error {
}

func TestWriteSampleHeader(t *testing.T) {
	//func writeSampleHeader(w io.Writer, instr *Instrument) error {
}

func TestWriteSampleHeaders(t *testing.T) {
	//func writeSampleHeaders(w io.Writer, slots [31]Instrument) error {
}

func TestWriteOrderList(t *testing.T) {
	//func writeOrderList(w io.Writer, orderList []uint8) error {
}

func TestWritePatterns(t *testing.T) {
	//func writePatterns(w io.Writer, pttns []Pattern) error {
}

func TestWriteSamples(t *testing.T) {
	//func writeSamples(w io.Writer, slots [31]Instrument) error {
}

func TestWriteHeaders(t *testing.T) {
	//func writeHeaders(w io.Writer,
}

func TestExtractAndInjectSamples(t *testing.T) {
	//func extractAndInjectSamples(proj *ModProject) ([]byte, error) {
}
