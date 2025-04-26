package writer

import (
	"errors"
	"fmt"
)

// writes from page to chunks (time chunk and value chunks) and then cleans the chunk-level memory
// time chunk writes timestamps, value chunk writes measurements values
// sent to tsfile_writer to combine multiple chunks into a chunk group

// MergedChunkWriter handles both time and value chunks
type MergedChunkWriter struct {
	ChunkType      string      // "time" or "value"
	PageWriter     *PageWriter // Unified PageWriter for either timestamps or values
	ChunkBuffer    []byte      // Buffer for storing chunk data
	NumOfPages     int         // Number of pages in the chunk
	ChunkStatistic *Statistic  // Statistics for the chunk
}

// NewMergedChunkWriter initializes a new MergedChunkWriter
// chunkType can be "time" or "value".
func NewMergedChunkWriter(chunkType string, dataType string, encoding string, compression string) (*MergedChunkWriter, error) {
	// Validate chunk type
	if chunkType != "time" && chunkType != "value" {
		return nil, errors.New("invalid chunk type: must be 'time' or 'value'")
	}

	// Create and initialize the PageWriter
	pageWriter, err := NewPageWriter(dataType, encoding, compression)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PageWriter: %v", err)
	}

	return &MergedChunkWriter{
		ChunkType:   chunkType,
		PageWriter:  pageWriter,
		ChunkBuffer: make([]byte, 0),
		NumOfPages:  0,
		// Initialize statistics factory for the specified data type
		ChunkStatistic: NewStatistic(dataType),
	}, nil
}

// Write writes timestamp or value to the chunk
func (mcw *MergedChunkWriter) Write(data interface{}) error {
	// Use PageWriter to write data (timestamp or value)
	switch mcw.ChunkType {
	case "time":
		// Cast data to int64 for timestamps
		timestamp, ok := data.(int64)
		if !ok {
			return errors.New("invalid data type for timestamp chunk")
		}
		return mcw.PageWriter.WriteTimestamp(timestamp)
	case "value":
		// WriteBuf value (should support multiple data types in PageWriter)
		return mcw.PageWriter.WriteValue(data)
	default:
		return errors.New("unknown chunk type")
	}
}

// FinalizeChunk finalizes the chunk and writes data to ChunkBuffer
func (mcw *MergedChunkWriter) FinalizeChunk(endChunk bool) error {
	// Merge statistics from the PageWriter into the ChunkStatistic
	if err := mcw.ChunkStatistic.MergeWith(mcw.PageWriter.GetStatistic()); err != nil {
		return err
	}

	// Handle the case of only one page in the chunk
	if mcw.NumOfPages == 0 {
		if endChunk {
			// Finalize chunk with one page
			err := mcw.PageWriter.WriteToChunk(mcw.ChunkBuffer, true, false, true)
			mcw.PageWriter.DestroyPageData()
			mcw.PageWriter.Destroy()
			return err
		} else {
			// Save the first page's data temporarily
			err := mcw.PageWriter.WriteToChunk(mcw.ChunkBuffer, true, false, false)
			if err == nil {
				mcw.PageWriter.Reset() // Prepare for the next page
			}
			return err
		}
	} else {
		// Flush the current page and handle multiple pages
		if mcw.NumOfPages == 1 {
			// WriteBuf the first saved page now
			mcw.FlushFirstPage()
		}
		// WriteBuf the current page
		err := mcw.PageWriter.WriteToChunk(mcw.ChunkBuffer, false, true, true)
		mcw.PageWriter.Reset()
		return err
	}
}

// FlushFirstPage writes the temporarily saved first page into the ChunkBuffer
func (mcw *MergedChunkWriter) FlushFirstPage() error {
	// TODO: Implement logic for saving and writing the first page's data
	return nil
}

// Destroy cleans up resources for the MergedChunkWriter
func (mcw *MergedChunkWriter) Destroy() {
	mcw.PageWriter.Destroy()
	if mcw.ChunkStatistic != nil {
		mcw.ChunkStatistic.Destroy()
	}
	mcw.ChunkBuffer = nil
}
