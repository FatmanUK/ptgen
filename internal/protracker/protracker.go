package protracker

import (
	"fmt"
	"io"
)

// Mostly written by Google Gemini Pro. Tweaked extensively by me.

// func WriteMod(w io.Writer, proj *ModProject) error
// WriteMod compiles proj into a binary ProTracker file stream.
// Populate proj, run WriteMod and the file appears in w.

const ERR_CELL_EFFT_SPEED = `Invalid initial speed. Maximum 31.`
const ERR_CELL_EFFT_BPM = `Invalid initial BPM. Minimum 32.`

// Magic bytes.
const MAGIC_BYTES = `M.K.`

func prepareCommands(s uint8, b uint8) ([]string, error) {
	var commands []string
	if s > 0 {
		if s >= 32 {
//			err := fmt.Errorf(ERR_CELL_EFFT_SPEED, s)
			err := fmt.Errorf(ERR_CELL_EFFT_SPEED)
			return commands, err
		}
		commands = append(commands, fmt.Sprintf("F%02X", s))
	}
	if b > 0 {
		if b < 32 {
//			err := fmt.Errorf(ERR_CELL_EFFT_BPM, b)
			err := fmt.Errorf(ERR_CELL_EFFT_BPM)
			return commands, err
		}
		commands = append(commands, fmt.Sprintf("F%02X", b))
	}
	return commands, nil
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

// 2. Write 31 Sample Headers (30 bytes each)
func writeSampleHeaders(w io.Writer, slots [31]Instrument) error {
	for _, i := range slots {
		err := i.WriteHeader(w)
		if err != nil {
			return err
		}
		//if i.ID > 0 {
		//	m := fmt.Sprintf("%d: %s", _, i.Save())
		//	log.Println(m)
		//}
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

// 5. Write Patterns: O(n^3) :(
// <129 patterns. =64 rows. =4 channels. Should never be too bad.
func writePatterns(w io.Writer, pttns []Pattern) error {
	for pid, p := range pttns {
		err := p.Write(w, uint8(pid))
		if err != nil {
			return err
		}
	}
	return nil
}

// 6. Write Raw PCM Sample Data block
func writeSamples(w io.Writer, slots [31]Instrument) error {
	for _, i := range slots {
		err := i.Write(w)
		if err != nil {
			return err
		}
	}
	return nil
}
