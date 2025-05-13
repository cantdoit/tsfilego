package utils

var ConfigValue = struct {
	MaxDegreeOfIndexNode     int
	PageWriterMaxPointNum    int
	PageWriterMaxMemoryBytes int
}{
	MaxDegreeOfIndexNode:     256,
	PageWriterMaxPointNum:    5,
	PageWriterMaxMemoryBytes: 128 * 1024,
}
