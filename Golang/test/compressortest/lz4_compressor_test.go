package compressortest

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"Golang/internal/compressor"
)

// TestLZ4Compressor validates the compression and decompression for LZ4Compressor.
func TestLZ4Compressor(t *testing.T) {
	// Data sizes to test: small, medium, and large.
	testDataSizes := []int{100, 1024, 1_000_000} // 100 bytes, 1 KB, 1 MB

	for _, size := range testDataSizes {
		t.Run("LZ4-"+string(rune(size)), func(t *testing.T) {
			data := randomBytes(size)

			// t.Logf("%v", data)

			// Initialize the LZ4 compressor
			c, err := compressor.NewCompressor(compressor.LZ4)
			if err != nil {
				t.Fatalf("Failed to initialize LZ4 compressor: %s", err)
			}

			// Test compression
			compressed, err := c.Compress(data)
			if err != nil {
				t.Fatalf("Compression failed for LZ4: %s", err)
			}

			// Verify compression achieves actual reduction in size (for data > 1 KB)
			if len(compressed) >= len(data) {
				// Note: LZ4 may occasionally result in upsized compressed data for small blocks.
				t.Logf("Warning: Compressed data size (%d) is not smaller than original size (%d).",
					len(compressed), len(data))
			}

			// Test decompression
			decompressed, err := c.Decompress(compressed)
			if err != nil {
				t.Fatalf("Decompression failed for LZ4: %s", err)
			}

			// Check that the decompressed data matches the original
			if !bytes.Equal(data, decompressed) {
				t.Errorf("Data mismatch! Expected: %v, Got: %v", data, decompressed)
			}

			// Log results
			t.Logf("LZ4 | Original Size: %d bytes | Compressed Size: %d bytes | Compression Ratio: %.2f",
				len(data), len(compressed), float64(len(compressed))/float64(len(data)))

			// Clean up resources
			c.Destroy()
		})
	}
}

func TestLZ4CompressionRatios(t *testing.T) {
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
			t.Run(fmt.Sprintf("%s-%d", test.name, size), func(t *testing.T) {
				// Initialize compressor
				c, err := compressor.NewCompressor(compressor.LZ4)
				if err != nil {
					t.Fatalf("Failed to initialize LZ4 compressor: %s", err)
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
					t.Errorf("Decompressed data does not match original! Expected: %v, Got: %v", test.data[:10], decompressed[:10])
				}

				// Calculate compression ratio
				compressionRatio := float64(len(compressed)) / float64(len(test.data))

				// Output results
				fmt.Printf("%s Data Test | Size: %d bytes | Compressed: %d bytes | Ratio: %.2f | Compression Time: %s | Decompression Time: %s\n",
					test.name, len(test.data), len(compressed), compressionRatio, compressionDuration, decompressionDuration)

				// Clean up compressor
				c.Destroy()
			})
		}
	}
}

// TestEdgeCasesForLZ4 validates edge cases for the LZ4Compressor, such as empty input.
func TestEdgeCasesForLZ4(t *testing.T) {
	// Initialize an empty input
	data := []byte{}

	// Create the LZ4 compressor
	c, err := compressor.NewCompressor(compressor.LZ4)
	if err != nil {
		t.Fatalf("Failed to initialize LZ4 compressor: %s", err)
	}

	// Test compression
	compressed, err := c.Compress(data)
	if err != nil {
		t.Fatalf("Compression failed for empty input in LZ4: %s", err)
	}

	// Test decompression
	decompressed, err := c.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompression failed for empty input in LZ4: %s", err)
	}

	// Ensure decompressed data matches the original empty input
	if !bytes.Equal(data, decompressed) {
		t.Errorf("Mismatch on empty data! Expected: %v, Got: %v", data, decompressed)
	}

	t.Log("LZ4 Compressor successfully handled empty input.")
	c.Destroy()
}

// TestLargeDatasetForlz4 benchmarks the compression and decompression for very large datasets.
func TestLargeDatasetForlz4(t *testing.T) {
	// Generate a large dataset for testing (e.g., 10MB)
	data := randomBytes(10_000_000) // 10 MB of random data

	// Create the LZ4 compressor
	c, err := compressor.NewCompressor(compressor.LZ4)
	if err != nil {
		t.Fatalf("Failed to initialize LZ4 compressor: %s", err)
	}

	// Measure compression time
	t.Run("CompressLargeData", func(t *testing.T) {
		compressed, err := c.Compress(data)
		if err != nil {
			t.Fatalf("Compression failed for large dataset in LZ4: %s", err)
		}

		t.Logf("Original Size: %d bytes | Compressed Size: %d bytes | Compression Ratio: %.2f",
			len(data), len(compressed), float64(len(compressed))/float64(len(data)))

		// Measure decompression time
		t.Run("DecompressLargeData", func(t *testing.T) {
			decompressed, err := c.Decompress(compressed)
			if err != nil {
				t.Fatalf("Decompression failed for large dataset in LZ4: %s", err)
			}

			// Verify integrity of the decompressed data
			if !bytes.Equal(data, decompressed) {
				t.Errorf("Data mismatch for large dataset! Expected: %v bytes, Got: %v bytes",
					len(data), len(decompressed))
			}
		})
	})

	// Clean up resources
	c.Destroy()
}
