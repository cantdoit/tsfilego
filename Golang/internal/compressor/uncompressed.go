package compressor

// UncompressedCompressor implements a No-Op Compressor
type UncompressedCompressor struct{}

func (c *UncompressedCompressor) Reset(forCompress bool) error {
	//TODO implement me
	panic("implement me")
}

func (c *UncompressedCompressor) Destroy() {
	//TODO implement me
	panic("implement me")
}

// newUncompressedCompressor creates and returns a new UncompressedCompressor
func newUncompressedCompressor() *UncompressedCompressor {
	return &UncompressedCompressor{}
}

// Compress simply returns the input data as-is (no compression)
func (c *UncompressedCompressor) Compress(data []byte) ([]byte, error) {
	// No compression logic; return the input data as-is
	return data, nil
}

// Decompress simply returns the input data as-is (no decompression)
func (c *UncompressedCompressor) Decompress(data []byte) ([]byte, error) {
	// No decompression logic; return the input data as-is
	return data, nil
}
