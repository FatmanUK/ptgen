package protracker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// A Pattern is a slice of ROWS_PER_PATTERN Rows.
type Pattern [ROWS_PER_PATTERN]Row

func PatternFactory() Pattern {
	p := Pattern{}
	for i := range p {
		p[i] = RowFactory()
	}
	return p
}

// Emplace a Row in a Pattern from an array of pre-Cell strings
func (p *Pattern) EmplaceRow(r uint16, m [4]string) (Pattern, error) {
	var err error
	for i := 0; i < CHANNELS_PER_ROW; i++ {
		(*p)[r][i] = CellFactory()
		err = (*p)[r][i].Load(m[i])
		if err != nil {
			return *p, err
		}
	}
	return *p, nil
}

// hmm. different semantic to Cell.Load. Think on it.
func (p *Pattern) Load(logs chan string, fileName string,
	isHexRows bool) error {
	*p = PatternFactory()
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	return p.Read(file, logs, isHexRows)
}

// Scan the 4 channels of Row 0 for empty effect fields
// Standard tracker empty signals are "", "000", or "---"
// Add speed and BPM commands there, if specified
// Return: number of commands injected
func (p *Pattern) InjectCommands(cmmds []string) uint8 {
	rowZero := &p[0]
	var cmdIdx uint8 = 0
	for ch := 0; ch < 4 && cmdIdx < uint8(len(cmmds)); ch++ {
		eff := strings.TrimSpace(rowZero[ch].Effect)
		if eff == "" || eff == "000" || eff == "---" {
			rowZero[ch].Effect = cmmds[cmdIdx]
			cmdIdx++
		}
	}
	return cmdIdx
}

func (p *Pattern) Write(w io.Writer, pid int) error {
	for rid, r := range *p {
		err := r.Write(w, pid, rid)
		if err != nil {
			return err
		}
	}
	return nil
}

type ScanData struct {
	scanner   *bufio.Scanner
	regex     *regexp.Regexp
	isHexRows bool
}

func (p *Pattern) Read(file io.Reader, logs chan string,
	isHexRows bool) error {
	data := ScanData{
		scanner:   bufio.NewScanner(file),
		regex:     regexp.MustCompile(RowRegexFactory()),
		isHexRows: isHexRows,
	}
	rowsRead, err := p.ScanLoop(logs, data)
	if err != nil {
		return err
	}
	logs <- fmt.Sprintf(MSG_TOTAL_ROWS_PROCESSED, rowsRead)
	return data.scanner.Err()
}

func (p *Pattern) ScanLoop(logs chan string, d ScanData) (uint, error) {
	var rowsRead uint
	for d.scanner.Scan() {
		m := d.regex.FindStringSubmatch(d.scanner.Text())
		if len(m) == 0 {
			continue
		}
		rowsRead++
		n, err := RowNumFromRowStr(m[1], d.isHexRows)
		if err != nil {
			return rowsRead, err
		}
		m2 := [4]string{m[2], m[3], m[4], m[5]}
		*p, err = p.EmplaceRow(n, m2)
		if err != nil {
			return rowsRead, err
		}
		logs <- fmt.Sprintf("%s", p[n].Save(n))
		if err != nil {
			return rowsRead, err
		}
	}
	return rowsRead, nil
}
