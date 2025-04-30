package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"

	"Golang/internal/common" // Import ByteStream from the common package
)

// GzipCompressor implements the Compressor interface, with ByteStream integration
type GzipCompressor struct {
	mu     sync.Mutex         // Mutex for thread-safe operation
	buffer bytes.Buffer       // Internal buffer for writing compressed data
	writer *gzip.Writer       // Gzip writer for compressing data
	reader *gzip.Reader       // Gzip reader for decompressing data
	stream *common.ByteStream // ByteStream object for buffering data
}

// newGzipCompressor creates a new GzipCompressor instance
func newGzipCompressor() *GzipCompressor {
	stream, err := common.NewByteStream(512) // Initialize ByteStream with 512-byte pages (custom value)
	if err != nil {
		panic("Failed to create ByteStream: " + err.Error()) // ByteStream initialization should not fail
	}
	return &GzipCompressor{
		stream: stream, // Attach the ByteStream instance
	}
}

// Compress compresses data using gzip and writes it into the ByteStream
func (c *GzipCompressor) Compress(data []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Reset internal buffer and gzip writer
	c.buffer.Reset()
	if c.writer == nil {
		c.writer = gzip.NewWriter(&c.buffer)
	} else {
		c.writer.Reset(&c.buffer)
	}

	// Write data to gzip.Writer
	if _, err := c.writer.Write(data); err != nil {
		return nil, err
	}

	// Close the writer to flush compressed data to the buffer
	if err := c.writer.Close(); err != nil {
		return nil, err
	}

	// Write the compressed data to ByteStream
	if err := c.stream.WriteBuf(c.buffer.Bytes(), uint32(c.buffer.Len())); err != nil {
		return nil, err
	}

	// Return the compressed data from ByteStream
	readbyte, _ := c.stream.GetBytesFromByteStream()
	return readbyte, nil
}

// Decompress reads from the ByteStream, decompresses the data, and returns the result
func (c *GzipCompressor) Decompress(data []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Short-circuit for empty input
	if len(data) == 0 {
		return []byte{}, nil
	}

	// Create a reader directly from the input data
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Decompress the data
	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return decompressedData, nil
}

// Reset resets the GzipCompressor for compressing or decompressing data
func (c *GzipCompressor) Reset(forCompress bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Reset ByteStream
	c.stream.Reset()

	// Reset the gzip writer or reader
	if forCompress && c.writer != nil {
		c.writer.Reset(&c.buffer)
	} else if !forCompress && c.reader != nil {
		if err := c.reader.Close(); err != nil {
			return err
		}
		c.reader = nil
	}
	return nil
}

// Destroy releases resources used by the GzipCompressor
func (c *GzipCompressor) Destroy() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close writer and reader
	if c.writer != nil {
		_ = c.writer.Close()
		c.writer = nil
	}
	if c.reader != nil {
		_ = c.reader.Close()
		c.reader = nil
	}
	c.buffer.Reset()

	// Reset ByteStream
	c.stream.Reset()
}
