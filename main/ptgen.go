package main

import (
	doh "github.com/FatmanUK/fatgo/docopt_helpers"
	"fmt"
	"log"
	pt "ptgen/internal/protracker"
	"strings"
	"github.com/FatmanUK/fatgo/utils"
)

const docoptFmt = `{{ .Name }} {{ .Version }}

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

func (d *DocOptVars) ProcessArgs() (map[string]any, error) {
	ds := string(utils.MustPrepTemplate("docopt", docoptFmt, *d))
	av := fmt.Sprintf("%s %s", d.Name, d.Version)
	return doh.NoExitParser.ParseArgs(ds, nil, av)
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

func threadGenerate(logs chan string, metaFile any, patternPath any) {
	defer close(logs)
	proj := pt.ModProjectFactory()
	p := utils.StringFromAny(patternPath)
	err := proj.PopulatePatterns(logs, p)
	panicIfNotNil(err)
	m := utils.StringFromAny(metaFile)
	err = proj.PopulateMetadata(logs, m)
	panicIfNotNil(err)
	err = proj.OutputEverything(logs)
	panicIfNotNil(err)
}

func threadLog(logs chan string) {
	for msg := range logs {
		msgs := strings.Split(msg, "\\n")
		for _, m := range utils.RemoveEmptyStrings(msgs) {
			log.Println(m)
		}
	}
}

func main() {
	dotv := DocOptVarsFactory()
	args, err := dotv.ProcessArgs()
	if args == nil { // All done, exit here.
		return
	}
	panicIfNotNil(err)
	logs := make(chan string)
	go threadGenerate(logs, args["-m"], args["-p"])
	threadLog(logs)
}
