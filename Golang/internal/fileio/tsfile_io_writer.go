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
	"Golang/internal/reader"
	"Golang/internal/utils"
	"fmt"
	"math"
)

// TsFileIOWriter handles the operations required for writing to .ts files.
type TsFileIOWriter struct {
	writeStream         *base.ByteStream                // Stream for writing data
	writeStreamConsumer *base.ByteStreamConsumer        // Consumer for the write stream
	curChunkMeta        *core.ChunkMeta                 // Current chunk metadata
	curChunkGroupMeta   *core.ChunkGroupMeta            // Current chunk group metadata
	chunkMetaCount      int                             // Count of chunks in the current group
	chunkGroupMetaList  []*core.ChunkGroupMeta          // All chunk group metadata
	usePrevAllocCgm     bool                            // Flag for reusing allocation
	curDeviceName       string                          // Current device name
	file                *WriteFile                      // The write file (.ts file)
	tsTimeIndexVector   []core.TimeseriesTimeIndexEntry // Time index vector for timeseries
	writeFileCreated    bool                            // Whether the write file has been created
}

// Constants
const (
	WriteStreamPageSize = 512
)

// NewTsFileIOWriter creates a new TsFileIOWriter instance.
func NewTsFileIOWriter() *TsFileIOWriter {
	// Create a new ByteStream for writing with a fixed page size.
	writeStream := &base.ByteStream{
		PageSize: WriteStreamPageSize,
		Pages:    []*base.Page{},
	}

	writeStreamConsumer := &base.ByteStreamConsumer{
		Host: writeStream,
	}

	return &TsFileIOWriter{
		writeStream:         writeStream,
		writeStreamConsumer: writeStreamConsumer, // Consumer uses the same ByteStream by default
		chunkMetaCount:      0,
		chunkGroupMetaList:  []*core.ChunkGroupMeta{},
		writeFileCreated:    false,
	}

}

// Init initializes the TsFileIOWriter with a WriteFile instance.
func (io *TsFileIOWriter) Init(writeFile *WriteFile) error {

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

	// Clean up allocated resources and metadata objects
	err := io.file.CloseFile()
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
func (io *TsFileIOWriter) EndFlushChunk() error {
	io.chunkMetaCount++
	// stat := statistic.Factory{}
	// Append the current chunk meta to the chunk group list
	io.curChunkGroupMeta.Push(io.curChunkMeta)
	return nil
}

// EndFlushChunkGroup finalizes the current chunk group.
func (io *TsFileIOWriter) EndFlushChunkGroup(isAligned bool) error {

	// Ensure all chunks in the group are written and metadata updated
	return nil
}

func (io *TsFileIOWriter) EndFile() error {
	wf := WriteFile{}

	// Write log index range
	if err := io.WriteLogIndexRange(); err != nil {
		fmt.Printf("writer range index error: %v\n", err)
		return err
	}

	// Write file index
	if err := io.WriteFileIndex(); err != nil {
		fmt.Printf("writer file index error: %v\n", err)
		return err
	}

	// Write file footer
	if err := io.WriteFileFooter(); err != nil {
		fmt.Printf("writer file footer error: %v\n", err)
		return err
	}

	// Sync file
	if err := wf.SyncFile(); err != nil {
		fmt.Printf("sync file error: %v\n", err)
		return err
	}

	// Close file
	if err := wf.CloseFile(); err != nil {
		return err
	}

	// Successfully completed
	return nil
}

func (io *TsFileIOWriter) WriteLogIndexRange() error {
	minPlanIndex := 0
	maxPlanIndex := 0
	err := io.WriteByte(core.OPERATION_INDEX_RANGE)
	if err != nil {
		return err
	}
	serial := base.SerializationUtil{}
	err = serial.WriteUint64(uint64(minPlanIndex), io.writeStream)
	if err != nil {
		return err
	}
	err = serial.WriteUint64(uint64(maxPlanIndex), io.writeStream)
	if err != nil {
		return err
	}
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
		buffer, length, err := io.writeStreamConsumer.GetNextBuf(io.writeStream)
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
	err := serial.WriteVarUint(uint32(len(str)), io.writeStream)
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

func (io *TsFileIOWriter) WriteFileIndex() error {
	var (
		TsfileMeta            *core.TsFileMeta
		entryCountInCurDevice int
		prevDeviceName        string
		curIndexNode          *core.MetaIndexNode
		curIndexNodeQueue     []*core.MetaIndexNode
		deviceMap             map[string]*core.MetaIndexNode
		writingMM             *FileIndexWritingMemManager
		filter                *reader.BloomFilter
	)

	tsmIter := core.TSMIterator{ChunkGroupMetaList: io.chunkGroupMetaList}
	/*
		err := tsmIter.TSMInit()
		if err != nil {
			return err
		}

	*/

	// Write the separator marker.
	if err := io.WriteSeperatorMarker(TsfileMeta); err != nil {
		return fmt.Errorf("failed to write separator marker: %w", err)
	}

	// Initialize the bloom filter.
	if err := filter.Init(0.05, int(io.GetPathCount(io.chunkGroupMetaList))); err != nil {
		return fmt.Errorf("failed to initialize Bloom filter: %w", err)
	}

	// Initialize the TSM iterator.
	if err := tsmIter.TSMInit(); err != nil {
		return fmt.Errorf("failed to initialize TSM iterator: %w", err)
	}

	for tsmIter.HasNext() {
		var (
			deviceName      string
			measurementName string
			tsIndex         *core.TimeseriesIndex
		)

		tsIndex.Reset()

		// Get the next metadata entry.

		var err error
		if deviceName, measurementName, tsIndex, err = tsmIter.GetNext(); err != nil {
			return fmt.Errorf("failed to get next TSM entry: %w", err)
		}

		// Handle a new device entry.
		if prevDeviceName != deviceName {
			if prevDeviceName != "" {
				// Add current index node to queue.
				if err := io.AddCurIndexNodeToQueue(curIndexNode, &curIndexNodeQueue); err != nil {
					return fmt.Errorf("failed to add current index node to queue: %w", err)
				}

				// Add device node to device map.
				if err := io.AddDeviceNode(deviceMap, prevDeviceName, curIndexNodeQueue, writingMM); err != nil {
					return fmt.Errorf("failed to add device node: %w", err)
				}
			}

			// Allocate new index nodes for the new device.
			if curIndexNodeQueue, err = io.AllocMetaIndexNodeQueue(writingMM); err != nil {
				return fmt.Errorf("failed to allocate meta index node queue: %w", err)
			}
			if curIndexNode, err = io.AllocAndInitMetaIndexNode(writingMM, core.LEAF_MEASUREMENT); err != nil {
				return fmt.Errorf("failed to allocate and initialize meta index node: %w", err)
			}

			// Update the previous device name and reset the entry counter.
			prevDeviceName = deviceName
			entryCountInCurDevice = 0
		}

		// Handle index entry creation.
		if entryCountInCurDevice%utils.ConfigValue.MaxDegreeOfIndexNode == 0 { // 256 is config of maxDegreeOfIndexNode
			if curIndexNode.IsFull() {
				if err := io.AddCurIndexNodeToQueue(curIndexNode, &curIndexNodeQueue); err != nil {
					return fmt.Errorf("failed to add current index node to queue: %w", err)
				}
				if curIndexNode, err = io.AllocAndInitMetaIndexNode(writingMM, core.LEAF_MEASUREMENT); err != nil {
					return fmt.Errorf("failed to allocate and initialize meta index node: %w", err)
				}
			}

			var metaIndexEntry *core.MetaIndexEntry
			if metaIndexEntry, err = io.AllocAndInitMetaIndexEntry(writingMM, measurementName); err != nil {
				return fmt.Errorf("failed to allocate and initialize meta index entry: %w", err)
			}
			if err := curIndexNode.PushEntry(metaIndexEntry); err != nil {
				return fmt.Errorf("failed to push entry to current index node: %w", err)
			}
		}

		// Add timeseries index to Bloom filter and serialize it.
		if tsIndex.GetDataType() != base.VECTOR {
			filter.AddPathEntry(deviceName, measurementName)

		}
		if err := tsIndex.SerializeTo(io.writeStream); err != nil {
			return fmt.Errorf("failed to serialize timeseries index: %w", err)
		}

		entryCountInCurDevice++
	}

	// Finalize processing for the last device.
	if curIndexNode != nil && curIndexNodeQueue != nil {
		if err := io.AddCurIndexNodeToQueue(curIndexNode, &curIndexNodeQueue); err != nil {
			return fmt.Errorf("failed to add current index node to queue: %w", err)
		}
		if err := io.AddDeviceNode(deviceMap, prevDeviceName, curIndexNodeQueue, writingMM); err != nil {
			return fmt.Errorf("failed to add device node: %w", err)
		}
	}

	// Build the device index level.
	var deviceIndexRootNode *core.MetaIndexNode
	if err := io.BuildDeviceLevel(deviceMap, &deviceIndexRootNode, writingMM); err != nil {
		return fmt.Errorf("failed to build device level: %w", err)
	}

	// Write the file metadata to the stream.
	tsFileMeta := core.TsFileMeta{}
	tsFileMeta.IndexNode = deviceIndexRootNode
	tsFileMeta.MetaOffset = int64(io.CurrentFilePosition())

	tsFileMetaOffset := io.CurrentFilePosition()
	if err := tsFileMeta.Serialize(io.writeStream); err != nil {
		return fmt.Errorf("failed to serialize file metadata: %w", err)
	}

	// Write the Bloom filter.
	if err := filter.SerializeTo(io.writeStream); err != nil {
		return fmt.Errorf("failed to serialize Bloom filter: %w", err)
	}

	// Write the metadata size.
	tsFileMetaEndOffset := io.CurrentFilePosition()
	metaSize := uint32(tsFileMetaEndOffset - tsFileMetaOffset)
	serial := base.SerializationUtil{}
	if err := serial.WriteUint32(metaSize, io.writeStream); err != nil {
		return fmt.Errorf("failed to write metadata size: %w", err)
	}

	return nil
}

func (io *TsFileIOWriter) BuildDeviceLevel(deviceMap map[string]*core.MetaIndexNode, retRoot **core.MetaIndexNode, wmm *FileIndexWritingMemManager) error {
	var err error

	// MetaIndexNode queue for managing intermediate nodes
	var nodeQueue []*core.MetaIndexNode

	var curIndexNode *core.MetaIndexNode
	if curIndexNode, err = io.AllocAndInitMetaIndexNode(wmm, (*retRoot).NodeType); err != nil {
		return fmt.Errorf("failed to allocate and initialize meta index node: %w", err)
	}

	// Iterate over the device map
	for deviceName, deviceNode := range deviceMap {
		// Check if the current node is full
		if curIndexNode.IsFull() {
			// Set the end offset of the current node
			curIndexNode.EndOffset = int64(io.curFilePosition())

			// Push the current node to the node queue
			nodeQueue = append(nodeQueue, curIndexNode)

			// Allocate and initialize a new MetaIndexNode
			if curIndexNode, err = io.AllocAndInitMetaIndexNode(wmm, curIndexNode.NodeType); err != nil {
				return fmt.Errorf("failed to allocate and initialize meta index node: %w", err)
			}
		}

		// Allocate and initialize a MetaIndexEntry
		var entry *core.MetaIndexEntry
		if entry, err = io.AllocAndInitMetaIndexEntry(wmm, deviceName); err != nil {
			return fmt.Errorf("failed to allocate and initialize meta index entry: %w", err)
		}

		// Serialize the device map into the write stream
		if err = deviceNode.Serialize(io.writeStream); err != nil {
			return fmt.Errorf("failed to serialize device node: %w", err)
		}

		// Push the entry into the current MetaIndexNode
		if err = curIndexNode.PushEntry(entry); err != nil {
			return fmt.Errorf("failed to push entry to curIndexNode: %w", err)
		}
	}

	// Check if the final current node is not empty
	if !curIndexNode.IsEmpty() {
		// Set the end offset and push it to the queue
		curIndexNode.EndOffset = int64(io.curFilePosition())
		nodeQueue = append(nodeQueue, curIndexNode)
	}

	// Generate root or set the root node
	if len(nodeQueue) > 0 {
		if *retRoot, err = io.GenerateRoot(nodeQueue, curIndexNode.NodeType, wmm); err != nil {
			return fmt.Errorf("failed to generate root: %w", err)
		}
	} else {
		*retRoot = curIndexNode
		(*retRoot).EndOffset = int64(io.curFilePosition())
		(*retRoot).NodeType = curIndexNode.NodeType
	}

	return nil
}

func (io *TsFileIOWriter) WriteSeperatorMarker(tsfIleMeta *core.TsFileMeta) error {
	tsfIleMeta.MetaOffset = int64(io.CurrentFilePosition())
	err := io.WriteByte(core.SEPARATOR_MARKER)
	if err != nil {
		return err
	}
	return nil
}

func (io *TsFileIOWriter) CurrentFilePosition() uint32 {
	return io.writeStream.TotalSize
}

func (io *TsFileIOWriter) AllocAndInitMetaIndexEntry(wmm *FileIndexWritingMemManager, name string) (*core.MetaIndexEntry, error) {
	// Allocate MetaIndexEntry
	entry := &core.MetaIndexEntry{
		Name:   name,                            // Copy the name
		Offset: int64(io.CurrentFilePosition()), // Get file position
	}
	return entry, nil
}

func (io *TsFileIOWriter) AllocAndInitMetaIndexNode(wmm *FileIndexWritingMemManager, nodeType core.MetaIndexNodeType) (*core.MetaIndexNode, error) {
	// Allocate MetaIndexNode
	node := &core.MetaIndexNode{
		NodeType:  nodeType,
		EndOffset: 0,
	}

	// Add it to the memory manager's list of all index nodes
	wmm.AllIndexNodes = append(wmm.AllIndexNodes, node)

	return node, nil
}

func (io *TsFileIOWriter) AddCurIndexNodeToQueue(node *core.MetaIndexNode, queue *[]*core.MetaIndexNode) error {
	if node == nil || queue == nil {
		return fmt.Errorf("invalid node or queue")
	}

	// Set the node's end offset to the current file position
	node.EndOffset = int64(io.CurrentFilePosition())

	// Add the node to the queue
	*queue = append(*queue, node)

	return nil
}

func (io *TsFileIOWriter) AllocMetaIndexNodeQueue(wmm *FileIndexWritingMemManager) ([]*core.MetaIndexNode, error) {
	// Allocate a new queue (slice of MetaIndexNode pointers)
	var queue []*core.MetaIndexNode

	return queue, nil
}

func (io *TsFileIOWriter) AddDeviceNode(deviceMap map[string]*core.MetaIndexNode, deviceName string, measurementIndexNodeQueue []*core.MetaIndexNode, wmm *FileIndexWritingMemManager) error {
	if len(measurementIndexNodeQueue) == 0 {
		return fmt.Errorf("measurementIndexNodeQueue is empty")
	}

	// Check if the device node already exists
	if _, exists := deviceMap[deviceName]; exists {
		return fmt.Errorf("device node already exists: %s", deviceName)
	}

	// Generate the root index node
	root, err := io.GenerateRoot(measurementIndexNodeQueue, core.INTERNAL_MEASUREMENT, wmm)
	if err != nil {
		return fmt.Errorf("failed to generate root for device node: %w", err)
	}

	// Add the root index node to the device map
	deviceMap[deviceName] = root

	return nil
}

func (io *TsFileIOWriter) GenerateRoot(nodeQueue []*core.MetaIndexNode, nodeType core.MetaIndexNodeType, wmm *FileIndexWritingMemManager) (*core.MetaIndexNode, error) {
	if len(nodeQueue) == 0 {
		return nil, fmt.Errorf("nodeQueue is empty")
	}

	// If the queue has only one node, it's already the root
	if len(nodeQueue) == 1 {
		return nodeQueue[0], nil
	}

	// Clone the node queue into a working list
	listX := append([]*core.MetaIndexNode{}, nodeQueue...)
	var listY []*core.MetaIndexNode

	var curIndexNode *core.MetaIndexNode

	// Start creating intermediate levels
	for {
		// Clear the next level list
		listY = listY[:0]

		for _, iterNode := range listX {
			// Allocate and initialize a new MetaIndexEntry for the current node
			name, err := iterNode.GetFirstChildName()
			if err != nil {
				return nil, fmt.Errorf("failed to get first child name: %w", err)
			}

			entry, err := io.AllocAndInitMetaIndexEntry(wmm, name)
			if err != nil {
				return nil, fmt.Errorf("failed to allocate MetaIndexEntry: %w", err)
			}

			// If the current index node is full, push it to the next level
			if curIndexNode != nil && curIndexNode.IsFull() {
				curIndexNode.EndOffset = int64(io.CurrentFilePosition())
				listY = append(listY, curIndexNode)

				// Allocate a new MetaIndexNode
				curIndexNode, err = io.AllocAndInitMetaIndexNode(wmm, nodeType)
				if err != nil {
					return nil, fmt.Errorf("failed to allocate new MetaIndexNode: %w", err)
				}
			}

			// Add the entry to the current index node
			if curIndexNode == nil {
				curIndexNode, err = io.AllocAndInitMetaIndexNode(wmm, nodeType)
				if err != nil {
					return nil, fmt.Errorf("failed to allocate initial MetaIndexNode: %w", err)
				}
			}

			if err := curIndexNode.PushEntry(entry); err != nil {
				return nil, fmt.Errorf("failed to push entry to MetaIndexNode: %w", err)
			}
		}

		// Process the last partially-full node
		if curIndexNode != nil && !curIndexNode.IsEmpty() {
			curIndexNode.EndOffset = int64(io.CurrentFilePosition())
			listY = append(listY, curIndexNode)
			curIndexNode = nil
		}

		// If the next level has only one node, it's the root
		if len(listY) == 1 {
			return listY[0], nil
		}

		// Swap lists for the next iteration
		listX, listY = listY, listX
	}
}

func (io *TsFileIOWriter) CloneNodeList(src []*core.MetaIndexNode) ([]*core.MetaIndexNode, error) {
	// Clone a slice of MetaIndexNodes
	dest := make([]*core.MetaIndexNode, len(src))
	copy(dest, src)
	return dest, nil
}

func (io *TsFileIOWriter) WriteFileFooter() error {
	err := io.WriteBuf(core.MAGIC_STRING_TSFILE, core.MAGIC_STRING_TSFILE_LEN)
	if err != nil {
		return err
	}
	err = io.FlushStreamToFile()
	if err != nil {
		return err
	}
	return nil
}

func (io *TsFileIOWriter) WriteBuf(tsfile string, tsfileLen int) error {
	err := io.writeStream.WriteBuf([]byte(tsfile), uint32(tsfileLen))
	if err != nil {
		return err
	}
	return nil
}

func (io *TsFileIOWriter) GetPathCount(cgmList []*core.ChunkGroupMeta) int32 {
	var pathCount int32 = 0
	var prevMeasurement string

	// Iterate through the ChunkGroupMeta list
	for _, cgm := range cgmList {
		// Iterate through the ChunkMeta list in each ChunkGroupMeta
		for _, cm := range cgm.ChunkMetaList {
			// Check if the measurement is different from the previous one
			if cm.MeasurementName != prevMeasurement {
				pathCount++
				prevMeasurement = cm.MeasurementName
			}
		}
	}

	return pathCount
}

// FileIndexWritingMemManager manages memory allocations and tracks allocated MetaIndexNodes.
type FileIndexWritingMemManager struct {
	ByteStream    *base.ByteStream      // Stream for managing allocations
	AllIndexNodes []*core.MetaIndexNode // List of all MetaIndexNodes
}

// NewFileIndexWritingMemManager initializes a new memory manager using ByteStream.
func NewFileIndexWritingMemManager(pageSize uint32) (*FileIndexWritingMemManager, error) {
	// Initialize a ByteStream to replace PageArena
	byteStream, err := base.NewByteStream(pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ByteStream: %w", err)
	}

	return &FileIndexWritingMemManager{
		ByteStream:    byteStream,
		AllIndexNodes: []*core.MetaIndexNode{},
	}, nil
}

// AddIndexNode adds a MetaIndexNode to the list of tracked nodes.
func (m *FileIndexWritingMemManager) AddIndexNode(node *core.MetaIndexNode) {
	m.AllIndexNodes = append(m.AllIndexNodes, node)
}

// Free releases resources for all tracked MetaIndexNodes.
func (m *FileIndexWritingMemManager) Free() {
	// Clear all tracked nodes
	for _, node := range m.AllIndexNodes {
		if node != nil {
			node.Children = nil // Allow garbage collection
		}
	}
	m.AllIndexNodes = nil
}
