package utils

var ConfigValue = struct {
	MaxDegreeOfIndexNode     int
	PageWriterMaxPointNum    int
	PageWriterMaxMemoryBytes int
	ChunkGroupSizeThreshold  int
}{
	MaxDegreeOfIndexNode:     256,
	PageWriterMaxPointNum:    5,
	PageWriterMaxMemoryBytes: 128 * 1024,
	ChunkGroupSizeThreshold:  128 * 1024 * 1024,
}
