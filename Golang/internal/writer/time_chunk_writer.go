package writer

import (
	"Golang/internal/common/base"
	"Golang/internal/common/core"
	"Golang/internal/common/statistic"
	_ "errors"
	"fmt"
)

// TimeChunkWriter handles writing a chunk of time-aligned data and manages its pages.
type TimeChunkWriter struct {
	Datatype           base.TSDataType
	TimePageWriter     *TimePageWriter     // Writer for individual time-aligned pages
	ChunkStatistic     statistic.Interface // Statistic for the entire chunk
	ChunkData          *base.ByteStream    // ByteStream containing chunk data
	FirstPageData      *TimePageData       // Data for the first page (if applicable)
	FirstPageStatistic statistic.Interface // Statistics for the first page
	ChunkHeader        *core.ChunkHeader   // Metadata for the chunk
	NumOfPages         int                 // Counter for the number of pages in the chunk
}

/*
// NewTimeChunkWriter creates a new instance of TimeChunkWriter.
func NewTimeChunkWriter() *TimeChunkWriter {
	return &TimeChunkWriter{
		ChunkData: base.NewByteStream(PagesDataPageSize),
	}
}

*/

// Initialize sets up the TimeChunkWriter with measurement details.
func (writer *TimeChunkWriter) Initialize(measurementName string, encoding base.TSEncoding, compressionType base.CompressionType) error {
	// Initialize Chunk Statistics
	stat := statistic.Factory{}
	stats, err := stat.AllocStatistic(base.VECTOR)
	if err != nil {
		return fmt.Errorf("failed to allocate chunk statistic: %w", err)
	}
	writer.ChunkStatistic = stats

	// Initialize Chunk Header
	writer.ChunkHeader = &core.ChunkHeader{
		MeasurementName: measurementName,
		DataType:        base.VECTOR, // Assumed type
		CompressionType: compressionType,
		EncodingType:    encoding,
	}

	// Initialize TimePageWriter
	writer.TimePageWriter = &TimePageWriter{}
	if err := writer.TimePageWriter.Initialize(encoding, compressionType); err != nil {
		return fmt.Errorf("failed to initialize time page writer: %w", err)
	}

	// Initialize First Page Statistic
	firstPageStats, err := stat.AllocStatistic(base.VECTOR)
	if err != nil {
		return fmt.Errorf("failed to allocate statistic for first page: %w", err)
	}
	writer.FirstPageStatistic = firstPageStats

	writer.NumOfPages = 0
	return nil
}

// Destroy releases all the resources held by the TimeChunkWriter.
func (writer *TimeChunkWriter) Destroy() {
	if writer.NumOfPages == 1 {
		writer.freeFirstWriterData()
	}

	if writer.TimePageWriter != nil {
		writer.TimePageWriter.Destroy()
		writer.TimePageWriter = nil
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
	writer.ChunkData.Reset()
}

// Write writes a single timestamp into the chunk.
func (writer *TimeChunkWriter) Write(timestamp int64) error {
	if err := writer.TimePageWriter.Write(timestamp); err != nil {
		return fmt.Errorf("failed to write timestamp: %w", err)
	}
	return nil
}

// GetChunkData retrieves all chunk data.
func (writer *TimeChunkWriter) GetChunkData() *base.ByteStream {
	return writer.ChunkData
}

// GetChunkStatistics retrieves the chunk statistics.
func (writer *TimeChunkWriter) GetChunkStatistics() statistic.Interface {
	return writer.ChunkStatistic
}

// EndEncodeChunk finalizes the encoding for the chunk.
func (writer *TimeChunkWriter) EndEncodeChunk() error {
	// Ensure the final page is sealed before ending the chunk
	err := writer.sealCurrentPage(true)
	if err != nil {
		return err
	}
	writer.ChunkHeader.DataSize = writer.ChunkData.TotalSize
	writer.ChunkHeader.NumOfPages = int32(writer.NumOfPages)
	return nil
}

// GetNumOfPages returns the number of pages in the chunk.
func (writer *TimeChunkWriter) GetNumOfPages() int32 {
	return int32(writer.NumOfPages)
}

// sealCurrentPageIfFull seals the current page if it exceeds any capacity limits.
func (writer *TimeChunkWriter) sealCurrentPageIfFull() error {
	if writer.TimePageWriter.IsPageFull() {
		return writer.sealCurrentPage(false)
	}
	return nil
}

// sealCurrentPage finalizes the current page and handles any chunk-specific logic.
func (writer *TimeChunkWriter) sealCurrentPage(endChunk bool) error {
	var err error

	// Merge current page statistics into chunk statistics
	if err = writer.ChunkStatistic.MergeWith(writer.TimePageWriter.GetStatistic()); err != nil {
		return fmt.Errorf("failed to merge page statistics into chunk statistics: %w", err)
	}

	// Handle logic for the first page
	if writer.NumOfPages == 0 {
		if endChunk {
			// If the current page is the only one, directly write to the chunk
			err = writer.TimePageWriter.WriteToChunk(writer.ChunkData, true, false, true)
			writer.TimePageWriter.DestroyPageData()
			writer.TimePageWriter.Destroy()
		} else {
			// Save first page data for possible reuse (deferred writing)
			err = writer.TimePageWriter.WriteToChunk(writer.ChunkData, true, false, false)
			if err == nil {
				writer.saveFirstPageData(writer.TimePageWriter)
				writer.TimePageWriter.Reset()
			}
		}
	} else {
		if writer.NumOfPages == 1 {
			// Write the first page if the chunk now has more than one page
			if err = writer.writeFirstPageData(writer.ChunkData); err != nil {
				return fmt.Errorf("failed to write first page: %w", err)
			}
			writer.freeFirstWriterData()
		}

		// Write the current page data to the chunk
		err = writer.TimePageWriter.WriteToChunk(writer.ChunkData, false, true, true)
		writer.TimePageWriter.DestroyPageData()
		writer.TimePageWriter.Reset()
	}

	// Increment page count
	writer.NumOfPages++
	return err
}

// saveFirstPageData saves the first page's data for deferred writing.
func (writer *TimeChunkWriter) saveFirstPageData(pageWriter *TimePageWriter) {
	writer.FirstPageData = pageWriter.GetTimeData()
	writer.FirstPageStatistic = pageWriter.GetStatistic().Clone()
}

// writeFirstPageData writes the saved first page's data to the chunk.
func (writer *TimeChunkWriter) writeFirstPageData(chunkData *base.ByteStream) error {
	stat := statistic.Factory{}
	firstPageStats, err := stat.AllocStatistic(base.VECTOR)
	if err != nil {
		return fmt.Errorf("failed to allocate statistic for first page: %w", err)
	}
	err = firstPageStats.SerializeTypedStat(chunkData)
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
func (writer *TimeChunkWriter) freeFirstWriterData() {
	if writer.FirstPageData != nil {
		writer.FirstPageData.Destroy()
		writer.FirstPageData = nil
	}
	if writer.FirstPageStatistic != nil {
		writer.FirstPageStatistic.Reset()
		writer.FirstPageStatistic = nil
	}
}
