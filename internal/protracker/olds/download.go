package protracker

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"github.com/FatmanUK/fatgo/utils"
	"github.com/FatmanUK/fatgo/xlha"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const SAFETY_PAUSE_S = 2

const FMT_HTTP_STATUS = `HTTP Status Code: %d`

const ERR_DL_HEX_INVALID = `Invalid hex string.`
const ERR_DL_CHECKSUM = `Checksum mismatch.`

const LOG_DL_REQUEST = `Requesting sample archive download.`
const LOG_DL_STATUS_OK = `Download request succeeded.`
const LOG_DL_IN_MEMORY = `Downloaded in memory.`
const LOG_DL_CHECKSUM_OK = `Checksum test passed.`
const LOG_DL_DOWNLOAD = `Downloading from %s. Saving as %s.`
const LOG_DL_PAUSE = `Pausing for %d seconds.`

var archiveMap = map[string]SampleArchive{
	"st01": {
		Url: "https://aminet.net/mods/inst/st-01.lha",
		Checksum: SHA256SumFactory(
			[]byte("8bd8c62d542de794a5f843b5351de1c6"),
			[]byte("f1bed8f7d7643a24a04d0ad8ce962553"),
		),
		File: "st-01.lha",
	},
	"st02": {
		Url: "https://aminet.net/mods/inst/st-02.lha",
		Checksum: SHA256SumFactory(
			[]byte("3ffbbf30e652a65aa47d08afd976990f"),
			[]byte("67e2f9de5cbce171706d2c813d8dbce9"),
		),
		File: "st-02.lha",
	},
}

type SHA256Sum struct {
	BytesHex []byte
}

func SHA256SumFactory(blob1 []byte, blob2 []byte) SHA256Sum {
	return SHA256Sum{BytesHex: append(blob1, blob2...)}
}

func (re *SHA256Sum) IsChecksumMatch(blob []byte) (bool, error) {
	calcSum := sha256.Sum256(blob)
	expected, err := hex.DecodeString(string(re.BytesHex))
	if err != nil {
		//return false, fmt.Errorf(ERR_DL_HEX_INVALID, err)
		return false, fmt.Errorf(ERR_DL_HEX_INVALID)
	}
	rv := (subtle.ConstantTimeCompare(calcSum[:], expected) == 1)
	return rv, nil
}

type SampleArchive struct {
	Url      string
	Checksum SHA256Sum
	File     string
	Data     []byte
	logs     chan string
}

func (re *SampleArchive) BlobFromURI() error {
	re.logs <- LOG_DL_REQUEST
	resp, err := http.Get(re.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	sc := resp.StatusCode
	if sc != 200 {
		return fmt.Errorf(FMT_HTTP_STATUS, sc)
	}
	re.logs <- LOG_DL_STATUS_OK
	blob, err := re.Validate(resp.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(re.File, blob, 0640)
	if err != nil {
		return err
	}
	return nil
}

// Has delay to prevent abuse/damage to Aminet. Also checks checksums.
func (re *SampleArchive) Download() error {
	exists, err := utils.IsFileExists(re.File, true)
	if err != nil {
		return err
	}
	if exists {
		file, err := os.Open(re.File)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = re.Validate(file)
		return err
	}
	re.logs <- fmt.Sprintf(LOG_DL_DOWNLOAD, re.Url, re.File)
	err = re.BlobFromURI()
	if err != nil {
		return err
	}
	re.logs <- fmt.Sprintf(LOG_DL_PAUSE, SAFETY_PAUSE_S)
	time.Sleep(SAFETY_PAUSE_S * time.Second)
	return nil
}

func (re *SampleArchive) Extract(n string) error {
	re.Data = []byte{0}
	cache, err := utils.GetUserAppCacheDir("ptgen")
	if err != nil {
		return err
	}
	file, err := os.Open(filepath.Join(cache, re.File))
	if err != nil {
		return fmt.Errorf(ERR_INST_OPEN, err)
	}
	defer file.Close()
	return xlha.ExtractFile(file, n, &(re.Data))
}

func (re *SampleArchive) SetLogs(logs chan string) {
	re.logs = logs
}

func (re *SampleArchive) Validate(r io.ReadCloser) ([]byte, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return blob, err
	}
	re.logs <- LOG_DL_IN_MEMORY
	isChecksumOK, err := re.Checksum.IsChecksumMatch(blob)
	if err != nil {
		return blob, err
	}
	if !isChecksumOK {
		return blob, fmt.Errorf(ERR_DL_CHECKSUM)
	}
	re.logs <- LOG_DL_CHECKSUM_OK
	return blob, nil
}
