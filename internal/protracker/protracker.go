package protracker

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Mostly written by Google Gemini Pro. Tweaked extensively by me.

// Just one function exported:
// func WriteMod(w io.Writer, proj *ModProject) error
// WriteMod compiles proj into a binary ProTracker file stream.
// Populate proj, run WriteMod and the file appears in w.

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
// Note: e.g., "C-4" or "---" or ""
// Instrument: 1-31 (0 means no instrument)
// Effect: e.g., "C40", "047", "F03", or ""
type Cell struct {
	Note       string `json:"note"`
	Instrument uint8  `json:"instrument"`
	Effect     string `json:"effect"`
}

// A Row contains exactly CHANNELS_PER_ROW channels.
type Row [CHANNELS_PER_ROW]Cell

// A Pattern is a slice of ROWS_PER_PATTERN Rows.
type Pattern [ROWS_PER_PATTERN]Row

type ModInfo struct {
	Title    string
	Speed    uint8
	BPM      uint8
	SequenceLen int
	PatternsLen int
}

// ModProject represents your input JSON structure.
// Speed: Optional: 1-31 (0: default, always 6)
// BPM: Optional: 32-255 (0: default, always 125)
// Patterns: a slice of Patterns
type ModProject struct {
	Title    string    `json:"title"`
	Speed    uint8     `json:"speed"`
	BPM      uint8     `json:"bpm"`
	Sequence []uint8   `json:"orderList"`
	Patterns []Pattern `json:"patterns"`
}

func ModProjectFactory() ModProject {
	return ModProject{
		Title:    "A Song With No Name",
		Speed:    DEFAULT_SPEED,
		BPM:      DEFAULT_BPM,
		Sequence: []uint8{0},
		Patterns: []Pattern{Pattern{}},
	}
}

func (re ModProject) ModInfoFactory() ModInfo {
	return ModInfo{
		Title: re.Title,
		Speed: re.Speed,
		BPM: re.BPM,
		SequenceLen: len(re.Sequence),
		PatternsLen: len(re.Patterns),
	}
}

func (re ModProject) IsOrderListValid() bool {
	// TODO: check that each entry in the order list refers to
	// a pattern.
	return true
}

// Parse Note - Validate and map note to period
func parseNote(n string) (uint16, bool) {
	isValid := true
	var p uint16 = 0
	if n != "" && n != "---" {
		p, isValid = periodMap[strings.ToUpper(n)]
	}
	return p, isValid
}

// Parse Hex Effect String
// Return: cmd, param, isValid, error
func parseEffect(e string) (uint8, uint8, bool, error) {
	var cmd uint8 = 0
	var param uint8 = 0
	var errEffect error
	// Ignore empty strings or std empty tracker representations
	if e != "" && e != "000" && e != "---" {
		if len(e) != 3 {
			errEffect = fmt.Errorf(MSG_INV_EFFT, e)
			return 0, 0, false, errEffect
		}
		cmdVal, err := strconv.ParseUint(e[0:1], 16, 8)
		if err != nil {
			errEffect = fmt.Errorf(MSG_INV_CMMD, e)
			return 0, 0, false, errEffect
		}
		cmd = uint8(cmdVal)
		paramVal, err := strconv.ParseUint(e[1:3], 16, 8)
		if err != nil {
			errEffect = fmt.Errorf(MSG_INV_PARM, e)
			return 0, 0, false, errEffect
		}
		param = uint8(paramVal)
	}
	return cmd, param, true, errEffect
}

// Pack cell data into exactly 4 bytes.
// Wait. Pack 5 bytes into 4? That doesn't work mathematically.
// 4 bytes is 32 bits. We have to throw away 8 bits somewhere.
// We're only using 12 bits of period. And we only need the lower
// 4 bits of cmd. Can lose the 4 MSB on both. Packing achieved.
func packBytes(i uint8, p uint16, ec uint8, ep uint8) [4]byte {
	var out [4]byte
	// Byte 0: Instrument (upper 4 bits) | Period (upper 4 bits)
	out[0] = (i & 0xF0) | uint8((p & 0x0F00)>>8)
	// Byte 1: Period (lower 8 bits)
	out[1] = uint8(p & 0x00FF)
	// Byte 2: Instrument (lower 4 bits) | Effect Command (4 bits)
	out[2] = ((i & 0x0F) << 4) | (ec & 0x0F)
	// Byte 3: Effect Parameter
	out[3] = ep
	return out
}

// encodeCell packs a JSON cell into the 4-byte ProTracker format,
// parsing MilkyTracker-style hex effect strings.
func encodeCell(c Cell) ([4]byte, error) {
	var out [4]byte
	var period uint16
	var cmd uint8
	var param uint8
	var isValid bool
	period, isValid = parseNote(c.Note)
	if !isValid {
		return out, fmt.Errorf(MSG_INV_NOTE, c.Note)
	}
	// Instrument must be at least 1. 0 is reserved internally,
	// but can be selected. (It just can't be set.)
	if c.Instrument > 31 {
		return out, fmt.Errorf(MSG_INV_INST, c.Instrument)
	}
	cmd, param, isValid, err := parseEffect(c.Effect)
	if !isValid {
		return out, err
	}
	out = packBytes(c.Instrument, period, cmd, param)
	return out, nil
}

func fstPttnValid(o []uint8, p []Pattern) (uint8, error) {
	if len(o) == 0 {
		return 0, errors.New(MSG_ERR_EMPTY_OLST)
	}
	// Identify the pattern that will play first in the order list
	firstPatternIdx := o[0]
	if int(firstPatternIdx) >= len(p) {
		err := fmt.Errorf(MSG_ERR_OOB_PTTN, firstPatternIdx)
		return firstPatternIdx, err
	}
	if len(p[firstPatternIdx]) == 0 {
		err := fmt.Errorf(MSG_ERR_EMPTY_PTTN, firstPatternIdx)
		return firstPatternIdx, err
	}
	return firstPatternIdx, nil
}

func prepareCommands(s uint8, b uint8) ([]string, error) {
	var commands []string
	if s > 0 {
		if s >= 32 {
			return commands, fmt.Errorf(MSG_INV_SPD, s)
		}
		commands = append(commands, fmt.Sprintf("F%02X", s))
	}
	if b > 0 {
		if b < 32 {
			return commands, fmt.Errorf(MSG_INV_BPM, b)
		}
		commands = append(commands, fmt.Sprintf("F%02X", b))
	}
	return commands, nil
}

// Scan the 4 channels of Row 0 for empty effect fields
// Standard tracker empty signals are "", "000", or "---"
// Add speed and BPM commands there, if specified
// Return: number of commands injected
func injectCommands(p *ModProject, i uint8, commands []string) int {
	pattern := p.Patterns[i]
	rowZero := &pattern[0]
	cmdIdx := 0
	for ch := 0; ch < 4 && cmdIdx < len(commands); ch++ {
		eff := strings.TrimSpace(rowZero[ch].Effect)
		if eff == "" || eff == "000" || eff == "---" {
			rowZero[ch].Effect = commands[cmdIdx]
			cmdIdx++
		}
	}
	return cmdIdx
}

// Preprocess metadata and inject F-commands into Row 0
// injectInitialTempo scans the first row of the starting pattern and
// attempts to inject Fxx speed/BPM commands into empty effect slots.
func injectInitialTempo(proj *ModProject) error {
	// If neither option is set, there's nothing to inject
	if proj.Speed == 0 && proj.BPM == 0 {
		return nil
	}
	fstPttnIdx, err := fstPttnValid(proj.Sequence, proj.Patterns)
	if err != nil {
		return err
	}
	// Prepare the targets we need to inject
	commandsToInject, err := prepareCommands(proj.Speed, proj.BPM)
	if err != nil {
		return err
	}
	cmdIdx := injectCommands(proj, fstPttnIdx, commandsToInject)
	// If we still have commands left over, row 0 was
	// too saturated with user effects
	if cmdIdx < len(commandsToInject) {
		return fmt.Errorf(MSG_ERR_EFFT_FULL, fstPttnIdx)
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

func writeHeader(w io.Writer, proj *ModProject) {
	// 1. Write song title
	titleBytes := make([]byte, 20)
	copy(titleBytes, proj.Title)
	w.Write(titleBytes)
	// 2. Write 31 Empty Sample Headers (30 bytes each)
	emptySample := make([]byte, 30)
	// 3. Repeat Length must be 1 word for empty samples
	emptySample[29] = 0x01
	for i := 0; i < 31; i++ {
		w.Write(emptySample)
	}
	songLength := uint8(len(proj.Sequence))
	w.Write([]byte{songLength, 0x7F}) // Song length, Restart byte
	sequenceTable := make([]byte, 128)
	copy(sequenceTable, proj.Sequence)
	w.Write(sequenceTable)
	// 4. Write Magic String
	w.Write([]byte(MAGIC_BYTES))
}

func writePatterns(w io.Writer, proj *ModProject) error {
	// 5. Write Patterns
	//var err error
	for pIdx, pattern := range proj.Patterns {
		if len(pattern) != 64 {
			return fmt.Errorf(MSG_INV_PTTN, pIdx)
		}
		for rIdx, row := range pattern {
			err := writeCell(w, row, pIdx, rIdx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteMod compiles proj into a binary ProTracker file stream.
func WriteMod(w io.Writer, proj *ModProject) error {
	err := injectInitialTempo(proj)
	if err != nil {
		return fmt.Errorf(MSG_ERR_INIT, err)
	}
	if len(proj.Sequence) == 0 || len(proj.Sequence) > 128 {
		return errors.New(MSG_INV_OLST)
	}
	writeHeader(w, proj)
	err = writePatterns(w, proj)
	return err
}
