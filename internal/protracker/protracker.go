package protracker

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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
const MSG_INV_SPD  = "invalid initial speed %d: must be less than 32"
const MSG_INV_BPM  = "invalid initial BPM %d: must be 32 or greater"

const MSG_LOCATION = "pattern %d, row %d, channel %d: %w"

const MSG_ERR_INIT = "tempo initialization error: %w"
const MSG_ERR_OOB_PTTN = "sequence references out-of-bounds pattern index %d"
const MSG_ERR_EMPTY_OLST = "cannot inject tempo parameters into an empty song sequence"
const MSG_ERR_EMPTY_PTTN = "starting pattern %d contains no rows"
const MSG_ERR_EFFT_FULL = "failed to inject initial song speed/BPM: row 0 of pattern %d does not have enough empty effect slots"

// Amiga PAL periods for MilkyTracker octaves 3, 4, and 5.
var periodMap = map[string]uint16{
	"C-3": 856, "C#3": 808, "D-3": 762, "D#3": 720,
	"E-3": 678, "F-3": 640, "F#3": 604, "G-3": 570,
	"G#3": 538, "A-3": 508, "A#3": 480, "B-3": 453,
	"C-4": 428, "C#4": 404, "D-4": 381, "D#4": 360,
	"E-4": 339, "F-4": 320, "F#4": 302, "G-4": 285,
	"G#4": 269, "A-4": 254, "A#4": 240, "B-4": 226,
	"C-5": 214, "C#5": 202, "D-5": 190, "D#5": 180,
	"E-5": 170, "F-5": 160, "F#5": 151, "G-5": 143,
	"G#5": 135, "A-5": 127, "A#5": 120, "B-5": 113,
}

// A Cell represents one channel on one row.
type Cell struct {
	Note       string `json:"note"`       // e.g., "C-4" or "---" or ""
	Instrument uint8  `json:"instrument"` // 1-31 (0 means no instrument)
	Effect     string `json:"effect"`     // e.g., "C40", "047", "F03", or ""
}

// A Row contains exactly CHANNELS_PER_ROW channels.
type Row [CHANNELS_PER_ROW]Cell

// A Pattern is a slice of ROWS_PER_PATTERN Rows.
type Pattern [ROWS_PER_PATTERN]Row

// ModProject represents your input JSON structure.
type ModProject struct {
	Title    string    `json:"title"`
	Speed    uint8     `json:"speed"` // Optional: 1-31 (0 means skip injection)
	BPM      uint8     `json:"bpm"`   // Optional: 32-255 (0 means skip injection)
	Sequence []uint8   `json:"sequence"`
	Patterns []Pattern `json:"patterns"` // Slice of patterns
}

// encodeCell packs a JSON cell into the 4-byte ProTracker format,
// parsing MilkyTracker-style hex effect strings.
func encodeCell(c Cell) ([4]byte, error) {
	var out [4]byte
	var period uint16 = 0

	// 1. Parse Note - Validate and map note to period
	if c.Note != "" && c.Note != "---" {
		p, exists := periodMap[strings.ToUpper(c.Note)]
		if !exists {
			return out, fmt.Errorf(MSG_INV_NOTE, c.Note)
		}
		period = p
	}

	// 2. Validate Instrument
	// TODO check low boundary is really 1 instead of 0
	if c.Instrument > 31 {
		return out, fmt.Errorf(MSG_INV_INST, c.Instrument)
	}

	// 3. Parse MilkyTracker Hex Effect String
	var cmd uint8 = 0
	var param uint8 = 0

	// Ignore empty strings or standard empty tracker representations
	if c.Effect != "" && c.Effect != "000" && c.Effect != "---" {
		if len(c.Effect) != 3 {
			return out, fmt.Errorf(MSG_INV_EFFT, c.Effect)
		}

		// Parse the 1-character command (Base 16, 8-bit size)
		cmdVal, err := strconv.ParseUint(c.Effect[0:1], 16, 8)
		if err != nil {
			return out, fmt.Errorf(MSG_INV_CMMD, c.Effect)
		}
		cmd = uint8(cmdVal)

		// Parse the 2-character parameter (Base 16, 8-bit size)
		paramVal, err := strconv.ParseUint(c.Effect[1:3], 16, 8)
		if err != nil {
			return out, fmt.Errorf(MSG_INV_PARM, c.Effect)
		}
		param = uint8(paramVal)
	}

	// 4. Pack into exactly 4 bytes
	// Reminder: we write period instead of c.Note
	// Byte 0: Instrument (upper 4 bits) | Period (upper 4 bits)
	out[0] = (c.Instrument & 0xF0) | uint8((period&0x0F00)>>8)
	// Byte 1: Period (lower 8 bits)
	out[1] = uint8(period & 0x00FF)
	// Byte 2: Instrument (lower 4 bits) | Effect Command (4 bits)
	out[2] = ((c.Instrument & 0x0F) << 4) | (cmd & 0x0F)
	// Byte 3: Effect Parameter
	out[3] = param

	return out, nil
}

// injectInitialTempo scans the first row of the starting pattern and attempts
// to inject Fxx speed/BPM commands into empty effect slots.
func injectInitialTempo(proj *ModProject) error {
	// If neither option is set, there's nothing to inject
	if proj.Speed == 0 && proj.BPM == 0 {
		return nil
	}

	if len(proj.Sequence) == 0 {
		return errors.New(MSG_ERR_EMPTY_OLST)
	}

	// Identify the pattern that will play first in the order list
	firstPatternIdx := proj.Sequence[0]
	if int(firstPatternIdx) >= len(proj.Patterns) {
		return fmt.Errorf(MSG_ERR_OOB_PTTN, firstPatternIdx)
	}

	firstPattern := proj.Patterns[firstPatternIdx]
	if len(firstPattern) == 0 {
		return fmt.Errorf(MSG_ERR_EMPTY_PTTN, firstPatternIdx)
	}

	// Prepare the targets we need to inject
	var commandsToInject []string
	if proj.Speed > 0 {
		if proj.Speed >= 32 {
			return fmt.Errorf(MSG_INV_SPD, proj.Speed)
		}
		fSpd := fmt.Sprintf("F%02X", proj.Speed)
		commandsToInject = append(commandsToInject, fSpd)
	}

	if proj.BPM > 0 {
		if proj.BPM < 32 {
			return fmt.Errorf(MSG_INV_BPM, proj.BPM)
		}
		fBPM := fmt.Sprintf("F%02X", proj.BPM)
		commandsToInject = append(commandsToInject, fBPM)
	}

	// Scan the 4 channels of Row 0 for empty effect fields
	// Standard tracker empty signals are "", "000", or "---"
	rowZero := &firstPattern[0]
	cmdIdx := 0

	for ch := 0; ch < 4 && cmdIdx < len(commandsToInject); ch++ {
		eff := strings.TrimSpace(rowZero[ch].Effect)
		if eff == "" || eff == "000" || eff == "---" {
			rowZero[ch].Effect = commandsToInject[cmdIdx]
			cmdIdx++
		}
	}

	// If we still have commands left over, row 0 was too saturated with user effects
	if cmdIdx < len(commandsToInject) {
		return fmt.Errorf(MSG_ERR_EFFT_FULL, firstPatternIdx)
	}

	return nil
}

func writeCell(w io.Writer, row Row, p int, r int) error {
	for cIdx, cell := range row {
		encodedBytes, err := encodeCell(cell)
		msgErr := fmt.Errorf(MSG_LOCATION, p, r, cIdx, err)
		if err != nil {
			return msgErr
		}
		w.Write(encodedBytes[:])
	}
	return nil
}

// WriteMod compiles the ModProject into a binary ProTracker file stream.
func WriteMod(w io.Writer, proj *ModProject) error {
	// 0. Preprocess metadata and inject initial F-commands into Row 0
	if err := injectInitialTempo(proj); err != nil {
		return fmt.Errorf(MSG_ERR_INIT, err)
	}

	if len(proj.Sequence) == 0 || len(proj.Sequence) > 128 {
		return errors.New(MSG_INV_OLST)
	}

	// 1. Write Title (20 bytes)
	titleBytes := make([]byte, 20)
	copy(titleBytes, proj.Title)
	w.Write(titleBytes)

	// 2. Write 31 Empty Sample Headers (30 bytes each)
	emptySample := make([]byte, 30)
	emptySample[29] = 0x01 // Repeat Length must be 1 word for empty samples
	for i := 0; i < 31; i++ {
		w.Write(emptySample)
	}

	// 3. Write Sequence Data
	songLength := uint8(len(proj.Sequence))
	w.Write([]byte{songLength, 0x7F}) // Song length and Restart byte

	sequenceTable := make([]byte, 128)
	copy(sequenceTable, proj.Sequence)
	w.Write(sequenceTable)

	// 4. Write Magic String
	w.Write([]byte(MAGIC_BYTES))

	// 5. Write Patterns
	var err error
	for pIdx, pattern := range proj.Patterns {
		if len(pattern) != 64 {
			return fmt.Errorf(MSG_INV_PTTN, pIdx)
		}
		for rIdx, row := range pattern {
			err = writeCell(w, row, pIdx, rIdx)
			if err != nil {
				return err
			}
		}
	}

	return err
}
