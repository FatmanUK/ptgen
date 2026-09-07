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
)

const DEFAULT_ROWS_PER_PATTERN = 64
const DEFAULT_CHANNELS_PER_ROW = 4

const DEFAULT_TITLE = `Nameless Song`
const DEFAULT_MESSAGE = ``
const DEFAULT_SPEED = 6
const DEFAULT_BPM = 125

const mdFmt = `Mod metadata:
Title:            {{ .Title }}
Message:          {{ .Msg }}
Initial Speed:    {{ .Speed }}
Initial BPM:      {{ .BPM }}
Rows per pattern: {{ .RowsPerPattern }}
Channels per row: {{ .ChannelsPerRow }}`

const RGX_PTTN_FILE = `^pattern[0-9A-Fa-f]{2}\.(?:txt|md)$`

// vars in mdFmt and yaml fields must be exported
type Mod struct {
	RowsPerPattern uint8        `yaml:"rowsPerPattern"`
	ChannelsPerRow uint8        `yaml:"channelsPerRow"`
	logs           chan string
	Title          string
	Msg            string       `yaml:"message"`
	Speed          uint8
	BPM            uint8
	Order          []uint8
	Patterns       []Pattern
	Instruments    []Instrument `yaml:"instruments"`
}

func ModFactory(l chan string) Mod {
	return Mod{
		RowsPerPattern: DEFAULT_ROWS_PER_PATTERN,
		ChannelsPerRow: DEFAULT_CHANNELS_PER_ROW,
		logs: l,
		Title: DEFAULT_TITLE,
		Msg: DEFAULT_MESSAGE,
		Speed: DEFAULT_SPEED,
		BPM: DEFAULT_BPM,
		Order: []uint8{},
		Patterns: []Pattern{},
		Instruments: []Instrument{},
	}
}

func (m Mod) Init(logs chan string, f string, p string) (Mod, error) {
	m, err := ModFactory(logs).loadMetadata(f)
	if err != nil {
		return m, err
	}
	m.logs <- string(utils.MustPrepTemplate("out", mdFmt, m))
	return m.enumeratePatterns(p)
}

func (m Mod) loadMetadata(metaFile string) (Mod, error) {
	if metaFile != "<nil>" {
		//m.logs <- LOG_LOADING_METADATA
		file, err := os.Open(metaFile)
		if err != nil {
			return m, err
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			return m, err
		}
		if _, err = file.Seek(0, io.SeekStart); err != nil {
			return m, err
		}
		if m, err = m.decodeMeta(file, content); err != nil {
			return m, err
		}
	}
	return m, m.validateMetadata()
}

func (m Mod) decodeMeta(file io.Reader, content []byte) (Mod, error) {
	if json.Valid(content) {
		//m.logs <- LOG_JSON_DETECTED
		return m, json.NewDecoder(file).Decode(&m)
	} else {
		var node yaml.Node
		if err := yaml.Unmarshal(content, &node); err != nil {
			return m, err
		}
		//m.logs <- LOG_YAML_DETECTED
		return m, yaml.NewDecoder(file).Decode(&m)
	}
	return m, nil
}

func (m Mod) validateMetadata() error {
	// look and see what original branch does
	//if len(m.Title) > 20 { // title less than 21 bytes
	//	return fmt.Errorf(ERR_MOD_TITLE_LONG)
	//}
	m.logs <- fmt.Sprintf("Instr: %d", len(m.Instruments))
	m.logs <- fmt.Sprintf("Patts: %d", len(m.Patterns))
	m.logs <- fmt.Sprintf("Order length: %d", len(m.Order))
	return nil
}

func (m Mod) enumeratePatterns(patternPath string) (Mod, error) {
	if patternPath != "<nil>" {
		//m.logs <- LOG_LOADING_PTTNS
		numPttns, err := utils.CountMatchingFiles(
			patternPath, RGX_PTTN_FILE)
		if err != nil {
			return m, err
		}
		m.logs <- fmt.Sprintf("Found %d patterns.", numPttns)
		m.Patterns = make([]Pattern, numPttns)
		for k := range m.Patterns {
			m.Patterns[k], err = m.Patterns[k].Init(
				m.logs,
				m.RowsPerPattern,
				m.ChannelsPerRow,
				patternPath,
				k,
			)
			if err != nil {
				return m, err
			}
		}
	}
	return m, m.validatePatterns()
}

func (m Mod) validatePatterns() error {
	// look and see what original branch does
	return nil
}

func (m Mod) Output() error {
	cache, err := utils.GetUserAppCacheDir("ptgen")
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err = m.downloadArchives(cache); err != nil {
		return err
	}
	if err = m.write(&buf); err != nil {
		return err
	}
/*
	errMsg := fmt.Sprintf(LOG_WRN_TOO_BIG_KB, MOD_TOO_BIG_KB)
	if buf.Len() > MOD_TOO_BIG_BYTES {
		p.logs <- errMsg
	}
*/
	return binary.Write(os.Stdout, binary.BigEndian, buf.Bytes())
}

func (m Mod) downloadArchives(cache string) error {
	for _, inst := range m.Instruments {
		_, err := Archive{}.Init(m.logs, cache, inst)
		if err != nil {
			return err
		}
	}
	return nil
}

// calculate loop points on the fly
func (m Mod) write(b *bytes.Buffer) error {
/*
	slots, err := p.preProcessInstruments()
	if err != nil {
		return slots, err
	}
	//err = p.injectInitialTempo()
	//if err != nil {
	//	return slots, fmt.Errorf(ERR_MOD_INJECT, err)
	//}
	//if !p.isOrderListValid() {
	//	return slots, fmt.Errorf(ERR_MOD_LIST)
	//}
	//if len(p.Instruments) > 31 {
	//	return slots, fmt.Errorf(ERR_INST_TOO_MANY)
	//}
	err = writeTitle(w, p.Title)
	if err != nil {
		return slots, err
	}
	err = writeSampleHeaders(w, slots)
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

*/
	return nil
}
