package protracker

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

// Mostly written by Google Gemini Pro. Tweaked extensively by me.

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
// Return: cmd, param, error
func parseEffect(e string) (uint8, uint8, error) {
	if len(e) != 3 {
		return 0, 0, fmt.Errorf(ERR_CELL_EFFT, e)
	}
	cmdVal, err := uint8FromHexString(e[0:1])
	if err != nil {
		err2 := fmt.Errorf(ERR_CELL_EFFT_CMMD, e[0:1], err)
		return 0, 0, err2
	}
	paramVal, err := uint8FromHexString(e[1:3])
	if err != nil {
		err2 := fmt.Errorf(ERR_CELL_EFFT_PARM, e[1:3], err)
		return 0, 0, err2
	}
	return uint8(cmdVal), uint8(paramVal), nil
}

// Pack cell data into exactly 4 bytes.
// Wait. Pack 5 bytes into 4? That doesn't work mathematically.
// 4 bytes is 32 bits. We have to throw away 8 bits somewhere.
// We're only using 12 bits of period. And we only need the lower
// 4 bits of cmd. Can lose the 4 MSB on both. Packing achieved.
// Byte 0: Instrument (upper 4 bits) | Period (upper 4 bits)
// Byte 1: Period (lower 8 bits)
// Byte 2: Instrument (lower 4 bits) | Effect Command (4 bits)
// Byte 3: Effect Parameter
func packBytes(i uint8, p uint16, ec uint8, ep uint8) [4]byte {
	var out [4]byte
	out[0] = (i & 0xF0) | uint8((p&0x0F00)>>8)
	out[1] = uint8(p & 0x00FF)
	out[2] = ((i & 0x0F) << 4) | (ec & 0x0F)
	out[3] = ep
	return out
}

// encodeCell packs a JSON cell into the 4-byte ProTracker format,
// parsing MilkyTracker-style hex effect strings.
// Instrument must be at least 1. 0 is reserved internally,
// but can be selected. (It just can't be set.)
func encodeCell(c Cell) ([4]byte, error) {
	var out [4]byte
	var cmd uint8
	var param uint8
	var err error
	period, isNoteValid := parseNote(c.Note)
	if !isNoteValid {
		return out, fmt.Errorf(ERR_CELL_NOTE, c.Note)
	}
	if c.Instr > 31 {
		return out, fmt.Errorf(ERR_CELL_INST, c.Instr)
	}
	if c.Effect != "" && c.Effect != "000" && c.Effect != "---" {
		cmd, param, err = parseEffect(c.Effect)
	}
	if err != nil {
		return out, err
	}
	out = packBytes(c.Instr, period, cmd, param)
	return out, nil
}

func prepareCommands(s uint8, b uint8) ([]string, error) {
	var commands []string
	if s > 0 {
		if s >= 32 {
			err := fmt.Errorf(ERR_CELL_EFFT_SPEED, s)
			return commands, err
		}
		commands = append(commands, fmt.Sprintf("F%02X", s))
	}
	if b > 0 {
		if b < 32 {
			err := fmt.Errorf(ERR_CELL_EFFT_BPM, b)
			return commands, err
		}
		commands = append(commands, fmt.Sprintf("F%02X", b))
	}
	return commands, nil
}

// Scan the 4 channels of Row 0 for empty effect fields
// Standard tracker empty signals are "", "000", or "---"
// Add speed and BPM commands there, if specified
// Return: number of commands injected
func injectCommands(p *ModProject, i uint8, commands []string) uint8 {
	pattern := &p.Patterns[i]
	rowZero := &pattern[0]
	var cmdIdx uint8 = 0
	for ch := 0; ch < 4 && cmdIdx < uint8(len(commands)); ch++ {
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
	if proj.Speed == 0 && proj.BPM == 0 {
		return nil
	}
	i, err := proj.IsFirstPatternValid()
	if err != nil {
		return err
	}
	commands, err := prepareCommands(proj.Speed, proj.BPM)
	if err != nil {
		return err
	}
	cmdIdx := injectCommands(proj, i, commands)
	if cmdIdx < uint8(len(commands)) {
		return fmt.Errorf(ERR_PTTN_ZERO_FULL, i)
	}
	return nil
}

// Validates and packs a 30-byte ProTracker sample header.
func encodeInstrumentHeader(inst Instrument) ([30]byte, error) {
	var out [30]byte
	// 1. Name (Padded/truncated to exactly 22 bytes)
	copy(out[0:22], inst.Name)
	dataLen := len(inst.Data)
	if dataLen > 131070 {
		return out, fmt.Errorf(ERR_SAMPLE_TOO_LONG, inst.Name)
	}
	if dataLen%2 != 0 {
		return out, fmt.Errorf(ERR_SAMPLE_LENGTH_ODD,
			inst.Name, dataLen)
	}
	// 2. Length (Stored in words)
	binary.BigEndian.PutUint16(out[22:24], uint16(dataLen/2))
	// 3. Finetune & Volume
	out[24] = inst.Finetune & 0x0F
	vol := inst.Volume
	if vol > 64 {
		vol = 64
	}
	out[25] = vol
	// 4. Loop Points (Stored in words)
	// If loop length is 2 bytes or less, treat it as unlooped
	if inst.Length <= 2 {
		// Start = 0, Length = 1 word (standard for no loop)
		binary.BigEndian.PutUint16(out[26:28], 0)
		binary.BigEndian.PutUint16(out[28:30], 1)
	} else {
		if inst.Start+inst.Length > uint32(dataLen) {
			m := fmt.Errorf(ERR_SAMPLE_LOOP_INVALID,
				inst.Name)
			return out, m
		}
		s := inst.Start / 2
		l := inst.Length / 2
		binary.BigEndian.PutUint16(out[26:28], uint16(s))
		binary.BigEndian.PutUint16(out[28:30], uint16(l))
	}
	return out, nil
}

func writeCell(w io.Writer, row Row, p int, r int) error {
	for c, cell := range row {
		encodedBytes, err := encodeCell(cell)
		if err != nil {
			return fmt.Errorf(ERR_LOCATION, p, r, c, err)
		}
		_, err = w.Write(encodedBytes[:])
		if err != nil {
			return err
		}
	}
	return nil
}

func injectAndCheck(proj *ModProject) error {
	err := injectInitialTempo(proj)
	if err != nil {
		return fmt.Errorf(ERR_INJECT, err)
	}
	if !proj.IsOrderListValid() {
		return fmt.Errorf(ERR_MOD_LIST)
	}
	if len(proj.Instruments) > 31 {
		return fmt.Errorf(ERR_INSTRUMENT_TOO_MANY)
	}
	return err
}

func isSlotPopulated(slot *Instrument) bool {
	//return (slot != nil)
	return (slot.ID > 0)
}

// 0. Pre-process sparse instruments into a rigid 31-slot lookup map
func preProcessInstruments(proj *ModProject) ([31]Instrument, error) {
	var slots [31]Instrument
	for _, i := range proj.Instruments {
		if i.ID < 1 || i.ID > 31 {
			m := fmt.Errorf(ERR_INSTRUMENT_INVALID_ID,
				i.Name, i.ID)
			return slots, m
		}
		if slots[i.ID].ID > 0 {
			m := fmt.Errorf(ERR_INSTRUMENT_DUPLICATE,
				i.ID)
			return slots, m
		}
		slots[i.ID] = i
	}
//	for k := range slots {
//		m := fmt.Sprintf("%d: %s", k, slots[k].Save())
//		log.Println(m)
//	}
	return slots, nil
}

// 1. Write song title
func writeTitle(w io.Writer, title string) error {
	titleBytes := make([]byte, 20)
	copy(titleBytes, title)
	_, err := w.Write(titleBytes)
	if err != nil {
		return err
	}
	return nil
}

func writeSampleHeader(w io.Writer, instr *Instrument) error {
	var err error
	if isSlotPopulated(instr) {
		headerBytes, err := encodeInstrumentHeader(*instr)
		if err != nil {
			//fmt.Errorf(ERR_INSTRUMENT_SLOT, i+1, err)
			return err
		}
		_, err = w.Write(headerBytes[:])
	} else {
		emptySample := make([]byte, 30)
		emptySample[29] = 0x01 // Repeat length = 1 word
		_, err = w.Write(emptySample)
	}
	return err
}

// 2. Write 31 Sample Headers (30 bytes each)
func writeSampleHeaders(w io.Writer, slots [31]Instrument) error {
	for i := 0; i < 31; i++ {
		err := writeSampleHeader(w, &(slots[i]))
		if err != nil {
			return err
		}
//		m := fmt.Sprintf("%d: %s", i, slots[i].Save())
//		log.Println(m)
	}
	return nil
}

// 3. Write Sequence Data
func writeOrderList(w io.Writer, orderList []uint8) error {
	songLength := uint8(len(orderList))
	_, err := w.Write([]byte{songLength, 0x7F})
	if err != nil {
		return err
	}
	orderTable := make([]byte, 128)
	copy(orderTable, orderList)
	_, err = w.Write(orderTable)
	if err != nil {
		return err
	}
	return nil
}

// 4. Write Magic String
func writeMagic(w io.Writer) error {
	_, err := w.Write([]byte(MAGIC_BYTES))
	if err != nil {
		return err
	}
	return nil
}

// 5. Write Patterns
func writePatterns(w io.Writer, pttns []Pattern) error {
	for pIdx, pattern := range pttns {
		for rIdx, row := range pattern {
			err := writeCell(w, row, pIdx, rIdx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 6. Write Raw PCM Sample Data block
func writeSamples(w io.Writer, slots [31]Instrument) error {
	for i := 0; i < 31; i++ {
		d := slots[i].Data
		if isSlotPopulated(&(slots[i])) && len(d) > 0 {
			if _, err := w.Write(d); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeHeaders(w io.Writer,
	proj *ModProject) ([31]Instrument, error) {
	slots, err := preProcessInstruments(proj)
	if err != nil {
		return slots, err
	}
	err = injectAndCheck(proj)
	if err != nil {
		return slots, err
	}
	err = writeTitle(w, proj.Title)
	if err != nil {
		return slots, err
	}
	return slots, writeSampleHeaders(w, slots)
}

// WriteMod compiles proj into a binary ProTracker file stream.
func WriteMod(w io.Writer, proj *ModProject) error {
	slots, err := writeHeaders(w, proj)
	if err != nil {
		return err
	}
	err = writeOrderList(w, proj.OrderList)
	if err != nil {
		return err
	}
	err = writeMagic(w)
	if err != nil {
		return err
	}
	err = writePatterns(w, proj.Patterns)
	if err != nil {
		return err
	}
	return writeSamples(w, slots)
}
