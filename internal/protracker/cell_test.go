package protracker

import (
	"fmt"
	"strings"
	"testing"
	"github.com/FatmanUK/fatgo/utils"
)

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

const TEST_CELL_STRING = `A#3 03 C38`

//const FMT_TST_STR = `Exp: "%s" Got: "%s"`
//const FMT_TST_NUM = `Exp:  %d  Got:  %d `
//const FMT_TST_VAR = `Exp: "%v" Got: "%v"`

const TST_OK_CELL_NOTE = `Note is ok.`
const TST_NO_CELL_NOTE = ERR_CELL_NOTE

const TST_OK_CELL_INST = `Instrument is ok.`
const TST_NO_CELL_INST = ERR_CELL_INST

const TST_OK_CELL_EFFT = `Effect is ok.`
const TST_NO_CELL_EFFT = ERR_CELL_EFFT

const TST_OK_CELL_REGEX = `Regex is ok.`
const TST_NO_CELL_REGEX = `Regex is wrong.`

const TST_OK_CELL_SAVE = `Save output is ok.`
const TST_NO_CELL_SAVE = `Save output is wrong.`

const TST_NO_CELL = `Generic Cell failure.`
const TST_NO_SUCCESS = `Should have failed but didn't.`
const TST_NO_WRONG_ERROR = `Wrong error message.`

func breakUpTestString(cellData *Cell) error {
	cellData.Note = TEST_CELL_STRING[0:3]
	inst := TEST_CELL_STRING[4:6]
	if inst == "--" {
		inst = "00"
	}
	cellData.Effect = TEST_CELL_STRING[7:10]
	i, err := utils.Uint8FromHexString(inst)
	cellData.Instr = i
	return err
}

func TestCellFactory_Must_Succeed(t *testing.T) {
	c := CellFactory()
	if c.Note == "---" {
		t.Logf(TST_OK_CELL_NOTE)
	} else {
		msgPair := [2]string{TST_NO_CELL_NOTE, FMT_TST_STR}
		valPair := [2]string{"---", c.Note}
		utils.TErr(t, msgPair, valPair)
	}
	if c.Instr == 0 {
		t.Logf(TST_OK_CELL_INST)
	} else {
		msgPair := [2]string{TST_NO_CELL_INST, FMT_TST_NUM}
		valPair := [2]uint8{0, c.Instr}
		utils.TErr(t, msgPair, valPair)
	}
	if c.Effect == "---" {
		t.Logf(TST_OK_CELL_EFFT)
	} else {
		msgPair := [2]string{TST_NO_CELL_EFFT, FMT_TST_STR}
		valPair := [2]string{"---", c.Effect}
		utils.TErr(t, msgPair, valPair)
	}
}

func TestCellFactory_Must_Fail(t *testing.T) {
//47:func CellFactory() Cell {
}

func TestCellRegexFactory_Must_Succeed(t *testing.T) {
	expectedRgx := fmt.Sprintf(RGX_CELL,
		RGX_NOTE, RGX_INST, RGX_EFFT,
		RGX_NOTE, RGX_INST)
	r := CellRegexFactory()
	if r == expectedRgx {
		t.Logf(TST_OK_CELL_REGEX)
	} else {
		msgPair := [2]string{TST_NO_CELL_REGEX, FMT_TST_STR}
		valPair := [2]string{expectedRgx, r}
		utils.TErr(t, msgPair, valPair)
	}
}

func TestCellRegexFactory_Must_Fail(t *testing.T) {
//55:func CellRegexFactory() string {
}

func TestCellLoad_Must_Succeed(t *testing.T) {
	c := CellFactory()
	cellData := CellFactory()
	err := breakUpTestString(&cellData)
	if err != nil {
		t.Errorf("%v", err)
	}
	err = c.Load(TEST_CELL_STRING)
	if err != nil {
		t.Errorf("%v", err)
	}
	if c.Note == cellData.Note {
		t.Logf(TST_OK_CELL_NOTE)
	} else {
		msgPair := [2]string{TST_NO_CELL_NOTE, FMT_TST_STR}
		valPair := [2]string{cellData.Note, c.Note}
		utils.TErr(t, msgPair, valPair)
	}
	if c.Instr == cellData.Instr {
		t.Logf(TST_OK_CELL_INST)
	} else {
		msgPair := [2]string{TST_NO_CELL_INST, FMT_TST_NUM}
		valPair := [2]uint8{cellData.Instr, c.Instr}
		utils.TErr(t, msgPair, valPair)
	}
	if c.Effect == cellData.Effect {
		t.Logf(TST_OK_CELL_EFFT)
	} else {
		msgPair := [2]string{TST_NO_CELL_EFFT, FMT_TST_STR}
		valPair := [2]string{cellData.Effect, c.Effect}
		utils.TErr(t, msgPair, valPair)
	}
}

func TestCellLoad_Must_Fail(t *testing.T) {
//61:func (c *Cell) Load(s string) error {
}

func TestCellSave_Must_Succeed(t *testing.T) {
	c := CellFactory()
	err := breakUpTestString(&c)
	if err != nil {
		t.Errorf("%v", err)
	}
	output := c.Save()
	if output == TEST_CELL_STRING {
		t.Logf(TST_OK_CELL_SAVE)
	} else {
		msgPair := [2]string{TST_NO_CELL_SAVE, FMT_TST_STR}
		valPair := [2]string{TEST_CELL_STRING, output}
		utils.TErr(t, msgPair, valPair)
	}
}

func TestCellSave_Must_Fail(t *testing.T) {
//86:func (c *Cell) Save() string {
}

func TestCellPack_Must_Succeed(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := tt.input.Pack()
			if err != nil {
				msgPair := [2]string{
					TST_NO_CELL, FMT_TST_STR}
				valPair := [2][4]byte{
					tt.expected, result}
				utils.TErr(t, msgPair, valPair)
			}
			if result != tt.expected {
				msgPair := [2]string{
					TST_NO_CELL, FMT_TST_STR}
				valPair := [2][4]byte{
					tt.expected, result}
				utils.TErr(t, msgPair, valPair)
			}
		})
	}
}

// TestEncodeCell_Must_Fail ensures validation blocks spec-breaking input
func TestCellPack_Must_Fail(t *testing.T) {
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
		expectedErr: "Note is wrong",
	})
	tests = append(tests, CellTestError{
		name:        "Octave above limit (C-6)",
		input:       cellTooHigh,
		expectedErr: "Note is wrong",
	})
	tests = append(tests, CellTestError{
		name:        "Invalid note format name",
		input:       cellRedDwarfError,
		expectedErr: "Note is wrong",
	})
	tests = append(tests, CellTestError{
		name:        "Instrument range overflow (>31)",
		input:       cellInstrTooHigh,
		expectedErr: "Instrument is wrong",
	})
	tests = append(tests, CellTestError{
		name:        "Effect string length underflow",
		input:       cellEffectTooShort,
		expectedErr: "Effect is wrong",
	})
	tests = append(tests, CellTestError{
		name:        "Effect string length overflow",
		input:       cellEffectTooLong,
		expectedErr: "Effect is wrong",
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := tt.input.Pack()
			if err == nil {
				t.Fatalf(TST_NO_SUCCESS)
			}
			hasSubstr := strings.Contains(err.Error(),
				tt.expectedErr)
			if !hasSubstr {
				msgPair := [2]string{
					TST_NO_WRONG_ERROR,
					FMT_TST_STR}
				valPair := [2]string{
					tt.expectedErr, err.Error()}
				utils.TErr(t, msgPair, valPair)
			}
		})
	}
}

func TestCellGetPeriod_Must_Succeed(t *testing.T) {
//95:func (c *Cell) getPeriod() uint16 {
}

func TestCellGetPeriod_Must_Fail(t *testing.T) {
//95:func (c *Cell) getPeriod() uint16 {
}

func TestCellGetEffectCommand_Must_Succeed(t *testing.T) {
//105:func (c *Cell) getEffectCommand() (uint8, error) {
}

func TestCellGetEffectCommand_Must_Fail(t *testing.T) {
//105:func (c *Cell) getEffectCommand() (uint8, error) {
}

func TestCellGetEffectParameter_Must_Succeed(t *testing.T) {
//123:func (c *Cell) getEffectParameter() (uint8, error) {
}

func TestCellGetEffectParameter_Must_Fail(t *testing.T) {
//123:func (c *Cell) getEffectParameter() (uint8, error) {
}

func TestCellWrite_Must_Succeed(t *testing.T) {
//180:func (c *Cell) Write(w io.Writer) error {
}

func TestCellWrite_Must_Fail(t *testing.T) {
//180:func (c *Cell) Write(w io.Writer) error {
}
