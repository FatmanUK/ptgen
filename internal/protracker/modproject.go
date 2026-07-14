package protracker

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

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
}

func ModProjectFactory() ModProject {
	return ModProject{
		Title:       DEFAULT_TITLE,
		Speed:       DEFAULT_SPEED,
		BPM:         DEFAULT_BPM,
		OrderList:   []uint8{0},
		Patterns:    []Pattern{PatternFactory()},
		Instruments: []Instrument{},
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

func (p *ModProject) OutputEverything(logs chan string) error {
	var buf bytes.Buffer
	var err error
	info := p.ModInfoFactory()
	output := MustPrepTemplate("output", MetadataString, info)
	logs <- string(output)
	if len(p.Title) > 20 { // title less than 21 bytes
		return fmt.Errorf(ERR_MOD_TITLE_LONG)
	}
	downloads, err := DecideSources(p.Instruments)
	if err != nil {
		return err
	}
	err = Download(downloads, logs)
	if err != nil {
		return err
	}
	err = p.WriteMod(&buf)
	if err != nil {
		return err
	}
	if buf.Len() > MOD_TOO_BIG_BYTES {
		logs <- fmt.Sprintf(MSG_WRN_TOO_BIG_KB,
			buf.Len()/1024, MOD_TOO_BIG_KB)
	}
	err = binary.Write(os.Stdout, binary.BigEndian, buf.Bytes())
	return err
}

// No monolithic pattern files. Too annoying.
// Requires the patterns formatted in plain text.
func (p *ModProject) PopulatePatterns(logs chan string,
	path string) error {
	logs <- MSG_LOADING_PTTNS
	numPttns, err := CountPatterns(path)
	if err != nil {
		return err
	}
	hex, err := IsHexDetected(path)
	if err != nil {
		return err
	}
	p.Patterns = make([]Pattern, numPttns)
	for i := range p.Patterns {
		fileName := fmt.Sprintf(MSG_PTTN_FILE, i)
		logs <- fmt.Sprintf(MSG_LOADING_FILE, fileName)
		fullPath := filepath.Join(path, fileName)
		err = p.Patterns[i].Load(logs, fullPath, hex)
		if err != nil {
			return err
		}
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
	for _, j := range p.OrderList {
		if int(j) >= len(p.Patterns) {
			return false
		}
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
			m := fmt.Errorf(ERR_INSTRUMENT_INVALID_ID,
				i.Name, i.ID)
			return slots, m
		}
		if slots[i.ID].ID > 0 {
			m := fmt.Errorf(ERR_INSTRUMENT_DUPLICATE,
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
		return slots, fmt.Errorf(ERR_INJECT, err)
	}
	if !p.isOrderListValid() {
		return slots, fmt.Errorf(ERR_MOD_LIST)
	}
	if len(p.Instruments) > 31 {
		return slots, fmt.Errorf(ERR_INSTRUMENT_TOO_MANY)
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

func (p *ModProject) ReadMetadata(logs chan string,
	file *os.File) error {
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
		logs <- MSG_JSON_DETECTED
		err = json.NewDecoder(file).Decode(p)
	} else {
		var node yaml.Node
		err = yaml.Unmarshal(content, &node)
		if err == nil {
			logs <- MSG_YAML_DETECTED
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
func (p *ModProject) PopulateMetadata(logs chan string,
	fileName string) error {
	var err error
	p.OrderList = p.DefaultOrderlist()
	if fileName != "<nil>" {
		logs <- MSG_LOADING_METADATA
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		err = p.ReadMetadata(logs, file)
		if err != nil {
			return err
		}
		p.CalculateLength()
	}
	if !p.isOrderListValid() {
		err = fmt.Errorf(ERR_MOD_OLST)
	}
	return err
}
