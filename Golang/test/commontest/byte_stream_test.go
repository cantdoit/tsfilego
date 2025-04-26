package commontest

import (
	"testing"

	"Golang/internal/common" // Import your `ByteStream` implementation
)

// ByteStreamTestSuite represents the test suite for ByteStream.
type ByteStreamTestSuite struct {
	byteStream *common.ByteStream
}

// Setup initializes a new ByteStream instance for tests.
func (suite *ByteStreamTestSuite) Setup(t *testing.T, pageSize uint32) {
	var err error
	suite.byteStream, err = common.NewByteStream(pageSize)
	if err != nil {
		t.Fatalf("Failed to set up ByteStream: %v", err)
	}
}

// TearDown releases the ByteStream instance.
func (suite *ByteStreamTestSuite) TearDown() {
	suite.byteStream = nil
}

// WriteToStream writes data to ByteStream.
func (suite *ByteStreamTestSuite) WriteToStream(t *testing.T, data []byte) {
	err := suite.byteStream.WriteBuf(data, uint32(len(data)))
	if err != nil {
		t.Fatalf("Failed to write to ByteStream: %v", err)
	}
}

// ReadFromStream reads data from ByteStream.
func (suite *ByteStreamTestSuite) ReadFromStream(t *testing.T, buffer []byte, wantLen uint32) uint32 {
	readLen, err := suite.byteStream.ReadBuf(buffer, wantLen)
	if err != nil {
		t.Fatalf("Failed to read from ByteStream: %v", err)
	}
	return readLen
}

// WrapExternalBuffer wraps an external buffer to be used by ByteStream.
func (suite *ByteStreamTestSuite) WrapExternalBuffer(buffer []byte) {
	suite.byteStream.WrapFrom(buffer, uint32(len(buffer)))
}

// TestWriteRead verifies writing to and reading from a ByteStream.
func TestWriteRead(t *testing.T) {
	suite := &ByteStreamTestSuite{}
	suite.Setup(t, 16) // Page size = 16 bytes
	defer suite.TearDown()

	data := []byte{0x01, 0x02, 0x03}
	suite.WriteToStream(t, data)

	readBuffer := make([]byte, len(data))
	readLen := suite.ReadFromStream(t, readBuffer, uint32(len(data)))

	if readLen != uint32(len(data)) {
		t.Fatalf("Read length mismatch. Expected: %d, Got: %d", len(data), readLen)
	}

	for i := 0; i < len(data); i++ {
		if readBuffer[i] != data[i] {
			t.Errorf("Data mismatch at index %d. Expected: %d, Got: %d", i, data[i], readBuffer[i])
		}
	}
}

// TestWriteReadLargeQuantities verifies writing and reading large data from a ByteStream.
func TestWriteReadLargeQuantities(t *testing.T) {
	suite := &ByteStreamTestSuite{}
	suite.Setup(t, 16)
	defer suite.TearDown()

	dataSize := 1024 * 1024
	data := make([]byte, dataSize)
	for i := 0; i < dataSize; i++ {
		data[i] = uint8(i & 0xff)
	}

	suite.WriteToStream(t, data)

	readBuffer := make([]byte, dataSize)
	for i := 0; i < dataSize; i++ {
		readLen := suite.ReadFromStream(t, readBuffer[i:], 1)
		if readLen != 1 {
			t.Fatalf("Failed to read one byte at index %d", i)
		}
	}

	for i := 0; i < dataSize; i++ {
		if readBuffer[i] != uint8(i&0xff) {
			t.Errorf("Data mismatch at index %d. Expected: %d, Got: %d", i, uint8(i&0xff), readBuffer[i])
		}
	}
}

// TestWrapExternalBuffer verifies wrapping of an external buffer.
func TestWrapExternalBuffer(t *testing.T) {
	suite := &ByteStreamTestSuite{}
	suite.Setup(t, 16)
	defer suite.TearDown()

	externalBuffer := []byte("Hello, World!")
	suite.WrapExternalBuffer(externalBuffer)

	if !suite.byteStream.IsWrappedBuffer() {
		t.Fatalf("Buffer wrapping failed. Expected the buffer to be wrapped.")
	}

	if string(suite.byteStream.WrappedBuffer) != string(externalBuffer) {
		t.Fatalf(
			"Wrapped buffer data mismatch. Expected: %s, Got: %s",
			string(externalBuffer),
			string(suite.byteStream.WrappedBuffer),
		)
	}
}

// TestSize checks the size calculations in ByteStream.
func TestSize(t *testing.T) {
	suite := &ByteStreamTestSuite{}
	suite.Setup(t, 16)
	defer suite.TearDown()

	data := []byte{0x01, 0x02, 0x03}
	suite.WriteToStream(t, data)

	totalSize := suite.byteStream.TotalSize
	if totalSize != uint32(len(data)) {
		t.Errorf("Total size mismatch. Expected: %d, Got: %d", len(data), totalSize)
	}

	remainingSize := suite.byteStream.RemainingSize()
	if remainingSize != uint32(len(data)) {
		t.Errorf("Remaining size mismatch after write. Expected: %d, Got: %d", len(data), remainingSize)
	}

	readBuffer := make([]byte, 2)
	suite.ReadFromStream(t, readBuffer, 2)

	remainingSize = suite.byteStream.RemainingSize()
	if remainingSize != uint32(len(data)-2) {
		t.Errorf("Remaining size mismatch after read. Expected: %d, Got: %d", len(data)-2, remainingSize)
	}
}

// TestMarkReadPosition verifies the marking of read positions in ByteStream.
func TestMarkReadPosition(t *testing.T) {
	suite := &ByteStreamTestSuite{}
	suite.Setup(t, 16)
	defer suite.TearDown()

	data := []byte{0x01, 0x02, 0x03}
	suite.WriteToStream(t, data)

	suite.byteStream.MarkReadPos()

	readBuffer := make([]byte, 2)
	suite.ReadFromStream(t, readBuffer, 2)

	if suite.byteStream.GetMarkLen() != 2 {
		t.Errorf("Mark length mismatch. Expected: %d, Got: %d", 2, suite.byteStream.GetMarkLen())
	}
}
