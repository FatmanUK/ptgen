package protracker

import (
	"testing"
)

func TestInstrumentFactory_Must_Succeed(t *testing.T) {
//45:func InstrumentFactory() Instrument {
}

func TestInstrumentFactory_Must_Fail(t *testing.T) {
//45:func InstrumentFactory() Instrument {
}

func TestInstrumentCalculateLength_Must_Succeed(t *testing.T) {
//59:func (i *Instrument) CalculateLength() {
}

func TestInstrumentCalculateLength_Must_Fail(t *testing.T) {
//59:func (i *Instrument) CalculateLength() {
}

func TestInstrumentEncodeInstrumentLoop_Must_Succeed(t *testing.T) {
//74:func (i *Instrument) encodeInstrumentLoop(out *[30]byte) error {
}

func TestInstrumentEncodeInstrumentLoop_Must_Fail(t *testing.T) {
//74:func (i *Instrument) encodeInstrumentLoop(out *[30]byte) error {
}

func TestInstrumentEncodeInstrumentHeader_Must_Succeed(t *testing.T) {
//93:func (i *Instrument) encodeInstrumentHeader() ([30]byte, error) {
}

func TestInstrumentEncodeInstrumentHeader_Must_Fail(t *testing.T) {
//93:func (i *Instrument) encodeInstrumentHeader() ([30]byte, error) {
}

func TestInstrumentIsSlotPopulated_Must_Succeed(t *testing.T) {
//116:func (i *Instrument) isSlotPopulated() bool {
}

func TestInstrumentIsSlotPopulated_Must_Fail(t *testing.T) {
//116:func (i *Instrument) isSlotPopulated() bool {
}

func TestInstrumentTemporaryWorkaroundWhileXlhaBroken_Must_Succeed(t *testing.T) {
//120:func (i *Instrument) temporaryWorkaroundWhileXlhaBroken(arch SampleArchive, data *[]byte) error {
}

func TestInstrumentTemporaryWorkaroundWhileXlhaBroken_Must_Fail(t *testing.T) {
//120:func (i *Instrument) temporaryWorkaroundWhileXlhaBroken(arch SampleArchive, data *[]byte) error {
}

func TestLhaLoop_Must_Succeed(t *testing.T) {
//142:func (i *Instrument) lhaLoop(lhaReader *xlha.Reader,
}

func TestLhaLoop_Must_Fail(t *testing.T) {
//142:func (i *Instrument) lhaLoop(lhaReader *xlha.Reader,
}

func TestLhaFile_Must_Succeed(t *testing.T) {
//169:func (i *Instrument) lhaFile(arch SampleArchive, data *[]byte) error {
}

func TestLhaFile_Must_Fail(t *testing.T) {
//169:func (i *Instrument) lhaFile(arch SampleArchive, data *[]byte) error {
}

func TestInstrumentExtractSample_Must_Succeed(t *testing.T) {
//182:func (i *Instrument) ExtractSample() ([]byte, error) {
}

func TestInstrumentExtractSample_Must_Fail(t *testing.T) {
//182:func (i *Instrument) ExtractSample() ([]byte, error) {
}

func TestInstrumentWriteHeader_Must_Succeed(t *testing.T) {
//198:func (i *Instrument) WriteHeader(w io.Writer) error {
}

func TestInstrumentWriteHeader_Must_Fail(t *testing.T) {
//198:func (i *Instrument) WriteHeader(w io.Writer) error {
}

func TestInstrumentWrite_Must_Succeed(t *testing.T) {
//214:func (i *Instrument) Write(w io.Writer) error {
}

func TestInstrumentWrite_Must_Fail(t *testing.T) {
//214:func (i *Instrument) Write(w io.Writer) error {
}
