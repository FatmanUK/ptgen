package main

import (
	"fmt"
	doh "github.com/FatmanUK/fatgo/docopt_helpers"
	"github.com/FatmanUK/fatgo/utils"
	"log"
	"ptgen/internal/protracker"
	"strings"
)

const docoptFmt = `{{ .Name }} {{ .Version }}

Outputs the binary data directly to stdout, so redirect it to a file.

Usage:
  {{ .Name }} -m <metadata> -p <patterns>
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
	mod := protracker.ModFactory(logs)
	mod, err := mod.Init(
		utils.StringFromAny(metaFile),
		utils.StringFromAny(patternPath),
	)
	panicIfNotNil(err)
	panicIfNotNil(mod.Output())
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
