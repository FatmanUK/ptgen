package main

const metadataString = `Mod metadata:
Song Title:      {{ .Title }}
Format:          ProTracker-compatible MOD
Channels:        4
Speed:           {{ .Speed }}
BPM:             {{ .BPM }}
Pattern length:  64 rows (Len. 0x40h)
Unique patterns: {{ .Patterns }}
Order length:    {{ .OrderLen }}
Instruments:     {{ .Instruments }}`

// Formatting patterns.
const FMT_TST_STRING = `Exp: "%s" Got: "%s"`
const FMT_TST_NUMBER = `Exp: %d Got: %d`
const FMT_TST_VAR = `Exp: "%v" Got: "%v"`
const FMT_HTTP_STATUS = `HTTP Status Code: %d`

// Test data.
const SAMPLE2_PTTN_COUNT = 4
const SAMPLE2_PTTN_DIR = `test_data/sample2`

// Regexes.
const RGX_PTTN_FILE = `^pattern[0-9A-Fa-f]{2}\.txt$`
const RGX_HEX_EVIDENCE = `[A-Fa-f]`

// Arbitrary constants.
const MOD_TOO_BIG_KB = 50
const MOD_TOO_BIG_BYTES = (MOD_TOO_BIG_KB * 1024)
const URL_ST01 = `https://aminet.net/mods/inst/st-01.lha`
const SHA256_ST01 = `8bd8c62d542de794a5f843b5351de1c6f1bed8f7d7643a24a04d0ad8ce962553`
const URL_ST02 = `https://aminet.net/mods/inst/st-02.lha`
const SHA256_ST02 = `3ffbbf30e652a65aa47d08afd976990f67e2f9de5cbce171706d2c813d8dbce9`
const SAFETY_PAUSE_S = 2

// Messages.
const MSG_WRN_TOO_BIG = `Mod is %d bytes. This is unusually large (>%d bytes).`
const MSG_WRN_TOO_BIG_KB = `Mod is %d kB. This is unusually large (>%d kB).`
const MSG_SLICES_OK = `Slices match.`
const MSG_SLICES_NOK = `Slices do not match, but that's ok.`
const MSG_ROW_NUM_OK = `Row num is ok.`
const MSG_ROW_NUM_NOK = `Row num is wrong, but that's ok.`
const MSG_MOD_OLST_OK = `Order list is ok.`
const MSG_MOD_PTTN_COUNT_OK = `Pattern count is ok.`
const MSG_TOTAL_ROWS_PROCESSED = `Total rows processed: %d`
const MSG_JSON_DETECTED = `JSON metadata detected.`
const MSG_YAML_DETECTED = `YAML metadata detected.`
const MSG_LOADING_METADATA = `Loading metadata and pattern order.`
const MSG_LOADING_PTTNS = `Loading patterns.`
const MSG_LOADING_FILE = `Loading %s.`
const MSG_PTTN_FILE = `pattern%02x.txt`
const MSG_PTTN_ARE_SEQ = `Patterns are sequential.`
const MSG_PTTN_NOT_SEQ = `Patterns aren't sequential.`
const MSG_PTTN_VALID = `Patterns are valid.`
const MSG_HEX_DETECTED = `Hex detected.`
const MSG_SAFETY_PAUSE = `Pausing for %d seconds.`
const MSG_DOWNLOAD = `Downloading from %s. Saving as %s.`
const MSG_DOWNLOAD_IN_MEM = `Downloaded in memory.`
const MSG_CHECKSUM_OK = `Checksum test passed.`
const MSG_REQUEST_DOWNLOAD = `Requesting sample archive download.`
const MSG_DOWNLOAD_STATUS_OK = `Download request succeeded.`

// Errors.
const ERR_SLICES = `Slices do not match.`
const NERR_SLICES = `Slices match, but that's wrong.`
const ERR_ROW_NUM = `Row num is wrong.`
const NERR_ROW_NUM = `Row num is ok, but that's wrong.`
const ERR_MOD_OLST = `Order list is wrong.`
const ERR_MOD_PTTN_COUNT = `Pattern count is wrong.`
const ERR_MOD_TITLE_LONG = `Title is too long.`
const ERR_HEX_NOT_DETECTED = `Hex not detected.`
const ERR_HEX_INVALID = `Invalid hex string: %w`
const ERR_CHECKSUM_MISMATCH = `Checksum mismatch.`
