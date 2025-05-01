package core

import (
	"Golang/internal/common/base"
	"Golang/internal/common/statistic"
	"errors"
	"sort"
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
		p.Statistic, _ = factory.AllocStatistic(dataType)
		if p.Statistic == nil {
			return errors.New("failed to allocate statistic")
		}

		// Deserialize the statistic
		if err = p.Statistic.(stream); err != nil {
			return errors.New("failed to deserialize statistic: " + err.Error())
		}
	}
	return nil
}

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
	if err := stream.WriteByte(c.ChunkType); err != nil {
		return ErrSerialization
	}
	if err := stream.WriteString(c.MeasurementName); err != nil {
		return ErrSerialization
	}
	if err := stream.WriteVarUint(c.DataSize); err != nil {
		return ErrSerialization
	}
	if err := stream.WriteByte(byte(c.DataType)); err != nil {
		return ErrSerialization
	}
	if err := stream.WriteByte(byte(c.CompressionType)); err != nil {
		return ErrSerialization
	}
	if err := stream.WriteByte(byte(c.EncodingType)); err != nil {
		return ErrSerialization
	}
	return nil
}

// DeserializeFrom deserializes a ChunkHeader from a ByteStream
func (c *ChunkHeader) DeserializeFrom(stream *base.ByteStream) error {
	var err error
	if c.ChunkType, err = stream.ReadByte(); err != nil {
		return ErrDeserialization
	}
	if c.MeasurementName, err = stream.ReadString(); err != nil {
		return ErrDeserialization
	}
	if c.DataSize, err = stream.ReadVarUint(); err != nil {
		return ErrDeserialization
	}
	if b, err := stream.ReadByte(); err != nil {
		return ErrDeserialization
	} else {
		c.DataType = base.TSDataType(b)
	}
	if b, err := stream.ReadByte(); err != nil {
		return ErrDeserialization
	} else {
		c.CompressionType = base.CompressionType(b)
	}
	if b, err := stream.ReadByte(); err != nil {
		return ErrDeserialization
	} else {
		c.EncodingType = base.TSEncoding(b)
	}
	return nil
}

// ChunkMeta represents metadata for a chunk
type ChunkMeta struct {
	MeasurementName string
	DataType        base.TSDataType
	OffsetOfHeader  int64
	Statistic       *statistic.Statistic
	Mask            byte
}

// Initialize creates a new ChunkMeta with the given parameters
func (c *ChunkMeta) Initialize(name string, dataType base.TSDataType, offset int64, stat *statistic.Statistic, mask byte) {
	c.MeasurementName = name
	c.DataType = dataType
	c.OffsetOfHeader = offset
	c.Statistic = stat
	c.Mask = mask
}

// CloneFrom duplicates data from another ChunkMeta
func (c *ChunkMeta) CloneFrom(other *ChunkMeta) error {
	c.MeasurementName = other.MeasurementName
	c.DataType = other.DataType
	c.OffsetOfHeader = other.OffsetOfHeader
	if other.Statistic != nil {
		c.Statistic = other.Statistic.Clone()
	}
	c.Mask = other.Mask
	return nil
}

// SerializeTo serializes ChunkMeta to a ByteStream
func (c *ChunkMeta) SerializeTo(stream *base.ByteStream, serializeStat bool) error {
	if err := stream.WriteInt64(c.OffsetOfHeader); err != nil {
		return ErrSerialization
	}
	if serializeStat && c.Statistic != nil {
		if err := c.Statistic.SerializeTo(stream); err != nil {
			return ErrSerialization
		}
	}
	return nil
}

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
