package main

import (
	"fmt"
	"path/filepath"
	pt "ptgen/internal/protracker"
	"testing"
)

func TestReadPattern(t *testing.T) {
	// func readPattern(file io.Reader, logs chan string, isHexRows bool) (pt.Pattern, error)
}

func TestReadMetadata(t *testing.T) {
	// func readMetadata(proj *pt.ModProject, logs chan string, file *os.File) (*pt.ModProject, error)
}

func TestDefaultOrderlist(t *testing.T) {
	proj := pt.ModProjectFactory()
	pttn := pt.PatternFactory()
	proj.Patterns = []pt.Pattern{pttn, pttn, pttn, pttn, pttn}
	proj.OrderList = defaultOrderlist(&proj)
	if proj.IsOrderListValid() {
		t.Logf(MSG_MOD_OLST_OK)
	} else {
		t.Errorf(ERR_MOD_OLST)
	}
}

func TestCountPatterns(t *testing.T) {
	adjustedPath := filepath.Join("..", SAMPLE2_PTTN_DIR)
	expectedCount := uint8(SAMPLE2_PTTN_COUNT)
	output, err := countPatterns(adjustedPath)
	if err != nil {
		t.Error(err)
	}
	if output == expectedCount {
		t.Log(MSG_MOD_PTTN_COUNT_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_MOD_PTTN_COUNT,
			FMT_TST_NUMBER)
		t.Errorf(m, expectedCount, output)
	}
}

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

func TestIsHexRowNotationDetected(t *testing.T) {
	adjustedPath := filepath.Join("..", SAMPLE2_PTTN_DIR)
	isHex, err := isHexRowNotationDetected(adjustedPath)
	if isHex && err == nil {
		t.Log(MSG_HEX_DETECTED)
	} else {
		m := fmt.Sprintf("%s %s", ERR_HEX_NOT_DETECTED,
			FMT_TST_VAR)
		t.Errorf(m, err)
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

/*
func scanLoop(logs chan string, d *ScanData) (pt.Pattern, error) {
func populateMetadata(proj *pt.ModProject, logs chan string,
func locateFile(expected string, entries []os.DirEntry) bool {
func prepareRegexes() (*regexp.Regexp, *regexp.Regexp,
func checkFileForHexRow(filePath string, reRow *regexp.Regexp,
func loadPattern(logs chan string, fileName string,
func populatePatterns(proj *pt.ModProject, logs chan string,
func decideSources(samples []pt.Instrument) (map[string]Download, error) {
func download(downloads map[string]Download, logs chan string) error {
func outputEverything(proj *pt.ModProject, logs chan string) error {
*/
