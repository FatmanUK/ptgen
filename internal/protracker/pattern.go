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

func IsTooManyPatterns(numPatterns uint8) error {
	if numPatterns > MAX_PATTERNS {
		return errors.New(ERR_PTTN_TOO_MANY)
	}
	return nil
}

// Emplace a Row in a Pattern from an array of pre-Cell strings
func (re Pattern) EmplaceRow(r uint16, m [4]string) (Pattern, error) {
	var err error
	for i := 0; i < CHANNELS_PER_ROW; i++ {
		re[r][i] = CellFactory()
		err = re[r][i].Load(m[i])
		if err != nil {
			return re, err
		}
	}
	return re, nil
}
