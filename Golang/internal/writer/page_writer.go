package writer

import (
	"Golang/internal/common/base"
	"Golang/internal/common/statistic"
	"Golang/internal/compressor"
	"Golang/internal/encoder"
	"errors"
	"fmt"
)

const OutStreamPageSize = 1024

// PageData represents metadata and buffers for page management.
type PageData struct {
	TimeBufSize      uint32 // Size of the time buffer
	ValueBufSize     uint32 // Size of the value buffer
	UncompressedSize uint32 // Total size of the uncompressed data
	CompressedSize   uint32 // Total size of the compressed data
	UncompressedBuf  []byte // Combined uncompressed buffer (time + value)
	CompressedBuf    []byte // Compressed buffer
	Compressor       compressor.Compressor
}

// Initialize initializes the PageData by combining the time and value byte streams into a single uncompressed buffer,
// and optionally compresses the data.
func (data *PageData) Initialize(timeBS, valueBS *base.ByteStream, compressor compressor.Compressor) error {
	// Save the sizes of the time and value buffers
	data.TimeBufSize = timeBS.PageSize
	data.ValueBufSize = valueBS.PageSize

	// Validate buffer sizes
	if data.TimeBufSize == 0 || data.ValueBufSize == 0 {
		return errors.New("time and value buffers must not be empty")
	}
	// Calculate the uncompressed size
	data.UncompressedSize = data.TimeBufSize + data.ValueBufSize

	// Allocate and populate the uncompressed buffer (time + value)
	data.UncompressedBuf = make([]byte, data.UncompressedSize)
	offset := uint32(0)

	// Copy all pages from the time buffer
	timeData, err := timeBS.GetBytesFromByteStream()
	if err != nil {
		return fmt.Errorf("failed to read time buffer: %w", err)
	}
	copy(data.UncompressedBuf[offset:], timeData)
	offset += data.TimeBufSize

	// Copy all pages from the value buffer
	valueData, err := valueBS.GetBytesFromByteStream()
	if err != nil {
		return fmt.Errorf("failed to read value buffer: %w", err)
	}
	copy(data.UncompressedBuf[offset:], valueData)

	// Set compressor
	data.Compressor = compressor

	// Compress the uncompressed buffer if a compressor is provided
	if compressor != nil {
		var err error
		data.CompressedBuf, err = compressor.Compress(data.UncompressedBuf)
		if err != nil {
			return fmt.Errorf("compression failed: %v", err)
		}
		data.CompressedSize = uint32(len(data.CompressedBuf))
	} else {
		data.CompressedSize = 0 // No compression
	}

	return nil
}

// Destroy releases the memory associated with the PageData buffers.
func (data *PageData) Destroy() {
	// Clear the uncompressed buffer
	data.UncompressedBuf = nil

	// Handle compressor-specific cleanup for the compressed buffer
	if data.Compressor != nil && len(data.CompressedBuf) > 0 {
		data.Compressor.Destroy()
		data.CompressedBuf = nil
	}
}

// PageWriter manages buffers for time and value streams, writes data, and prepares pages.
type PageWriter struct {
	DataType       base.TSDataType  // Data type
	TimeOutStream  *base.ByteStream // Output stream for time data
	ValueOutStream *base.ByteStream // Output stream for value data
	TimeEncoder    encoder.Encoder  // time encoder
	ValueEncoder   encoder.Encoder  // value encoder
	Compressor     compressor.Compressor
	Statistic      statistic.Interface // Interface for maintaining statistics
	PointCount     int                 // Number of data points written
	PageData       PageData
}

// Initialize sets up the PageWriter with the required type, encoding, and compression.
func (writer *PageWriter) Initialize(dataType base.TSDataType, encodingType base.TSEncoding, compressionType base.CompressionType) error {

	writer.DataType = dataType
	var err error

	// Create the ByteStreams
	writer.TimeOutStream, err = base.NewByteStream(OutStreamPageSize)
	// fmt.Printf("time out stream size %v", writer.TimeOutStream.TotalSize)
	if err != nil {
		return fmt.Errorf("failed to initialize time stream: %w", err)
	}

	writer.ValueOutStream, err = base.NewByteStream(OutStreamPageSize)
	if err != nil {
		return fmt.Errorf("failed to initialize value stream: %w", err)
	}

	// Initialize the statistic component
	statFactory := statistic.Factory{}
	writer.Statistic, err = statFactory.AllocStatistic(dataType)
	if err != nil {
		return fmt.Errorf("failed to initialize statistics: %w", err)
	}

	// Allocate time encoder
	writer.TimeEncoder, err = encoder.NewEncoder(dataType, encodingType)
	if err != nil {
		return fmt.Errorf("failed to allocate time encoder: %w", err)
	}

	// Allocate value encoder
	writer.ValueEncoder, err = encoder.NewEncoder(dataType, encodingType)
	if err != nil {
		return fmt.Errorf("failed to allocate value encoder: %w", err)
	}

	// Allocate statistics
	factory := statistic.Factory{}
	writer.Statistic, err = factory.AllocStatistic(dataType)
	if err != nil {
		return fmt.Errorf("failed to allocate statistic: %w", err)
	}

	// Allocate compressor
	writer.Compressor, err = compressor.NewCompressor(compressionType)
	if err != nil {
		return fmt.Errorf("failed to allocate compressor: %w", err)
	}

	writer.PointCount = 0

	return nil
}

// Write encodes and writes a data point with a timestamp and value into the PageWriter's output streams.
func (writer *PageWriter) Write(timestamp int64, value interface{}) error {
	// Validate data type
	if !writer.isDataTypeMatch(value) {
		return fmt.Errorf("data type mismatch %v", value)
	}

	// Encode timestamp
	err := writer.TimeEncoder.Encode(timestamp, writer.TimeOutStream)
	if err != nil {
		return fmt.Errorf("failed to encode timestamp: %w", err)
	}

	// Encode value
	err = writer.ValueEncoder.Encode(value, writer.ValueOutStream)
	if err != nil {
		return fmt.Errorf("failed to encode value: %w", err)
	}

	// Update statistics
	err = writer.Statistic.Update(timestamp, value)
	if err != nil {
		return err
	}
	// fmt.Printf("Statistic (%v)", writer.Statistic)

	// Increment point count
	writer.PointCount++

	return nil
}

// FinalizePage finalizes the current page and returns the constructed PageData.
func (writer *PageWriter) FinalizePage() (*PageData, error) {
	// Create a new PageData instance
	pageData := &PageData{}

	// Initialize the page with the current time and value buffers
	err := pageData.Initialize(writer.TimeOutStream, writer.ValueOutStream, writer.Compressor)
	if err != nil {
		return nil, fmt.Errorf("failed to finalize page: %w", err)
	}

	// Reset the PageWriter after finalizing a page
	writer.Reset()

	return pageData, nil
}

// WriteToChunk write the current page data into a chunk
func (writer *PageWriter) WriteToChunk(currPageData *base.ByteStream, writeHeader bool, writeStatistic bool, writeDataToChunk bool) error {
	// prepare to end the current page

	//init the page data
	pageData := &PageData{}
	err := pageData.Initialize(writer.TimeOutStream, writer.ValueOutStream, writer.Compressor)
	if err != nil {
		return fmt.Errorf("failed to finalize page: %w", err)
	}
	serial := base.SerializationUtil{}
	if writeHeader {
		err := serial.WriteVarUint(uint32(pageData.UncompressedSize), currPageData)
		if err != nil {
			return err
		}

		err = serial.WriteVarUint(pageData.CompressedSize, currPageData)
		if err != nil {
			return err
		}
	}
	if writeStatistic {
		err = writer.Statistic.SerializeTypedStat(currPageData)
		if err != nil {
			return err
		}
	}
	if writeDataToChunk {
		err := currPageData.WriteBuf(pageData.CompressedBuf, pageData.CompressedSize)
		if err != nil {
			return err
		}
	}
	return nil
}

// Reset clears the internal state of the PageWriter, including output streams and statistics.
func (writer *PageWriter) Reset() {
	writer.TimeOutStream.Reset()
	writer.ValueOutStream.Reset()
	writer.Statistic.Reset()
	// writer.PointCount = 0

}

// GetPointCount retrieves the current number of points written to the writer.
func (writer *PageWriter) GetPointCount() int {
	return writer.PointCount
}

// GetTimeOutStreamSize retrieves the current size of the time output stream.
func (writer *PageWriter) GetTimeOutStreamSize() int {
	return int(writer.TimeOutStream.TotalSize)
}

// Destroy releases resources allocated by the PageWriter.
func (writer *PageWriter) Destroy() {
	if writer.TimeEncoder != nil {
		writer.TimeEncoder = nil
	}
	if writer.ValueEncoder != nil {
		writer.ValueEncoder = nil
	}
	if writer.Statistic != nil {
		writer.Statistic = nil
	}
	if writer.TimeOutStream != nil {
		writer.TimeOutStream = nil
	}
	if writer.ValueOutStream != nil {
		writer.ValueOutStream = nil
	}
	if writer.Compressor != nil {
		writer.Compressor.Destroy()
		writer.Compressor = nil
	}
}

// Helper method to validate data type compatibility.
func (writer *PageWriter) isDataTypeMatch(value interface{}) bool {
	switch value.(type) {
	case bool:
		return writer.DataType == base.BOOLEAN
	case int32:
		return writer.DataType == base.INT32
	case int64:
		return writer.DataType == base.INT64
	case float32:
		return writer.DataType == base.FLOAT
	case float64:
		return writer.DataType == base.DOUBLE
	default:
		return false
	}
}

// GetCurrPageData returns the current page data
func (writer *PageWriter) GetCurrPageData() PageData {
	return writer.PageData
}
