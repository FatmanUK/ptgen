package protracker

// These are just defaults.
const DEFAULT_TITLE = `A Song With No Name`
const DEFAULT_SPEED = 6
const DEFAULT_BPM = 125
const MAX_PATTERNS = 128 // or 64 ? Gemini seems undecided.

// These can be changed, if you don't want PT compatibility.
const ROWS_PER_PATTERN = 64
const CHANNELS_PER_ROW = 4

// Magic bytes.
const MAGIC_BYTES = `M.K.`

// Regexes.
const RGX_ROWNUM = `[0-9A-Fa-f]{2}`
const RGX_NOTE = `[A-G-][-#][3-5-]`
const RGX_INSTR = `[0-9A-F-]{2}`
const RGX_EFFECT = `[0-9A-F-]{3}`
const RGX_CELL = `%s %s %s|%s %s|\.\.\.|---|-`
const RGX_CELL_SEPARATOR = `[:|]`
const RGX_ROW = `(%s) *%s *(%s) *%s *(%s) *%s *(%s) *%s *(%s) *%s?`

// Formatting patterns.
const FMT_ROW_PRETTY = `(%02x) %s | %s | %s | %s |`
const FMT_TST_STRING = `Exp: "%s" Got: "%s"`
const FMT_TST_NUMBER = `Exp: %d Got: %d`
const FMT_TST_VAR = `Exp: "%v" Got: "%v"`
const FMT_TST_VAR_X = `\nExpected: [% X]\nGot:      [% X]`

// Test data
const TEST_ROW_ROWNUM = 42
const TEST_CELL_STRING = `A#3 03 C38`

// Error strings
const ERR_PTTN_TOO_MANY = `Too many patterns.`
const ERR_LOCATION = `Pattern %d, Row %d, Channel %d: %w`
const ERR_INJECT = `Tempo injection error: %w`

// Test error strings
const ERR_CELL_NOTE = `Note "%s" is wrong. Must be a valid note in octaves 3-5.`
const ERR_CELL_INST = `Instrument %d is wrong. Minimum 1, maximum 31.`
const ERR_CELL_EFFT = `Effect "%s" is wrong. Must be 3 hex characters.`
const ERR_CELL_EFFT_CMMD = `Invalid effect command "%s". %v`
const ERR_CELL_EFFT_PARM = `Invalid effect parameter "%s". %v`
const ERR_CELL_EFFT_SPEED = `Invalid initial speed %d. Maximum 31.`
const ERR_CELL_EFFT_BPM = `Invalid initial BPM %d. Minimum 32.`
const ERR_CELL_REGEX = `Regex is wrong.`
const ERR_CELL_SAVE = `Save output is wrong.`

const ERR_CELL_FAILED = `Generic Cell failure: %v`
const ERR_CELL_SUCCESS = `Should have failed but didn't.`
const ERR_CELL_WRONG_ERROR = `Expected error message containing '%s', got '%v'`

const ERR_ROW_REGEX = `Regex is wrong.`
const ERR_ROW_SAVE = `Save output is wrong.`
const ERR_ROW_CELLS = `Number of Row Cells is wrong.`
const ERR_ROW_CELLX = `Cell %d is wrong.`

const ERR_PTTN_NUM = `Pattern quantity (%d) is wrong.`
const NERR_PTTN_NUM = `Pattern quantity (%d) is ok.`
const ERR_PTTN_ZERO_FULL = `Injection failed. Row 0 of first pattern %d has insufficient empty effect slots`

const ERR_MOD_TITLE = `Title is wrong.`
const ERR_MOD_SPEED = `Speed is wrong.`
const ERR_MOD_BPM = `BPM is wrong.`
const ERR_MOD_ORDER = `OrderLen is wrong.`
const ERR_MOD_PTTN = `Patterns is wrong.`
const ERR_MOD_LIST = `OrderList is wrong.`
const ERR_MOD_LIST_EMPTY = ERR_MOD_LIST + ` Must not be empty.`
const ERR_MOD_LIST_OOB = ERR_MOD_LIST + ` Invalid pattern reference %d.`
const ERR_MOD_LIST_OVERFLOW = ERR_MOD_LIST + ` Maximum 128 items.` // TODO: is this right? 128 is max patterns, also max orderlist?

const ERR_MOD_CORRUPT = `Failed to notice a corrupt payload setup.`
const ERR_MOD_FAILED = `Failed with valid payload. %x`
const ERR_MOD_SIZE = `Binary payload footprint mismatch. Expected %d bytes, got %d.`
const ERR_MOD_TITLE_PAD = `Title field error. Expected custom byte-padding format layout.`
const ERR_MOD_MAGIC = `Magic bytes missing from target offset 1080 (got %q).`
const ERR_MOD_INSTR = `Instrument sample slot header #%d failed constraints.`
const ERR_SAMPLE_TOO_LONG = `sample '%s' exceeds max length of 131070 bytes`
const ERR_SAMPLE_LENGTH_ODD = `sample '%s' must have an even byte length (got %d)`
const ERR_SAMPLE_LOOP_INVALID = `sample '%s' loop points exceed total sample length`
const ERR_INSTRUMENT_TOO_MANY = `a maximum of 31 instruments are supported`
const ERR_INSTRUMENT_INVALID_ID = `instrument '%s' has invalid ID %d (must be 1-31)`
const ERR_INSTRUMENT_DUPLICATE = `duplicate instrument ID %d detected`
const ERR_INSTRUMENT_SLOT = `instrument slot %d error: %w`

// Test ok strings
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
