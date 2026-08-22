package protracker

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/FatmanUK/fatgo/utils"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

const MetadataFmt = `Mod metadata:
Song Title:      {{ .Title }}
Format:          ProTracker-compatible MOD
Channels:        4
Speed:           {{ .Speed }}
BPM:             {{ .BPM }}
Pattern length:  64 rows (Len. 0x40h)
Unique patterns: {{ .Patterns }}
Order length:    {{ .OrderLen }}
Instruments:     {{ .Instruments }}`

// Limits.
const MAX_PATTERNS = 128 // or 64 ? Gemini seems undecided.
const MAX_ORDER_LIST = MAX_PATTERNS

// These are just defaults.
const DEFAULT_TITLE = `A Song With No Name`
const DEFAULT_SPEED = 6
const DEFAULT_BPM = 125

// Warnings.
const MOD_TOO_BIG_KB = 50
const MOD_TOO_BIG_BYTES = (MOD_TOO_BIG_KB * 1024)

const LOG_WRN_TOO_BIG = `Mod is unusually large (>%d bytes).`
const LOG_WRN_TOO_BIG_KB = `Mod is unusually large (>%d kB).`
const LOG_LOADING_FILE = `Loading %s.`
const LOG_LOADING_PTTNS = `Loading patterns.`
const LOG_JSON_DETECTED = `JSON metadata detected.`
const LOG_YAML_DETECTED = `YAML metadata detected.`
const LOG_LOADING_METADATA = `Loading metadata and pattern order.`

const FMT_PTTN_FILE = `pattern%02x.txt`
const FMT_PTTN_MD = `pattern%02x.md`

const RGX_PTTN_FILE = `^pattern[0-9A-Fa-f]{2}\.(?:txt|md)$`

const ERR_MOD_INJECT = `Tempo injection error: %w`
const ERR_MOD_ORDER_LIST = `Order list is wrong.`
const ERR_MOD_TITLE_LONG = `Title is too long.`
const ERR_MOD_LIST = `OrderList is wrong.`
const ERR_MOD_LIST_EMPTY = ERR_MOD_LIST + ` Must not be empty.`
const ERR_MOD_LIST_OOB = ERR_MOD_LIST + ` Invalid pattern reference %d.`
const ERR_MOD_LIST_OVERFLOW = ERR_MOD_LIST + ` Maximum 128 items.` // TODO: is this right? 128 is max patterns, also max orderlist?

const ERR_INST_TOO_MANY = `a maximum of 31 instruments are supported`
const ERR_INST_INVALID_ID = `instrument '%s' has invalid ID %d (must be 1-31)`
const ERR_INST_DUPE = `duplicate instrument ID %d detected`
const ERR_INST_SLOT = `instrument slot %d error: %w`

const ERR_PTTN_TOO_MANY = `Too many patterns.`
const ERR_PTTN_ZERO_FULL = `Injection failed. Row 0 of first pattern %d has insufficient empty effect slots`

type ModInfo struct {
	Title       string
	Speed       uint8
	BPM         uint8
	OrderLen    uint8
	Patterns    uint8
	Instruments uint8
}

// TODO: enforce these limits.
// ModProject represents your input structure.
// Speed: Optional: 1-31 (0: default, always 6)
// BPM: Optional: 32-255 (0: default, always 125)
// Patterns: a slice of Patterns
// Up to 31 instruments
type ModProject struct {
	Title       string       `json:"title"`
	Speed       uint8        `json:"speed"`
	BPM         uint8        `json:"bpm"`
	OrderList   []uint8      `json:"orderList" yaml:"orderList"`
	Patterns    []Pattern    `json:"patterns"`
	Instruments []Instrument `json:"instruments"`
	logs        chan string
}

func ModProjectFactory(logs chan string) ModProject {
	return ModProject{
		Title:       DEFAULT_TITLE,
		Speed:       DEFAULT_SPEED,
		BPM:         DEFAULT_BPM,
		OrderList:   []uint8{0},
		Patterns:    []Pattern{PatternFactory()},
		Instruments: []Instrument{},
		logs:        logs,
	}
}

func (p *ModProject) ModInfoFactory() ModInfo {
	return ModInfo{
		Title:       p.Title,
		Speed:       p.Speed,
		BPM:         p.BPM,
		OrderLen:    uint8(len(p.OrderList)),
		Patterns:    uint8(len(p.Patterns)),
		Instruments: uint8(len(p.Instruments)),
	}
}

func (p *ModProject) OutputEverything() error {
	var buf bytes.Buffer
	var err error
	info := p.ModInfoFactory()
	output := utils.MustPrepTemplate("output", MetadataFmt, info)
	p.logs <- string(output)
	if len(p.Title) > 20 { // title less than 21 bytes
		return fmt.Errorf(ERR_MOD_TITLE_LONG)
	}
	cache, err := utils.GetUserAppCacheDir("ptgen")
	if err != nil {
		return err
	}
	for _, i := range p.Instruments {
		arch := archiveMap[i.Source]
		arch.File = filepath.Join(cache, arch.File)
		err = download(arch, p.logs)
		if err != nil {
			return err
		}
	}
	err = p.WriteMod(&buf)
	if err != nil {
		return err
	}
	// bufKb := buf.Len()/1024
	errMsg := fmt.Sprintf(LOG_WRN_TOO_BIG_KB, MOD_TOO_BIG_KB)
	if buf.Len() > MOD_TOO_BIG_BYTES {
		p.logs <- errMsg
	}
	err = binary.Write(os.Stdout, binary.BigEndian, buf.Bytes())
	return err
}

func findFile(path string, i int) (string, error) {
	fileName := fmt.Sprintf(FMT_PTTN_FILE, i)
	fullPath := filepath.Join(path, fileName)
	ex, err := utils.IsFileExists(fullPath, false)
	if err != nil {
		fmt.Println("File .txt found but error")
		return fileName, err
	}
	if !ex {
		fileName = fmt.Sprintf(FMT_PTTN_MD, i)
		fullPath = filepath.Join(path, fileName)
		ex, err = utils.IsFileExists(fullPath, false)
		if err != nil {
			fmt.Println("File .md found but error")
			return fileName, err
		}
		if !ex {
			m := fmt.Errorf("Pattern %d not found.", i)
			return fileName, m
		}
	}
	return fileName, nil
}

func (p *ModProject) loadPttns(path string, isHex bool) error {
	for i := range p.Patterns {
		// load txt or md file
		fileName, err := findFile(path, i)
		if err != nil {
			return err
		}
		p.logs <- fmt.Sprintf(LOG_LOADING_FILE, fileName)
		fullPath := filepath.Join(path, fileName)
		err = p.Patterns[i].Load(p.logs, fullPath, isHex)
		if err != nil {
			return err
		}
	}
	return nil
}

// No monolithic pattern files. Too annoying.
// Requires the patterns formatted in plain text.
func (p *ModProject) PopulatePatterns(path string) error {
	p.logs <- LOG_LOADING_PTTNS
	numPttns, err := utils.CountMatchingFiles(path, RGX_PTTN_FILE)
	if err != nil {
		return err
	}
	p.logs <- fmt.Sprintf("Found %d patterns.", numPttns)
	hex, err := IsHexDetected(path)
	if err != nil {
		return err
	}
	p.logs <- fmt.Sprintf("Is hex rows: %v", hex)
	p.Patterns = make([]Pattern, numPttns)
	err = p.loadPttns(path, hex)
	if err != nil {
		return err
	}
	return p.isTooManyPatterns()
}

func (p *ModProject) isTooManyPatterns() error {
	if len(p.Patterns) > MAX_PATTERNS {
		return fmt.Errorf(ERR_PTTN_TOO_MANY)
	}
	return nil
}

func (p *ModProject) isOrderListValid() bool {
	l := len(p.OrderList)
	if l == 0 || l > MAX_PATTERNS {
		return false
	}
	orderMax := uint8(0)
	for _, j := range p.OrderList {
		if j > orderMax {
			orderMax = j
		}
		if int(j) >= len(p.Patterns) {
			return false
		}
	}
	if len(p.Patterns) != int(orderMax+1) {
		return false
	}
	return true
}

func (p *ModProject) getFirstPatternIndex() (uint8, error) {
	var err error
	if len(p.OrderList) == 0 {
		return 0, fmt.Errorf(ERR_MOD_LIST_EMPTY)
	}
	i := p.OrderList[0]
	if int(i) >= len(p.Patterns) {
		err = fmt.Errorf(ERR_MOD_LIST_OOB, i)
	}
	return i, err
}

// Preprocess metadata and inject F-commands into Row 0
// injectInitialTempo scans the first row of the starting pattern and
// attempts to inject Fxx speed/BPM commands into empty effect slots.
func (p *ModProject) injectInitialTempo() error {
	if p.Speed == 0 && p.BPM == 0 {
		return nil
	}
	firstIndex, err := p.getFirstPatternIndex()
	if err != nil {
		return err
	}
	commands, err := prepareCommands(p.Speed, p.BPM)
	if err != nil {
		return err
	}
	firstPattern := &(p.Patterns[firstIndex])
	cmdIdx := firstPattern.InjectCommands(commands)
	if cmdIdx < uint8(len(commands)) {
		return fmt.Errorf(ERR_PTTN_ZERO_FULL, firstIndex)
	}
	return nil
}

// "If any of those raw files happen to have an odd byte length (which
// occasionally happened with manual rips of those old Amiga disks),
// you can simply append a single 0x00 byte to bassData before
// assigning it to the struct to satisfy the strict word-length
// constraint we built into encodeInstrumentHeader."
//
// TODO: split
// 0. Pre-process sparse instruments into a rigid 31-slot lookup map
func (p *ModProject) preProcessInstruments() ([31]Instrument, error) {
	// TODO: hmm... do we need p.Instruments /and/ slots?
	var slots [31]Instrument
	for _, i := range p.Instruments {
		d, err := i.ExtractSample()
		if len(d)%2 == 1 { // see comment above
			d = append(d, 0x00)
		}
		i.Data = d
		if err != nil {
			return slots, err
		}
		if i.ID < 1 || i.ID > 31 {
			m := fmt.Errorf(ERR_INST_INVALID_ID,
				i.Name, i.ID)
			return slots, m
		}
		if slots[i.ID].ID > 0 {
			m := fmt.Errorf(ERR_INST_DUPE,
				i.ID)
			return slots, m
		}
		// TODO: test i.ID ranges here. One too high so -1?
		// Suspicious. Check instruments 1 and 31.
		slots[i.ID-1] = i
	}
	//for k := range slots {
	//	m := fmt.Sprintf("%d: %s", k, slots[k].Save())
	//	log.Println(m)
	//}
	return slots, nil
}

func (p *ModProject) writeHeaders(
	w io.Writer) ([31]Instrument, error) {
	slots, err := p.preProcessInstruments()
	if err != nil {
		return slots, err
	}
	err = p.injectInitialTempo()
	if err != nil {
		return slots, fmt.Errorf(ERR_MOD_INJECT, err)
	}
	if !p.isOrderListValid() {
		return slots, fmt.Errorf(ERR_MOD_LIST)
	}
	if len(p.Instruments) > 31 {
		return slots, fmt.Errorf(ERR_INST_TOO_MANY)
	}
	err = writeTitle(w, p.Title)
	if err != nil {
		return slots, err
	}
	return slots, writeSampleHeaders(w, slots)
}

// WriteMod compiles proj into a binary ProTracker file stream.
func (p *ModProject) WriteMod(w io.Writer) error {
	slots, err := p.writeHeaders(w)
	if err != nil {
		return err
	}
	err = writeOrderList(w, p.OrderList)
	if err != nil {
		return err
	}
	err = writeMagic(w)
	if err != nil {
		return err
	}
	err = writePatterns(w, p.Patterns)
	if err != nil {
		return err
	}
	return writeSamples(w, slots)
}

func (p *ModProject) ReadMetadata(file *os.File) error {
	var err error
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	if json.Valid(content) {
		p.logs <- LOG_JSON_DETECTED
		err = json.NewDecoder(file).Decode(p)
	} else {
		var node yaml.Node
		err = yaml.Unmarshal(content, &node)
		if err == nil {
			p.logs <- LOG_YAML_DETECTED
			err = yaml.NewDecoder(file).Decode(p)
		}
	}
	return err
}

func (p *ModProject) DefaultOrderlist() []uint8 {
	var orderList []uint8
	for n := uint8(0); n < uint8(len(p.Patterns)); n++ {
		orderList = append(orderList, n)
	}
	return orderList
}

func (p *ModProject) CalculateLength() {
	for k := range p.Instruments {
		p.Instruments[k].CalculateLength()
	}
}

// The song metadata in JSON or YAML format, including order list.
func (p *ModProject) PopulateMetadata(fileName string) error {
	var err error
	p.OrderList = p.DefaultOrderlist()
	if fileName != "<nil>" {
		p.logs <- LOG_LOADING_METADATA
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		err = p.ReadMetadata(file)
		if err != nil {
			return err
		}
		p.CalculateLength()
	}
	if !p.isOrderListValid() {
		err = fmt.Errorf(ERR_MOD_ORDER_LIST)
	}
	return err
}
