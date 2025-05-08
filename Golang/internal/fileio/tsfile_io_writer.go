package fileio

/*
Handle file writing index, chunk, flush
writing into chunks for buffer memory into stream
Chunk gets passed to writefile to be flushed into disk
*/

import (
	"Golang/internal/common/base"
	"Golang/internal/common/core"
	"Golang/internal/common/statistic"
	"Golang/internal/utils"
	"fmt"
	"math"
	"sync"
)

// TsFileIOWriter handles the operations required for writing to .ts files.
type TsFileIOWriter struct {
	writeStream         *base.ByteStream                // Stream for writing data
	writeStreamConsumer *base.ByteStream                // Consumer for the write stream
	curChunkMeta        *core.ChunkMeta                 // Current chunk metadata
	curChunkGroupMeta   *core.ChunkGroupMeta            // Current chunk group metadata
	chunkMetaCount      int                             // Count of chunks in the current group
	chunkGroupMetaList  []*core.ChunkGroupMeta          // All chunk group metadata
	usePrevAllocCgm     bool                            // Flag for reusing allocation
	curDeviceName       string                          // Current device name
	file                *WriteFile                      // The write file (.ts file)
	tsTimeIndexVector   []core.TimeseriesTimeIndexEntry // Time index vector for timeseries
	writeFileCreated    bool                            // Whether the write file has been created
	lock                sync.Mutex                      // To ensure thread safety
}

// Constants
const (
	WriteStreamPageSize = 512 // FIXME: Adjust as per requirements
)

// NewTsFileIOWriter creates a new TsFileIOWriter instance.
func NewTsFileIOWriter() *TsFileIOWriter {
	// Create a new ByteStream for writing with a fixed page size.
	writeStream := &base.ByteStream{
		PageSize: WriteStreamPageSize,
		Pages:    []*base.Page{},
	}

	return &TsFileIOWriter{
		writeStream:         writeStream,
		writeStreamConsumer: writeStream, // Consumer uses the same ByteStream by default
		chunkMetaCount:      0,
		chunkGroupMetaList:  []*core.ChunkGroupMeta{},
		writeFileCreated:    false,
	}

}

// Init initializes the TsFileIOWriter with a WriteFile instance.
func (io *TsFileIOWriter) Init(writeFile *WriteFile) error {
	io.lock.Lock()
	defer io.lock.Unlock()

	if io.writeFileCreated {
		return fmt.Errorf("write file already initialized")
	}

	// Set up the WriteFile
	io.file = writeFile
	io.writeFileCreated = true

	// Initialize other metadata as needed here...
	return nil
}

// Destroy cleans up resources used by TsFileIOWriter.
func (io *TsFileIOWriter) Destroy() {
	io.lock.Lock()
	defer io.lock.Unlock()

	// Clean up allocated resources and metadata objects
	err := io.file.Close()
	if err != nil {
		return
	} // Ensure the file is properly closed
	io.file = nil

	io.writeFileCreated = false
	// Reset all fields
	io.curChunkMeta = nil
	io.curChunkGroupMeta = nil
	io.chunkGroupMetaList = nil
	io.tsTimeIndexVector = nil
	io.chunkMetaCount = 0
	io.curDeviceName = ""

}

// StartFile initializes the writing process for the TsFile.
func (io *TsFileIOWriter) StartFile() error {
	// Write the magic string to the file
	if err := io.writeStream.WriteBuf([]uint8(core.MAGIC_STRING_TSFILE), uint32(core.MAGIC_STRING_TSFILE_LEN)); err != nil {
		return fmt.Errorf("%w: writing magic string", err)
	}

	// Write the version number
	if err := io.WriteByte(core.VERSION_NUM_BYTE); err != nil {
		return fmt.Errorf("%w: writing version number", err)
	}

	// Flush the stream data to file
	if err := io.FlushStreamToFile(); err != nil {
		return fmt.Errorf("%w: flushing to file", err)
	}
	return nil
}

// StartFlushChunkGroup begins flushing a new chunk group into the file.
func (io *TsFileIOWriter) StartFlushChunkGroup(deviceName string, isAligned bool) error {
	// Write the marker for a chunk group header
	if err := io.WriteByte(core.CHUNK_GROUP_HEADER_MARKER); err != nil {
		return fmt.Errorf("%w: writing chunk group header marker", err)
	}

	// Write the device name
	if err := io.WriteString(deviceName); err != nil {
		return fmt.Errorf("%w: writing device name", err)
	}

	// Set the current device name
	io.curDeviceName = deviceName

	// Check for a reusable chunk group metadata object
	io.usePrevAllocCgm = false
	for _, cgMeta := range io.chunkGroupMetaList {
		if cgMeta.DeviceName == deviceName {
			io.usePrevAllocCgm = true
			io.curChunkGroupMeta = cgMeta
			break
		}
	}

	// If no reusable metadata object is found, create a new one
	if !io.usePrevAllocCgm {
		// Create a new ChunkGroupMeta instance and initialize it
		newMeta := &core.ChunkGroupMeta{}
		if err := core.NewChunkGroupMeta(deviceName); err != nil {
			return fmt.Errorf("%v: initializing new ChunkGroupMetadata", err)
		}

		io.curChunkGroupMeta = newMeta
		io.chunkGroupMetaList = append(io.chunkGroupMetaList, newMeta)
	}

	return nil

}

// StartFlushChunk starts the flush process for a chunk with specific metadata.
func (io *TsFileIOWriter) StartFlushChunk(
	chunkData *base.ByteStream,
	measurementName string,
	dataType base.TSDataType,
	encoding base.TSEncoding,
	compression base.CompressionType,
	numOfPages int32,
	TsID utils.TsID,
) error {
	const mask = 0 // For common chunk

	// Step 1: Record chunk meta
	if io.curChunkMeta != nil {
		return fmt.Errorf("current chunk metadata is not nil")
	}

	// Allocate memory for chunk metadata and statistics creation
	curChunkMeta := &core.ChunkMeta{}
	StatisticFactory := statistic.Factory{}
	chunkStatistic, err := StatisticFactory.AllocStatistic(dataType)
	if err != nil {
		return fmt.Errorf("failed to create statistics: %w", err)
	}

	// Initialize chunk metadata
	err = curChunkMeta.Initialize(
		measurementName,
		dataType,
		int64(io.curFilePosition()), // offsetOfHeader
		chunkStatistic,
		TsID,
		mask,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize chunk metadata: %w", err)
	}
	io.curChunkMeta = curChunkMeta

	var chunkTy byte = 0
	if numOfPages <= 1 {
		chunkTy = core.ONLY_ONE_PAGE_CHUNK_HEADER_MARKER
	} else {
		chunkTy = core.CHUNK_HEADER_MARKER
	}

	// Step 2: Serialize chunk header to writeStream
	chunkHeader := core.ChunkHeader{
		MeasurementName: measurementName,
		DataSize:        chunkData.TotalSize, // Retrieved from ByteStream
		DataType:        dataType,
		CompressionType: compression,
		EncodingType:    encoding,
		NumOfPages:      numOfPages,
		ChunkType:       chunkTy,
	}
	err = chunkHeader.SerializeTo(io.writeStream)
	if err != nil {
		return err
	}

	return nil

}

// FlushChunk flushes the contents of the current chunk to the file.
func (io *TsFileIOWriter) FlushChunk(chunkData *base.ByteStream) error {
	err := io.WriteChunkData(chunkData)
	if err != nil {
		return err
	}
	err = io.FlushStreamToFile()
	if err != nil {
		return err
	}
	return nil
}

func (io *TsFileIOWriter) WriteChunkData(chunkData *base.ByteStream) error {
	err := chunkData.MergeByteStream(io.writeStream, chunkData, true)
	if err != nil {
		return err
	}
	return nil
}

// EndFlushChunk finalizes and closes the current chunk.
func (io *TsFileIOWriter) EndFlushChunk(chunkStatistics *statistic.Interface) error {
	io.chunkMetaCount++
	stat := statistic.Factory{}
	err := stat.CloneStatistic(io.curChunkMeta, io.curChunkMeta, io.curChunkMeta.DataType)
	if err != nil {
		return err
	}
	return nil
}

// EndFlushChunkGroup finalizes the current chunk group.
func (io *TsFileIOWriter) EndFlushChunkGroup(isAligned bool) error {
	io.lock.Lock()
	defer io.lock.Unlock()

	// Ensure all chunks in the group are written and metadata updated
	return nil
}

// EndFile finalizes the TsFile writing process.
func (io *TsFileIOWriter) EndFile() error {
	io.lock.Lock()
	defer io.lock.Unlock()

	// Flush any remaining data and write final file metadata
	return nil
}

// GetFilePath returns the path of the associated file.
func (io *TsFileIOWriter) GetFilePath() string {
	if io.file == nil {
		return ""
	}
	return io.file.GetFilePath()
}

func (io *TsFileIOWriter) WriteByte(written byte) error {
	serial := base.SerializationUtil{}
	err := serial.WriteUint8(written, io.writeStream)
	if err != nil {
		return err
	}
	return nil
}

func (io *TsFileIOWriter) FlushStreamToFile() error {
	for {
		// Retrieve the next buffer from the writeStreamConsumer
		buffer, length, err := io.writeStreamConsumer.GetNextBuffer()
		if err != nil {
			// If there's an error retrieving the buffer, stop the process
			return fmt.Errorf("failed to get next buffer: %w", err)
		}

		// If no buffer is available (end of stream), break the loop
		if buffer == nil {
			break
		}

		// Write the buffer content to the file
		if err = io.file.Write(buffer, length); err != nil {
			return fmt.Errorf("failed to write buffer to file: %w", err)
		}
	}

	// Purge previous pages in the writeStream to free memory
	io.writeStream.PurgePrevPages(math.MaxInt32)

	return nil
}

func (io *TsFileIOWriter) WriteString(str string) error {
	serial := base.SerializationUtil{}
	err := serial.WriteString(str, io.writeStream)
	if err != nil {
		return err
	}
	err = io.writeStream.WriteBuf([]byte(str), uint32(len(str)))
	if err != nil {
		return err
	}
	return nil
}

func (io *TsFileIOWriter) curFilePosition() uint32 {
	return io.writeStream.TotalSize
}
