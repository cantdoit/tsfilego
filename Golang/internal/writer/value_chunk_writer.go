package writer

import (
	"Golang/internal/common/base"
	"Golang/internal/common/core"
	"Golang/internal/common/statistic"
	"fmt"
)

// ValueChunkWriter manages multiple pages and writes them into a value chunk.
type ValueChunkWriter struct {
	DataType           base.TSDataType     // The data type for the chunk
	ValuePageWriter    *ValuePageWriter    // The PageWriter managing individual pages
	ChunkStatistic     statistic.Interface // Overall statistics for the chunk
	ChunkData          *base.ByteStream    // ByteStream where the chunk data is written
	FirstPageData      *ValuePageData      // Data for the first page, if applicable
	FirstPageStatistic statistic.Interface // Statistics for the first page
	ChunkHeader        core.ChunkHeader    // Metadata for the chunk
	NumOfPages         int                 // Total number of pages in the chunk
}

// Initialize sets up the ValueChunkWriter with measurement metadata.
func (writer *ValueChunkWriter) Initialize(measurementName string, dataType base.TSDataType, encoding base.TSEncoding, compressionType base.CompressionType) error {
	writer.DataType = dataType
	writer.ChunkHeader = core.ChunkHeader{
		MeasurementName: measurementName,
		DataType:        dataType,
		CompressionType: compressionType,
		EncodingType:    encoding,
	}

	// Initialize the chunk statistic
	var err error
	stat := statistic.Factory{}
	writer.ChunkStatistic, err = stat.AllocStatistic(dataType)
	if err != nil {
		return fmt.Errorf("failed to initialize chunk statistic: %w", err)
	}

	// Initialize the value page writer
	err = writer.ValuePageWriter.Initialize(dataType, encoding, compressionType)
	if err != nil {
		return fmt.Errorf("failed to initialize value page writer: %w", err)
	}

	// Initialize the first page statistic
	writer.FirstPageStatistic, err = stat.AllocStatistic(dataType)
	if err != nil {
		return fmt.Errorf("failed to initialize first page statistic: %w", err)
	}

	// Initialize the ByteStream for the chunk data
	writer.ChunkData, err = base.NewByteStream(OutStreamPageSize)
	if err != nil {
		return fmt.Errorf("failed to create chunk data ByteStream: %w", err)
	}

	writer.NumOfPages = 0
	return nil
}

// Destroy releases resources used by the ValueChunkWriter.
func (writer *ValueChunkWriter) Destroy() {
	if writer.NumOfPages == 1 {
		writer.freeFirstWriterData()
	}
	if writer.ValuePageWriter != nil {
		writer.ValuePageWriter.Destroy()
		writer.ValuePageWriter = nil
	}
	if writer.ChunkStatistic != nil {
		writer.ChunkStatistic.Reset()
		writer.ChunkStatistic = nil
	}
	if writer.FirstPageStatistic != nil {
		writer.FirstPageStatistic.Reset()
		writer.FirstPageStatistic = nil
	}
	writer.NumOfPages = 0
	writer.ChunkData = nil // Free memory
}

// Write handles writing a timestamp, value, and null indicator to the chunk.
func (writer *ValueChunkWriter) Write(timestamp int64, value interface{}, isNull bool) error {
	// Ensure the data type matches
	if writer.ValuePageWriter.DataType != writer.DataType {
		return fmt.Errorf("data type mismatch: %d, %d", base.TSDataType.TSDataTypeToEnum(writer.ValuePageWriter.DataType), base.TSDataType.TSDataTypeToEnum(writer.DataType))
	}

	err := writer.ValuePageWriter.Write(timestamp, value, isNull)
	if err != nil {
		return fmt.Errorf("failed to write value: %w", err)
	}

	// Seal the current page if it is full
	if writer.shouldSealCurrentPage() {
		return writer.SealCurrentPage(false)
	}
	return nil
}

// shouldSealCurrentPage checks if the current page exceeds limits and needs sealing.
func (writer *ValueChunkWriter) shouldSealCurrentPage() bool {
	return writer.ValuePageWriter.PointCount >= OutStreamPageSize
}

// SealCurrentPage finalizes the current page and handles chunk-specific logic.
func (writer *ValueChunkWriter) SealCurrentPage(endChunk bool) error {
	// Merge statistics from the current page into the chunk's statistics
	err := writer.ChunkStatistic.MergeWith(writer.ValuePageWriter.Statistic)
	if err != nil {
		return fmt.Errorf("failed to merge page statistics with chunk statistics: %w", err)
	}

	// Handle the first page logic
	if writer.NumOfPages == 0 {
		if endChunk {
			// If this is the only page, finalize the chunk by writing the data
			err = writer.ValuePageWriter.WriteToChunk(writer.ChunkData, true, false, true)
			if err != nil {
				return fmt.Errorf("failed to write first page to chunk: %w", err)
			}
			writer.ValuePageWriter.Destroy()
		} else {
			// Save the data for the first page
			err = writer.ValuePageWriter.WriteToChunk(writer.ChunkData, true, false, false)
			if err != nil {
				return fmt.Errorf("failed to save first page data: %w", err)
			}
			writer.saveFirstPageData(writer.ValuePageWriter)
			writer.ValuePageWriter.Reset()
		}
	} else {
		// Handle subsequent pages
		if writer.NumOfPages == 1 {
			// Write the first page if this is the second page
			err = writer.writeFirstPageData(writer.ChunkData)
			if err != nil {
				return fmt.Errorf("failed to write saved first page: %w", err)
			}
			writer.freeFirstWriterData()
		}

		err = writer.ValuePageWriter.WriteToChunk(writer.ChunkData, false, true, true)
		if err != nil {
			return fmt.Errorf("failed to write page data to chunk: %w", err)
		}
		writer.ValuePageWriter.Reset()
	}

	writer.NumOfPages++
	return nil
}

// saveFirstPageData saves the data of the first page for deferred writing.
func (writer *ValueChunkWriter) saveFirstPageData(pageWriter *ValuePageWriter) {
	writer.FirstPageData = &pageWriter.PageData
	writer.FirstPageStatistic = pageWriter.Statistic.Clone()
}

// writeFirstPageData writes the saved first page's data to the chunk.
func (writer *ValueChunkWriter) writeFirstPageData(chunkData *base.ByteStream) error {
	err := writer.FirstPageStatistic.SerializeTypedStat(chunkData)
	if err != nil {
		return err
	}
	err = chunkData.WriteBuf(writer.FirstPageData.CompressedBuf, writer.FirstPageData.CompressedSize)
	if err != nil {
		return err
	}
	return nil
}

// freeFirstWriterData frees the resources associated with the first page.
func (writer *ValueChunkWriter) freeFirstWriterData() {
	if writer.FirstPageData != nil {
		writer.FirstPageData.Destroy()
		writer.FirstPageData = nil
	}
	if writer.FirstPageStatistic != nil {
		writer.FirstPageStatistic.Reset()
		writer.FirstPageStatistic = nil
	}
}

// GetChunkData returns the byte stream containing the chunk's data.
func (writer *ValueChunkWriter) GetChunkData() *base.ByteStream {
	return writer.ChunkData
}

// GetChunkStatistics returns the chunk's statistics.
func (writer *ValueChunkWriter) GetChunkStatistics() statistic.Interface {
	return writer.ChunkStatistic
}

// EndEncodeChunk finalizes the encoding process for the chunk.
func (writer *ValueChunkWriter) EndEncodeChunk() error {
	return writer.SealCurrentPage(true)
}
