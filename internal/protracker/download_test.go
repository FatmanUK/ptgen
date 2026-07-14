package protracker

import (
	"testing"
)

func TestSHA256SumFactory_Must_Succeed(t *testing.T) {
	//24:func SHA256SumFactory(blob []byte) SHA256Sum {
}

func TestSHA256SumFactory_Must_Fail(t *testing.T) {
	//24:func SHA256SumFactory(blob []byte) SHA256Sum {
}

func TestSHA256SumIsChecksumMatch_Must_Succeed(t *testing.T) {
	//28:func (re *SHA256Sum) IsChecksumMatch(blob *[]byte) (bool, error) {
}

func TestSHA256SumIsChecksumMatch_Must_Fail(t *testing.T) {
	//28:func (re *SHA256Sum) IsChecksumMatch(blob *[]byte) (bool, error) {
}

func TestDownloadFactory_Must_Succeed(t *testing.T) {
	//47:func DownloadFileFactory(url string, e SHA256Sum, cache string) DownloadFile {
}

func TestDownloadFactory_Must_Fail(t *testing.T) {
	//47:func DownloadFileFactory(url string, e SHA256Sum, cache string) DownloadFile {
}

func TestDownloadIsChecksumMatch_Must_Succeed(t *testing.T) {
	//57:func (re *DownloadFile) IsChecksumMatch(blob *[]byte) (bool, error) {
}

func TestDownloadIsChecksumMatch_Must_Fail(t *testing.T) {
	//57:func (re *DownloadFile) IsChecksumMatch(blob *[]byte) (bool, error) {
}

func TestDownloadValidate_Must_Succeed(t *testing.T) {
	//61:func (re *DownloadFile) Validate(r io.ReadCloser,
}

func TestDownloadValidate_Must_Fail(t *testing.T) {
	//61:func (re *DownloadFile) Validate(r io.ReadCloser,
}

func TestDownloadBlobFromURI_Must_Succeed(t *testing.T) {
	//79:func (re *DownloadFile) BlobFromURI(logs chan string) error {
}

func TestDownloadBlobFromURI_Must_Fail(t *testing.T) {
	//79:func (re *DownloadFile) BlobFromURI(logs chan string) error {
}

func TestDecideSources_Must_Succeed(t *testing.T) {
	// func DecideSources(samples []Instrument) (map[string]DownloadFile, error) {
}

func TestDecideSources_Must_Fail(t *testing.T) {
	// func DecideSources(samples []Instrument) (map[string]DownloadFile, error) {
}

func TestDownload_Must_Succeed(t *testing.T) {
	// func Download(downloads map[string]DownloadFile, logs chan string) error {
}

func TestDownload_Must_Fail(t *testing.T) {
	// func Download(downloads map[string]DownloadFile, logs chan string) error {
}
