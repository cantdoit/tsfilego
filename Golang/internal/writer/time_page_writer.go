package writer

import (
	"Golang/internal/common/base"
	"Golang/internal/common/statistic"
	"Golang/internal/compressor"
	"Golang/internal/encoder"
	"Golang/internal/utils"
	"errors"
	"fmt"
)

// TimePageData represents metadata and buffers for a time page.
type TimePageData struct {
	TimeBufSize      uint32                // Buffer size of the time data
	UncompressedSize uint32                // Total size of uncompressed data
	CompressedSize   uint32                // Total size of compressed data
	UncompressedBuf  []byte                // Uncompressed buffer
	CompressedBuf    []byte                // Compressed buffer
	Compressor       compressor.Compressor // Compressor instance
}

// Initialize sets up the TimePageData by combining and optionally compressing the time buffer.
func (data *TimePageData) Initialize(timeBS *base.ByteStream, compressor compressor.Compressor) error {
	// Save the size of the time buffer
	data.TimeBufSize = timeBS.TotalSize

	// Calculate the uncompressed size (var_uint_size + buffer size)
	varSize := base.GetVarUintSize(data.TimeBufSize) // You should implement this utility if not already present
	data.UncompressedSize = varSize + data.TimeBufSize

	// Allocate memory for the uncompressed buffer
	data.UncompressedBuf = make([]byte, data.UncompressedSize)
	if data.UncompressedBuf == nil {
		return errors.New("failed to allocate memory for uncompressed buffer")
	}

	// Validate the time buffer size
	if data.TimeBufSize == 0 {
		return errors.New("time buffer size cannot be zero")
	}

	// Write the variable-length size of the time buffer into the uncompressed buffer
	serial := base.SerializationUtil{}
	err := serial.WriteVarUint(data.TimeBufSize, timeBS)
	if err != nil {
		return fmt.Errorf("failed to write var uint: %w", err)
	}

	// Copy the time buffer into the uncompressed buffer
	err = timeBS.CopyBSToBuffer(timeBS, data.UncompressedBuf[varSize:], varSize)
	if err != nil {
		return fmt.Errorf("failed to copy time buffer into uncompressed buffer: %w", err)
	}

	// Set the compressor
	data.Compressor = compressor

	// Compress the uncompressed buffer if a compressor is provided
	if compressor != nil {
		var err error
		data.CompressedBuf, err = compressor.Compress(data.UncompressedBuf)
		if err != nil {
			return fmt.Errorf("compression failed: %w", err)
		}
		data.CompressedSize = uint32(len(data.CompressedBuf))
	} else {
		data.CompressedSize = 0 // No compression
	}

	return nil
}

// Destroy releases the memory associated with the buffers in TimePageData.
func (data *TimePageData) Destroy() {
	// Clear the uncompressed buffer
	data.UncompressedBuf = nil

	// Handle cleanup for the compressed buffer and the compressor
	if data.Compressor != nil && len(data.CompressedBuf) > 0 {
		data.Compressor.Destroy()
		data.CompressedBuf = nil
	}
}

// TimePageWriter handles the writing of time-aligned pages into a chunk.
type TimePageWriter struct {
	DataType      base.TSDataType       // The data type of the page
	TimeEncoder   encoder.Encoder       // Encoder for time data
	Statistic     statistic.Interface   // Statistics for the page data (e.g., min, max)
	TimeOutStream *base.ByteStream      // Output stream for the page's time data
	Compressor    compressor.Compressor // Compressor to handle compression of the page
	CurPageData   TimePageData          // Metadata and buffer management for the current page
	IsInitialized bool                  // Indicates whether the writer was initialized
}

// Initialize sets up the TimePageWriter with encoding and compression.
func (writer *TimePageWriter) Initialize(encoding base.TSEncoding, compression base.CompressionType) error {
	// Set default data type (VECTOR in this case, based on the C++ code)
	writer.DataType = base.VECTOR

	var err error

	// Allocate Encoder
	// TODO implement more encoding types
	writer.TimeEncoder = encoder.NewPlainEncoder(writer.DataType)

	// Allocate Statistic
	statFactory := statistic.Factory{}
	writer.Statistic, err = statFactory.AllocStatistic(writer.DataType)
	if err != nil {
		return fmt.Errorf("failed to allocate statistic: %w", err)
	}

	// Allocate Compressor
	writer.Compressor, err = compressor.NewCompressor(compression)
	if err != nil {
		return fmt.Errorf("failed to allocate compressor: %w", err)
	}

	// Allocate ByteStream
	writer.TimeOutStream, err = base.NewByteStream(OutStreamPageSize)
	if err != nil {
		return fmt.Errorf("failed to allocate time output stream: %w", err)
	}

	// Mark initialization complete
	writer.IsInitialized = true
	return nil
}

// Reset reinitializes the TimePageWriter for a new page.
func (writer *TimePageWriter) Reset() {
	if writer.TimeEncoder != nil {
		writer.TimeEncoder.Destroy()
	}
	if writer.Statistic != nil {
		writer.Statistic.Reset()
	}
	if writer.TimeOutStream != nil {
		writer.TimeOutStream.Reset()
	}
}

// Destroy releases resources used by the TimePageWriter.
func (writer *TimePageWriter) Destroy() {
	if writer.IsInitialized {
		writer.IsInitialized = false
		if writer.TimeEncoder != nil {
			writer.TimeEncoder.Destroy()
		}
		if writer.Statistic != nil {
			writer.Statistic.Reset()
		}
		if writer.Compressor != nil {
			writer.Compressor.Destroy()
		}
	}
}

// Write encodes a timestamp and updates statistics.
func (writer *TimePageWriter) Write(timestamp int64) error {
	// Encode the timestamp
	err := writer.TimeEncoder.Encode(timestamp, writer.TimeOutStream)
	if err != nil {
		return fmt.Errorf("failed to encode timestamp: %w", err)
	}

	// Update statistics with the timestamp
	err = writer.Statistic.Update(timestamp, nil)
	if err != nil {
		return err
	}
	return nil
}

func (writer *TimePageWriter) IsPageFull() bool {
	return (int(writer.GetPointNumber()) >= utils.ConfigValue.PageWriterMaxPointNum) || (int(writer.GetPageMemorySize()) >= utils.ConfigValue.PageWriterMaxMemoryBytes)
}

// GetPointNumber returns the number of points stored in the page's statistics.
func (writer *TimePageWriter) GetPointNumber() uint32 {
	return uint32(writer.Statistic.GetCount())
}

// GetTimeOutStreamSize returns the size of the time output stream.
func (writer *TimePageWriter) GetTimeOutStreamSize() uint32 {
	return writer.TimeOutStream.TotalSize
}

// GetPageMemorySize returns the memory size used by the page.
func (writer *TimePageWriter) GetPageMemorySize() uint32 {
	return writer.TimeOutStream.TotalSize
}

/*
// EstimateMaxMemSize estimates the maximum memory size for the page.
func (writer *TimePageWriter) EstimateMaxMemSize() uint32 {
	if writer.TimeEncoder == nil {
		return 0
	}
	return writer.TimeOutStream.TotalSize + writer.TimeEncoder.GetMaxByteSize()
}
*/

// PrepareEndPage finalizes the page by flushing the encoder.
func (writer *TimePageWriter) PrepareEndPage() error {
	return writer.TimeEncoder.Flush(writer.TimeOutStream)
}

// WriteToChunk handles writing the page data to the chunk data stream.
func (writer *TimePageWriter) WriteToChunk(chunkData *base.ByteStream, writeHeader, writeStatistic, writeData bool) error {
	// Finalize the current page
	err := writer.PrepareEndPage()
	if err != nil {
		return fmt.Errorf("failed to finalize page: %w", err)
	}

	// Copy the page data into the chunk
	if writeData {
		err = writer.copyPageDataTo(writer.TimeOutStream, chunkData)
		if err != nil {
			return fmt.Errorf("failed to copy page data to chunk: %w", err)
		}
	}

	// If requested, write statistics
	if writeStatistic {
		err = writer.Statistic.SerializeTypedStat(chunkData)
		if err != nil {
			return fmt.Errorf("failed to write statistics to chunk: %w", err)
		}
	}

	return nil
}

// GetTimeData retrieves the ByteStream containing the page's time data.
func (writer *TimePageWriter) GetTimeData() TimePageData {
	return writer.CurPageData
}

// GetStatistic retrieves statistics for the page's data.
func (writer *TimePageWriter) GetStatistic() statistic.Interface {
	return writer.Statistic
}

// DestroyPageData destroys the resources associated with the current page.
func (writer *TimePageWriter) DestroyPageData() {
	writer.CurPageData.Destroy()
}

// copyPageDataTo copies the page data into a specified chunk data stream.
func (writer *TimePageWriter) copyPageDataTo(pageData, chunkData *base.ByteStream) error {
	data, err := pageData.GetBytesFromByteStream()
	if err != nil {
		return fmt.Errorf("failed to access page data: %w", err)
	}

	if err := chunkData.WriteBuf(data, uint32(len(data))); err != nil {
		return fmt.Errorf("failed to write page data to chunk: %w", err)
	}

	return nil
}
