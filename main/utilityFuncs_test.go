package main

import (
	"testing"
)

type NumTests struct {
	number string
	isHex bool
	expected uint16
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
		t.Logf(MSG_SLICES_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_SLICES, FMT_TST_VAR)
		t.Errorf(m, output, solution)
	}
}

func TestRowNumFromRowStr(t *testing.T) {
	tests := []NumTests{
		{ "1e", true, 30 },
	}
	for _, test := range tests {
		out, _ := rowNumFromRowStr(test.number, test.isHex)
		if out == test.expected {
			t.Logf(MSG_ROW_NUM_OK)
		} else {
			m := fmt.Sprintf("%s %s", ERR_ROW_NUM,
				FMT_TST_NUMBER)
			t.Errorf(m, out, test.expected)
		}
	}
}

// Need to think about this some more.
func TestTemplateFactory(t *testing.T) {
	// func TemplateFactory(name string, format string) *template.Template
}

// Need to think about this some more.
func TestMustPrepTemplate(t *testing.T) {
	// func mustPrepTemplate(name string, formatString string, data any) []byte
}

func TestCompareSlices(t *testing.T) {
	// func compareSlices[T comparable](a []T, b []T) bool
}
