package main

const metadataString = `Mod metadata:
Song Title:      {{ .Title }}
Format:          ProTracker-compatible MOD
Channels:        4
Speed:           {{ .Speed }}
BPM:             {{ .BPM }}
Pattern length:  64 rows (Len. 0x40h)
Unique patterns: {{ .Patterns }}
Order length:    {{ .OrderLen }}`

// Formatting patterns.
const FMT_TST_STRING = `Exp: "%s" Got: "%s"`
const FMT_TST_NUMBER = `Exp: %d Got: %d`
const FMT_TST_VAR = `Exp: "%v" Got: "%v"`

// Arbitrary constants.
const MOD_TOO_BIG_KB = 50
const MOD_TOO_BIG_BYTES = (MOD_TOO_BIG_KB * 1024)

// Messages.
const MSG_WRN_TOO_BIG = `Mod is %d bytes. This is unusually large (>%d bytes).`
const MSG_WRN_TOO_BIG_KB = `Mod is %d kB. This is unusually large (>%d kB).`
const MSG_SLICES_OK = `Slices match.`
const MSG_ROW_NUM_OK = `Row num is ok.`
const MSG_MOD_OLST_OK = `Order list is ok.`

// Errors.
const ERR_SLICES = `Slices do not match.`
const ERR_ROW_NUM = `Row num is wrong.`
const ERR_MOD_OLST = `Order list is wrong.`
