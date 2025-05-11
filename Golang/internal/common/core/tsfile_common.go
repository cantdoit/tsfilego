package core

import (
	"Golang/internal/common/base"
	"Golang/internal/common/statistic"
	"Golang/internal/reader"
	"Golang/internal/utils"
	"errors"
	"sort"
	"strconv"
	"time"
)

// TsFileID type alias for consistency
type TsFileID int64

// Constants defined in header file
const (
	MAGIC_STRING_TSFILE               = "TsFile"
	MAGIC_STRING_TSFILE_LEN           = 6
	VERSION_NUM_BYTE                  = byte(0x03)
	CHUNK_GROUP_HEADER_MARKER         = byte(0x00)
	CHUNK_HEADER_MARKER               = byte(0x01)
	ONLY_ONE_PAGE_CHUNK_HEADER_MARKER = byte(0x05)
	SEPARATOR_MARKER                  = byte(0x02)
	OPERATION_INDEX_RANGE             = byte(0x04)
)

// Errors used for handling invalid cases
var (
	ErrInvalidArg           = errors.New("invalid argument")
	ErrOutOfMemory          = errors.New("out of memory")
	ErrSerialization        = errors.New("serialization failed")
	ErrDeserialization      = errors.New("deserialization failed")
	ErrDeserializationField = errors.New("deserialization failed: field")
)

/////////////////
// page header //
/////////////////

// PageHeader struct represents the page metadata
type PageHeader struct {
	UncompressedSize uint32
	CompressedSize   uint32
	Statistic        statistic.Interface
}

// NewPageHeader initializes a new PageHeader with default values
func NewPageHeader() *PageHeader {
	return &PageHeader{
		UncompressedSize: 0,
		CompressedSize:   0,
		Statistic:        nil,
	}
}

// Reset clears the PageHeader fields, releasing the associated statistic.
func (p *PageHeader) Reset(factory statistic.Factory) {
	if p.Statistic != nil {
		factory.FreeStatistic(p.Statistic)
		p.Statistic = nil
	}
	p.UncompressedSize = 0
	p.CompressedSize = 0
}

// DeserializeFrom deserializes a PageHeader from ByteStream
func (p *PageHeader) DeserializeFrom(stream *base.ByteStream, deserializeStat bool, dataType base.TSDataType, serial base.SerializationUtil, factory statistic.Factory) error {
	var err error

	// Read UncompressedSize
	if p.UncompressedSize, err = serial.ReadVarUint(stream); err != nil {
		return errors.New("failed to read uncompressed size: " + err.Error())
	}

	// Read CompressedSize
	if p.CompressedSize, err = serial.ReadVarUint(stream); err != nil {
		return errors.New("failed to read compressed size: " + err.Error())
	}

	if deserializeStat {
		// Allocate a statistic object
		stat, err := factory.AllocStatistic(dataType)
		if err != nil {
			return errors.New("failed to allocate statistic: " + err.Error())
		}

		// Deserialize the statistic
		if err = stat.DeserializeTypedStat(stream); err != nil {
			return errors.New("failed to deserialize statistic: " + err.Error())
		}

		p.Statistic = stat // Assign the deserialized statistic

	}
	return nil
}

// EstimateMaxPageHeaderSizeWithoutStatistics used for memory estimate
func EstimateMaxPageHeaderSizeWithoutStatistics() int {
	return 2 * (4 + 1)
}

//////////////////
// chunk header //
//////////////////

// ChunkHeader struct represents chunk metadata
type ChunkHeader struct {
	MeasurementName string
	DataSize        uint32
	DataType        base.TSDataType
	CompressionType base.CompressionType
	EncodingType    base.TSEncoding
	NumOfPages      int32
	SerializedSize  int32
	ChunkType       byte
}

// NewChunkHeader initializes a new ChunkHeader with default values
func NewChunkHeader() *ChunkHeader {
	return &ChunkHeader{
		MeasurementName: "",
		DataSize:        0,
		DataType:        base.INVALID_TS,
		CompressionType: base.INVALID_C,
		EncodingType:    base.INVALID_E,
		NumOfPages:      0,
		SerializedSize:  0,
		ChunkType:       0,
	}
}

// Reset clears the ChunkHeader fields
func (c *ChunkHeader) Reset() {
	c.MeasurementName = ""
	c.DataSize = 0
	c.DataType = base.INVALID_TS
	c.CompressionType = base.INVALID_C
	c.EncodingType = base.INVALID_E
	c.NumOfPages = 0
	c.SerializedSize = 0
	c.ChunkType = 0
}

// SerializeTo serializes the ChunkHeader to a ByteStream
func (c *ChunkHeader) SerializeTo(stream *base.ByteStream) error {
	util := base.SerializationUtil{}
	err := util.WriteUint8(c.ChunkType, stream) // Write_char direct convert
	if err != nil {
		return ErrSerialization
	}
	if err := util.WriteString(c.MeasurementName, stream); err != nil {
		return ErrSerialization
	}
	if err := util.WriteVarUint(c.DataSize, stream); err != nil {
		return ErrSerialization
	}
	if err := util.WriteUint8(c.DataType.TSDataTypeToEnum(), stream); err != nil {
		return ErrSerialization
	}
	if err := util.WriteUint8(c.CompressionType.CompressionTypeToEnum(), stream); err != nil {
		return ErrSerialization
	}
	if err := util.WriteUint8(c.EncodingType.TSEncodingToEnum(), stream); err != nil {
		return ErrSerialization
	}
	return nil
}

// DeserializeFrom deserializes a ChunkHeader from a ByteStream
func (c *ChunkHeader) DeserializeFrom(stream *base.ByteStream, util *base.SerializationUtil) error {
	var err error
	if c.ChunkType, err = util.ReadUint8(stream); err != nil {
		return ErrDeserialization
	}
	if c.MeasurementName, err = util.ReadString(stream); err != nil {
		return ErrDeserialization
	}
	if c.DataSize, err = util.ReadVarUint(stream); err != nil {
		return ErrDeserialization
	}
	if b, err := util.ReadUint8(stream); err != nil {
		return ErrDeserialization
	} else {
		c.DataType = base.TSDataType(b)
	}
	if b, err := util.ReadUint8(stream); err != nil {
		return ErrDeserialization
	} else {
		c.CompressionType = base.CompressionType(b)
	}
	if b, err := util.ReadUint8(stream); err != nil {
		return ErrDeserialization
	} else {
		c.EncodingType = base.TSEncoding(b)
	}
	return nil
}

////////////////
// chunk meta //
////////////////

// ChunkMeta represents metadata for a chunk
type ChunkMeta struct {
	MeasurementName string
	DataType        base.TSDataType
	OffsetOfHeader  int64
	Statistic       statistic.Interface
	TsID            utils.TsID
	Mask            byte
}

// Initialize creates a new ChunkMeta with the given parameters
func (c *ChunkMeta) Initialize(name string, dataType base.TSDataType, offset int64, stat statistic.Interface, TsID utils.TsID, mask byte) error {
	c.MeasurementName = name
	c.DataType = dataType
	c.OffsetOfHeader = offset
	c.Statistic = stat
	c.TsID = TsID
	c.Mask = mask
	return nil
}

// CloneStatisticFrom duplicates data from another ChunkMeta
func (c *ChunkMeta) CloneStatisticFrom(stat statistic.Interface) error {
	err := statistic.CloneStatistic(stat, c.Statistic, c.DataType)
	if err != nil {
		return err
	}
	return nil
}

// CloneFrom duplicates the data from another ChunkMeta.
func (c *ChunkMeta) CloneFrom(that *ChunkMeta, stat statistic.Factory) error {
	// Clone MeasurementName
	c.MeasurementName = that.MeasurementName

	// Clone the basic fields
	c.DataType = that.DataType
	c.OffsetOfHeader = that.OffsetOfHeader
	c.TsID = that.TsID
	c.Mask = that.Mask

	// Clone the statistic if available in the source
	if that.Statistic != nil {
		// Ensure the target Statistic exists; initialize if nil
		if c.Statistic == nil {
			c.Statistic, _ = stat.AllocStatistic(c.DataType) // Create a new statistic based on DataType
		}
		// Call the clone method for the statistic
		if err := statistic.CloneStatistic(that.Statistic, c.Statistic, c.DataType); err != nil {
			return err
		}
	}

	return nil
}

// SerializeTo serializes ChunkMeta to a ByteStream
func (c *ChunkMeta) SerializeTo(stream *base.ByteStream, serializeStat bool, util *base.SerializationUtil) error {
	if err := util.WriteUint64(uint64(c.OffsetOfHeader), stream); err != nil {
		return ErrSerialization
	}
	if serializeStat && c.Statistic != nil {
		if err := c.Statistic.SerializeTypedStat(stream); err != nil {
			return ErrSerialization
		}
	}
	return nil
}

// DeserializeFrom deserializes the ChunkMeta from a ByteStream.
func (c *ChunkMeta) DeserializeFrom(stream *base.ByteStream, deserializeStat bool, util *base.SerializationUtil) error {
	// Read the OffsetOfHeader from the ByteStream
	offset, err := util.ReadUint64(stream)
	if err != nil {
		return ErrDeserialization
	}
	c.OffsetOfHeader = int64(offset)

	// Deserialize the statistic if needed
	if deserializeStat {
		// Deserialize into the Statistic
		if err := c.Statistic.DeserializeTypedStat(stream); err != nil {
			return ErrDeserialization
		}
	}

	return nil
}

//////////////////////
// chunk group meta //
//////////////////////

// ChunkGroupMeta represents metadata for a group of chunks
type ChunkGroupMeta struct {
	DeviceName    string
	ChunkMetaList []*ChunkMeta
}

// NewChunkGroupMeta creates an empty ChunkGroupMeta
func NewChunkGroupMeta(deviceName string) *ChunkGroupMeta {
	return &ChunkGroupMeta{
		DeviceName:    deviceName,
		ChunkMetaList: []*ChunkMeta{},
	}
}

// Push adds a ChunkMeta to the group
func (cg *ChunkGroupMeta) Push(meta *ChunkMeta) {
	cg.ChunkMetaList = append(cg.ChunkMetaList, meta)
}

// SortChunks organizes the ChunkMetaList by measurement name and offset
func (cg *ChunkGroupMeta) SortChunks() {
	sort.SliceStable(cg.ChunkMetaList, func(i, j int) bool {
		if cg.ChunkMetaList[i].MeasurementName == cg.ChunkMetaList[j].MeasurementName {
			return cg.ChunkMetaList[i].OffsetOfHeader < cg.ChunkMetaList[j].OffsetOfHeader
		}
		return cg.ChunkMetaList[i].MeasurementName < cg.ChunkMetaList[j].MeasurementName
	})
}

///////////////////////
// time series index //
///////////////////////

// TimeseriesIndex represents a time-series index containing chunk metadata.
type TimeseriesIndex struct {
	TimeseriesMetaType      byte
	ChunkMetaListDataSize   uint32
	MeasurementName         string
	TsID                    utils.TsID
	DataType                base.TSDataType
	Statistic               statistic.Interface
	StatisticFromSerialized bool
	ChunkMetaListSerialized []byte
	ChunkMetaList           []*ChunkMeta
}

// Initialize creates a new TimeseriesIndex with the default values.
func (idx *TimeseriesIndex) Initialize() {
	idx.TimeseriesMetaType = 255
	idx.ChunkMetaListDataSize = 0
	idx.MeasurementName = ""
	idx.TsID = utils.TsID{}
	idx.DataType = base.INVALID_TS // Represents an invalid data type.
	idx.Statistic = nil
	idx.StatisticFromSerialized = false
	idx.ChunkMetaListSerialized = nil
	idx.ChunkMetaList = nil
}

// SortChunkMetaList sorts the ChunkMetaList by MeasurementName, and by OffsetOfHeader if names are equal.
func (idx *TimeseriesIndex) SortChunkMetaList() {
	sort.Slice(idx.ChunkMetaList, func(i, j int) bool {
		if idx.ChunkMetaList[i].MeasurementName != idx.ChunkMetaList[j].MeasurementName {
			return idx.ChunkMetaList[i].MeasurementName < idx.ChunkMetaList[j].MeasurementName
		}
		return idx.ChunkMetaList[i].OffsetOfHeader < idx.ChunkMetaList[j].OffsetOfHeader
	})
}

// Reset resets the TimeseriesIndex to default values.
func (idx *TimeseriesIndex) Reset() {
	idx.TimeseriesMetaType = 0
	idx.ChunkMetaListDataSize = 0
	idx.MeasurementName = ""
	idx.TsID.Reset()
	idx.DataType = base.VECTOR
	idx.Statistic = nil
	idx.StatisticFromSerialized = false
	idx.ChunkMetaListSerialized = nil
	idx.ChunkMetaList = nil
}

// AddChunkMeta adds a ChunkMeta to the TimeseriesIndex and optionally serializes it.
func (idx *TimeseriesIndex) AddChunkMeta(chunkMeta *ChunkMeta, serializeStatistic bool, stream *base.ByteStream) error {
	// Check if the provided chunkMeta is nil
	if chunkMeta == nil {
		return errors.New("chunkMeta cannot be nil")
	}

	// Serialize the chunkMeta to the internal serialized buffer if required
	if serializeStatistic {
		if err := chunkMeta.SerializeTo(stream, serializeStatistic, nil); err != nil {
			return err
		}

		// Append the serialized bytes to the serialized buffer
		serialisedBytes, _ := stream.GetBytesFromByteStream()
		idx.ChunkMetaListSerialized = append(idx.ChunkMetaListSerialized, serialisedBytes...)
		idx.ChunkMetaListDataSize = uint32(len(idx.ChunkMetaListSerialized))
	}

	// Add the chunkMeta to the list of chunk metas
	idx.ChunkMetaList = append(idx.ChunkMetaList, chunkMeta)

	// Merge the statistics of the TimeseriesIndex with the chunkMeta's statistics
	if idx.Statistic != nil && chunkMeta.Statistic != nil {
		// Merge the statistics
		if err := idx.Statistic.MergeWith(chunkMeta.Statistic); err != nil {
			return err
		}
	}

	return nil
}

// SetMeasurementName sets the measurement name for the TimeseriesIndex.
func (idx *TimeseriesIndex) SetMeasurementName(name string) {
	idx.MeasurementName = name
}

// GetMeasurementName retrieves the measurement name.
func (idx *TimeseriesIndex) GetMeasurementName() string {
	return idx.MeasurementName
}

// SetDataType sets the data type for the TimeseriesIndex.
func (idx *TimeseriesIndex) SetDataType(dataType base.TSDataType) {
	idx.DataType = dataType
}

// GetDataType retrieves the data type of the TimeseriesIndex.
func (idx *TimeseriesIndex) GetDataType() base.TSDataType {
	return idx.DataType
}

// InitStatistic initializes the Statistic field based on the data type.
func (idx *TimeseriesIndex) InitStatistic(dataType base.TSDataType) error {
	factory := statistic.Factory{}
	stat, err := factory.AllocStatistic(dataType)
	if err != nil {
		return err
	}
	idx.Statistic = stat
	idx.Statistic.Reset()
	return nil
}

// SerializeTo serializes the TimeseriesIndex into a ByteStream.
func (idx *TimeseriesIndex) SerializeTo(stream *base.ByteStream) error {
	util := base.SerializationUtil{}

	if err := util.WriteUint8(idx.TimeseriesMetaType, stream); err != nil {
		return ErrSerialization
	}
	if err := util.WriteString(idx.MeasurementName, stream); err != nil {
		return ErrSerialization
	}
	if err := util.WriteUint8(base.TSDataType.TSDataTypeToEnum(idx.DataType), stream); err != nil {
		return ErrSerialization
	}
	if err := util.WriteUint32(idx.ChunkMetaListDataSize, stream); err != nil {
		return ErrSerialization
	}
	if idx.Statistic != nil {
		if err := idx.Statistic.SerializeTypedStat(stream); err != nil {
			return ErrSerialization
		}
	}
	if len(idx.ChunkMetaListSerialized) > 0 {
		err := stream.WriteBuf(idx.ChunkMetaListSerialized, uint32(len(idx.ChunkMetaListSerialized)))
		if err != nil {
			return err
		}
	}
	return nil
}

// DeserializeFrom deserializes the TimeseriesIndex from a ByteStream.
func (idx *TimeseriesIndex) DeserializeFrom(stream *base.ByteStream) error {
	util := base.SerializationUtil{}

	var err error
	// Deserialize fields
	idx.TimeseriesMetaType, err = util.ReadUint8(stream)
	if err != nil {
		return ErrDeserializationField
	}
	idx.MeasurementName, err = util.ReadString(stream)
	if err != nil {
		return ErrDeserializationField
	}
	datatype, err := util.ReadUint8(stream)
	if err != nil {
		return ErrDeserializationField
	}
	idx.DataType = base.TSDataType.EnumToTSDataType(base.TSDataType(strconv.Itoa(int(datatype))))
	idx.ChunkMetaListDataSize, err = util.ReadUint32(stream)
	if err != nil {
		return ErrDeserializationField
	}

	// Deserialize statistic object
	statFactory := statistic.Factory{}
	stat, err := statFactory.AllocStatistic(idx.DataType)
	if err != nil {
		return err
	}
	if err := stat.DeserializeTypedStat(stream); err != nil {
		return err
	}
	idx.Statistic = stat
	idx.StatisticFromSerialized = true

	// Deserialize chunk meta list
	for pos := uint32(0); pos < idx.ChunkMetaListDataSize; {
		chunkMeta := &ChunkMeta{}
		if err := chunkMeta.DeserializeFrom(stream, true, nil); err != nil {
			return ErrDeserialization
		}
		idx.ChunkMetaList = append(idx.ChunkMetaList, chunkMeta)
		pos += stream.ReadPos
	}
	return nil
}

// CloneFrom clones the data from another TimeseriesIndex.
func (idx *TimeseriesIndex) CloneFrom(that *TimeseriesIndex) error {
	idx.TimeseriesMetaType = that.TimeseriesMetaType
	idx.ChunkMetaListDataSize = that.ChunkMetaListDataSize
	idx.MeasurementName = that.MeasurementName
	idx.TsID = that.TsID
	idx.DataType = that.DataType

	// Clone the statistic
	if that.Statistic != nil {
		statFactory := statistic.Factory{}
		stat, err := statFactory.AllocStatistic(idx.DataType)
		if err != nil {
			return err
		}
		if err := statistic.CloneStatistic(that.Statistic, stat, idx.DataType); err != nil {
			return err
		}
		idx.Statistic = stat
	}

	// Clone the chunkMeta list
	for _, cm := range that.ChunkMetaList {
		clonedChunkMeta := &ChunkMeta{}
		if err := clonedChunkMeta.CloneFrom(cm, statistic.Factory{}); err != nil {
			return err
		}
		idx.ChunkMetaList = append(idx.ChunkMetaList, clonedChunkMeta)
	}
	return nil
}

// SetTsID updates the TsID of the TimeseriesIndex and propagates it to all ChunkMeta in the ChunkMetaList.
func (idx *TimeseriesIndex) SetTsID(tsID utils.TsID) {
	idx.TsID = tsID

	/*
		// Debug behavior: propagate the tsID to all ChunkMeta in the ChunkMetaList.
		if idx.ChunkMetaList != nil {
			for _, chunkMeta := range idx.ChunkMetaList {
				chunkMeta.TsID = tsID
			}
		}
	*/
}

// GetTsId returns the TsID of the TimeseriesIndex.
func (idx *TimeseriesIndex) GetTsId() utils.TsID {
	return idx.TsID
}

// Finish finalizes the TimeseriesIndex by updating the serialized buffer size.
func (idx *TimeseriesIndex) Finish() {
	idx.ChunkMetaListDataSize = uint32(len(idx.ChunkMetaListSerialized))
}

// SetTsMetaType is set as 1 if there are more than 1 chunk meta and 0 if there are only 1
func (idx *TimeseriesIndex) SetTsMetaType(metaType byte) {
	idx.TimeseriesMetaType = metaType
}

///////////////////////////////
// aligned time series index //
///////////////////////////////

type AlignedTimeseriesIndex struct {
	TimeTsIdx  *TimeseriesIndex
	ValueTsIdx *TimeseriesIndex
}

func NewAlignedTimeseriesIndex() *AlignedTimeseriesIndex {
	return &AlignedTimeseriesIndex{
		TimeTsIdx:  nil,
		ValueTsIdx: nil,
	}
}

// GetTimeChunkMetaList returns the chunk meta list from the time index.
func (a *AlignedTimeseriesIndex) GetTimeChunkMetaList() []*ChunkMeta {
	if a.TimeTsIdx == nil {
		return nil
	}
	return a.TimeTsIdx.ChunkMetaList
}

// GetValueChunkMetaList returns the chunk meta list from the value index.
func (a *AlignedTimeseriesIndex) GetValueChunkMetaList() []*ChunkMeta {
	if a.ValueTsIdx == nil {
		return nil
	}
	return a.ValueTsIdx.ChunkMetaList
}

// GetMeasurementName returns the measurement name from the value index.
func (a *AlignedTimeseriesIndex) GetMeasurementName() string {
	if a.ValueTsIdx == nil {
		return ""
	}
	return a.ValueTsIdx.GetMeasurementName()
}

// GetDataType returns the data type from the time index.
func (a *AlignedTimeseriesIndex) GetDataType() base.TSDataType {
	if a.TimeTsIdx == nil {
		return base.INVALID_TS
	}
	return a.TimeTsIdx.GetDataType()
}

// GetStatistic returns the statistic from the value index.
func (a *AlignedTimeseriesIndex) GetStatistic() *statistic.Statistic {
	if a.ValueTsIdx == nil {
		return nil
	}
	return a.GetStatistic()

}

/////////////////
// TSMIterator //
/////////////////

type TSMIterator struct {
	ChunkGroupMetaList []*ChunkGroupMeta
	ChunkGroupMetaIter int
	ChunkMetaIter      int

	TSMChunkMetaInfo  map[string]map[string][]*ChunkMeta // device_name -> measurement_name -> chunk_meta
	DeviceIterator    []string                           // ordered list of device names
	MeasurementIter   int
	CurrentDeviceIter int
}

// TSMInit initializes and organizes the TSMIterator by grouping metadata and sorting chunks.
func (iter *TSMIterator) TSMInit() error {
	// Build the TSMChunkMetaInfo structure: group by device name and measurement name
	iter.TSMChunkMetaInfo = make(map[string]map[string][]*ChunkMeta)
	for _, groupMeta := range iter.ChunkGroupMetaList {
		deviceName := groupMeta.DeviceName
		if _, exists := iter.TSMChunkMetaInfo[deviceName]; !exists {
			iter.TSMChunkMetaInfo[deviceName] = make(map[string][]*ChunkMeta)
		}

		// Group chunks by measurement name and sort them
		for _, chunkMeta := range groupMeta.ChunkMetaList {
			measurementName := chunkMeta.MeasurementName
			iter.TSMChunkMetaInfo[deviceName][measurementName] = append(iter.TSMChunkMetaInfo[deviceName][measurementName], chunkMeta)
		}

		// Sort chunks by OffsetOfHeader for every measurement under this device
		for measurementName, chunks := range iter.TSMChunkMetaInfo[deviceName] {
			sort.Slice(chunks, func(i, j int) bool {
				return chunks[i].OffsetOfHeader < chunks[j].OffsetOfHeader
			})
			iter.TSMChunkMetaInfo[deviceName][measurementName] = chunks
		}
	}

	// Prepare iterators for traversal
	iter.DeviceIterator = make([]string, 0, len(iter.TSMChunkMetaInfo))
	for device := range iter.TSMChunkMetaInfo {
		iter.DeviceIterator = append(iter.DeviceIterator, device)
	}
	sort.Strings(iter.DeviceIterator) // Sort device names to iterate in a consistent order

	iter.CurrentDeviceIter = 0
	iter.MeasurementIter = 0
	return nil
}

// HasNext checks if there are more entries to iterate over
func (iter *TSMIterator) HasNext() bool {
	return iter.CurrentDeviceIter < len(iter.DeviceIterator)
}

// GetNext retrieves the next device, measurement, and timeseries index
func (iter *TSMIterator) GetNext() (string, string, *TimeseriesIndex, error) {
	if !iter.HasNext() {
		return "", "", nil, utils.GetError(utils.ErrNoMoreData) // Example error for no more data
	}

	deviceName := iter.DeviceIterator[iter.CurrentDeviceIter]
	measurements := iter.TSMChunkMetaInfo[deviceName]

	// Get the current measurement
	measurementNames := make([]string, 0, len(measurements))
	for measurement := range measurements {
		measurementNames = append(measurementNames, measurement)
	}
	sort.Strings(measurementNames) // Ensure consistent order for measurements
	if iter.MeasurementIter >= len(measurementNames) {
		// Move to the next device

		iter.MeasurementIter = 0
		return iter.GetNext()
	}

	measurementName := measurementNames[iter.MeasurementIter]
	chunkMetaList := measurements[measurementName]
	iter.MeasurementIter++

	if len(chunkMetaList) == 0 {
		return "", "", nil, utils.GetError(utils.ErrMeasurementNotExist) // Example error for missing metadata
	}

	// Create a TimeseriesIndex with metadata from the chunks
	tsIndex := &TimeseriesIndex{}
	multiChunks := len(chunkMetaList) > 1
	firstChunk := chunkMetaList[0]

	// Fill in the basic metadata for the timeseries index
	metaType := byte(0)
	if multiChunks {
		metaType = 1
	}
	metaType |= firstChunk.Mask

	tsIndex.SetTsMetaType(metaType)
	tsIndex.SetMeasurementName(measurementName)
	tsIndex.SetDataType(firstChunk.DataType)
	err := tsIndex.InitStatistic(firstChunk.DataType)
	if err != nil {
		return "", "", nil, err
	}
	tsIndex.SetTsID(firstChunk.TsID)

	// Add all chunk metadata to the TimeseriesIndex
	for _, chunkMeta := range chunkMetaList {
		if err := tsIndex.AddChunkMeta(chunkMeta, multiChunks, nil); err != nil {
			return "", "", nil, err
		}
	}
	tsIndex.Finish()

	if deviceName == "" {
		return "", "", nil, utils.GetError(utils.ErrTsFileWriterMetaErr) // Example error for invalid metadata
	}
	iter.CurrentDeviceIter++

	return deviceName, measurementName, tsIndex, nil
}

//////////////////
// tsfile index //
//////////////////

type MetaIndexEntry struct {
	Name   string
	Offset int64
}

// Init initializes the MetaIndexEntry with a name and offset
func (entry *MetaIndexEntry) Init(name string, offset int64) {
	entry.Name = name
	entry.Offset = offset
}

// Serialize serializes the MetaIndexEntry into a ByteStream.
func (entry *MetaIndexEntry) Serialize(stream *base.ByteStream) error {
	serial := base.SerializationUtil{}
	if err := serial.WriteString(entry.Name, stream); err != nil {
		return err
	}
	if err := serial.WriteUint64(uint64(entry.Offset), stream); err != nil {
		return err
	}
	return nil
}

// Deserialize deserializes the MetaIndexEntry from a ByteStream.
func (entry *MetaIndexEntry) Deserialize(stream *base.ByteStream) error {
	serial := base.SerializationUtil{}
	name, err := serial.ReadString(stream)
	if err != nil {
		return errors.New("readString error: " + err.Error())
	}
	entry.Name = name

	offset, err := serial.ReadUint64(stream)
	if err != nil {
		return err
	}
	entry.Offset = int64(offset)

	return nil
}

type MetaIndexNode struct {
	Children  []*MetaIndexEntry
	EndOffset int64
	NodeType  MetaIndexNodeType
}

// GetFirstChildName retrieves the name of the first child matching the specified node type or returns an error if unavailable.
func (node *MetaIndexNode) GetFirstChildName() (string, error) {
	if len(node.Children) == 0 {
		return "", errors.New("node has no children")
	}
	return node.Children[0].Name, nil
}

// Serialize serializes the MetaIndexNode into a ByteStream.
func (node *MetaIndexNode) Serialize(stream *base.ByteStream) error {
	serial := base.SerializationUtil{}
	// Write the number of children
	if err := serial.WriteVarUint(uint32(len(node.Children)), stream); err != nil {
		return err
	}

	// Serialize each child entry
	for _, entry := range node.Children {
		if err := entry.Serialize(stream); err != nil {
			return err
		}
	}

	// Write the end offset
	if err := serial.WriteUint64(uint64(node.EndOffset), stream); err != nil {
		return err
	}

	// Write the node type
	if err := serial.WriteUint8(uint8(node.NodeType), stream); err != nil {
		return err
	}

	return nil
}

// Deserialize deserializes the MetaIndexNode from a ByteStream.
func (node *MetaIndexNode) Deserialize(stream *base.ByteStream) error {
	serial := base.SerializationUtil{}
	// Read the number of children
	childrenSize, err := serial.ReadVarUint(stream)
	if err != nil {
		return err
	}

	// Deserialize each child entry
	node.Children = make([]*MetaIndexEntry, 0, childrenSize)
	for i := 0; i < int(childrenSize); i++ {
		entry := &MetaIndexEntry{}
		if err := entry.Deserialize(stream); err != nil {
			return err
		}
		node.Children = append(node.Children, entry)
	}

	// Read the end offset
	offset, err := serial.ReadUint64(stream)
	if err != nil {
		return err
	}

	node.EndOffset = int64(offset)

	// Read the node type
	nodeType, err := serial.ReadUint8(stream)
	if err != nil {
		return err
	}
	node.NodeType = MetaIndexNodeType(nodeType)

	return nil
}

// BinarySearchChildren Returns nil on success, or an error if the entry is not found.
func (node *MetaIndexNode) BinarySearchChildren(name string, exactSearch bool, retIndexEntry *MetaIndexEntry) (int64, error) {
	// Edge case: Leaf measurement with a single empty child
	isAligned := false
	if node.NodeType == LEAF_MEASUREMENT && len(node.Children) == 1 && node.Children[0].Name == "" {
		isAligned = true
	}

	// Binary search variables
	l := -1
	if isAligned {
		l = 0
	} else {
		h := len(node.Children)
		found := false
		for l < h-1 {
			m := (l + h) / 2
			cmp := compareStrings(node.Children[m].Name, name)
			if cmp == 0 { // Exact match
				l = m
				found = true
				break
			} else if cmp < 0 { // Children[m].Name < name
				l = m
			} else { // Children[m].Name > name
				h = m
			}
		}

		if !found && exactSearch {
			return 0, errors.New("child not found")
		}
	}

	// If not exact search, `l` will point to the largest entry <= name
	if l == -1 || l >= len(node.Children) {
		return 0, errors.New("child not found")
	}

	// Populate result fields for the matched entry
	*retIndexEntry = *node.Children[l]
	return node.Children[l].Offset, nil
}

// Helper function to compare strings, simulating C++'s compare
func compareStrings(a, b string) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func (node *MetaIndexNode) IsFull() bool {
	return len(node.Children) >= utils.ConfigValue.MaxDegreeOfIndexNode
}

func (node *MetaIndexNode) PushEntry(entry *MetaIndexEntry) error {
	node.Children = append(node.Children, entry)
	return nil
}

func (node *MetaIndexNode) IsEmpty() bool {
	return len(node.Children) == 0
}

// MetaIndexNodeType represents the type of the meta index node.
type MetaIndexNodeType int

const (
	INTERNAL_DEVICE = MetaIndexNodeType(iota)
	LEAF_DEVICE
	INTERNAL_MEASUREMENT
	LEAF_MEASUREMENT
	INVALID_META_NODE_TYPE
)

type TsFileMeta struct {
	IndexNode   *MetaIndexNode
	MetaOffset  int64
	BloomFilter *reader.BloomFilter
}

// Serialize serializes the TsFileMeta into a ByteStream.
func (meta *TsFileMeta) Serialize(stream *base.ByteStream) error {
	serial := base.SerializationUtil{}
	if err := meta.IndexNode.Serialize(stream); err != nil {
		return err
	}
	if err := serial.WriteUint64(uint64(meta.MetaOffset), stream); err != nil {
		return err
	}
	if err := meta.BloomFilter.SerializeTo(stream); err != nil {
		return err
	}
	return nil
}

// Deserialize deserializes the TsFileMeta from a ByteStream.
func (meta *TsFileMeta) Deserialize(stream *base.ByteStream) error {
	serial := base.SerializationUtil{}
	meta.IndexNode = &MetaIndexNode{}
	if err := meta.IndexNode.Deserialize(stream); err != nil {
		return err
	}

	offset, err := serial.ReadUint64(stream)
	if err != nil {
		return err
	}
	meta.MetaOffset = int64(offset)

	meta.BloomFilter = &reader.BloomFilter{}
	if err := meta.BloomFilter.DeserializeFrom(stream); err != nil {
		return err
	}

	return nil
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

type TimeseriesTimeIndexEntry struct {
	TsID      utils.TsID
	TimeRange TimeRange
}
