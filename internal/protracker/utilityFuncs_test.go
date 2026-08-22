package protracker

import (
	"fmt"
	"testing"
)

const SAMPLE2_PTTN_COUNT = 4
const SAMPLE2_PTTN_DIR = `../../test_data/sample2`

const SAMPLE1_PTTN_DIR = `../../test_data/sample1`

const MSG_HEX_DETECTED = `Hex detected.`
const ERR_HEX_NOT_DETECTED = `Hex not detected.`

const MSG_HEX_NOT_DETECTED = `Hex not detected.`
const ERR_HEX_DETECTED = `Hex detected.`

const MSG_MOD_PTTN_COUNT_OK = `Pattern count is ok.`
const ERR_MOD_PTTN_COUNT = `Pattern count is wrong.`

func TestIsHexDetected_Must_Succeed(t *testing.T) {
	isHex, err := IsHexDetected(SAMPLE2_PTTN_DIR)
	if isHex && err == nil {
		t.Log(MSG_HEX_DETECTED)
	} else {
		m := fmt.Sprintf("%s %s", ERR_HEX_NOT_DETECTED,
			FMT_TST_VAR)
		t.Errorf(m, err)
	}
}

func TestIsHexDetected_Must_Fail(t *testing.T) {
	isHex, err := IsHexDetected(SAMPLE1_PTTN_DIR)
	if isHex && err == nil {
		t.Errorf(ERR_HEX_DETECTED)
	} else {
		t.Log(MSG_HEX_NOT_DETECTED)
	}
}
