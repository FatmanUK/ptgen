package protracker

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/FatmanUK/fatgo/utils"
	"github.com/FatmanUK/fatgo/xlha"
)

const TMP_WRKRND_CMDLINE = `/usr/bin/lha evifw=%s %s %s`

const ERR_INST_TOO_BIG = `sample '%s' exceeds max length of 131070 bytes`
const ERR_INST_ODD = `sample '%s' must have an even byte length (got %d)`
const ERR_INST_LOOP = `sample '%s' loop points exceed total sample length`
const ERR_INST_OPEN = `Failed to open archive: %v`

// Start, Length and End do not refer to the Instrument but rather the
// forward loop (if present). Obviously, Start < End,
// Start + Length = End, and End < file size.
// Max file size = 65535*2 = 131070, needs a uint32 to store them.
// Start/Length = 0 for no loop.
// ID: The target slot (1-31)
// Volume: 0-64
// Finetune: 0-15 (0=0, 1=1... 8=-8... 15=-1)
// Start: Offset in BYTES
// Length: Length in BYTES (0 means no loop)
// Data: Raw 8-bit SIGNED PCM
type Instrument struct {
	ID       uint8  `json:"id"`
	Source   string `json:"source"`
	Name     string `json:"name"`
	Volume   uint8  `json:"volume"`
	Finetune uint8  `json:"finetune"`
	Start    uint32 `json:"start"`
	Length   uint32 `json:"length"`
	End      uint32 `json:"end"`
	Data     []byte `json:"data"`
}

func InstrumentFactory() Instrument {
	return Instrument{
		ID:       0,
		Source:   "",
		Name:     "",
		Volume:   64,
		Finetune: 0,
		Start:    0,
		Length:   0,
		End:      0,
		Data:     []byte{0},
	}
}

func (i *Instrument) CalculateLength() {
	if i.Start == 0 {
		i.Length = 0
		i.End = 0
	}
	if i.Start > 0 && i.Length == 0 {
		i.Length = (i.End - i.Start)
	}
	if i.Start > 0 && i.End == 0 {
		i.End = (i.Length + i.Start)
	}
}

// 4. Loop Points (Stored in words)
// If loop length is 2 bytes or less, treat it as unlooped
func (i *Instrument) encodeInstrumentLoop(out *[30]byte) error {
	if i.Length <= 2 {
		// Start = 0, Length = 1 word (standard for no loop)
		binary.BigEndian.PutUint16((*out)[26:28], 0)
		binary.BigEndian.PutUint16((*out)[28:30], 1)
	} else {
		if i.Start+i.Length > uint32(len(i.Data)) {
			m := fmt.Errorf(ERR_INST_LOOP, i.Name)
			return m
		}
		s := i.Start / 2
		l := i.Length / 2
		binary.BigEndian.PutUint16((*out)[26:28], uint16(s))
		binary.BigEndian.PutUint16((*out)[28:30], uint16(l))
	}
	return nil
}

// Validates and packs a 30-byte ProTracker sample header.
func (i *Instrument) encodeInstrumentHeader() ([30]byte, error) {
	var out [30]byte
	// 1. Name (Padded/truncated to exactly 22 bytes)
	copy(out[0:22], i.Name)
	dataLen := len(i.Data)
	if dataLen > 131070 {
		return out, fmt.Errorf(ERR_INST_TOO_BIG, i.Name)
	}
	if dataLen%2 != 0 {
		return out, fmt.Errorf(ERR_INST_ODD, i.Name, dataLen)
	}
	// 2. Length (Stored in words)
	binary.BigEndian.PutUint16(out[22:24], uint16(dataLen/2))
	// 3. Finetune & Volume
	out[24] = i.Finetune & 0x0F
	vol := i.Volume
	if vol > 64 {
		vol = 64
	}
	out[25] = vol
	return out, i.encodeInstrumentLoop(&out)
}

func (i *Instrument) isSlotPopulated() bool {
	return (i.ID > 0)
}

func (i *Instrument) temporaryWorkaroundWhileXlhaBroken(arch SampleArchive, data *[]byte) error {
	var err error
	cache := filepath.Dir(arch.File)
	wholeCmd := fmt.Sprintf(TMP_WRKRND_CMDLINE, cache,
		arch.File, i.Name)
	cmdArray := strings.Split(wholeCmd, " ")
	// assumes lhasa lha command installed
	cmd := exec.Command(cmdArray[0], cmdArray[1:]...)
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	sampleName := filepath.Base(i.Name)
	sampleFile, err := os.Open(filepath.Join(cache, sampleName))
	if err != nil {
		return err
	}
	defer sampleFile.Close()
	*data, err = io.ReadAll(sampleFile)
	return err
}

func (i *Instrument) lhaLoop(lhaReader *xlha.Reader,
		data *[]byte) (bool, error) {
	h, err := lhaReader.Next()
	if err == io.EOF {
		return true, err
	}
	if err != nil {
		//return fmt.Errorf(ERR_ARCH_HEADER_PARSE, err)
		return false, err
	}
	//log.Println(fmt.Sprintf(MSG_ARCH_HEADER_PARSE_OK, h.Name, h.Method, h.OriginalSize))
	*data, err = io.ReadAll(lhaReader)
	if err != nil {
		//return fmt.Errorf(ERR_ARCH_EXTRACTION, h.Name, err)
		return false, err
	}
	written := uint32(len(*data))
	if written != h.OriginalSize {
		return false, fmt.Errorf(ERR_ARCH_SIZE_MISMATCH,
			h.Name, h.OriginalSize, written)
	}
	if i.Name == h.Name { // found our file
		return true, nil
	}
	return false, nil
}

func (i *Instrument) lhaFile(arch SampleArchive, data *[]byte) error {
	file, err := os.Open(arch.File)
	if err != nil {
		return fmt.Errorf(ERR_INST_OPEN, err)
	}
	defer file.Close()
	lhaReader := xlha.NewReader(file)
	for isDone := false; !isDone; {
		isDone, err = i.lhaLoop(lhaReader, data)
	}
	return err
}

func (i *Instrument) ExtractSample() ([]byte, error) {
	data := []byte{0}
	arch := archiveMap[i.Source]
	cache, err := utils.GetUserAppCacheDir("ptgen")
	if err != nil {
		return data, err
	}
	arch.File = filepath.Join(cache, arch.File)
	//*
	err = i.temporaryWorkaroundWhileXlhaBroken(arch, &data)
	/*/
	err = i.lhaFile(arch, &data)
	//*/
	return data, err
}

func (i *Instrument) WriteHeader(w io.Writer) error {
	var err error
	if i.isSlotPopulated() {
		headerBytes, err := i.encodeInstrumentHeader()
		if err != nil {
			return err
		}
		_, err = w.Write(headerBytes[:])
	} else {
		emptySample := make([]byte, 30)
		emptySample[29] = 0x01 // Repeat length = 1 word
		_, err = w.Write(emptySample)
	}
	return err
}

func (i *Instrument) Write(w io.Writer) error {
	if i.isSlotPopulated() && len(i.Data) > 0 {
		_, err := w.Write(i.Data)
		if err != nil {
			return err
		}
	}
	return nil
}
