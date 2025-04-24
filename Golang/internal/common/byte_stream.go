package common

import (
	"bytes"
	_ "errors"
	"os"
)

// ByteStream provides a mechanism to write data to memory or file.
type ByteStream struct {
	buffer *bytes.Buffer // Buffer to hold data before it's written to a file
	file   *os.File      // File to write data to, if applicable
}

// NewByteStream initializes a new ByteStream with an in-memory buffer.
// To write data into a file, use OpenFile.
func NewByteStream() *ByteStream {
	return &ByteStream{
		buffer: bytes.NewBuffer([]byte{}), // Initialize empty buffer
		file:   nil,                       // No file associated initially
	}
}

// OpenFile opens a file for writing and associates it with the ByteStream.
func (bs *ByteStream) OpenFile(filePath string) error {
	file, err := os.Create(filePath) // Open (or create) file for writing
	if err != nil {
		return err
	}
	bs.file = file
	return nil
}

// Write writes raw bytes into the ByteStream (in-memory buffer or file).
func (bs *ByteStream) Write(data []byte) error {
	// If no file is open, write to the in-memory buffer
	if bs.file == nil {
		_, err := bs.buffer.Write(data)
		return err
	}

	// If a file is open, write to the file
	_, err := bs.file.Write(data)
	return err
}

// WriteByte writes a single byte to the ByteStream.
func (bs *ByteStream) WriteByte(data byte) error {
	// If no file is open, write to the in-memory buffer
	if bs.file == nil {
		return bs.buffer.WriteByte(data)
	}

	// If a file is open, write to the file
	_, err := bs.file.Write([]byte{data})
	return err
}

// WriteUint16 writes a uint16 value to the ByteStream in big-endian format.
func (bs *ByteStream) WriteUint16(value uint16) error {
	data := []byte{
		byte(value >> 8), // Most significant byte
		byte(value),      // Least significant byte
	}
	return bs.Write(data)
}

// WriteUint32 writes a uint32 value to the ByteStream in big-endian format.
func (bs *ByteStream) WriteUint32(value uint32) error {
	data := []byte{
		byte(value >> 24), // MSB (most significant byte)
		byte(value >> 16),
		byte(value >> 8),
		byte(value), // LSB (least significant byte)
	}
	return bs.Write(data)
}

// WriteUint64 writes a uint64 value to the ByteStream in big-endian format.
func (bs *ByteStream) WriteUint64(value uint64) error {
	data := []byte{
		byte(value >> 56), // MSB
		byte(value >> 48),
		byte(value >> 40),
		byte(value >> 32),
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value), // LSB
	}
	return bs.Write(data)
}

// Close closes the ByteStream and its associated file, if any.
func (bs *ByteStream) Close() error {
	if bs.file != nil {
		err := bs.file.Close()
		bs.file = nil
		return err
	}
	return nil
}

// Reset resets the ByteStream by clearing the in-memory buffer.
// Note: If the ByteStream is writing to a file, this does NOT affect the file.
func (bs *ByteStream) Reset() {
	bs.buffer.Reset()
}

// GetBufferBytes returns the contents of the in-memory buffer as a byte slice.
// This is useful for testing or if no file was used.
func (bs *ByteStream) GetBufferBytes() []byte {
	return bs.buffer.Bytes()
}

// TODO: Implement Read methods (e.g., ReadByte, ReadUint16, ReadUint32, etc.).
