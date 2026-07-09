package main

import (
	"testing"
)

func compareSlices(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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

type NumTests struct {
	number string
	isHex bool
	expected uint16
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
