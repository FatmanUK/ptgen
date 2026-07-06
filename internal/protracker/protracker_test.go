package protracker

import (
	"bytes"
	"strings"
	"testing"
)

// Mostly written by Google Gemini Pro. Tweaked extensively by me.

const MSG_ERR_FAILED = "Failed handling valid input: %v"
const MSG_ERR_EXPECTED = "\nExpected: [% X]\nGot:      [% X]"
const MSG_ERR_FAIL_FAIL = "Expected constraint violation error, but execution passed"
const MSG_ERR_WRONG_ERROR = "Expected error message containing '%s', got '%v'"
const MSG_ERR_FAIL_CRRPT = "Failed to notice a corrupt payload setup"
const MSG_ERR_FAIL_VALID = "Compilation pipeline failed on valid payload setup: %v"
const MSG_ERR_SIZE = "Binary payload footprint mismatch. Expected %d bytes, got %d"
const MSG_ERR_TITLE = "Title field error. Expected custom byte-padding format layout"
const MSG_ERR_BAD_FMT = "Format validation failure. 'M.K.' magic bytes missing from target offset 1080 (got %q)"
const MSG_ERR_BAD_INS = "Instrument sample slot header #%d failed configuration constraints setup"

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
	cellBlank := Cell{Note: "", Instrument: 0, Effect: ""}
	cellEmpty := Cell{Note: "---", Instrument: 0, Effect: "---"}
	cellLowBound := Cell{Note: "C-3", Instrument: 0, Effect: ""}
	cellMidNote := Cell{Note: "C-4", Instrument: 0, Effect: ""}
	cellHighNote := Cell{Note: "B-5", Instrument: 0, Effect: ""}
	cellCase := Cell{Note: "g#4", Instrument: 0, Effect: ""}
	cellEffectOnly := Cell{Note: "", Instrument: 0, Effect: "C40"}
	cellEffectCase := Cell{Note: "", Instrument: 0, Effect: "f03"}
	cellInstrBound := Cell{Note: "", Instrument: 31, Effect: ""}
	cellPacked := Cell{Note: "C-4", Instrument: 17, Effect: "C40"}
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
				t.Fatalf(MSG_ERR_FAILED,
					err,
				)
			}
			if result != tt.expected {
				t.Errorf(MSG_ERR_EXPECTED,
					tt.expected,
					result,
				)
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
	cellInstrTooHigh := Cell{Instrument: 32}
	cellEffectTooShort := Cell{Effect: "C4"}
	cellEffectTooLong := Cell{Effect: "C400"}
	cellEffectWtf := Cell{Effect: "G00"}
	cellEffectMangled := Cell{Effect: "C0X"}
	tests = append(tests, CellTestError{
		name:        "Octave below limit (B-2)",
		input:       cellTooLow,
		expectedErr: "only octaves 3-5 allowed",
	})
	tests = append(tests, CellTestError{
		name:        "Octave above limit (C-6)",
		input:       cellTooHigh,
		expectedErr: "only octaves 3-5 allowed",
	})
	tests = append(tests, CellTestError{
		name:        "Invalid note format name",
		input:       cellRedDwarfError,
		expectedErr: "invalid note",
	})
	tests = append(tests, CellTestError{
		name:        "Instrument range overflow (>31)",
		input:       cellInstrTooHigh,
		expectedErr: "instrument 32 out of bounds",
	})
	tests = append(tests, CellTestError{
		name:        "Effect string length underflow",
		input:       cellEffectTooShort,
		expectedErr: "must be 3 hex chars",
	})
	tests = append(tests, CellTestError{
		name:        "Effect string length overflow",
		input:       cellEffectTooLong,
		expectedErr: "must be 3 hex chars",
	})
	tests = append(tests, CellTestError{
		name:        "Non-hex effect command byte",
		input:       cellEffectWtf,
		expectedErr: "invalid effect command",
	})
	tests = append(tests, CellTestError{
		name:        "Non-hex effect parameter bytes",
		input:       cellEffectMangled,
		expectedErr: "invalid effect parameter",
	})

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := encodeCell(tt.input)
			if err == nil {
				t.Fatalf(MSG_ERR_FAIL_FAIL)
			}
			hasSubstr := strings.Contains(err.Error(),
				tt.expectedErr,
			)
			if !hasSubstr {
				t.Errorf(MSG_ERR_WRONG_ERROR,
					tt.expectedErr,
					err,
				)
			}
		})
	}
}

// TestWriteMod_Integration ensures structural math and pattern
// layouts align perfectly
func TestWriteMod_Integration(t *testing.T) {
	// Generate pattern mock payload (64 rows, 4 channels)
	mockPattern := Pattern{}
	cellMock := Cell{Note: "C-4", Instrument: 1, Effect: "000"}
	for i := 0; i < 64; i++ {
		mockPattern[i] = Row{
			cellMock, Cell{}, Cell{}, Cell{},
		}
	}

	var tests []ModProjectTestError
	mpCompliant := ModProject{
		Title:    "Pro Validation",
		Sequence: []uint8{0, 0, 1},
		Patterns: []Pattern{mockPattern, mockPattern},
	}
	mpEmptyOrderList := ModProject{
		Title:    "Order List Too Short (Empty)",
		Sequence: []uint8{},
		Patterns: []Pattern{mockPattern},
	}
	mpOrderListTooLong := ModProject{
		Title:    "Order List Too Long",
		Sequence: make([]uint8, 129), // Cap is 128
		Patterns: []Pattern{mockPattern},
	}
	tests = append(tests, ModProjectTestError{
		name:        "Standard compliant export compilation",
		project:     mpCompliant,
		shouldFail:  false,
		expectedLen: 1084 + (2 * 1024), // Header + 2 patterns
	})
	tests = append(tests, ModProjectTestError{
		name:       "Reject sequence tracking length underflow",
		project:    mpEmptyOrderList,
		shouldFail: true,
	})
	tests = append(tests, ModProjectTestError{
		name:       "Reject sequence tracking length overflow",
		project:    mpOrderListTooLong,
		shouldFail: true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteMod(&buf, &tt.project)

			if tt.shouldFail {
				if err == nil {
					t.Errorf(MSG_ERR_FAIL_CRRPT)
				}
				return
			}

			if err != nil {
				t.Fatalf(MSG_ERR_FAIL_VALID, err)
			}

			// Size Assertion
			if buf.Len() != tt.expectedLen {
				t.Errorf(MSG_ERR_SIZE,
					tt.expectedLen,
					buf.Len(),
				)
			}

			outBytes := buf.Bytes()
			// Check title padding execution
			title := []byte("Pro Validation")
			expTitle := append(title, make([]byte, 6)...)
			if !bytes.Equal(outBytes[0:20], expTitle) {
				t.Errorf(MSG_ERR_TITLE)
			}

			// Validate Magic Marker Offset placement
			// 20 (Title) + 930 (Samples) + 1 (Len)
			// + 1 (Restart) + 128 (Seq Array)
			// = Offset 1080
			magicMarker := string(outBytes[1080:1084])
			if magicMarker != MAGIC_BYTES {
				t.Errorf(MSG_ERR_BAD_FMT, magicMarker)
			}

			// Sample header integrity loop check
			// Ensure empty sample arrays maintain default
			// loop size lengths (Byte 29 must be 0x01)
			for i := 0; i < 31; i++ {
				offset := 20 + (i * 30) + 29
				if outBytes[offset] != 0x01 {
					t.Errorf(MSG_ERR_BAD_INS, i+1)
				}
			}
		})
	}
}
