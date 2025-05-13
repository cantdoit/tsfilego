package writer

import (
	"Golang/internal/common/base"
	"Golang/internal/common/statistic"
	"Golang/internal/compressor"
	"Golang/internal/encoder"
	"errors"
	"fmt"
)

// ValuePageData represents metadata and buffers for a value page.
type ValuePageData struct {
	ColNotNullBitmapBufSize uint32                // Buffer size of the column not-null bitmap
	ValueBufSize            uint32                // Buffer size of the value data
	UncompressedSize        uint32                // Total size of uncompressed data
	CompressedSize          uint32                // Total size of compressed data
	UncompressedBuf         []byte                // Uncompressed buffer
	CompressedBuf           []byte                // Compressed buffer
	Compressor              compressor.Compressor // Compressor instance
}

// Initialize sets up the ValuePageData by combining not-null bitmap and value data streams.
func (data *ValuePageData) Initialize(colNotNullBitmapBS *base.ByteStream, valueBS *base.ByteStream, compressor compressor.Compressor, size uint32) error {
	// Save the sizes of the buffers
	data.ColNotNullBitmapBufSize = colNotNullBitmapBS.TotalSize
	data.ValueBufSize = valueBS.TotalSize

	// Validate buffer sizes
	if data.ColNotNullBitmapBufSize == 0 || data.ValueBufSize == 0 {
		return errors.New("not-null bitmap or value buffer must not be empty")
	}
	// Calculate the uncompressed size (size + bitmap buffer + value buffer)
	data.UncompressedSize = 4 + data.ColNotNullBitmapBufSize + data.ValueBufSize

	// Allocate and populate the uncompressed buffer
	data.UncompressedBuf = make([]byte, data.UncompressedSize)
	offset := uint32(0)

	// Add the size (4-byte big-endian integer)
	data.UncompressedBuf[offset+0] = byte((size >> 24) & 0xFF)
	data.UncompressedBuf[offset+1] = byte((size >> 16) & 0xFF)
	data.UncompressedBuf[offset+2] = byte((size >> 8) & 0xFF)
	data.UncompressedBuf[offset+3] = byte(size & 0xFF)
	offset += 4

	// Copy the column not-null bitmap buffer
	bitmapData, err := colNotNullBitmapBS.GetBytesFromByteStream()
	if err != nil {
		return fmt.Errorf("failed to read bitmap buffer: %w", err)
	}
	copy(data.UncompressedBuf[offset:], bitmapData)
	offset += data.ColNotNullBitmapBufSize

	// Copy the value buffer
	valueData, err := valueBS.GetBytesFromByteStream()
	if err != nil {
		return fmt.Errorf("failed to read value buffer: %w", err)
	}
	copy(data.UncompressedBuf[offset:], valueData)

	// Set compressor and compress the data
	data.Compressor = compressor
	if compressor != nil {
		data.CompressedBuf, err = compressor.Compress(data.UncompressedBuf)
		if err != nil {
			return fmt.Errorf("compression failed: %w", err)
		}
		data.CompressedSize = uint32(len(data.CompressedBuf))
	} else {
		// If no compressor is provided, we still track sizes
		data.CompressedSize = 0
	}

	return nil
}

// Destroy cleans up memory and releases compressor resources.
func (data *ValuePageData) Destroy() {
	data.UncompressedBuf = nil
	if data.Compressor != nil && len(data.CompressedBuf) > 0 {
		data.Compressor.Destroy()
		data.CompressedBuf = nil
	}
}

// ValuePageWriter manages value-based page data and streams for encoding.
type ValuePageWriter struct {
	DataType               base.TSDataType     // Data type for the values
	ValueEncoder           encoder.Encoder     // Encoder for value data
	Statistic              statistic.Interface // Statistics for the page
	ValueOutStream         *base.ByteStream    // Output stream for value buffer
	ColNotNullBitmapStream *base.ByteStream    // Column not-null bitmap stream
	ColNotNullBitmap       []byte
	PointCount             int                   // Number of written points
	Compressor             compressor.Compressor // Compressor for the page
	PageData               ValuePageData         // Metadata and buffers for the value page
	size                   uint32
}

// Initialize sets up the ValuePageWriter with the provided type, encoding, and compression.
func (writer *ValuePageWriter) Initialize(dataType base.TSDataType, encodingType base.TSEncoding, compressionType base.CompressionType) error {
	writer.DataType = dataType

	// Initialize the ValueEncoder
	var err error
	writer.ValueEncoder = encoder.NewPlainEncoder(dataType)

	// Allocate a ByteStream for value output
	writer.ValueOutStream, err = base.NewByteStream(OutStreamPageSize)
	if err != nil {
		return fmt.Errorf("failed to initialize value stream: %w", err)
	}

	writer.ColNotNullBitmapStream, err = base.NewByteStream(OutStreamPageSize)
	if err != nil {
		return fmt.Errorf("failed to initialize Bitmap stream: %w", err)
	}
	writer.ColNotNullBitmap = make([]byte, 0)

	// Initialize the statistics for the page
	statFactory := statistic.Factory{}
	writer.Statistic, err = statFactory.AllocStatistic(dataType)
	if err != nil {
		return fmt.Errorf("failed to allocate statistic: %w", err)
	}

	// Allocate a compressor
	writer.Compressor, err = compressor.NewCompressor(compressionType)
	if err != nil {
		return fmt.Errorf("failed to allocate compressor: %w", err)
	}

	return nil
}

// Write encodes a timestamp and value into the ValuePageWriter along with null handling.
func (writer *ValuePageWriter) Write(timestamp int64, value interface{}, isNull bool) error {
	if base.TSDataType.TSDataTypeToEnum(writer.DataType) >= 254 { //NULLTYPE or ERROR
		return errors.New("data type is not initialized")
	}

	// Handle the column not-null bitmap
	if !isNull {
		bitmapIdx := writer.PointCount / 8
		if bitmapIdx >= len(writer.ColNotNullBitmap) {
			writer.ColNotNullBitmap = append(writer.ColNotNullBitmap, 0)
		}
		writer.ColNotNullBitmap[bitmapIdx] |= 1 << (7 - (writer.PointCount % 8))
	}

	writer.PointCount++

	if isNull {
		return nil // Skip value encoding if it's null
	}

	// Encode the value into the output stream
	err := writer.ValueEncoder.Encode(value, writer.ValueOutStream)
	if err != nil {
		return fmt.Errorf("value encoding failed: %w", err)
	}

	// Update the statistic
	err = writer.Statistic.Update(timestamp, value)
	if err != nil {
		return err
	}
	return nil
}

func (writer *ValuePageWriter) ColNotNullBitmapSize() uint32 {
	return writer.ColNotNullBitmapStream.TotalSize
}

func (writer *ValuePageWriter) GetPageMemorySize() uint32 {
	return writer.ValueOutStream.TotalSize + writer.ColNotNullBitmapStream.TotalSize
}

func (writer *ValuePageWriter) WriteToChunk(chunkData *base.ByteStream, writeHeader bool, writeStatistic bool, writeDataToChunk bool) error {
	err := writer.PrepareEndPage()
	if err != nil {
		return err
	}
	err = writer.PageData.Initialize(writer.ColNotNullBitmapStream, writer.ValueOutStream, writer.Compressor, writer.size)
	if err != nil {
		return err
	}
	writer.ColNotNullBitmap = nil //clear
	writer.size = 0
	writer.ColNotNullBitmapStream = nil
	serial := base.SerializationUtil{}

	if writeHeader {
		err := serial.WriteVarUint(writer.PageData.UncompressedSize, chunkData)
		if err != nil {
			return err
		}
		err = serial.WriteVarUint(writer.PageData.CompressedSize, chunkData)
		if err != nil {
			return err
		}
	}
	if writeStatistic {
		err := writer.Statistic.SerializeTypedStat(chunkData)
		if err != nil {
			return err
		}
	}
	if writeDataToChunk {
		err := chunkData.WriteBuf(writer.PageData.UncompressedBuf, writer.PageData.UncompressedSize)
		if err != nil {
			return err
		}
	}

	return nil
}

func (writer *ValuePageWriter) PrepareEndPage() error {
	// Flush the value encoder
	err := writer.ValueEncoder.Flush(writer.ValueOutStream)
	if err != nil {
		return fmt.Errorf("failed to flush value encoder: %w", err)
	}

	// Write each byte of the colNotNullBitmap into the corresponding stream
	for range writer.ColNotNullBitmap {
		if err := writer.ColNotNullBitmapStream.WriteBuf(writer.ColNotNullBitmap, uint32(len(writer.ColNotNullBitmap))); err != nil {
			return fmt.Errorf("failed to write not-null bitmap byte: %w", err)
		}
	}
	return nil
}

func (writer *ValuePageWriter) EstimateMaxMemSize() uint32 {
	return /*sizeof(int32_t)*/ 4 + 1 + writer.ColNotNullBitmapStream.TotalSize + writer.ValueOutStream.TotalSize // + encoder.GetMaxByteSize(valueEncoder) Returns 0
}

/*
// SealPage finalizes the page by compressing its data.
func (writer *ValuePageWriter) SealPage() error {
	colNotNullBitmapBS := base.by(writer.ColNotNullBitmapStream)
	return writer.PageData.Initialize(colNotNullBitmapBS, writer.ValueOutStream, writer.Compressor, uint32(writer.PointCount))
}

*/

// Reset clears the current page, preparing the writer for the next page.
func (writer *ValuePageWriter) Reset() {
	writer.ValueOutStream.Reset()
	writer.ColNotNullBitmapStream = nil
	writer.PointCount = 0
	writer.Statistic.Reset()
}

// Destroy releases the resources used by the ValuePageWriter.
func (writer *ValuePageWriter) Destroy() {
	if writer.ValueEncoder != nil {
		writer.ValueEncoder.Destroy()
	}
	if writer.Statistic != nil {
		writer.Statistic.Reset()
	}
	if writer.Compressor != nil {
		writer.Compressor.Destroy()
	}
	writer.PageData.Destroy()
	writer.ValueOutStream = nil
	writer.ColNotNullBitmapStream = nil
}

// GetPageData retrieves the metadata and buffers for the current page.
func (writer *ValuePageWriter) GetPageData() ValuePageData {
	return writer.PageData
}

// GetStatistic returns the current statistics for the page.
func (writer *ValuePageWriter) GetStatistic() statistic.Interface {
	return writer.Statistic
}
