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
)

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
		return false, fmt.Errorf(ERR_HEX_INVALID, err)
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

func DownloadFileFactory(url string, e SHA256Sum, cache string) DownloadFile {
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
	logs <- MSG_DOWNLOAD_IN_MEM
	isChecksumOK, err := re.IsChecksumMatch(&blob)
	if err != nil {
		return blob, err
	}
	if !isChecksumOK {
		return blob, fmt.Errorf(ERR_CHECKSUM_MISMATCH)
	}
	logs <- MSG_CHECKSUM_OK
	return blob, nil
}

func (re *DownloadFile) BlobFromURI(logs chan string) error {
	logs <- MSG_REQUEST_DOWNLOAD
	resp, err := http.Get(re.Remote)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	sc := resp.StatusCode
	if sc != 200 {
		return fmt.Errorf(FMT_HTTP_STATUS, sc)
	}
	logs <- MSG_DOWNLOAD_STATUS_OK
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

// Find the user cache dir. Create source mappings
// ("ref"=>("https://url","local/file")). Decide on mappings depending
// which sample refs are selected.
func DecideSources(samples []Instrument) (map[string]DownloadFile, error) {
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
	for _, s := range samples {
		chosen[s.Source] = downloads[s.Source]
	}
	return chosen, nil
}

// Has delay to prevent abuse/damage to Aminet. Also checks checksums.
func Download(downloads map[string]DownloadFile, logs chan string) error {
	for _, v := range downloads {
		exists, err := CheckFileExistsWithMkdir(v.Local)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		logs <- fmt.Sprintf(MSG_DOWNLOAD, v.Remote, v.Local)
		err = v.BlobFromURI(logs)
		if err != nil {
			return err
		}
		logs <- fmt.Sprintf(MSG_SAFETY_PAUSE, SAFETY_PAUSE_S)
		time.Sleep(SAFETY_PAUSE_S * time.Second)
	}
	return nil
}
