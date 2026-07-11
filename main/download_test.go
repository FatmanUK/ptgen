package main

import (
	"testing"
)

func TestThing(t *testing.T) {
}

/*

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

type Download struct {
	Remote string
	Local string
	ExpectedSum SHA256Sum
}

func DownloadFactory(url string, e SHA256Sum, cache string) Download {
func (re *Download) IsChecksumMatch(blob *[]byte) (bool, error) {
	return re.ExpectedSum.IsChecksumMatch(blob)
}

func (re *Download) Validate(r io.ReadCloser,
	logs chan string) ([]byte, error) {
func (re *Download) BlobFromURI(logs chan string) error {

*/
