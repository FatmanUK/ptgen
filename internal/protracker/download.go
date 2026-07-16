package protracker

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"github.com/FatmanUK/fatgo/utils"
)

const FMT_HTTP_STATUS = `HTTP Status Code: %d`

const ERR_DL_HEX_INVALID = `Invalid hex string. %v`
const ERR_DL_CHECKSUM = `Checksum mismatch`

const LOG_DL_IN_MEMORY = `Downloaded in memory.`
const LOG_DL_CHECKSUM_OK = `Checksum test passed.`
const LOG_DL_REQUEST = `Requesting sample archive download.`
const LOG_DL_STATUS_OK = `Download request succeeded.`
const LOG_DL_DOWNLOAD = `Downloading from %s. Saving as %s.`
const LOG_DL_PAUSE = `Pausing for %d seconds.`

// Arbitrary constants.
const URL_ST01 = `https://aminet.net/mods/inst/st-01.lha`
const URL_ST02 = `https://aminet.net/mods/inst/st-02.lha`

const SHA256_ST01 = `8bd8c62d542de794a5f843b5351de1c6f1bed8f7d7643a24a04d0ad8ce962553`
const SHA256_ST02 = `3ffbbf30e652a65aa47d08afd976990f67e2f9de5cbce171706d2c813d8dbce9`

const SAFETY_PAUSE_S = 2

type FileSource struct {
	Ref         string
	URL         string
	ExpectedSum SHA256Sum
}

type SHA256Sum struct {
	BytesHex []byte
}

func SHA256SumFactory(blob []byte) SHA256Sum {
	return SHA256Sum{BytesHex: blob}
}

func (re *SHA256Sum) IsChecksumMatch(blob *[]byte) (bool, error) {
	actualSum := sha256.Sum256(*blob)
	expectedBytes, err := hex.DecodeString(string(re.BytesHex))
	if err != nil {
		return false, fmt.Errorf(ERR_DL_HEX_INVALID, err)
	}
	if subtle.ConstantTimeCompare(actualSum[:],
		expectedBytes) == 1 {
		return true, nil
	}
	return false, nil
}

type DownloadFile struct {
	Remote      string
	Local       string
	ExpectedSum SHA256Sum
}

func DownloadFileFactory(url string, e SHA256Sum,
		cache string) DownloadFile {
	cache = filepath.Join(cache, "ptgen")
	d := filepath.Join(cache, filepath.Base(url))
	return DownloadFile{
		Remote:      url,
		Local:       d,
		ExpectedSum: e,
	}
}

func (re *DownloadFile) IsChecksumMatch(blob *[]byte) (bool, error) {
	return re.ExpectedSum.IsChecksumMatch(blob)
}

func (re *DownloadFile) Validate(r io.ReadCloser,
		logs chan string) ([]byte, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return blob, err
	}
	logs <- LOG_DL_IN_MEMORY
	isChecksumOK, err := re.IsChecksumMatch(&blob)
	if err != nil {
		return blob, err
	}
	if !isChecksumOK {
		return blob, fmt.Errorf(ERR_DL_CHECKSUM)
	}
	logs <- LOG_DL_CHECKSUM_OK
	return blob, nil
}

func (re *DownloadFile) BlobFromURI(logs chan string) error {
	logs <- LOG_DL_REQUEST
	resp, err := http.Get(re.Remote)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	sc := resp.StatusCode
	if sc != 200 {
		return fmt.Errorf(FMT_HTTP_STATUS, sc)
	}
	logs <- LOG_DL_STATUS_OK
	blob, err := re.Validate(resp.Body, logs)
	if err != nil {
		return err
	}
	err = os.WriteFile(re.Local, blob, 0640)
	if err != nil {
		return err
	}
	return nil
}

// TODO: is this a bit kak-handed? Try and improve.
// Find the user cache dir. Create source mappings
// ("ref"=>("https://url","local/file")). Decide on mappings depending
// which sample refs are selected.
func DecideSources(s []Instrument) (map[string]DownloadFile, error) {
	st01 := []byte(SHA256_ST01)
	st02 := []byte(SHA256_ST02)
	sources := []FileSource{
		{"st01", URL_ST01, SHA256SumFactory(st01)},
		{"st02", URL_ST02, SHA256SumFactory(st02)},
	}
	downloads := map[string]DownloadFile{}
	chosen := map[string]DownloadFile{}
	cache, err := os.UserCacheDir()
	if err != nil {
		return chosen, err
	}
	for _, m := range sources {
		downloads[m.Ref] = DownloadFileFactory(m.URL,
			m.ExpectedSum, cache)
	}
	for _, sample := range s {
		chosen[sample.Source] = downloads[sample.Source]
	}
	return chosen, nil
}

// Has delay to prevent abuse/damage to Aminet. Also checks checksums.
func Download(downloads map[string]DownloadFile,
		logs chan string) error {
	for _, v := range downloads {
		exists, err := utils.IsFileExists(v.Local, true)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		logs <- fmt.Sprintf(LOG_DL_DOWNLOAD, v.Remote, v.Local)
		err = v.BlobFromURI(logs)
		if err != nil {
			return err
		}
		logs <- fmt.Sprintf(LOG_DL_PAUSE, SAFETY_PAUSE_S)
		time.Sleep(SAFETY_PAUSE_S * time.Second)
	}
	return nil
}
