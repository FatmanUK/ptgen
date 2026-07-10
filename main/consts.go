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

const MOD_TOO_BIG_BYTES = (50 * 1024)

const MSG_WRN_TOO_BIG = "Warning: Mod is %d bytes. Typically, they're smaller than %d bytes. This is a very large mod."
