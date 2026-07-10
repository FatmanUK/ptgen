package main

import (
	"fmt"
	"testing"
)

type NumTests struct {
	number     string
	isHex      bool
	expected   uint16
	shouldFail bool
}

func TestRemoveEmptyStrings(t *testing.T) {
	testStrings := []string{
		"", "test", "1234", "", "",
		"5678", "", "blah", "", "",
	}
	solution := []string{"test", "1234", "5678", "blah"}
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
		{"1e", true, 30, false},
		{"20", true, 32, false},
		{"20", false, 20, false},

		{"1e", true, 31, true},
		{"20", true, 31, true},
		{"2x", false, 20, true},
	}
	for _, test := range tests {
		out, _ := rowNumFromRowStr(test.number, test.isHex)
		if test.shouldFail {
			if out != test.expected {
				t.Logf(MSG_ROW_NUM_NOK)
			} else {
				m := fmt.Sprintf("%s %s",
					NERR_ROW_NUM, FMT_TST_NUMBER)
				t.Errorf(m, out, test.expected)
			}
		} else {
			if out == test.expected {
				t.Logf(MSG_ROW_NUM_OK)
			} else {
				m := fmt.Sprintf("%s %s",
					ERR_ROW_NUM, FMT_TST_NUMBER)
				t.Errorf(m, out, test.expected)
			}
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
	as := []string{"abc", "123", "xyz"}
	ai := []int{1, -8, 43}

	bs := as
	bi := ai

	fs := []string{}
	fi := []int{1, 9, -4}

	passes := []bool{}
	passes = append(passes, compareSlices(as, bs))
	passes = append(passes, compareSlices(ai, bi))

	fails := []bool{}
	fails = append(fails, compareSlices(as, fs))
	fails = append(fails, compareSlices(ai, fi))

	for _, pass := range passes {
		if pass {
			t.Logf(MSG_SLICES_OK)
		} else {
			t.Errorf(ERR_SLICES)
		}
	}

	for _, fail := range fails {
		if !fail {
			t.Logf(MSG_SLICES_NOK)
		} else {
			t.Errorf(NERR_SLICES)
		}
	}
}
