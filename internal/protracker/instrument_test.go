package protracker

import (
	"testing"
)

func TestInstrumentFactory_Must_Succeed(t *testing.T) {
	//36:func InstrumentFactory() Instrument {
}

func TestInstrumentFactory_Must_Fail(t *testing.T) {
	//36:func InstrumentFactory() Instrument {
}

func TestInstrumentCalculateLength_Must_Succeed(t *testing.T) {
	//50:func (i *Instrument) CalculateLength() {
}

func TestInstrumentCalculateLength_Must_Fail(t *testing.T) {
	//50:func (i *Instrument) CalculateLength() {
}

func TestInstrumentSave_Must_Succeed(t *testing.T) {
	//64:func (i *Instrument) Save() string {
}

func TestInstrumentSave_Must_Fail(t *testing.T) {
	//64:func (i *Instrument) Save() string {
}

func TestInstrumentEncodeInstrumentTitle_Must_Succeed(t *testing.T) {
	//72:func (i *Instrument) encodeInstrumentTitle(out *[30]byte) error {
}

func TestInstrumentEncodeInstrumentTitle_Must_Fail(t *testing.T) {
	//72:func (i *Instrument) encodeInstrumentTitle(out *[30]byte) error {
}

func TestInstrumentEncodeInstrumentLoop_Must_Succeed(t *testing.T) {
	//87:func (i *Instrument) encodeInstrumentLoop(out *[30]byte) error {
}

func TestInstrumentEncodeInstrumentLoop_Must_Fail(t *testing.T) {
	//87:func (i *Instrument) encodeInstrumentLoop(out *[30]byte) error {
}

func TestInstrumentEncodeInstrumentHeader_Must_Succeed(t *testing.T) {
	//107:func (i *Instrument) encodeInstrumentHeader() ([30]byte, error) {
}

func TestInstrumentEncodeInstrumentHeader_Must_Fail(t *testing.T) {
	//107:func (i *Instrument) encodeInstrumentHeader() ([30]byte, error) {
}

func TestInstrumentIsSlotPopulated_Must_Succeed(t *testing.T) {
	//126:func (i *Instrument) isSlotPopulated() bool {
}

func TestInstrumentIsSlotPopulated_Must_Fail(t *testing.T) {
	//126:func (i *Instrument) isSlotPopulated() bool {
}

func TestInstrumentTemporaryWorkaroundWhileXlhaBroken_Must_Succeed(t *testing.T) {
	//130:func (i *Instrument) temporaryWorkaroundWhileXlhaBroken(data *[]byte) error {
}

func TestInstrumentTemporaryWorkaroundWhileXlhaBroken_Must_Fail(t *testing.T) {
	//130:func (i *Instrument) temporaryWorkaroundWhileXlhaBroken(data *[]byte) error {
}

func TestInstrumentExtractSample_Must_Succeed(t *testing.T) {
	//151:func (i *Instrument) ExtractSample() ([]byte, error) {
}

func TestInstrumentExtractSample_Must_Fail(t *testing.T) {
	//151:func (i *Instrument) ExtractSample() ([]byte, error) {
}

func TestInstrumentWriteHeader_Must_Succeed(t *testing.T) {
	//182:func (i *Instrument) WriteHeader(w io.Writer) error {
}

func TestInstrumentWriteHeader_Must_Fail(t *testing.T) {
	//182:func (i *Instrument) WriteHeader(w io.Writer) error {
}

func TestInstrumentWrite_Must_Succeed(t *testing.T) {
	//198:func (i *Instrument) Write(w io.Writer) error {
}

func TestInstrumentWrite_Must_Fail(t *testing.T) {
	//198:func (i *Instrument) Write(w io.Writer) error {
}
