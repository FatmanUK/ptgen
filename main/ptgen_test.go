package main

import (
	"fmt"
	"testing"
)

// Formatting patterns.
const FMT_TST_VAR = `Expected: "%v" Got: "%v"`

const TST_OK_NAME = `Name init'd ok`
const TST_NO_NAME = `Name init'd wrong. Exp: "%s" Got: "%s"`

const TST_OK_VERSION = `Version init'd ok`
const TST_NO_VERSION = `Name init'd wrong. Exp: "%s" Got: "%s"`

const TST_OK_PROCESS_ARGS = `Processed args ok.`
const TST_NO_PROCESS_ARGS = `Error: %v`

const TST_OK_ARGS = `Args map is ok.`
const TST_NO_ARGS = `Args map not empty.`

// Only defining succeed as it's too simple for fail conditions.
func TestDocOptVars_Must_Succeed(t *testing.T) {
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

// Only defining succeed as it's too simple for fail conditions.
func TestDocOptVarsProcessArgs_Must_Succeed(t *testing.T) {
	// Outputs "badvalue badvalue". This uninitialised
	// name/version string is correct, but I don't know why it's
	// output at all. Doesn't happen in normal operation. Might
	// just be some quirk of docopt.
	dotv := DocOptVarsFactory()
	args, err := dotv.ProcessArgs()
	if err == nil {
		t.Logf(TST_OK_PROCESS_ARGS)
	} else {
		t.Errorf(TST_NO_PROCESS_ARGS, err)
	}
	if len(args) == 0 {
		t.Logf(TST_OK_ARGS)
	} else {
		m := fmt.Sprintf("%s %s", TST_NO_ARGS, FMT_TST_VAR)
		t.Errorf(m, map[string]any{}, args)
	}
}
