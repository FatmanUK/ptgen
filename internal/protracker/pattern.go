package protracker

import (
	"errors"
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
func (re Pattern) EmplaceRow(r uint16, m []string) (Pattern, error) {
	var err error
	re[r][0], err = CellFromString(m[2])
	if err != nil {
		return re, err
	}
	re[r][1], err = CellFromString(m[3])
	if err != nil {
		return re, err
	}
	re[r][2], err = CellFromString(m[4])
	if err != nil {
		return re, err
	}
	re[r][3], err = CellFromString(m[5])
	if err != nil {
		return re, err
	}
	return re, nil
}

func IsTooManyPatterns(numPatterns uint8) error {
	if numPatterns > MAX_PATTERNS {
		return errors.New("too many patterns")
	}
	return nil
}
