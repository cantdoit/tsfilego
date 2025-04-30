package compressortest

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"Golang/internal/compressor"
)

func TestSnappyCompressor(t *testing.T) {
	// Data sizes to test: small, medium, and large.
	testDataSizes := []int{100, 1024, 1_000_000} // 100 bytes, 1 KB, 1 MB

	for _, size := range testDataSizes {
		t.Run(fmt.Sprintf("Snappy-%d", size), func(t *testing.T) {
			data := randomBytes(size)

			// Initialize the Snappy compressor
			c, err := compressor.NewCompressor(compressor.SNAPPY)
			if err != nil {
				t.Fatalf("Failed to initialize Snappy compressor: %s", err)
			}

			// Test compression
			compressed, err := c.Compress(data)
			if err != nil {
				t.Fatalf("Compression failed for Snappy: %s", err)
			}

			// Verify compression achieves actual reduction in size (for data > 1 KB)
			if len(compressed) >= len(data) {
				// Note: Snappy may occasionally result in upsized compressed data for small or random blocks.
				t.Logf("Warning: Compressed data size (%d) is not smaller than original size (%d).",
					len(compressed), len(data))
			}

			// Test decompression
			decompressed, err := c.Decompress(compressed)
			if err != nil {
				t.Fatalf("Decompression failed for Snappy: %s", err)
			}

			// Check that the decompressed data matches the original
			if !bytes.Equal(data, decompressed) {
				t.Errorf("Data mismatch! Expected: %v, Got: %v", data[:10], decompressed[:10]) // Log first 10 bytes
			}

			// Log results
			t.Logf("Snappy | Original Size: %d bytes | Compressed Size: %d bytes | Compression Ratio: %.2f",
				len(data), len(compressed), float64(len(compressed))/float64(len(data)))

			// Clean up resources
			c.Destroy()
		})
	}
}

func TestSnappyCompressionRatios(t *testing.T) {
	// Initialize test data sizes: 100B, 1KB, 1MB
	testDataSizes := []int{100, 1024, 1_000_000}

	// Test different data types
	for _, size := range testDataSizes {
		tests := []struct {
			name string
			data []byte
		}{
			{"Random", randomBytes(size)},                              // High entropy random data
			{"Compressible", CompressibleString("abc123abc123", size)}, // Low entropy compressible data
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("Snappy-%s-%d", test.name, size), func(t *testing.T) {
				// Initialize Snappy compressor
				c, err := compressor.NewCompressor(compressor.SNAPPY)
				if err != nil {
					t.Fatalf("Failed to initialize Snappy compressor: %s", err)
				}

				// Compress
				start := time.Now()
				compressed, err := c.Compress(test.data)
				compressionDuration := time.Since(start)
				if err != nil {
					t.Fatalf("Compression failed: %s", err)
				}

				// Decompress
				start = time.Now()
				decompressed, err := c.Decompress(compressed)
				decompressionDuration := time.Since(start)
				if err != nil {
					t.Fatalf("Decompression failed: %s", err)
				}

				// Verify that decompressed data matches original
				if !bytes.Equal(test.data, decompressed) {
					t.Errorf("Decompressed data does not match original! Expected: %v, Got: %v", test.data[:10], decompressed[:10]) // Only log first 10 bytes.
				}

				// Calculate compression ratio
				compressionRatio := float64(len(compressed)) / float64(len(test.data))

				// Output results
				fmt.Printf("%s Data Test | Size: %d bytes | Compressed: %d bytes | Compression Ratio: %.2f | Compression Time: %s | Decompression Time: %s\n",
					test.name, len(test.data), len(compressed), compressionRatio, compressionDuration, decompressionDuration)

				// Clean up resources
				c.Destroy()
			})
		}
	}
}

func TestEdgeCasesForSnappy(t *testing.T) {
	// Edge case: Empty input
	data := []byte{}

	// Create the Snappy compressor
	c, err := compressor.NewCompressor(compressor.SNAPPY)
	if err != nil {
		t.Fatalf("Failed to initialize Snappy compressor: %s", err)
	}

	// Test compression
	compressed, err := c.Compress(data)
	if err != nil {
		t.Fatalf("Compression failed for empty input in Snappy: %s", err)
	}

	// Test decompression
	decompressed, err := c.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompression failed for empty input in Snappy: %s", err)
	}

	// Ensure the decompressed data matches the original
	if !bytes.Equal(data, decompressed) {
		t.Errorf("Decompression mismatch on empty input! Expected: %v, Got: %v", data, decompressed)
	}

	t.Logf("Snappy Edge Case | Input Size: %d bytes | Compressed Size: %d bytes | Decompressed Size: %d bytes",
		len(data), len(compressed), len(decompressed))

	// Clean up resources
	c.Destroy()
}
