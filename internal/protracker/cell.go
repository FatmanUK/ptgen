package protracker

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"github.com/FatmanUK/fatgo/utils"
)

const RGX_NOTE = `[A-G-\.][-#\.][3-5-\.]`
const RGX_INST = `[0-9A-F-\.]{2}`
const RGX_EFFT = `[0-9A-F-\.]{3}`
const RGX_CELL = `%s %s %s|%s %s|\.{3}|-{3}|-`

const ERR_STR = `Actual value: "%s"`

const ERR_CELL_NOTE = `Note is wrong. Must be A-G# in octaves 3-5.`
const ERR_CELL_INST = `Instrument is wrong. Minimum 1, maximum 31.`
const ERR_CELL_EFFT = `Effect is wrong. Must be 3 hex characters.`
const ERR_CELL_CMMD = `Invalid effect command.`
const ERR_CELL_PARM = `Invalid effect parameter.`

// Amiga PAL periods for octaves 3, 4, and 5.
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

func isEmpty(e string) bool {
	return e == "" || e == "000" || e == "---" || e == "..."
}

// A Cell represents one channel on one row.
// Note: e.g., "C-4" or "---" or ""
// Instrument: 1-31 (0 means no instrument)
// Effect: e.g., "C40", "047", "F03", or ""
type Cell struct {
	Note   string `json:"note"`
	Instr  uint8  `json:"instr"`
	Effect string `json:"effect"`
}

func CellFactory() Cell {
	return Cell{
		Note:   "---",
		Instr:  0,
		Effect: "---",
	}
}

func CellRegexFactory() string {
	return fmt.Sprintf(RGX_CELL, RGX_NOTE, RGX_INST, RGX_EFFT,
		RGX_NOTE, RGX_INST)
}

// Load: deserialise struct from a string
func (c *Cell) Load(s string) error {
	if isEmpty(s) {
		return nil
	}
	l := len(s)
	if l >= 3 {
		c.Note = s[0:3]
	}
	if l >= 6 {
		i := s[4:6]
		if i != "--" {
			u, err := strconv.ParseUint(i, 16, 16)
			if err != nil {
				return err
			}
			c.Instr = uint8(u)
		}
	}
	if l >= 10 {
		c.Effect = s[7:10]
	}
	return nil
}

// Save: serialise struct to a string
func (c *Cell) Save() string {
	i := fmt.Sprintf("%02x", c.Instr)
	if i == "00" {
		i = "--"
	}
	return fmt.Sprintf("%s %s %s", c.Note, i, c.Effect)
}

// Validate and map note to period
func (c *Cell) getPeriod() uint16 {
	var p uint16 = 0
	if !isEmpty(c.Note) {
		p = periodMap[strings.ToUpper(c.Note)]
	}
	return p
}

// Parse Hex Effect String
// Return: cmd, param, error
func (c *Cell) getEffectCommand() (uint8, error) {
	var cmd uint8
	var err error
	if !isEmpty(c.Effect) {
		if len(c.Effect) != 3 {
			msgPair := [2]string{ERR_CELL_EFFT, ERR_STR}
			return 0, utils.Err(msgPair, c.Effect, nil)
		}
		cmdStr := c.Effect[0:1]
		cmd, err = utils.Uint8FromHexString(cmdStr)
		if err != nil {
			msgPair := [2]string{ERR_CELL_CMMD, ERR_STR}
			return 0, utils.Err(msgPair, cmdStr, err)
		}
	}
	return cmd, err
}

func (c *Cell) getEffectParameter() (uint8, error) {
	var param uint8
	var err error
	if !isEmpty(c.Effect) {
		if len(c.Effect) != 3 {
			msgPair := [2]string{ERR_CELL_EFFT, ERR_STR}
			return 0, utils.Err(msgPair, c.Effect, nil)
		}
		paramStr := c.Effect[1:3]
		param, err = utils.Uint8FromHexString(paramStr)
		if err != nil {
			msgPair := [2]string{ERR_CELL_PARM, ERR_STR}
			return 0, utils.Err(msgPair, paramStr, err)
		}
	}
	return param, err
}

// encodeCell packs a JSON cell into the 4-byte ProTracker format,
// parsing MilkyTracker-style hex effect strings.
// Instrument must be at least 1. 0 is reserved internally,
// but can be selected. (It just can't be set.)
// Pack cell data into exactly 4 bytes.
// Wait. Pack 5 bytes into 4? That doesn't work mathematically.
// 4 bytes is 32 bits. We have to throw away 8 bits somewhere.
// We're only using 12 bits of period. And we only need the lower
// 4 bits of cmd. Can lose the 4 MSB on both. Packing achieved.
// Byte 0: Instrument (upper 4 bits) | Period (upper 4 bits)
// Byte 1: Period (lower 8 bits)
// Byte 2: Instrument (lower 4 bits) | Effect Command (4 bits)
// Byte 3: Effect Parameter
func (c *Cell) Pack() ([4]byte, error) {
	var out [4]byte
	period := c.getPeriod()
	if period == 0 && !isEmpty(c.Note) {
		msgPair := [2]string{ERR_CELL_NOTE, ERR_STR}
		return out, utils.Err(msgPair, c.Note, nil)
	}
	if c.Instr > 31 {
		msgPair := [2]string{ERR_CELL_INST, ERR_STR}
		return out, utils.Err(msgPair, c.Instr, nil)
	}
	cmd, err := c.getEffectCommand()
	if err != nil {
		return out, err
	}
	param, err := c.getEffectParameter()
	if err != nil {
		return out, err
	}
	out[0] = (c.Instr & 0xF0) | uint8((period&0x0F00)>>8)
	out[1] = uint8(period & 0x00FF)
	out[2] = ((c.Instr & 0x0F) << 4) | (cmd & 0x0F)
	out[3] = param
	return out, nil
}

func (c *Cell) Write(w io.Writer) error {
	encodedBytes, err := c.Pack()
	if err != nil {
		return err
	}
	_, err = w.Write(encodedBytes[:])
	if err != nil {
		return err
	}
	return nil
}
