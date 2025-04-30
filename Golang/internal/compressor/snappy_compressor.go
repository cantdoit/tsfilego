package compressor

import (
	"errors"
	"github.com/golang/snappy"
	"sync"
)

var (
	ErrSnappyDecompressionFailure = errors.New("failed to decompress snappy data")
	ErrSnappyMemoryFailure        = errors.New("memory allocation failed while compressing snappy data")
)

// SnappyCompressor implements the `Compressor` interface for Snappy
type SnappyCompressor struct {
	mu              sync.Mutex // Mutex for thread-safe usage
	compressedBuf   []byte     // Buffer used for compressed data
	uncompressedBuf []byte     // Buffer used for decompressed data
}

// newSnappyCompressor creates and returns a new instance of SnappyCompressor
func newSnappyCompressor() *SnappyCompressor {
	return &SnappyCompressor{}
}

// Compress compresses the provided data using Snappy and returns the compressed data
func (c *SnappyCompressor) Compress(data []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Perform Snappy compression
	compressed := snappy.Encode(nil, data)
	if compressed == nil {
		return nil, ErrSnappyMemoryFailure
	}

	// Update internal buffer for compressed data (optional, based on need)
	c.compressedBuf = compressed

	return compressed, nil
}

// Decompress decompresses Snappy-compressed data and returns the original uncompressed data
func (c *SnappyCompressor) Decompress(data []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Attempt to decompress using Snappy
	decompressed, err := snappy.Decode(nil, data)
	if err != nil {
		return nil, ErrSnappyDecompressionFailure
	}

	// Update internal buffer for uncompressed data (optional, based on need)
	c.uncompressedBuf = decompressed

	return decompressed, nil
}

// Reset resets the compressor's internal buffers, allowing reuse
func (c *SnappyCompressor) Reset(forCompress bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clear internal buffers
	c.compressedBuf = nil
	c.uncompressedBuf = nil

	return nil
}

// Destroy frees up resources used by this SnappyCompressor
func (c *SnappyCompressor) Destroy() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clear buffers to release memory
	c.compressedBuf = nil
	c.uncompressedBuf = nil
}
