package common

// DefaultEncodings maps data types to their default encoding methods
var DefaultEncodings = map[TSDataType]TSEncoding{
	BOOLEAN: PLAIN,
	INT32:   PLAIN,
	INT64:   PLAIN,
	FLOAT:   PLAIN,
	DOUBLE:  PLAIN,
	TEXT:    PLAIN,
}

// DefaultCompressions maps data types to their default compression methods
var DefaultCompressions = map[TSDataType]CompressionType{
	BOOLEAN: UNCOMPRESSED,
	INT32:   UNCOMPRESSED,
	INT64:   UNCOMPRESSED,
	FLOAT:   UNCOMPRESSED,
	DOUBLE:  UNCOMPRESSED,
	TEXT:    UNCOMPRESSED,
}

// GetDefaultEncoding returns the default encoding for a given data type
func GetDefaultEncoding(dataType TSDataType) TSEncoding {
	if encoding, exists := DefaultEncodings[dataType]; exists {
		return encoding
	}
	return INVALID_E
}

// GetDefaultCompression returns the default compression method for a given data type
func GetDefaultCompression(dataType TSDataType) CompressionType {
	if compression, exists := DefaultCompressions[dataType]; exists {
		return compression
	}
	return INVALID_C
}

// IsValidDataType verifies if a given TSDataType is valid
func IsValidDataType(dataType TSDataType) bool {
	switch dataType {
	case BOOLEAN, INT32, INT64, FLOAT, DOUBLE, TEXT:
		return true
	default:
		return false
	}
}

// IsValidEncoding verifies if a given TSEncoding is valid
func IsValidEncoding(encoding TSEncoding) bool {
	switch encoding {
	case PLAIN, DICTIONARY, RLE, DIFF, TS_2DIFF, BITMAP, REGULAR:
		return true
	default:
		return false
	}
}

// IsValidCompression verifies if a given CompressionType is valid
func IsValidCompression(compression CompressionType) bool {
	switch compression {
	case UNCOMPRESSED, SNAPPY, GZIP, LZO, LZ4:
		return true
	default:
		return false
	}
}
