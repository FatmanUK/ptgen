package main

const docoptString = `{{ .Name }} {{ .Version }}

Usage:
  {{ .Name }} [-m <metadata>] [-o <orderlist>] -p <patterns>
  {{ .Name }} -h | --help
  {{ .Name }} -v | --version

Options:
  -h --help       Show this screen
  -v --version    Show version
  -m <metadata>   Metadata input file
  -o <orderlist>  Order list input file
  -p <patterns>   Patterns file or directory
`

var APP_NAME string = "badvalue"
var VERSION  string = "badvalue"

type DocOptVars struct {
	Name string
	Version string
}

func DocOptVarsFactory() DocOptVars {
	return DocOptVars{
		Name: APP_NAME,
		Version: VERSION,
	}
}
