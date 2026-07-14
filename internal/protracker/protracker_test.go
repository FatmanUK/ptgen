package protracker

import (
	//	"bytes"
	"testing"
)

// Mostly written by Google Gemini Pro. Tweaked extensively by me.

/*
type ModProjectTestError struct {
	name        string
	project     ModProject
	shouldFail  bool
	expectedLen int
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

func TestInjectAndCheck(t *testing.T) {
	//func injectAndCheck(proj *ModProject) error {
}

func TestIsSlotPopulated(t *testing.T) {
	//func isSlotPopulated(slot *Instrument) bool {
}

func TestPreProcessInstruments(t *testing.T) {
	//func preProcessInstruments(proj *ModProject) ([31]Instrument, error) {
}

func TestWriteHeaders(t *testing.T) {
	//func writeHeaders(w io.Writer,
}

func TestExtractAndInjectSamples(t *testing.T) {
	//func extractAndInjectSamples(proj *ModProject) ([]byte, error) {
}
*/

func TestPrepareCommands_Must_Succeed(t *testing.T) {
	//14:func prepareCommands(s uint8, b uint8) ([]string, error) {
}

func TestPrepareCommands_Must_Fail(t *testing.T) {
	//14:func prepareCommands(s uint8, b uint8) ([]string, error) {
}

func TestWriteTitle_Must_Succeed(t *testing.T) {
	//34:func writeTitle(w io.Writer, title string) error {
}

func TestWriteTitle_Must_Fail(t *testing.T) {
	//34:func writeTitle(w io.Writer, title string) error {
}

func TestWriteSampleHeaders_Must_Succeed(t *testing.T) {
	//45:func writeSampleHeaders(w io.Writer, slots [31]Instrument) error {
}

func TestWriteSampleHeaders_Must_Fail(t *testing.T) {
	//45:func writeSampleHeaders(w io.Writer, slots [31]Instrument) error {
}

func TestWriteOrderList_Must_Succeed(t *testing.T) {
	//60:func writeOrderList(w io.Writer, orderList []uint8) error {
}

func TestWriteOrderList_Must_Fail(t *testing.T) {
	//60:func writeOrderList(w io.Writer, orderList []uint8) error {
}

func TestWriteMagic_Must_Succeed(t *testing.T) {
	//76:func writeMagic(w io.Writer) error {
}

func TestWriteMagic_Must_Fail(t *testing.T) {
	//76:func writeMagic(w io.Writer) error {
}

func TestWritePatterns_Must_Succeed(t *testing.T) {
	//86:func writePatterns(w io.Writer, pttns []Pattern) error {
}

func TestWritePatterns_Must_Fail(t *testing.T) {
	//86:func writePatterns(w io.Writer, pttns []Pattern) error {
}

func TestWriteSamples_Must_Succeed(t *testing.T) {
	//97:func writeSamples(w io.Writer, slots [31]Instrument) error {
}

func TestWriteSamples_Must_Fail(t *testing.T) {
	//97:func writeSamples(w io.Writer, slots [31]Instrument) error {
}

/*
func TestIsSequentialPatterns(t *testing.T) {
	adjustedPath := filepath.Join("..", SAMPLE2_PTTN_DIR)
	c, err := countPatterns(adjustedPath)
	if err != nil {
		t.Error(err)
	}
	err = isSequentialPatterns(adjustedPath, c)
	if err == nil {
		t.Log(MSG_PTTN_ARE_SEQ)
	} else {
		t.Error(err)
	}
}

func TestValidatePatterns(t *testing.T) {
	adjustedPath := filepath.Join("..", SAMPLE2_PTTN_DIR)
	_, err := validatePatterns(adjustedPath)
	if err == nil {
		t.Log(MSG_PTTN_VALID)
	} else {
		t.Error(err)
	}
}
*/
