package protracker

const MetadataString = `Mod metadata:
Song Title:      {{ .Title }}
Format:          ProTracker-compatible MOD
Channels:        4
Speed:           {{ .Speed }}
BPM:             {{ .BPM }}
Pattern length:  64 rows (Len. 0x40h)
Unique patterns: {{ .Patterns }}
Order length:    {{ .OrderLen }}
Instruments:     {{ .Instruments }}`

const TMP_WRKRND_CMDLINE = `/usr/bin/lha evifw=%s %s %s`

// These are just defaults.
const DEFAULT_TITLE = `A Song With No Name`
const DEFAULT_SPEED = 6
const DEFAULT_BPM = 125
const MAX_PATTERNS = 128 // or 64 ? Gemini seems undecided.

// These can be changed, if you don't want PT compatibility.
const ROWS_PER_PATTERN = 64
const CHANNELS_PER_ROW = 4

// Regexes.
const RGX_NOTE = `[A-G-][-#][3-5-]`
const RGX_INSTR = `[0-9A-F-]{2}`
const RGX_EFFECT = `[0-9A-F-]{3}`
const RGX_CELL = `%s %s %s|%s %s|\.\.\.|---|-`
const RGX_ROWNUM = `[0-9A-Fa-f]{2}`
const RGX_CELL_SEPARATOR = `[:|]`
const RGX_ROW = `(%s) *%s *(%s) *%s *(%s) *%s *(%s) *%s *(%s) *%s?`
const RGX_PTTN_FILE = `^pattern[0-9A-Fa-f]{2}\.txt$`
const RGX_HEX_EVIDENCE = `[A-Fa-f]`

// Magic bytes.
const MAGIC_BYTES = `M.K.`

// Warnings.
const MOD_TOO_BIG_KB = 50
const MOD_TOO_BIG_BYTES = (MOD_TOO_BIG_KB * 1024)
const MSG_WRN_TOO_BIG = `Mod is %d bytes. This is unusually large (>%d bytes).`
const MSG_WRN_TOO_BIG_KB = `Mod is %d kB. This is unusually large (>%d kB).`

// Test ok strings.
const MSG_CELL_NOTE_OK = `Note is ok.`
const MSG_CELL_INST_OK = `Instrument is ok.`
const MSG_CELL_EFFT_OK = `Effect is ok.`
const MSG_CELL_REGEX_OK = `Regex is ok.`
const MSG_CELL_SAVE_OK = `Save output is ok.`

const MSG_ROW_REGEX_OK = `Regex is ok.`
const MSG_ROW_SAVE_OK = `Save output is ok.`
const MSG_ROW_CELLS_OK = `Row has %d Cells.`
const MSG_ROW_CELLX_OK = `Cell %d is ok.`

const MSG_PTTN_NUM_OK = `Pattern quantity (%d) is ok.`      // awkward; can't find synonym
const MSG_PTTN_NUM_NOK = `Pattern quantity (%d) is not ok.` // awkward; can't find synonym

const MSG_MOD_TITLE_OK = `Title is ok.`
const MSG_MOD_SPEED_OK = `Speed is ok.`
const MSG_MOD_BPM_OK = `BPM is ok.`
const MSG_MOD_ORDER_OK = `OrderLen is ok.`
const MSG_MOD_PTTN_OK = `Patterns is ok.`
const MSG_MOD_LIST_OK = `OrderList is ok.`

// Test error strings
const ERR_CELL_NOTE = `Note "%s" is wrong. Must be a valid note in octaves 3-5.`
const ERR_CELL_INST = `Instrument %d is wrong. Minimum 1, maximum 31.`
const ERR_CELL_EFFT = `Effect "%s" is wrong. Must be 3 hex characters.`
const ERR_CELL_EFFT_CMMD = `Invalid effect command "%s". %v`
const ERR_CELL_EFFT_PARAM = `Invalid effect parameter "%s". %v`
const ERR_HEX_INVALID = `Invalid hex string: %w`
const ERR_CHECKSUM_MISMATCH = `Checksum mismatch.`
const ERR_SAMPLE_TOO_LONG = `sample '%s' exceeds max length of 131070 bytes`
const ERR_SAMPLE_LENGTH_ODD = `sample '%s' must have an even byte length (got %d)`
const ERR_SAMPLE_LOOP_INVALID = `sample '%s' loop points exceed total sample length`
const ERR_ARCH_CACHE_DIR = `Didn't get user cache directory: %v`
const ERR_ARCH_OPEN = `Failed to open archive: %v`
const ERR_ARCH_HEADER_PARSE = `failed header parse: %v`
const ERR_ARCH_EXTRACTION = `Extraction failed for file %s: %v`
const ERR_ARCH_SIZE_MISMATCH = `Size mismatch for %s: expected %d bytes, extracted %d`
const ERR_MOD_TITLE_LONG = `Title is too long.`
const ERR_PTTN_TOO_MANY = `Too many patterns.`
const ERR_MOD_LIST = `OrderList is wrong.`
const ERR_MOD_LIST_EMPTY = ERR_MOD_LIST + ` Must not be empty.`
const ERR_MOD_LIST_OOB = ERR_MOD_LIST + ` Invalid pattern reference %d.`
const ERR_MOD_LIST_OVERFLOW = ERR_MOD_LIST + ` Maximum 128 items.` // TODO: is this right? 128 is max patterns, also max orderlist?
const ERR_INSTRUMENT_TOO_MANY = `a maximum of 31 instruments are supported`
const ERR_INSTRUMENT_INVALID_ID = `instrument '%s' has invalid ID %d (must be 1-31)`
const ERR_INSTRUMENT_DUPLICATE = `duplicate instrument ID %d detected`
const ERR_INSTRUMENT_SLOT = `instrument slot %d error: %w`
const ERR_PTTN_ZERO_FULL = `Injection failed. Row 0 of first pattern %d has insufficient empty effect slots`
const ERR_INJECT = `Tempo injection error: %w`
const ERR_MOD_OLST = `Order list is wrong.`
const ERR_CELL_EFFT_SPEED = `Invalid initial speed %d. Maximum 31.`
const ERR_CELL_EFFT_BPM = `Invalid initial BPM %d. Minimum 32.`
const ERR_CELL_REGEX = `Regex is wrong.`
const ERR_CELL_SAVE = `Save output is wrong.`
const ERR_LOCATION = `Pattern %d, Row %d, Channel %d: %w`
const ERR_MOD_TITLE = `Title is wrong.`
const ERR_MOD_SPEED = `Speed is wrong.`
const ERR_MOD_BPM = `BPM is wrong.`
const ERR_MOD_ORDER = `OrderLen is wrong.`
const ERR_MOD_PTTN = `Patterns is wrong.`
const ERR_MOD_CORRUPT = `Failed to notice a corrupt payload setup.`
const ERR_MOD_FAILED = `Failed with valid payload. %x`
const ERR_MOD_SIZE = `Binary payload footprint mismatch. Expected %d bytes, got %d.`
const ERR_MOD_TITLE_PAD = `Title field error. Expected custom byte-padding format layout.`
const ERR_MOD_MAGIC = `Magic bytes missing from target offset 1080 (got %q).`
const ERR_MOD_INSTR = `Instrument sample slot header #%d failed constraints.`
const ERR_ROW_NUM = `Row num is wrong.`
const NERR_ROW_NUM = `Row num is ok, but that's wrong.`
const ERR_MOD_PTTN_COUNT = `Pattern count is wrong.`
const ERR_ROW_CELLS = `Number of Row Cells is wrong.`
const ERR_ROW_CELLX = `Cell %d is wrong.`
const ERR_ROW_REGEX = `Regex is wrong.`
const ERR_ROW_SAVE = `Save output is wrong.`
const ERR_HEX_NOT_DETECTED = `Hex not detected.`
const ERR_PTTN_NUM = `Pattern quantity (%d) is wrong.`
const NERR_PTTN_NUM = `Pattern quantity (%d) is ok.`

// Messages.
const MSG_CHECKSUM_OK = `Checksum test passed.`
const MSG_DOWNLOAD_IN_MEM = `Downloaded in memory.`
const MSG_REQUEST_DOWNLOAD = `Requesting sample archive download.`
const MSG_DOWNLOAD_STATUS_OK = `Download request succeeded.`
const MSG_LOADING_PTTNS = `Loading patterns.`
const MSG_LOADING_FILE = `Loading %s.`
const MSG_JSON_DETECTED = `JSON metadata detected.`
const MSG_YAML_DETECTED = `YAML metadata detected.`
const MSG_LOADING_METADATA = `Loading metadata and pattern order.`
const MSG_TOTAL_ROWS_PROCESSED = `Total rows processed: %d`
const MSG_SAFETY_PAUSE = `Pausing for %d seconds.`
const MSG_DOWNLOAD = `Downloading from %s. Saving as %s.`
const MSG_ARCH_HEADER_PARSE_OK = `Processing file: %s (Method: %s, Original Size: %d)`
const MSG_ROW_NUM_OK = `Row num is ok.`
const MSG_ROW_NUM_NOK = `Row num is wrong, but that's ok.`
const MSG_MOD_OLST_OK = `Order list is ok.`
const MSG_MOD_PTTN_COUNT_OK = `Pattern count is ok.`
const MSG_HEX_DETECTED = `Hex detected.`

// Formatting.
const FMT_HTTP_STATUS = `HTTP Status Code: %d`
const MSG_PTTN_FILE = `pattern%02x.txt`
const FMT_ROW_PRETTY = `(%02x) %s | %s | %s | %s |`
const FMT_TST_STRING = `Exp: "%s" Got: "%s"`
const FMT_TST_NUMBER = `Exp: %d Got: %d`
const FMT_TST_VARIANT = `Exp: "%v" Got: "%v"`

// Arbitrary constants.
const URL_ST01 = `https://aminet.net/mods/inst/st-01.lha`
const SHA256_ST01 = `8bd8c62d542de794a5f843b5351de1c6f1bed8f7d7643a24a04d0ad8ce962553`
const URL_ST02 = `https://aminet.net/mods/inst/st-02.lha`
const SHA256_ST02 = `3ffbbf30e652a65aa47d08afd976990f67e2f9de5cbce171706d2c813d8dbce9`
const SAFETY_PAUSE_S = 2
const FILE_ST01 = `st-01.lha`
const FILE_ST02 = `st-02.lha`

// Test data.
const TEST_ROW_ROWNUM = 42
const TEST_CELL_STRING = `A#3 03 C38`
const SAMPLE2_PTTN_COUNT = 4
const SAMPLE2_PTTN_DIR = `test_data/sample2`

/*
// Formatting patterns.
const FMT_TST_VAR_X = `\nExpected: [% X]\nGot:      [% X]`

// Test error strings
const ERR_CELL_FAILED = `Generic Cell failure: %v`
const ERR_CELL_SUCCESS = `Should have failed but didn't.`
const ERR_CELL_WRONG_ERROR = `Expected error message containing '%s', got '%v'`

// Test data.
const MSG_PTTN_ARE_SEQ = `Patterns are sequential.`
const MSG_PTTN_NOT_SEQ = `Patterns aren't sequential.`
const MSG_PTTN_VALID = `Patterns are valid.`
*/
