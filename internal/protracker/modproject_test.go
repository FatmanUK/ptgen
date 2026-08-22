package protracker

import (
	"fmt"
	"testing"
)

const TST_OK_MOD_TITLE = `Title is ok.`
const TST_NO_MOD_TITLE = `Title is wrong.`

const TST_OK_MOD_SPEED = `Speed is ok.`
const TST_NO_MOD_SPEED = `Speed is wrong.`

const TST_OK_MOD_BPM = `BPM is ok.`
const TST_NO_MOD_BPM = `BPM is wrong.`

const TST_OK_MOD_OLEN = `OrderLen is ok.`
const TST_NO_MOD_OLEN = `OrderLen is wrong.`

const TST_OK_MOD_PTTN = `Patterns is ok.`
const TST_NO_MOD_PTTN = `Patterns is wrong.`

const TST_OK_MOD_OLST = `OrderList is ok.`
const TST_NO_MOD_OLST = `OrderList is wrong.`

func TestModProjectFactory_Must_Succeed(t *testing.T) {
	p := ModProjectFactory(make(chan string))
	if p.Title == DEFAULT_TITLE {
		t.Logf(TST_OK_MOD_TITLE)
	} else {
		m := fmt.Sprintf("%s %s", TST_NO_MOD_TITLE, FMT_TST_STR)
		t.Errorf(m, DEFAULT_TITLE, p.Title)
	}
	if p.Speed == DEFAULT_SPEED {
		t.Logf(TST_OK_MOD_SPEED)
	} else {
		m := fmt.Sprintf("%s %s", TST_NO_MOD_SPEED, FMT_TST_NUM)
		t.Errorf(m, DEFAULT_SPEED, p.Speed)
	}
	if p.BPM == DEFAULT_BPM {
		t.Logf(TST_OK_MOD_BPM)
	} else {
		m := fmt.Sprintf("%s %s", TST_NO_MOD_BPM, FMT_TST_NUM)
		t.Errorf(m, DEFAULT_BPM, p.BPM)
	}
}

func TestModInfoFactory_Must_Succeed(t *testing.T) {
	p := ModProjectFactory(make(chan string))
	i := p.ModInfoFactory()
	if i.Title == DEFAULT_TITLE {
		t.Logf(TST_OK_MOD_TITLE)
	} else {
		m := fmt.Sprintf("%s %s", TST_NO_MOD_TITLE, FMT_TST_STR)
		t.Errorf(m, DEFAULT_TITLE, i.Title)
	}
	if i.Speed == DEFAULT_SPEED {
		t.Logf(TST_OK_MOD_SPEED)
	} else {
		m := fmt.Sprintf("%s %s", TST_NO_MOD_SPEED, FMT_TST_NUM)
		t.Errorf(m, DEFAULT_SPEED, i.Speed)
	}
	if i.BPM == DEFAULT_BPM {
		t.Logf(TST_OK_MOD_BPM)
	} else {
		m := fmt.Sprintf("%s %s", TST_NO_MOD_BPM, FMT_TST_NUM)
		t.Errorf(m, DEFAULT_BPM, i.BPM)
	}
	if i.OrderLen == 1 {
		t.Logf(TST_OK_MOD_OLEN)
	} else {
		m := fmt.Sprintf("%s %s", TST_NO_MOD_OLEN, FMT_TST_NUM)
		t.Errorf(m, 1, i.OrderLen)
	}
	if i.Patterns == 1 {
		t.Logf(TST_OK_MOD_PTTN)
	} else {
		m := fmt.Sprintf("%s %s", TST_NO_MOD_PTTN, FMT_TST_NUM)
		t.Errorf(m, 1, i.Patterns)
	}
}

func TestIsTooManyPatterns_Must_Succeed(t *testing.T) {
//173:func (p *ModProject) isTooManyPatterns() error {
/*
	var err error
	n := uint8(MAX_PATTERNS - 1)
	err = p.isTooManyPatterns(n)
	if err == nil {
		t.Logf(MSG_PTTN_NUM_OK, n)
	} else {
		t.Errorf(ERR_PTTN_NUM, n)
	}
	n = uint8(MAX_PATTERNS)
	err = p.isTooManyPatterns(n)
	if err == nil {
		t.Logf(MSG_PTTN_NUM_OK, n)
	} else {
		t.Errorf(ERR_PTTN_NUM, n)
	}
	n = uint8(MAX_PATTERNS + 1)
	err = p.isTooManyPatterns(n)
	if err != nil {
		t.Logf(MSG_PTTN_NUM_NOK, n)
	} else {
		t.Errorf(NERR_PTTN_NUM, n)
	}
*/
}

func TestIsTooManyPatterns_Must_Fail(t *testing.T) {
//173:func (p *ModProject) isTooManyPatterns() error {
}

func TestIsOrderListValid_Must_Succeed(t *testing.T) {
	p := ModProjectFactory(make(chan string))
	if p.isOrderListValid() {
		t.Logf(TST_OK_MOD_OLST)
	} else {
		t.Errorf(TST_NO_MOD_OLST)
	}
}

func TestIsOrderListValid_Must_Fail(t *testing.T) {
//180:func (p *ModProject) isOrderListValid() bool {
}

func TestIsGetFirstPatternIndex_Must_Succeed(t *testing.T) {
//193:func (p *ModProject) getFirstPatternIndex() (uint8, error) {
}

func TestIsGetFirstPatternIndex_Must_Fail(t *testing.T) {
//193:func (p *ModProject) getFirstPatternIndex() (uint8, error) {
}

func TestIsInjectInitialTempo_Must_Succeed(t *testing.T) {
//208:func (p *ModProject) injectInitialTempo() error {
}

func TestIsInjectInitialTempo_Must_Fail(t *testing.T) {
//208:func (p *ModProject) injectInitialTempo() error {
}

func TestPreProcessInstruments_Must_Succeed(t *testing.T) {
//236:func (p *ModProject) preProcessInstruments() ([31]Instrument, error) {
}

func TestPreProcessInstruments_Must_Fail(t *testing.T) {
//236:func (p *ModProject) preProcessInstruments() ([31]Instrument, error) {
}

func TestDefaultOrderlist_Must_Succeed(t *testing.T) {
	proj := ModProjectFactory(make(chan string))
	pttn := PatternFactory()
	proj.Patterns = []Pattern{pttn, pttn, pttn, pttn, pttn}
	proj.OrderList = proj.DefaultOrderlist()
	if proj.isOrderListValid() {
		t.Logf(TST_OK_MOD_OLST)
	} else {
		t.Errorf(TST_NO_MOD_OLST)
	}
}

func TestDefaultOrderlist_Must_Fail(t *testing.T) {
//338:func (p *ModProject) DefaultOrderlist() []uint8 {
}

func TestCalculateLength_Must_Succeed(t *testing.T) {
//346:func (p *ModProject) CalculateLength() {
}

func TestCalculateLength_Must_Fail(t *testing.T) {
//346:func (p *ModProject) CalculateLength() {
}
