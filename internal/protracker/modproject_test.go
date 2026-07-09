package protracker

import (
	"fmt"
	"testing"
)

func TestModProjectFactory(t *testing.T) {
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

func TestModInfoFactory(t *testing.T) {
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

func TestIsOrderListValid(t *testing.T) {
	p := ModProjectFactory()
	if p.IsOrderListValid() {
		t.Logf(MSG_MOD_LIST_OK)
	} else {
		t.Errorf(ERR_MOD_LIST)
	}
}
