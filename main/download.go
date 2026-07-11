package main

import (
	"os"
	"io"
	"fmt"
	"net/http"
	"encoding/hex"
	"path/filepath"
	"crypto/sha256"
	"crypto/subtle"
)

type FileSource struct {
	Ref string
	URL string
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

type Download struct {
	Remote string
	Local string
	ExpectedSum SHA256Sum
}

func DownloadFactory(url string, e SHA256Sum, cache string) Download {
	cache = filepath.Join(cache, "ptgen")
	d := filepath.Join(cache, filepath.Base(url))
	return Download{
		Remote: url,
		Local: d,
		ExpectedSum: e,
	}
}

func (re *Download) IsChecksumMatch(blob *[]byte) (bool, error) {
	return re.ExpectedSum.IsChecksumMatch(blob)
}

func (re *Download) Validate(r io.ReadCloser,
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

func (re *Download) BlobFromURI(logs chan string) error {
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
