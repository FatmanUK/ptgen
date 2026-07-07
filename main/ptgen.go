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

func readJSON(proj interface{}, file io.Reader) (interface{}, error) {
	err := json.NewDecoder(file).Decode(proj)
	return proj, err
}

func readYAML(proj interface{}, file io.Reader) (interface{}, error) {
	err := yaml.NewDecoder(file).Decode(proj)
	return proj, err
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
		rowNum, err := rowNumFromRowStr(matches[1], isHexRows)
		if err != nil {
			return pattern, err
		}
		pattern, err = pattern.EmplaceRow(rowNum, matches)
		if err != nil {
			return pattern, err
		}
		logs <- "Row " + pattern[rowNum].StringFromRow(rowNum)
	}
	logs <- fmt.Sprintf("Total rows processed: %d", rowsRead)
	return pattern, scanner.Err()
}

func isFileJSON(fileName string) (bool, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return false, err
	}
	defer file.Close()
	buf := make([]byte, 1)
	n, err := file.Read(buf)
	if err != nil || n != 1 {
		return false, errors.New("error reading file")
	}
	return buf[0] == '{', nil
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
		var p interface{}
		isJSON, err := isFileJSON(fileName)
		if err != nil {
			return err
		}
		if isJSON {
			p, err = readJSON(proj, file)
		} else {
			p, err = readYAML(proj, file)
		}
		proj = p.(*pt.ModProject)
	}
	return err
}

func countPatterns(path string) int {
	// count patternNN (txt or json) files in dir
	// TODO
	return 4
}

func isSequentialPatterns(path string) bool {
	// check all pattern files run in a sequence
	// pattern0a not pattern10
	// TODO
	return true
}

func isHexRowNotationDetected(path string) bool {
	// is row number hex or dec? can detect hex, but not dec.
	// must be consistent in this (why wouldn't you be?)
	// detect across all patterns. unlikely we wouldn't find at
	// least one hex row
	// can use empty row eg. 5e to denote hex
	// TODO
	return true
}

// TODO: could split this fn
func loadPattern(logs chan string, fileBase string,
		isHexRows bool) (pt.Pattern, error) {
	var err error
	var pattern pt.Pattern
	var isJSON bool = true
	fileNameJSON := fmt.Sprintf("%s.json", fileBase)
	file, err := os.Open(fileNameJSON)
	if err != nil {
		isJSON = false
		fileNameTXT := fmt.Sprintf("%s.txt", fileBase)
		file, err = os.Open(fileNameTXT)
		if err != nil {
			return pt.PatternFactory(), err
		}
	}
	defer file.Close()
	var p interface{}
	if isJSON {
		isFileJSON, err := isFileJSON(fileNameJSON)
		if err != nil {
			return pattern, err
		}
		if isFileJSON {
			p, err = readJSON(pattern, file)
		}
	} else {
		p, err = readTXT(file, logs, isHexRows)
	}
	pattern = p.(pt.Pattern)
	return pattern, err
}

// No monolithic pattern files. Too annoying.
// Requires the patterns formatted in plain text or JSON.
func populatePatterns(proj *pt.ModProject, logs chan string,
		path string) error {
	var err error
	logs <- "Loading patterns."
	numPatterns := countPatterns(path)
	if pt.IsTooManyPatterns(numPatterns) {
		return errors.New("too many patterns")
	}
	if !isSequentialPatterns(path) {
		return errors.New("patterns not sequential")
	}
	isHexRows := isHexRowNotationDetected(path)
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
	var p interface{}
	if fileName != "<nil>" {
		logs <- "Loading order list."
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		p, err = readJSON(proj, file)
		proj = p.(*pt.ModProject)
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
	if err != nil {
		panic(err)
	}
	pttns := stringFromInterface(patternsPath)
	err = populatePatterns(&proj, logs, pttns)
	if err != nil {
		panic(err)
	}
	orderList := stringFromInterface(orderFile)
	err = populateOrderlist(&proj, logs, orderList)
	if err != nil {
		panic(err)
	}
	err = outputEverything(&proj, logs)
	if err != nil {
		panic(err)
	}
}

func main() {
	dotv := DocOptVarsFactory()
	ds := string(mustPrepTemplate("docopt", docoptString, dotv))
	appVer := dotv.Name + " " + dotv.Version
	args, err := doh.NoExitParser.ParseArgs(ds, nil, appVer)
	if args == nil { // All done, exit here.
		return
	}
	if err != nil {
		panic(err)
	}
	logs := make(chan string)
	go threadGenerate(logs,	args["-m"], args["-o"], args["-p"])
	for msg := range logs {
		msgs := strings.Split(msg, "\\n")
		for _, m := range removeEmptyStrings(msgs) {
			log.Println(m)
		}
	}
}
