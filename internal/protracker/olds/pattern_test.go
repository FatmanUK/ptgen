package protracker

import (
	"fmt"
	"testing"
)

func TestPatternEmplaceRow_Must_Succeed(t *testing.T) {
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
			t.Logf(TST_OK_ROW_CELL, n)
		} else {
			m := fmt.Sprintf("%s %s", TST_NO_ROW_CELL, FMT_TST_STR)
			t.Errorf(m, n, testRow[n], s)
		}
	}
}

func TestPatternEmplaceRow_Must_Fail(t *testing.T) {
	//26:func (p *Pattern) EmplaceRow(r uint8, m [4]string) (Pattern, error) {
}

func TestPatternInjectCommands_Must_Succeed(t *testing.T) {
	//54:func (p *Pattern) InjectCommands(cmmds []string) uint8 {
}

func TestPatternInjectCommands_Must_Fail(t *testing.T) {
	//54:func (p *Pattern) InjectCommands(cmmds []string) uint8 {
}
