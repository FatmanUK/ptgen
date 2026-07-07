package main

import (
	"encoding/binary"
	"encoding/json"
	"path/filepath"
	"gopkg.in/yaml.v3"
	"os"
	"io"
	"fmt"
	"log"
	"bufio"
	"bytes"
	"regexp"
	"errors"
	"strings"
	pt "ptgen/internal/protracker"
	doh "github.com/FatmanUK/fatgo/docopt_helpers"
)

const metadataString = `Mod metadata:
Song Title:      {{ .Title }}
Format:          4-channel MOD / ProTracker-compatible
Speed:           {{ .Speed }}
BPM:             {{ .BPM }}
Channels:        4
Pattern length:  64 rows (Len. 0x40h)
Unique patterns: {{ .PatternsLen }}
Order length:    {{ .SequenceLen }}`

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

func readJSON(proj *pt.ModProject,
		file io.Reader) (*pt.ModProject, error) {
	err := json.NewDecoder(file).Decode(proj)
	return proj, err
}

func readYAML(proj *pt.ModProject,
		file io.Reader) (*pt.ModProject, error) {
	err := yaml.NewDecoder(file).Decode(proj)
	return proj, err
}

func emplaceRow(pattern *pt.Pattern, logs chan string,
		matches []string, isHexRows bool) error {
	rowNum, err := rowNumFromRowStr(matches[1], isHexRows)
	if err != nil {
		return err
	}
	*pattern, err = pattern.EmplaceRow(rowNum, matches)
	if err != nil {
		return err
	}
	logs <- "Row " + pattern[rowNum].StringFromRow(rowNum)
	return nil
}

func readTXT(file io.Reader, logs chan string,
		isHexRows bool) (interface{}, error) {
	pattern := pt.PatternFactory()
	re := regexp.MustCompile(pt.RowRegexFactory())
	scanner := bufio.NewScanner(file)
	rowsRead := 0
	for scanner.Scan() {
		matches := re.FindStringSubmatch(scanner.Text())
		if len(matches) == 0 {
			continue
		}
		rowsRead++
		err := emplaceRow(&pattern, logs, matches, isHexRows)
		if err != nil {
			return pattern, err
		}
	}
	logs <- fmt.Sprintf("Total rows processed: %d", rowsRead)
	return pattern, scanner.Err()
}

// The song metadata in JSON or YAML format.
func populateMetadata(proj *pt.ModProject, logs chan string,
		fileName string) error {
	var err error
	if fileName != "<nil>" {
		logs <- "Loading metadata."
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		if strings.HasSuffix(fileName, ".json") {
			proj, err = readJSON(proj, file)
		} else if strings.HasSuffix(fileName, ".yaml") {
			proj, err = readYAML(proj, file)
		}
	}
	return err
}

func countPatterns(path string) (uint, error) {
	var count uint = 0
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

// check all required pattern files exist
func isSequentialPatterns(path string, count uint) bool {
	var err error
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for c := 0; uint(c) < count; c++ {
		expectedName := fmt.Sprintf("pattern%02x.txt", c)
		isFound := false
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if entry.Name() == expectedName {
				isFound = true
				break
			}
		}
		if !isFound {
			return false
		}
	}
	return true
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
			matchesRow := reRow.FindStringSubmatch(line)
			if len(matchesRow) == 0 {
				continue
			}
			matchesHex := reHex.FindStringSubmatch(
					matchesRow[1])
			if len(matchesHex) > 0 {
				return true, nil
			}
		}
	}
	return false, nil
}

func loadPattern(logs chan string, fileBase string,
		isHexRows bool) (pt.Pattern, error) {
	var err error
	var pattern pt.Pattern
	fileName := fmt.Sprintf("%s.txt", fileBase)
	file, err := os.Open(fileName)
	if err != nil {
		return pt.PatternFactory(), err
	}
	defer file.Close()
	var p interface{}
	p, err = readTXT(file, logs, isHexRows)
	pattern = p.(pt.Pattern)
	return pattern, err
}

// No monolithic pattern files. Too annoying.
// Requires the patterns formatted in plain text. Screw JSON.
func populatePatterns(proj *pt.ModProject, logs chan string,
		path string) error {
	var err error
	logs <- "Loading patterns."
	numPatterns, err := countPatterns(path)
	if err != nil {
		return err
	}
	if pt.IsTooManyPatterns(numPatterns) {
		return errors.New("too many patterns")
	}
	if !isSequentialPatterns(path, numPatterns) {
		return errors.New("patterns not sequential")
	}
	isHexRows, err := isHexRowNotationDetected(path)
	if err != nil {
		return err
	}
	patterns := make([]pt.Pattern, numPatterns)
	for pIdx := range patterns {
		fileBase := fmt.Sprintf("pattern%02x", pIdx)
		logs <- fmt.Sprintf("Loading %s.", fileBase)
		filePathBase := filepath.Join(path, fileBase)
		patterns[pIdx], err = loadPattern(logs, filePathBase,
				isHexRows)
		if err != nil {
			return err
		}
	}
	proj.Patterns = patterns
	return err
}

func defaultOrderlist(proj *pt.ModProject) []uint8 {
	var orderList []uint8
	p := len(proj.Patterns)
	for n := 0; n < p; n++ {
		orderList = append(orderList, uint8(n))
	}
	return orderList
}

// A pattern order list formatted in JSON.
func populateOrderlist(proj *pt.ModProject, logs chan string,
		fileName string) error {
	var err error
	if fileName != "<nil>" {
		logs <- "Loading order list."
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		proj, err = readJSON(proj, file)
		if err != nil {
			return err
		}
	} else {
		proj.Sequence = defaultOrderlist(proj)
	}
	if !proj.IsOrderListValid() {
		err = errors.New("Invalid order list")
	}
	return err
}

func outputEverything(proj *pt.ModProject, logs chan string) error {
	var buf bytes.Buffer
	var err error
	info := proj.ModInfoFactory()
	output := mustPrepTemplate("output", metadataString, info)
	logs <- string(output)
	// TODO: add tests from testsuite ... here? or after WriteMod?
	// not writing speed/bpm? not writing the defaults either?
	// add a test for this
	err = pt.WriteMod(&buf, proj)
	if err != nil {
		return err
	}
	// TODO: add tests from testsuite ... here? or after WriteMod?
	err = binary.Write(os.Stdout, binary.BigEndian, buf.Bytes())
	return err
}

func threadGenerate(logs chan string, metaFile interface{},
		orderFile interface{}, patternsPath interface{}) {
	defer close(logs)
	var err error
	proj := pt.ModProjectFactory()
	meta := stringFromInterface(metaFile)
	err = populateMetadata(&proj, logs, meta)
	panicIfNotNil(err)
	pttns := stringFromInterface(patternsPath)
	err = populatePatterns(&proj, logs, pttns)
	panicIfNotNil(err)
	orderList := stringFromInterface(orderFile)
	err = populateOrderlist(&proj, logs, orderList)
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
	go threadGenerate(logs,	args["-m"], args["-o"], args["-p"])
	for msg := range logs {
		msgs := strings.Split(msg, "\\n")
		for _, m := range removeEmptyStrings(msgs) {
			log.Println(m)
		}
	}
}
