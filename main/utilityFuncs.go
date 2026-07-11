package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		// doesn't exist
		return false, nil
	}
	// permissions issue
	return false, err
}

func checkFileExistsWithMkdir(path string) (bool, error) {
	cacheDir := filepath.Dir(path)
	err := os.MkdirAll(cacheDir, 0750)
	if err != nil {
		return false, err
	}
	exists, err := fileExists(path)
	if err != nil {
		return false, err
	}
	return exists, nil
}
