package protracker

import (
	"fmt"
	"testing"
)

func TestModProjectFactory_Must_Succeed(t *testing.T) {
	p := ModProjectFactory()
	if p.Title == DEFAULT_TITLE {
		t.Logf(MSG_MOD_TITLE_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_MOD_TITLE,
			FMT_TST_STRING)
		t.Errorf(m, DEFAULT_TITLE, p.Title)
	}
	if p.Speed == DEFAULT_SPEED {
		t.Logf(MSG_MOD_SPEED_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_MOD_SPEED,
			FMT_TST_NUMBER)
		t.Errorf(m, DEFAULT_SPEED, p.Speed)
	}
	if p.BPM == DEFAULT_BPM {
		t.Logf(MSG_MOD_BPM_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_MOD_BPM,
			FMT_TST_NUMBER)
		t.Errorf(m, DEFAULT_BPM, p.BPM)
	}
}

func TestModProjectFactory_Must_Fail(t *testing.T) {
}

func TestModInfoFactory_Must_Succeed(t *testing.T) {
	p := ModProjectFactory()
	i := p.ModInfoFactory()
	if i.Title == DEFAULT_TITLE {
		t.Logf(MSG_MOD_TITLE_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_MOD_TITLE,
			FMT_TST_STRING)
		t.Errorf(m, DEFAULT_TITLE, i.Title)
	}
	if i.Speed == DEFAULT_SPEED {
		t.Logf(MSG_MOD_SPEED_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_MOD_SPEED,
			FMT_TST_NUMBER)
		t.Errorf(m, DEFAULT_SPEED, i.Speed)
	}
	if i.BPM == DEFAULT_BPM {
		t.Logf(MSG_MOD_BPM_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_MOD_BPM,
			FMT_TST_NUMBER)
		t.Errorf(m, DEFAULT_BPM, i.BPM)
	}
	if i.OrderLen == 1 {
		t.Logf(MSG_MOD_ORDER_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_MOD_ORDER,
			FMT_TST_NUMBER)
		t.Errorf(m, 1, i.OrderLen)
	}
	if i.Patterns == 1 {
		t.Logf(MSG_MOD_PTTN_OK)
	} else {
		m := fmt.Sprintf("%s %s", ERR_MOD_PTTN,
			FMT_TST_NUMBER)
		t.Errorf(m, 1, i.Patterns)
	}
}

func TestModInfoFactory_Must_Fail(t *testing.T) {
}

func TestOutputEverything_Must_Succeed(t *testing.T) {
	//60:func (p *ModProject) OutputEverything(logs chan string) error {
}

func TestOutputEverything_Must_Fail(t *testing.T) {
	//60:func (p *ModProject) OutputEverything(logs chan string) error {
}

func TestPopulatePatterns_Must_Succeed(t *testing.T) {
	//91:func (p *ModProject) PopulatePatterns(logs chan string,
}

func TestPopulatePatterns_Must_Fail(t *testing.T) {
	//91:func (p *ModProject) PopulatePatterns(logs chan string,
}

func TestIsTooManyPatterns_Must_Succeed(t *testing.T) {
	//115:func (p *ModProject) isTooManyPatterns() error {
	/*
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
	*/
}

func TestIsTooManyPatterns_Must_Fail(t *testing.T) {
	//115:func (p *ModProject) isTooManyPatterns() error {
}

func TestIsOrderListValid_Must_Succeed(t *testing.T) {
	p := ModProjectFactory()
	if p.isOrderListValid() {
		t.Logf(MSG_MOD_LIST_OK)
	} else {
		t.Errorf(ERR_MOD_LIST)
	}
}

func TestIsOrderListValid_Must_Fail(t *testing.T) {
	//122:func (p *ModProject) isOrderListValid() bool {
}

func TestIsGetFirstPatternIndex_Must_Succeed(t *testing.T) {
	//135:func (p *ModProject) getFirstPatternIndex() (uint8, error) {
}

func TestIsGetFirstPatternIndex_Must_Fail(t *testing.T) {
	//135:func (p *ModProject) getFirstPatternIndex() (uint8, error) {
}

func TestIsInjectInitialTempo_Must_Succeed(t *testing.T) {
	//150:func (p *ModProject) injectInitialTempo() error {
}

func TestIsInjectInitialTempo_Must_Fail(t *testing.T) {
	//150:func (p *ModProject) injectInitialTempo() error {
}

func TestPreProcessInstruments_Must_Succeed(t *testing.T) {
	//178:func (p *ModProject) preProcessInstruments() ([31]Instrument, error) {
}

func TestPreProcessInstruments_Must_Fail(t *testing.T) {
	//178:func (p *ModProject) preProcessInstruments() ([31]Instrument, error) {
}

func TestWriteHeaders_Must_Succeed(t *testing.T) {
	//211:func (p *ModProject) writeHeaders(
}

func TestWriteHeaders_Must_Fail(t *testing.T) {
	//211:func (p *ModProject) writeHeaders(
}

func TestWriteMod_Must_Succeed(t *testing.T) {
	//235:func (p *ModProject) WriteMod(w io.Writer) error {
}

func TestWriteMod_Must_Fail(t *testing.T) {
	//235:func (p *ModProject) WriteMod(w io.Writer) error {
}

func TestReadMetadata_Must_Succeed(t *testing.T) {
	// func readMetadata(proj *pt.ModProject, logs chan string, file *os.File) (*pt.ModProject, error)
}

func TestReadMetadata_Must_Fail(t *testing.T) {
	// func readMetadata(proj *pt.ModProject, logs chan string, file *os.File) (*pt.ModProject, error)
}

func TestDefaultOrderlist_Must_Succeed(t *testing.T) {
	proj := ModProjectFactory()
	pttn := PatternFactory()
	proj.Patterns = []Pattern{pttn, pttn, pttn, pttn, pttn}
	proj.OrderList = proj.DefaultOrderlist()
	if proj.isOrderListValid() {
		t.Logf(MSG_MOD_OLST_OK)
	} else {
		t.Errorf(ERR_MOD_OLST)
	}
}

func TestDefaultOrderlist_Must_Fail(t *testing.T) {
	//280:func (p *ModProject) DefaultOrderlist() []uint8 {
}

func TestCalculateLength_Must_Succeed(t *testing.T) {
	//288:func (p *ModProject) CalculateLength(logs chan string) {
}

func TestCalculateLength_Must_Fail(t *testing.T) {
	//288:func (p *ModProject) CalculateLength(logs chan string) {
}

func TestPopulateMetadata_Must_Succeed(t *testing.T) {
	//296:func (p *ModProject) PopulateMetadata(logs chan string,
}

func TestPopulateMetadata_Must_Fail(t *testing.T) {
	//296:func (p *ModProject) PopulateMetadata(logs chan string,
}
