package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	doh "github.com/FatmanUK/fatgo/docopt_helpers"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path/filepath"
	pt "ptgen/internal/protracker"
	"regexp"
	"strings"
)

const metadataString = `Mod metadata:
Song Title:      {{ .Title }}
Format:          ProTracker-compatible MOD
Channels:        4
Speed:           {{ .Speed }}
BPM:             {{ .BPM }}
Pattern length:  64 rows (Len. 0x40h)
Unique patterns: {{ .PatternsLen }}
Order length:    {{ .SequenceLen }}`

const MOD_TOO_BIG_BYTES = (50 * 1024)

const MSG_WRN_TOO_BIG = "Warning: Mod is %d bytes. Typically, they're smaller than %d bytes. This is a very large mod."

type PatternScanData struct {
	pattern   pt.Pattern
	scanner   *bufio.Scanner
	regex     *regexp.Regexp
	isHexRows bool
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

func scanLoop(logs chan string, d *PatternScanData,
	rowsRead *uint) (pt.Pattern, error) {
	for d.scanner.Scan() {
		m := d.regex.FindStringSubmatch(d.scanner.Text())
		if len(m) != 0 {
			*rowsRead++
			n, err := rowNumFromRowStr(m[1], d.isHexRows)
			if err != nil {
				return d.pattern, err
			}
			d.pattern, err = d.pattern.EmplaceRow(n, m)
			if err != nil {
				return d.pattern, err
			}
			logs <- ">" + d.pattern[n].StringFromRow(n)
			if err != nil {
				return d.pattern, err
			}
		}
	}
	return d.pattern, nil
}

func readPattern(file io.Reader, logs chan string,
	isHexRows bool) (pt.Pattern, error) {
	var rowsRead uint = 0
	var err error
	data := PatternScanData{
		pattern:   pt.PatternFactory(),
		scanner:   bufio.NewScanner(file),
		regex:     regexp.MustCompile(pt.RowRegexFactory()),
		isHexRows: isHexRows,
	}
	data.pattern, err = scanLoop(logs, &data, &rowsRead)
	if err != nil {
		return data.pattern, err
	}
	logs <- fmt.Sprintf("Total rows processed: %d", rowsRead)
	return data.pattern, data.scanner.Err()
}

func readMetadata(proj *pt.ModProject, logs chan string,
	file *os.File) (*pt.ModProject, error) {
	var err error
	content, err := io.ReadAll(file)
	if err != nil {
		return proj, err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return proj, err
	}
	if json.Valid(content) {
		logs <- "JSON metadata detected"
		err = json.NewDecoder(file).Decode(proj)
	} else {
		var node yaml.Node
		err = yaml.Unmarshal(content, &node)
		if err == nil {
			logs <- "YAML metadata detected"
			err = yaml.NewDecoder(file).Decode(proj)
		}
	}
	return proj, err
}

func defaultOrderlist(proj *pt.ModProject) []uint8 {
	var orderList []uint8
	p := len(proj.Patterns)
	for n := 0; n < p; n++ {
		orderList = append(orderList, uint8(n))
	}
	return orderList
}

// The song metadata in JSON or YAML format, including order list.
func populateMetadata(proj *pt.ModProject, logs chan string,
	fileName string) error {
	var err error
	proj.Sequence = defaultOrderlist(proj)
	if fileName != "<nil>" {
		logs <- "Loading metadata and pattern order."
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		proj, err = readMetadata(proj, logs, file)
		if err != nil {
			return err
		}
	}
	if !proj.IsOrderListValid() {
		err = errors.New("Invalid order list")
	}
	return err
}

func countPatterns(path string) (uint8, error) {
	var count uint8 = 0
	var err error
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		re := regexp.MustCompile(`pattern[0-F][0-F].txt`)
		if len(re.FindStringSubmatch(entry.Name())) == 0 {
			continue
		}
		count++
	}
	return count, nil
}

func locateFile(expected string, entries []os.DirEntry) bool {
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == expected {
			return true
		}
	}
	return false
}

// check all required pattern files exist
func isSequentialPatterns(path string, count uint8) error {
	var err error
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for c := 0; uint8(c) < count; c++ {
		expected := fmt.Sprintf("pattern%02x.txt", c)
		if !locateFile(expected, entries) {
			return errors.New("patterns not sequential")
		}
	}
	return nil
}

// Row notation must be consistent across pattern files, but I don't
// see a way to enforce it.
func isHexRowNotationDetected(path string) (bool, error) {
	var err error
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	reRow := regexp.MustCompile(pt.RowRegexFactory())
	reHex := regexp.MustCompile(`[A-Fa-f]`)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		rePttn := regexp.MustCompile(`pattern[0-F][0-F].txt`)
		if len(rePttn.FindStringSubmatch(entry.Name())) == 0 {
			continue
		}
		fileName := filepath.Join(path, entry.Name())
		file, err := os.Open(fileName)
		if err != nil {
			return false, err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			mRow := reRow.FindStringSubmatch(line)
			if len(mRow) == 0 {
				continue
			}
			mHex := reHex.FindStringSubmatch(mRow[1])
			if len(mHex) > 0 {
				return true, nil
			}
		}
	}
	return false, nil
}

func loadPattern(logs chan string, fileName string,
	isHexRows bool) (pt.Pattern, error) {
	var err error
	var pattern pt.Pattern
	file, err := os.Open(fileName)
	if err != nil {
		return pt.PatternFactory(), err
	}
	defer file.Close()
	var p interface{}
	p, err = readPattern(file, logs, isHexRows)
	pattern = p.(pt.Pattern)
	return pattern, err
}

func validatePatterns(path string) (uint8, error) {
	numPttns, err := countPatterns(path)
	if err != nil {
		return numPttns, err
	}
	err = pt.IsTooManyPatterns(numPttns)
	if err != nil {
		return numPttns, err
	}
	err = isSequentialPatterns(path, numPttns)
	if err != nil {
		return numPttns, err
	}
	return numPttns, nil
}

// No monolithic pattern files. Too annoying.
// Requires the patterns formatted in plain text. Screw JSON.
func populatePatterns(proj *pt.ModProject, logs chan string,
	path string) error {
	logs <- "Loading patterns."
	numPatterns, err := validatePatterns(path)
	isHex, err := isHexRowNotationDetected(path)
	if err != nil {
		return err
	}
	patterns := make([]pt.Pattern, numPatterns)
	for i := range patterns {
		fileName := fmt.Sprintf("pattern%02x.txt", i)
		fullPath := filepath.Join(path, fileName)
		logs <- fmt.Sprintf("Loading %s.", fileName)
		patterns[i], err = loadPattern(logs, fullPath, isHex)
		if err != nil {
			return err
		}
	}
	proj.Patterns = patterns
	return err
}

func outputEverything(proj *pt.ModProject, logs chan string) error {
	var buf bytes.Buffer
	var err error
	info := proj.ModInfoFactory()
	output := mustPrepTemplate("output", metadataString, info)
	logs <- string(output)
	if len(proj.Title) > 20 { // title less than 21 bytes
		return errors.New("title too long")
	}
	err = pt.WriteMod(&buf, proj)
	if err != nil {
		return err
	}
	if buf.Len() > MOD_TOO_BIG_BYTES {
		logs <- fmt.Sprintf(MSG_WRN_TOO_BIG,
			buf.Len(), MOD_TOO_BIG_BYTES)
	}
	err = binary.Write(os.Stdout, binary.BigEndian, buf.Bytes())
	return err
}

func threadGenerate(logs chan string, metaFile interface{},
	patternsPath interface{}) {
	defer close(logs)
	var err error
	proj := pt.ModProjectFactory()
	pttns := stringFromInterface(patternsPath)
	err = populatePatterns(&proj, logs, pttns)
	panicIfNotNil(err)
	meta := stringFromInterface(metaFile)
	err = populateMetadata(&proj, logs, meta)
	panicIfNotNil(err)
	err = outputEverything(&proj, logs)
	panicIfNotNil(err)
}

func main() {
	dotv := DocOptVarsFactory()
	ds := string(mustPrepTemplate("docopt", docoptString, dotv))
	appVer := dotv.Name + " " + dotv.Version
	args, err := doh.NoExitParser.ParseArgs(ds, nil, appVer)
	if args == nil { // All done, exit here.
		return
	}
	panicIfNotNil(err)
	logs := make(chan string)
	go threadGenerate(logs, args["-m"], args["-p"])
	for msg := range logs {
		msgs := strings.Split(msg, "\\n")
		for _, m := range removeEmptyStrings(msgs) {
			log.Println(m)
		}
	}
}
