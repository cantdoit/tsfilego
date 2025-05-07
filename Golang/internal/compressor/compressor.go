package compressor

import (
	"Golang/internal/common/base"
	"errors"
)

// Compressor is the interface for all compressions
type Compressor interface {
	Compress(data []byte) ([]byte, error)   // Compress data
	Decompress(data []byte) ([]byte, error) // Decompress data
	Reset(forCompress bool) error           // Reset compressor/decompressor state
	Destroy()                               // Clean up resources
}

// Available compression types
const (
	UNCOMPRESSED = "uncompressed"
	GZIP         = "gzip"
	SNAPPY       = "snappy"
	LZ4          = "lz4"
)

// NewCompressor - Factory method to create specific compressors
func NewCompressor(compressionType base.CompressionType) (Compressor, error) {
	switch compressionType {
	case GZIP:
		return newGzipCompressor(), nil
	case LZ4:
		return newLz4Compressor(), nil
	case SNAPPY:
		return newSnappyCompressor(), nil
	case base.UNCOMPRESSED:
		return newUncompressedCompressor(), nil
	default:
		return nil, errors.New(string("unsupported compression type: " + compressionType))
	}
}
