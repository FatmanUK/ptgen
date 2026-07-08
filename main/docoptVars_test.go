package main

import (
	"testing"
)

func TestDocOptVars(t *testing.T) {
	dov := DocOptVarsFactory()
	if dov.Name == APP_NAME {
		t.Logf(`Name init'd ok`)
	} else {
		t.Errorf(`Name init'd wrong. Exp: "%s" Got: "%s"`,
			APP_NAME,
			dov.Name,
		)
	}
	if dov.Version == VERSION {
		t.Logf(`Version init'd ok`)
	} else {
		t.Errorf(`Name init'd wrong. Exp: "%s" Got: "%s"`,
			VERSION,
			dov.Version,
		)
	}
}
