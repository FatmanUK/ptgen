package protracker

import (
	"fmt"
	"testing"
)

func TestPatternFactory_Must_Succeed(t *testing.T) {
//17:func PatternFactory() Pattern {
}

func TestPatternFactory_Must_Fail(t *testing.T) {
//17:func PatternFactory() Pattern {
}

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

func TestPatternLoad_Must_Succeed(t *testing.T) {
//39:func (p *Pattern) Load(logs chan string, fileName string,
}

func TestPatternLoad_Must_Fail(t *testing.T) {
//39:func (p *Pattern) Load(logs chan string, fileName string,
}

func TestPatternInjectCommands_Must_Succeed(t *testing.T) {
//54:func (p *Pattern) InjectCommands(cmmds []string) uint8 {
}

func TestPatternInjectCommands_Must_Fail(t *testing.T) {
//54:func (p *Pattern) InjectCommands(cmmds []string) uint8 {
}

func TestPatternWrite_Must_Succeed(t *testing.T) {
//67:func (p *Pattern) Write(w io.Writer, pid uint8) error {
}

func TestPatternWrite_Must_Fail(t *testing.T) {
//67:func (p *Pattern) Write(w io.Writer, pid uint8) error {
}

func TestPatternRead_Must_Succeed(t *testing.T) {
//83:func (p *Pattern) Read(file io.Reader, logs chan string,
	/*
	   func TestReadPattern(t *testing.T) {
	   	// func readPattern(file io.Reader, logs chan string, isHexRows bool) (pt.Pattern, error)
	   }
	*/
}

func TestPatternRead_Must_Fail(t *testing.T) {
//83:func (p *Pattern) Read(file io.Reader, logs chan string,
}

func TestPatternScanLoop_Must_Succeed(t *testing.T) {
//98:func (p *Pattern) ScanLoop(logs chan string, d ScanData) (uint8, error) {
}

func TestPatternScanLoop_Must_Fail(t *testing.T) {
//98:func (p *Pattern) ScanLoop(logs chan string, d ScanData) (uint8, error) {
}
