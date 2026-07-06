package main

import (
	"bytes"
	"text/template"
)

func TemplateFactory(name string, format string) *template.Template {
	t := template.New(name)
	t = template.Must(t.Parse(format))
	return t
}

func mustPrepTemplate(name string, formatString string, data interface{}) []byte {
	var wr bytes.Buffer
	err := TemplateFactory(name, formatString).Execute(&wr, data)
	if err != nil {
		panic(err)
	}
	return wr.Bytes()
}
