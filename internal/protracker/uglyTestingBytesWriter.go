package protracker

// Apologies for this horror. It's only for testing, I promise.

type AwfulUnsafeWriter []byte

// Not thread-safe due to this. I can't work out a good way to persist
// data when you also need an interface with non-pointer receivers.
var internalTestingBuffer []byte

func (b AwfulUnsafeWriter) Write(p []byte) (int, error) {
	oldLen := len(internalTestingBuffer)
	newLen := oldLen + len(p)
	newData := make([]byte, newLen)
	if oldLen > 0 {
		copy(newData[0:oldLen], internalTestingBuffer)
	}
	copy(newData[oldLen:], p)
	internalTestingBuffer = newData
	return len(p), nil
}

func (b AwfulUnsafeWriter) Bytes() []byte {
	b = internalTestingBuffer
	return b
}
