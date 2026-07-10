package main

import (
	"testing"
)

type NumTests struct {
	number string
	isHex bool
	expected uint16
}

func TestCompareSlices(t *testing.T) {
	//func compareSlices(a []string, b []string) bool
}

func TestRemoveEmptyStrings(t *testing.T) {
	testStrings := []string{
		"",
		"test",
		"1234",
		"",
		"",
		"5678",
		"",
		"blah",
		"",
		"",
	}
	solution := []string{
		"test",
		"1234",
		"5678",
		"blah",
	}
	output := removeEmptyStrings(testStrings)
	if compareSlices(output, solution) {
		t.Logf("Output matches expected solution.")
	} else {
		t.Errorf("Output doesn't match expected solution.")
	}
}

func TestRowNumFromRowStr(t *testing.T) {
	tests := []NumTests{
		{ "1e", true, 30 },
	}
	for _, test := range tests {
		out, _ := rowNumFromRowStr(test.number, test.isHex)
		if out == test.expected {
			t.Logf("Output matches expected solution.")
		} else {
			t.Errorf("Output doesn't match expected solution. %d != %d", out, test.expected)
		}
	}
}

// Need to think about this some more.
func TestTemplateFactory(t *testing.T) {
	// func TemplateFactory(name string, format string) *template.Template
}

func TestMustPrepTemplate(t *testing.T) {
	// func mustPrepTemplate(name string, formatString string, data interface{}) []byte
}
