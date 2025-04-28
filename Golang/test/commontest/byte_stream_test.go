package commontest

import (
	"Golang/internal/common"
	"Golang/internal/utils"
	"errors"
	"testing"
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

// TestReset verifies that ByteStream correctly resets its state.
func TestReset(t *testing.T) {
	suite := &ByteStreamTestSuite{}
	suite.Setup(t, 16) // Page size = 16 bytes
	defer suite.TearDown()

	data := []byte{0x01, 0x02, 0x03}
	suite.WriteToStream(t, data)

	// Reset the ByteStream
	suite.byteStream.Reset()

	// Verify that the stream is cleared and all values are reset
	if suite.byteStream.TotalSize != 0 {
		t.Errorf("Total size mismatch after reset. Expected: 0, Got: %d", suite.byteStream.TotalSize)
	}
	if suite.byteStream.RemainingSize() != 0 {
		t.Errorf("Remaining size mismatch after reset. Expected: 0, Got: %d", suite.byteStream.RemainingSize())
	}
	if suite.byteStream.ReadPos != uint32(0) {
		t.Errorf("Read position mismatch after reset. Expected: 0, Got: %v", suite.byteStream.ReadPos)
	}
	if suite.byteStream.IsWrappedBuffer() {
		t.Fatalf("Wrapped buffer was not cleared after reset.")
	}
}

// TestWriteBeyondCapacity verifies writing beyond the capacity of the ByteStream.
func TestWriteBeyondCapacity(t *testing.T) {
	suite := &ByteStreamTestSuite{}
	suite.Setup(t, 8) // Page size = 8 bytes
	defer suite.TearDown()

	// Write data larger than the initial capacity
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}
	suite.WriteToStream(t, data)

	if suite.byteStream.TotalSize != uint32(len(data)) {
		t.Errorf("Total size mismatch after write. Expected: %d, Got: %d", len(data), suite.byteStream.TotalSize)
	}

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

// TestWriteMoreThanPageSize verifies writing data larger than the page size into ByteStream.
func TestWriteMoreThanPageSize(t *testing.T) {
	suite := &ByteStreamTestSuite{}
	suite.Setup(t, 16) // Page size = 16 bytes
	defer suite.TearDown()

	// Write 20 bytes of data into the ByteStream
	data := make([]byte, 20)
	for i := 0; i < 20; i++ {
		data[i] = uint8(i)
	}
	suite.WriteToStream(t, data)

	if suite.byteStream.TotalSize != uint32(len(data)) {
		t.Errorf("Total size mismatch. Expected: %d, Got: %d", len(data), suite.byteStream.TotalSize)
	}
	if suite.byteStream.RemainingSize() != uint32(len(data)) {
		t.Errorf("Remaining size mismatch. Expected: %d, Got: %d", len(data), suite.byteStream.RemainingSize())
	}

	// Attempt to read back the data
	readBuffer := make([]byte, 20)
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

// TestReadMoreThanAvailable verifies reading more data than available in ByteStream.
func TestReadMoreThanAvailable(t *testing.T) {
	suite := &ByteStreamTestSuite{}
	suite.Setup(t, 16) // Page size = 16 bytes
	defer suite.TearDown()

	data := []byte{0x01, 0x02, 0x03}
	suite.WriteToStream(t, data)

	// Attempt to read more data than available
	readBuffer := make([]byte, 4) // Asking for 4 bytes, but only 3 are available
	readLen, err := suite.byteStream.ReadBuf(readBuffer, 4)

	if err == nil || !errors.Is(utils.GetError(utils.ErrPartialRead), err) {
		t.Fatalf("Expected a partial read error, but got: %v", err)
	}
	if readLen != uint32(len(data)) {
		t.Errorf("Partial read length mismatch. Expected: %d, Got: %d", len(data), readLen)
	}

	// Verify the data read matches the original data
	for i := 0; i < len(data); i++ {
		if readBuffer[i] != data[i] {
			t.Errorf("Data mismatch at index %d. Expected: %d, Got: %d", i, data[i], readBuffer[i])
		}
	}
}

// TestEmptyStreamRead verifies behavior when trying to read from an empty stream.
func TestEmptyStreamRead(t *testing.T) {
	suite := &ByteStreamTestSuite{}
	suite.Setup(t, 16)
	defer suite.TearDown()

	readBuffer := make([]byte, 4)
	readLen, err := suite.byteStream.ReadBuf(readBuffer, 4)

	if err == nil {
		t.Fatalf("Read should fail when stream is empty, but no error occurred.")
	}
	if readLen != 0 {
		t.Errorf("Read length mismatch when stream is empty. Expected: 0, Got: %d", readLen)
	}
}

// TestWrapAndClear verifies wrapping an external buffer and clearing it in ByteStream.
func TestWrapAndClear(t *testing.T) {
	suite := &ByteStreamTestSuite{}
	suite.Setup(t, 16)
	defer suite.TearDown()

	// Wrap an external buffer
	externalBuffer := []byte("Hello, World!")
	suite.WrapExternalBuffer(externalBuffer)

	if !suite.byteStream.IsWrappedBuffer() {
		t.Fatalf("Expected ByteStream to be using a wrapped buffer, but it wasn't.")
	}
	if string(suite.byteStream.WrappedBuffer) != string(externalBuffer) {
		t.Fatalf("Wrapped buffer content mismatch. Expected: %s, Got: %s",
			string(externalBuffer), string(suite.byteStream.WrappedBuffer))
	}

	// Clear the wrapped buffer
	suite.byteStream.ClearWrappedBuffer()
	if suite.byteStream.IsWrappedBuffer() {
		t.Fatalf("Expected ByteStream to no longer be using a wrapped buffer, but it still is.")
	}
	if suite.byteStream.WrappedBuffer != nil {
		t.Fatalf("Expected wrapped buffer to be nil after clearing, but got: %v", suite.byteStream.WrappedBuffer)
	}
}

func TestBufferIterator(t *testing.T) {
	// Create a new ByteStream with a page size of 16 bytes
	suite := &ByteStreamTestSuite{}
	suite.Setup(t, 4)
	defer suite.TearDown()

	// Write data into the ByteStream, which creates multiple pages
	data1 := []byte{0x01, 0x02, 0x03, 0x04} // Page 1
	data2 := []byte{0x05, 0x06, 0x07, 0x08} // Page 2
	data3 := []byte{0x09, 0x10}             // Page 3

	suite.WriteToStream(t, data1)
	suite.WriteToStream(t, data2)
	suite.WriteToStream(t, data3)

	// Expected data from all pages
	expectedBuffers := [][]byte{
		{0x01, 0x02, 0x03, 0x04},
		{0x05, 0x06, 0x07, 0x08},
		{0x09, 0x10},
	}

	// Ensure the iterator starts properly from the first page
	it := suite.byteStream.BufferIterator()
	for i, expectedBuffer := range expectedBuffers {
		buf, length, ok := it()
		if !ok {
			t.Errorf("Iterator stopped early, expected %d pages but got fewer", len(expectedBuffers))
			break
		}

		// Verify buffer length
		if length != uint32(len(expectedBuffer)) {
			t.Errorf("Page %d: expected length %d, got %d", i+1, len(expectedBuffer), length)
		}

		// Verify buffer contents
		for j, expectedByte := range expectedBuffer {
			if buf[j] != expectedByte {
				t.Errorf(
					"Page %d: expected byte at index %d to be %d, got %d",
					i+1, j, expectedByte, buf[j],
				)
			}
		}
	}

	// Check if iterator correctly returns false after all pages
	_, _, ok := it()
	if ok {
		t.Error("Iterator did not return false after traversing all buffers")
	}
}
