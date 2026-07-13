package main

import (
	doh "github.com/FatmanUK/fatgo/docopt_helpers"
	"log"
	pt "ptgen/internal/protracker"
	"strings"
)

const docoptString = `{{ .Name }} {{ .Version }}

Outputs the binary data directly to stdout, so redirect it to a file.

Usage:
  {{ .Name }} [-m <metadata>] -p <patterns>
  {{ .Name }} -h | --help
  {{ .Name }} -v | --version

Options:
  -h --help       Show this screen
  -v --version    Show version
  -m <metadata>   Metadata input file
  -p <patterns>   Patterns directory
`

var APP_NAME string = "badvalue"
var VERSION string = "badvalue"

type DocOptVars struct {
	Name    string
	Version string
}

func DocOptVarsFactory() DocOptVars {
	return DocOptVars{
		Name:    APP_NAME,
		Version: VERSION,
	}
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

func threadGenerate(logs chan string, metaFile any, patternPath any) {
	defer close(logs)
	var err error
	pttns := stringFromAny(patternPath)
	meta := stringFromAny(metaFile)
	proj := pt.ModProjectFactory()
	err = proj.PopulatePatterns(logs, pttns)
	panicIfNotNil(err)
	err = proj.PopulateMetadata(logs, meta)
	panicIfNotNil(err)
	err = proj.OutputEverything(logs)
	panicIfNotNil(err)
}

func main() {
	dotv := DocOptVarsFactory()
	ds := string(MustPrepTemplate("docopt", docoptString, dotv))
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
