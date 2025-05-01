package compressor

import (
	"Golang/internal/common/base"
	"bytes"
	"errors"
	"github.com/pierrec/lz4"
	"io"
	"sync"
)

// Errors specific to the LZ4 compressor
var (
	ErrCompressionFailed   = errors.New("lz4: compression failed")
	ErrDecompressionFailed = errors.New("lz4: decompression failed")
)

// LZ4Compressor handles compression and decompression for LZ4 with ByteStream integration.
type LZ4Compressor struct {
	mu     sync.Mutex       // Mutex for thread-safe operations
	buffer bytes.Buffer     // Temporary buffer for compression data
	writer *lz4.Writer      // LZ4 Writer for compression
	stream *base.ByteStream // ByteStream for efficient memory management
}

// newLz4Compressor initializes a new LZ4Compressor.
func newLz4Compressor() *LZ4Compressor {
	stream, err := base.NewByteStream(512) // ByteStream with 512-byte page size
	if err != nil {
		panic("Failed to initialize ByteStream: " + err.Error())
	}

	return &LZ4Compressor{
		stream: stream,
		writer: nil, // Lazily initialized during compression
	}
}

// Compress compresses the input data using the LZ4 algorithm.
func (c *LZ4Compressor) Compress(data []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Reset internal buffer
	c.buffer.Reset()

	// Lazily initialize the LZ4 writer
	if c.writer == nil {
		c.writer = lz4.NewWriter(&c.buffer)
	} else {
		c.writer.Reset(&c.buffer)
	}

	// Write data to LZ4 writer for compression
	if _, err := c.writer.Write(data); err != nil {
		return nil, ErrCompressionFailed
	}

	// Close the writer to flush compressed data to the internal buffer
	if err := c.writer.Close(); err != nil {
		return nil, ErrCompressionFailed
	}

	// Write the compressed data to ByteStream
	compressedData := c.buffer.Bytes()
	if err := c.stream.WriteBuf(compressedData, uint32(len(compressedData))); err != nil {
		return nil, err
	}

	// Return the compressed data as a slice
	return compressedData, nil
}

// Decompress decompresses LZ4-compressed data.
func (c *LZ4Compressor) Decompress(data []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Short-circuit for empty input
	if len(data) == 0 {
		return []byte{}, nil
	}

	// Create a new LZ4 reader for the input data
	reader := lz4.NewReader(bytes.NewReader(data))

	// Decompress the data into a buffer
	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, ErrDecompressionFailed
	}

	return decompressedData, nil
}

// Reset clears the state of the compressor, preparing it for reuse.
func (c *LZ4Compressor) Reset(forCompress bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Reset internal buffer and ByteStream
	c.buffer.Reset()
	c.stream.Reset()

	// Reset the LZ4 writer only if used for compression
	if forCompress && c.writer != nil {
		c.writer.Reset(&c.buffer)
	}

	return nil
}

// Destroy cleans up the resources used by the LZ4Compressor.
func (c *LZ4Compressor) Destroy() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ensure all resources are released
	c.buffer.Reset()
	c.stream.Reset()
	c.writer = nil // Dereference the LZ4 writer
}
