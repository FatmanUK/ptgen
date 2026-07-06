package protracker

// These are just defaults.
const DEFAULT_SPEED = 6
const DEFAULT_BPM   = 125

// These can be changed, if you don't want PT compatibility.
const ROWS_PER_PATTERN = 64
const CHANNELS_PER_ROW = 4

// Magic bytes.
const MAGIC_BYTES = "M.K."

// Message strings.
const MSG_INV_NOTE = "invalid note '%s': only octaves 3-5 allowed"
const MSG_INV_INST = "instrument %d out of bounds (1-31)"
const MSG_INV_EFFT = "effect '%s' must be 3 hex chars (e.g., 'C40')"
const MSG_INV_CMMD = "invalid effect command in '%s'"
const MSG_INV_PARM = "invalid effect parameter in '%s'"
const MSG_INV_OLST = "order list needs at least 1 and at most 128 patterns"
const MSG_INV_PTTN = "pattern %d must have exactly 64 rows"
const MSG_INV_SPD = "invalid initial speed %d: must be less than 32"
const MSG_INV_BPM = "invalid initial BPM %d: must be 32 or greater"

const MSG_LOCATION = "pattern %d, row %d, channel %d: %w"

const MSG_ERR_INIT = "tempo initialization error: %w"
const MSG_ERR_OOB_PTTN = "sequence references out-of-bounds pattern index %d"
const MSG_ERR_EMPTY_OLST = "cannot inject tempo parameters into an empty song sequence"
const MSG_ERR_EMPTY_PTTN = "starting pattern %d contains no rows"
const MSG_ERR_EFFT_FULL = "failed to inject initial song speed/BPM: row 0 of pattern %d does not have enough empty effect slots"
const MSG_ERR_FAILED = "Failed handling valid input: %v"
const MSG_ERR_EXPECTED = "\nExpected: [% X]\nGot:      [% X]"
const MSG_ERR_FAIL_FAIL = "Expected constraint violation error, but execution passed"
const MSG_ERR_WRONG_ERROR = "Expected error message containing '%s', got '%v'"
const MSG_ERR_FAIL_CRRPT = "Failed to notice a corrupt payload setup"
const MSG_ERR_FAIL_VALID = "Compilation pipeline failed on valid payload setup: %v"
const MSG_ERR_SIZE = "Binary payload footprint mismatch. Expected %d bytes, got %d"
const MSG_ERR_TITLE = "Title field error. Expected custom byte-padding format layout"
const MSG_ERR_BAD_FMT = "Format validation failure. 'M.K.' magic bytes missing from target offset 1080 (got %q)"
const MSG_ERR_BAD_INS = "Instrument sample slot header #%d failed configuration constraints setup"
