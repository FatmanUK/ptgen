package protracker

import (
	"fmt"
	"strconv"
	"testing"
)

func makeVars() (string, uint8, string, error) {
	note := TEST_CELL_STRING[0:3]
	effect := TEST_CELL_STRING[7:10]
	instStr := TEST_CELL_STRING[4:6]
	if instStr == "--" {
		instStr = "00"
	}
	inst, err := strconv.ParseUint(instStr, 16, 8)
	return note, uint8(inst), effect, err
}

func TestCellFactory(t *testing.T) {
	c := CellFactory()
	if c.Note == "---" {
		t.Logf(MSG_CELL_NOTE_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_CELL_NOTE,
			FMT_TST_STRING)
		t.Errorf(m, "---", c.Note)
	}
	if c.Instr == 0 {
		t.Logf(MSG_CELL_INST_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_CELL_INST,
			FMT_TST_NUMBER)
		t.Errorf(m, 0, c.Instr)
	}
	if c.Effect == "---" {
		t.Logf(MSG_CELL_EFFT_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_CELL_EFFT,
			FMT_TST_STRING)
		t.Errorf(m, "---", c.Effect)
	}
}

func TestCellRegexFactory(t *testing.T) {
	expectedRgx := fmt.Sprintf(RGX_CELL,
		RGX_NOTE, RGX_INSTR, RGX_EFFECT,
		RGX_NOTE, RGX_INSTR)
	r := CellRegexFactory()
	if r == expectedRgx {
		t.Logf(MSG_CELL_REGEX_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_CELL_REGEX,
			FMT_TST_STRING)
		t.Errorf(m, expectedRgx, r)
	}
}

func TestCellLoad(t *testing.T) {
	var err error
	c := CellFactory()
	note, inst, effect, err := makeVars()
	if err != nil {
		t.Errorf("%v", err)
	}
	err = c.Load(TEST_CELL_STRING)
	if err != nil {
		t.Errorf("%v", err)
	}
	if c.Note == note {
		t.Logf(MSG_CELL_NOTE_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_CELL_NOTE,
			FMT_TST_STRING)
		t.Errorf(m, note, c.Note)
	}
	if c.Instr == inst {
		t.Logf(MSG_CELL_INST_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_CELL_INST,
			FMT_TST_NUMBER)
		t.Errorf(m, inst, c.Instr)
	}
	if c.Effect == effect {
		t.Logf(MSG_CELL_EFFT_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_CELL_EFFT,
			FMT_TST_STRING)
		t.Errorf(m, effect, c.Effect)
	}
}

func TestCellSave(t *testing.T) {
	var err error
	c := CellFactory()
	c.Note, c.Instr, c.Effect, err = makeVars()
	if err != nil {
		t.Errorf("%v", err)
	}
	output := c.Save()
	if output == TEST_CELL_STRING {
		t.Logf(MSG_CELL_SAVE_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_CELL_SAVE,
			FMT_TST_STRING)
		t.Errorf(m, TEST_CELL_STRING, output)
	}
}
