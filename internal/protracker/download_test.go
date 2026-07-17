package protracker

import (
	"testing"
)

func TestSHA256SumFactory_Must_Succeed(t *testing.T) {
//33:func SHA256SumFactory(blob1 []byte, blob2 []byte) SHA256Sum {
}

func TestSHA256SumFactory_Must_Fail(t *testing.T) {
//33:func SHA256SumFactory(blob1 []byte, blob2 []byte) SHA256Sum {
}

func TestSHA256SumIsChecksumMatch_Must_Succeed(t *testing.T) {
//37:func (re *SHA256Sum) IsChecksumMatch(blob []byte) (bool, error) {
}

func TestSHA256SumIsChecksumMatch_Must_Fail(t *testing.T) {
//37:func (re *SHA256Sum) IsChecksumMatch(blob []byte) (bool, error) {
}

func TestDownloadValidate_Must_Succeed(t *testing.T) {
//104:func (re *SampleArchive) Validate(r io.ReadCloser,
}

func TestDownloadValidate_Must_Fail(t *testing.T) {
//104:func (re *SampleArchive) Validate(r io.ReadCloser,
}

func TestDownloadBlobFromURI_Must_Succeed(t *testing.T) {
//56:func (re *SampleArchive) BlobFromURI(logs chan string) error {
}

func TestDownloadBlobFromURI_Must_Fail(t *testing.T) {
//56:func (re *SampleArchive) BlobFromURI(logs chan string) error {
}

func TestDownload_Must_Succeed(t *testing.T) {
//80:func download(arch SampleArchive, logs chan string) error {
}

func TestDownload_Must_Fail(t *testing.T) {
//80:func download(arch SampleArchive, logs chan string) error {
}
