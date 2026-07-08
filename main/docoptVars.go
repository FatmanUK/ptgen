package main

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
