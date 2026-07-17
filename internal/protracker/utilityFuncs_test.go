package protracker

import (
	"fmt"
	"path/filepath"
	"testing"
)

const SAMPLE2_PTTN_COUNT = 4
const SAMPLE2_PTTN_DIR = `test_data/sample2`

const MSG_HEX_DETECTED = `Hex detected.`
const ERR_HEX_NOT_DETECTED = `Hex not detected.`

const MSG_MOD_PTTN_COUNT_OK = `Pattern count is ok.`
const ERR_MOD_PTTN_COUNT = `Pattern count is wrong.`

func TestIsHexDetected_Must_Succeed(t *testing.T) {
	adjustedPath := filepath.Join("..", "..", SAMPLE2_PTTN_DIR)
	isHex, err := IsHexDetected(adjustedPath)
	if isHex && err == nil {
		t.Log(MSG_HEX_DETECTED)
	} else {
		m := fmt.Sprintf("%s %s", ERR_HEX_NOT_DETECTED,
			FMT_TST_VAR)
		t.Errorf(m, err)
	}
}

func TestIsHexDetected_Must_Fail(t *testing.T) {
//20:func IsHexDetected(path string) (bool, error) {
}

func TestCheckHex_Must_Succeed(t *testing.T) {
//44:func CheckHex(filePath string, reRow *regexp.Regexp,
}

func TestCheckHex_Must_Fail(t *testing.T) {
//44:func CheckHex(filePath string, reRow *regexp.Regexp,
}

func TestPrepareRegexes_Must_Succeed(t *testing.T) {
//61:func PrepareRegexes() []*regexp.Regexp {
}

func TestPrepareRegexes_Must_Fail(t *testing.T) {
//61:func PrepareRegexes() []*regexp.Regexp {
}
