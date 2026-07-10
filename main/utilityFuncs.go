package main

import (
	"fmt"
	"strconv"
	"bytes"
	"text/template"
)

func compareSlices[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !(a[i] == b[i]) {
			return false
		}
	}
	return true
}

// loop and copy nonempty - O(n)
// join and split - O(???)
func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func stringFromAny(any any) string {
	return fmt.Sprintf("%v", any)
}

func rowNumFromRowStr(rowStr string, isHex bool) (uint16, error) {
	rowBase := 10
	if isHex {
		rowBase = 16
	}
	u, err := strconv.ParseUint(rowStr, rowBase, 16)
	if err != nil {
		return 0, err
	}
	return uint16(u), nil
}

func TemplateFactory(name string, format string) *template.Template {
	t := template.New(name)
	t = template.Must(t.Parse(format))
	return t
}

func mustPrepTemplate(name string, formatString string,
	data any) []byte {
	var wr bytes.Buffer
	err := TemplateFactory(name, formatString).Execute(&wr, data)
	if err != nil {
		panic(err)
	}
	return wr.Bytes()
}
