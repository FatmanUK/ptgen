package protracker

import (
	"fmt"
	"testing"
)

func TestEmplaceRow(t *testing.T) {
	var err error
	testRow := [4]string{
		"A#4 01 C40",
		"G-3 02 C38",
		"--- -- F06",
		"C-4 04 C30",
	}
	p := PatternFactory()
	p, err = p.EmplaceRow(0, testRow)
	if err != nil {
		t.Errorf("%v", err)
	}
	for n := 0; n < CHANNELS_PER_ROW; n++ {
		s := p[0][n].Save()
		if s == testRow[n] {
			t.Logf(MSG_ROW_CELLX_OK, n)
		} else {
			m := fmt.Sprintf("%s %s", ERR_ROW_CELLX,
				FMT_TST_STRING)
			t.Errorf(m, n, testRow[n], s)
		}
	}
}

func TestIsTooManyPatterns(t *testing.T) {
	var err error
	n := uint8(MAX_PATTERNS - 1)
	err = IsTooManyPatterns(n)
	if err == nil {
		t.Logf(MSG_PTTN_NUM_OK, n)
	} else {
		t.Errorf(ERR_PTTN_NUM, n)
	}
	n = uint8(MAX_PATTERNS)
	err = IsTooManyPatterns(n)
	if err == nil {
		t.Logf(MSG_PTTN_NUM_OK, n)
	} else {
		t.Errorf(ERR_PTTN_NUM, n)
	}
	n = uint8(MAX_PATTERNS + 1)
	err = IsTooManyPatterns(n)
	if err != nil {
		t.Logf(MSG_PTTN_NUM_NOK, n)
	} else {
		t.Errorf(NERR_PTTN_NUM, n)
	}
}
