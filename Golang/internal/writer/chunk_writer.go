package writer

import (
	"Golang/internal/common/base"
	"Golang/internal/common/core"
	"Golang/internal/common/statistic"
	"errors"
	"fmt"
)

// ChunkWriter is responsible for managing multiple pages and writing them into a chunk.
type ChunkWriter struct {
	DataType           base.TSDataType     // The data type for the chunk
	PageWriter         *PageWriter         // The PageWriter managing individual pages
	ChunkStatistic     statistic.Interface // Overall statistic for the chunk
	ChunkData          *base.ByteStream    // ByteStream where the chunk data is written
	FirstPageData      *PageData           // Data for the first page (if the chunk has one page only)
	FirstPageStatistic statistic.Interface // Statistic for the first page
	ChunkHeader        core.ChunkHeader    // Metadata for the chunk
	NumOfPages         int                 // Total number of pages in the chunk
}

// Initialize sets up the ChunkWriter given measurement metadata.
func (writer *ChunkWriter) Initialize(measurementName string, dataType base.TSDataType, encoding base.TSEncoding, compressionType base.CompressionType) error {
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

	// Initialize the PageWriter
	writer.PageWriter = &PageWriter{}
	err = writer.PageWriter.Initialize(dataType, encoding, compressionType)
	if err != nil {
		return fmt.Errorf("failed to initialize page writer: %w", err)
	}

	// Initialize the ByteStream for the chunk data
	writer.ChunkData, err = base.NewByteStream(OutStreamPageSize)
	if err != nil {
		return fmt.Errorf("failed to create chunk data ByteStream: %w", err)
	}

	writer.NumOfPages = 0
	return nil
}

// Destroy releases resources used by the ChunkWriter.
func (writer *ChunkWriter) Destroy() {
	if writer.NumOfPages == 1 && writer.FirstPageData != nil {
		writer.PageWriter.Destroy()
	}
	if writer.PageWriter != nil {
		writer.PageWriter.Destroy()
		writer.PageWriter = nil
	}
	if writer.ChunkStatistic != nil {
		writer.ChunkStatistic.Reset()
		writer.ChunkStatistic = nil
	}
	if writer.FirstPageStatistic != nil {
		writer.FirstPageStatistic.Reset()
		writer.FirstPageStatistic = nil
	}
	/*
		if writer.ChunkData != nil {
			writer.ChunkData = nil
		}

	*/
	writer.NumOfPages = 0
}

// Write writes a data point to the ChunkWriter.
func (writer *ChunkWriter) Write(timestamp int64, value interface{}) error {
	// Ensure the data type matches
	if writer.PageWriter.DataType != writer.DataType {
		return fmt.Errorf("data type mismatch: %d, %d", base.TSDataType.TSDataTypeToEnum(writer.PageWriter.DataType), base.TSDataType.TSDataTypeToEnum(writer.DataType))
	}

	// Write the data to the current page
	err := writer.PageWriter.Write(timestamp, value)
	if err != nil {
		return err
	}

	// Seal the current page if it is full
	return writer.SealCurrentPageIfFull()
}

// SealCurrentPageIfFull seals the current page if it is full.
func (writer *ChunkWriter) SealCurrentPageIfFull() error {
	// Use the PageWriter's statistics to check if the current page is full
	if writer.PageWriter.PointCount >= int(OutStreamPageSize) {
		return writer.SealCurrentPage(false)
	}
	return nil

}

// SealCurrentPage seals the current page and adds it to the chunk.
func (writer *ChunkWriter) SealCurrentPage(endChunk bool) error {
	// pagewriter := writer.PageWriter
	// Merge the current page's statistics into the chunk's statistics
	err := writer.ChunkStatistic.MergeWith(writer.PageWriter.Statistic)
	if err != nil {
		return fmt.Errorf("failed to merge page statistics with chunk statistics: %w", err)
	}

	if writer.NumOfPages == 0 {
		if endChunk {
			// Write the entire page into the chunk if this is the only page
			err := writer.PageWriter.WriteToChunk(writer.ChunkData, true, false, true)
			if err != nil {
				return fmt.Errorf("failed to write page to chunk: %w", err)
			}
		} else {
			// Save the first page for potential later writing
			err := writer.PageWriter.WriteToChunk(writer.ChunkData, true, false, false)
			if err != nil {
				return err
			}
			err = writer.SaveFirstPageData()
			if err != nil {
				return err
			}
		}
	} else {
		// If this is the first page, flush its data to the chunk
		if writer.NumOfPages == 1 {
			// Write the first page's data
			if err := writer.WriteFirstPageData(); err != nil {
				return err
			}
		}

		// Write the current page's data
		if err := writer.PageWriter.WriteToChunk(writer.ChunkData, true, true, true); err != nil {
			return err
		}
		writer.PageWriter.Reset()
		// pagewriter.Destroy()
	}

	writer.NumOfPages++
	return nil
}

// SaveFirstPageData saves the first page's data and statistics in memory
func (writer *ChunkWriter) SaveFirstPageData() error {
	pageData, err := writer.PageWriter.FinalizePage()
	if err != nil {
		return fmt.Errorf("failed to finalize first page: %w", err)
	}
	writer.FirstPageData = pageData

	// Clone current page statistics into the first page
	writer.FirstPageStatistic = writer.PageWriter.Statistic.Clone()
	return nil
}

// WriteFirstPageData writes the saved first page data to the chunk.
func (writer *ChunkWriter) WriteFirstPageData() error {
	if writer.FirstPageData == nil {
		return errors.New("no first page data to write")
	}

	if writer.FirstPageStatistic == nil {
		return errors.New("no first page statistic to serialize")
	}

	// Serialize the statistics from the first page
	err := writer.FirstPageStatistic.SerializeTypedStat(writer.ChunkData)
	if err != nil {
		return fmt.Errorf("failed to serialize first page statistics: %w", err)
	}

	// Write the data for the first page to the chunk
	err = writer.ChunkData.WriteBuf(writer.FirstPageData.CompressedBuf, writer.FirstPageData.CompressedSize)
	if err != nil {
		return err
	}
	return nil
}

func (writer *ChunkWriter) EndEncodeChunk() error {
	if writer.PageWriter.PointCount > 0 {
		err := writer.SealCurrentPage(true)
		if err != nil {
			return err
		}
		writer.ChunkHeader.DataSize = writer.ChunkData.TotalSize
		writer.ChunkHeader.NumOfPages = int32(writer.NumOfPages)
	} else {
		return errors.New("no data to write")
	}
	return nil
}
