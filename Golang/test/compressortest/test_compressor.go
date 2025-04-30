package compressortest

import (
	"bytes"
	"crypto/rand"
	"io"
)

// Helper function to generate random byte slices for testing.
func randomBytes(size int) []byte {
	data := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, data) // Safer than rand.Read()
	if err != nil {
		return nil // Handle error (e.g., insufficient entropy)
	}
	return data
}

// CompressibleString generates a large string based on a repeating pattern
func CompressibleString(pattern string, totalLength int) []byte {
	// Repeat the pattern until the total length is achieved
	var result bytes.Buffer
	for result.Len() < totalLength {
		result.WriteString(pattern)
	}
	return result.Bytes()[:totalLength] // Trim to the exact total length
}
