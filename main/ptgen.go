package main

import (
	"encoding/binary"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"io"
	"log"
//	"bufio"
	"bytes"
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

func readJSON(proj interface{}, file io.Reader) error {
	return json.NewDecoder(file).Decode(proj)
}

func readYAML(proj interface{}, file io.Reader) error {
	return yaml.NewDecoder(file).Decode(proj)
}

func readText(proj interface{}, file io.Reader) error {
	// TODO
/*
// look for lines with hex/dec number then colon.
// roll through, look for lines like "Pattern XX:"
// if not enough patterns, add until enough
*/
	return nil
}

func readDirectory(proj *pt.ModProject, path string) error {
	// TODO
	// look for files named pattern00.txt, pattern01.txt etc.
/*
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if isFileJSON(file) {
		err = readJSON(proj, file)
	} else {
		err = readText(proj, file)
*/
	return nil
}

func isFileJSON(file io.Reader) bool {
	// TODO
	// detect, not validate
	// first char open brace will do
	return true
}

func isDirectory(path string) bool {
	// TODO
	return false
}

// The song metadata in JSON or YAML format.
func populateMetadata(proj *pt.ModProject,
		logs chan string,
		fileName string) error {
	var err error
	if fileName != "<nil>" {
		logs <- "Loading metadata."
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		if isFileJSON(file) {
			err = readJSON(proj, file)
		} else {
			err = readYAML(proj, file)
		}
	}
	return err
}

// Requires the patterns formatted in plain text or JSON.
func populatePatterns(proj *pt.ModProject,
		logs chan string,
		path string) error {
	var err error
	logs <- "Loading patterns."
	if isDirectory(path) {
		err = readDirectory(proj, path)
	} else {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		if isFileJSON(file) {
			err = readJSON(proj, file)
		} else {
			err = readText(proj, file)
		}
	}
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
func populateOrderlist(proj *pt.ModProject,
		logs chan string,
		fileName string) error {
	var err error
	if fileName != "<nil>" {
		logs <- "Loading order list."
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		err = readJSON(proj, file)
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
