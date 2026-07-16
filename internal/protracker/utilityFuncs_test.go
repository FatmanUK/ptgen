package protracker

import (
	"fmt"
	"path/filepath"
	"testing"
	"github.com/FatmanUK/fatgo/utils"
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
	//33:func IsHexDetected(path string) (bool, error) {
}

func TestCheckHex_Must_Succeed(t *testing.T) {
	//57:func CheckHex(filePath string, reRow *regexp.Regexp,
}

func TestCheckHex_Must_Fail(t *testing.T) {
	//57:func CheckHex(filePath string, reRow *regexp.Regexp,
}

func TestCountPatterns_Must_Succeed(t *testing.T) {
	path := filepath.Join("..", "..", SAMPLE2_PTTN_DIR)
	expectedCount := uint8(SAMPLE2_PTTN_COUNT)
	output, err := utils.CountMatchingFiles(path, RGX_PTTN_FILE)
	if err != nil {
		t.Log(err)
	}
	if output == expectedCount {
		t.Log(MSG_MOD_PTTN_COUNT_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_MOD_PTTN_COUNT,
			FMT_TST_NUM)
		t.Errorf(m, expectedCount, output)
	}
}

func TestCountPatterns_Must_Fail(t *testing.T) {
	//74:func CountPatterns(path string) (uint8, error) {
}

func TestPrepareRegexes_Must_Succeed(t *testing.T) {
	//90:func PrepareRegexes() ([]*regexp.Regexp) {
}

func TestPrepareRegexes_Must_Fail(t *testing.T) {
	//90:func PrepareRegexes() ([]*regexp.Regexp) {
}

func TestUint8FromHexString_Must_Succeed(t *testing.T) {
	//97:func uint8FromHexString(s string) (uint8, error) {
}

func TestUint8FromHexString_Must_Fail(t *testing.T) {
	//97:func uint8FromHexString(s string) (uint8, error) {
}

func TestCompareSlices_Must_Succeed(t *testing.T) {
	//103:func compareSlices[T comparable](a []T, b []T) bool {
}

func TestCompareSlices_Must_Fail(t *testing.T) {
	//103:func compareSlices[T comparable](a []T, b []T) bool {
}

func TestCheckFileExistsWithMkdir_Must_Succeed(t *testing.T) {
	//162:func CheckFileExistsWithMkdir(path string) (bool, error) {
}

func TestCheckFileExistsWithMkdir_Must_Fail(t *testing.T) {
	//162:func CheckFileExistsWithMkdir(path string) (bool, error) {
}

func TestFileExists_Must_Succeed(t *testing.T) {
	//175:func fileExists(path string) (bool, error) {
}

func TestFileExists_Must_Fail(t *testing.T) {
	//175:func fileExists(path string) (bool, error) {
}

func TestPrepareArchives_Must_Succeed(t *testing.T) {
	//188:func prepareArchives() (map[string]string, error) {
}

func TestPrepareArchives_Must_Fail(t *testing.T) {
	//188:func prepareArchives() (map[string]string, error) {
}

func TestExtractSample_Must_Succeed(t *testing.T) {
	//199:func extractSample(n string, lr *xlha.Reader,
}

func TestExtractSample_Must_Fail(t *testing.T) {
	//199:func extractSample(n string, lr *xlha.Reader,
}
