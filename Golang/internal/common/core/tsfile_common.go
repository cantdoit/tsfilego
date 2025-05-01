package core

import (
	"Golang/internal/common/base"
	"Golang/internal/common/statistic"
	"Golang/internal/utils"
	"errors"
)

// TsFileID type alias for consistency
type TsFileID int64

// Constants defined in header file
const (
	MAGIC_STRING_TSFILE               = "TsFile"
	MAGIC_STRING_TSFILE_LEN           = len(MAGIC_STRING_TSFILE)
	VERSION_NUM_BYTE                  = byte(0x03)
	CHUNK_GROUP_HEADER_MARKER         = byte(0x00)
	CHUNK_HEADER_MARKER               = byte(0x01)
	ONLY_ONE_PAGE_CHUNK_HEADER_MARKER = byte(0x05)
	SEPARATOR_MARKER                  = byte(0x02)
	OPERATION_INDEX_RANGE             = byte(0x04)
)

// Errors used for handling invalid cases
var (
	ErrInvalidArg      = errors.New("invalid argument")
	ErrOutOfMemory     = errors.New("out of memory")
	ErrSerialization   = errors.New("serialization failed")
	ErrDeserialization = errors.New("deserialization failed")
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
func (c *ChunkHeader) SerializeTo(stream *base.ByteStream, util *base.SerializationUtil) error {
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

/*
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
*/
