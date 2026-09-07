package protracker

import (
)

type Pattern struct {
}
// maybe just have pattern this time, don't do row and cell?

func PatternFactory(l chan string, r uint8, c uint8) Pattern {
	//v = PatternFactory(m.logs, m.RowsPerPattern, m.ChannelsPerRow)
	return Pattern{
	}
}

func (p Pattern) Init(logs chan string, r uint8, c uint8, path string,
	n int) (Pattern, error) {
	return PatternFactory(logs, r, c).loadFile(path, n)
}

// get hex on the fly
func (p Pattern) loadFile(rootPath string, idx int) (Pattern, error) {
	return p, nil
}

//hex, err := IsHexDetected(path)
//if err != nil {
//	return err
//}
//m.logs <- fmt.Sprintf("Is hex rows: %v", hex)
//func (p Pattern) loadFile(i uint8, rootPath string) {
/*
	// load txt or md file
	fileName, err := findFile(rootPath, i)
	if err != nil {
		return err
	}
	p.logs <- fmt.Sprintf(LOG_LOADING_FILE, fileName)
	fullPath := filepath.Join(path, fileName)
	err = p.Patterns[i].Load(fullPath, isHex, p.Channels)
*/
//}
