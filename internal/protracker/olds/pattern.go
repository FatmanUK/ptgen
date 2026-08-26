package protracker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
//	"strings"
)

const LOG_TOTAL_ROWS_PROCESSED = `Total rows processed: %d`

type ScanData struct {
	scanner   *bufio.Scanner
	regex     *regexp.Regexp
	isHexRows bool
}

type Pattern struct {
	Rows []Row
	logs chan string
}

func PatternFactory(p uint8, c uint8, logs chan string) Pattern {
	pttn := Pattern{
		Rows: make([]Row, p),
		logs: logs,
	}
	for i := range pttn.Rows {
		pttn.Rows[i] = RowFactory(c)
	}
	return pttn
}

// hmm. different semantic to Cell.Load. Think on it.
func (p *Pattern) Load(fileName string, isHexRows bool, c uint8) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	data := ScanData{
		scanner:   bufio.NewScanner(file),
		regex:     regexp.MustCompile(RowRegexFactory(c)),
		isHexRows: isHexRows,
	}
	rowsRead, err := p.ScanLoop(data)
	if err != nil {
		return err
	}
	p.logs <- fmt.Sprintf(LOG_TOTAL_ROWS_PROCESSED, rowsRead)
	return data.scanner.Err()
}

func (p *Pattern) InjectCommands(cmmds []string) uint8 {
	return 0
}

func (p *Pattern) Write(w io.Writer, pid uint8) error {
	return nil
}

func (p *Pattern) ScanLoop(d ScanData) (uint8, error) {
	var rowsRead uint8
	for d.scanner.Scan() {
		m := d.regex.FindStringSubmatch(d.scanner.Text())
		if len(m) == 0 {
			continue
		}
		rowsRead++
		row, err := RowNumFromRowStr(m[1], d.isHexRows)
		if err != nil {
			return rowsRead, err
		}
		err = p.EmplaceRow(row, m[2:])
		if err != nil {
			return rowsRead, err
		}
		r := p.Rows[row].Save(row, uint8(len(m[2:])))
		p.logs <- fmt.Sprintf("%s", r)
		if err != nil {
			return rowsRead, err
		}
	}
	return rowsRead, nil
}

// Emplace a Row in a Pattern from an array of pre-Cell strings
func (p *Pattern) EmplaceRow(rowNum uint8, m []string) error {
	var err error
	for i := 0; i < len(m); i++ {
		p.Rows[rowNum].Cells[i] = CellFactory()
		err = p.Rows[rowNum].Cells[i].Load(m[i])
		if err != nil {
			break
		}
	}
	return err
}

/*
// Scan the 4 channels of Row 0 for empty effect fields
// Standard tracker empty signals are "", "000", or "---"
// Add speed and BPM commands there, if specified
// Return: number of commands injected
func (p *Pattern) InjectCommands(cmmds []string) uint8 {
	rowZero := &p[0]
	var cmdIdx uint8 = 0
	for ch := 0; ch < 4 && cmdIdx < uint8(len(cmmds)); ch++ {
		eff := strings.TrimSpace(rowZero[ch].Effect)
		if isEmpty(eff) {
			rowZero[ch].Effect = cmmds[cmdIdx]
			cmdIdx++
		}
	}
	return cmdIdx
}

func (p *Pattern) Write(w io.Writer, pid uint8) error {
	for rid, r := range *p {
		err := r.Write(w, pid, uint8(rid))
		if err != nil {
			return err
		}
	}
	return nil
}
*/
