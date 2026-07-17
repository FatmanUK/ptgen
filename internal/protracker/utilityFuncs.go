package protracker

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
)

const RGX_HEX_EVIDENCE = `[A-Fa-f]`

const ERR_ARCH_HEADER_PARSE = `failed header parse: %v`
const MSG_ARCH_HEADER_PARSE_OK = `Processing file: %s (Method: %s, Original Size: %d)`

const ERR_ARCH_EXTRACTION = `Extraction failed for file %s: %v`
const ERR_ARCH_SIZE_MISMATCH = `Size mismatch for %s: expected %d bytes, extracted %d`

// Row notation must be consistent across pattern files, but I don't
// see a way to enforce it.
func IsHexDetected(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	rgxs := PrepareRegexes()
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !rgxs[2].MatchString(name) {
			continue
		}
		filePath := filepath.Join(path, name)
		detected, err := CheckHex(filePath, rgxs[0], rgxs[1])
		if err != nil {
			return false, err
		}
		if detected {
			return true, nil
		}
	}
	return false, nil
}

// TODO: row number '5x' or '6x' is an instant return false, nil.
func CheckHex(filePath string, reRow *regexp.Regexp,
	reHex *regexp.Regexp) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mRow := reRow.FindStringSubmatch(scanner.Text())
		if len(mRow) > 1 && reHex.MatchString(mRow[1]) {
			return true, nil
		}
	}
	return false, scanner.Err()
}

func PrepareRegexes() []*regexp.Regexp {
	reRow := regexp.MustCompile(RowRegexFactory())
	reHex := regexp.MustCompile(RGX_HEX_EVIDENCE)
	rePttn := regexp.MustCompile(RGX_PTTN_FILE)
	return []*regexp.Regexp{reRow, reHex, rePttn}
}
