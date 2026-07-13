package main

import (
	"testing"
)

func TestDocOptVars(t *testing.T) {
	dov := DocOptVarsFactory()
	if dov.Name == APP_NAME {
		t.Logf(TST_OK_NAME)
	} else {
		t.Errorf(TST_NO_NAME, APP_NAME, dov.Name)
	}
	if dov.Version == VERSION {
		t.Logf(TST_OK_VERSION)
	} else {
		t.Errorf(TST_NO_VERSION, VERSION, dov.Version)
	}
}
