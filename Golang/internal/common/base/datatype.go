package base

import (
	"errors"
	"fmt"
	"strconv"
)

// TSEncoding defines the available encoding types
type TSEncoding string

// CompressionType defines the available compression types
type CompressionType string

// TSDataType defines the supported data types for timeseries values
type TSDataType string

// Supported timeseries data types
const (
	BOOLEAN    TSDataType = "BOOLEAN"
	INT32      TSDataType = "INT32"
	INT64      TSDataType = "INT64"
	FLOAT      TSDataType = "FLOAT"
	DOUBLE     TSDataType = "DOUBLE"
	TEXT       TSDataType = "TEXT"
	VECTOR     TSDataType = "VECTOR"
	NULL_TYPE  TSDataType = "NULL_TYPE"
	INVALID_TS TSDataType = "INVALID_DATATYPE"
)

func (t TSDataType) TSDataTypeToEnum() uint8 {
	switch t {
	case BOOLEAN:
		return 0
	case INT32:
		return 1
	case INT64:
		return 2
	case FLOAT:
		return 3
	case DOUBLE:
		return 4
	case TEXT:
		return 5
	case VECTOR:
		return 6
	case NULL_TYPE:
		return 254
	default:
		return 255 // or some error value
	}
}

func (t TSDataType) EnumToTSDataType() TSDataType {
	switch t {
	case TSDataType(strconv.Itoa(0)):
		return BOOLEAN
	case TSDataType(strconv.Itoa(1)):
		return INT32
	case TSDataType(strconv.Itoa(2)):
		return INT64
	case TSDataType(strconv.Itoa(3)):
		return FLOAT
	case TSDataType(strconv.Itoa(4)):
		return DOUBLE
	case TSDataType(strconv.Itoa(5)):
		return TEXT
	case TSDataType(strconv.Itoa(6)):
		return VECTOR
	case TSDataType(strconv.Itoa(254)):
		return NULL_TYPE
	default:
		return INVALID_TS // or some error value
	}
}

// Supported encoding types
const (
	PLAIN      TSEncoding = "PLAIN"
	DICTIONARY TSEncoding = "DICTIONARY"
	RLE        TSEncoding = "RLE"
	DIFF       TSEncoding = "DIFF"
	TS_2DIFF   TSEncoding = "TS_2DIFF"
	BITMAP     TSEncoding = "BITMAP"
	REGULAR    TSEncoding = "REGULAR"
	INVALID_E  TSEncoding = "INVALID_ENCODING"
)

func (e TSEncoding) TSEncodingToEnum() uint8 {
	switch e {
	case PLAIN:
		return 0
	case DICTIONARY:
		return 1
	case RLE:
		return 2
	//TODO add the remaining ones
	default:
		return 255
	}
}

// Supported compression types
const (
	UNCOMPRESSED CompressionType = "UNCOMPRESSED"
	SNAPPY       CompressionType = "SNAPPY"
	GZIP         CompressionType = "GZIP"
	LZO          CompressionType = "LZO"
	LZ4          CompressionType = "LZ4"
	INVALID_C    CompressionType = "INVALID_COMPRESSION"
)

func (c CompressionType) CompressionTypeToEnum() uint8 {
	switch c {
	case UNCOMPRESSED:
		return 0
	// TODO add the remaining enum
	default:
		return 255
	}
}

// Value represents a type-agnostic data holder
type Value struct {
	Type      TSDataType // Data type of the value
	BoolVal   bool
	Int32Val  int32
	Int64Val  int64
	FloatVal  float32
	DoubleVal float64
	StringVal string
}

// NewValue creates a new Value instance and sets its initial value and type
func NewValue(dataType TSDataType, rawValue interface{}) (*Value, error) {
	val := &Value{Type: dataType}

	switch dataType {
	case BOOLEAN:
		if v, ok := rawValue.(bool); ok {
			val.BoolVal = v
		} else {
			return nil, fmt.Errorf("invalid value type for BOOLEAN: %T", rawValue)
		}
	case INT32:
		if v, ok := rawValue.(int32); ok {
			val.Int32Val = v
		} else {
			return nil, fmt.Errorf("invalid value type for INT32: %T", rawValue)
		}
	case INT64:
		if v, ok := rawValue.(int64); ok {
			val.Int64Val = v
		} else {
			return nil, fmt.Errorf("invalid value type for INT64: %T", rawValue)
		}
	case FLOAT:
		if v, ok := rawValue.(float32); ok {
			val.FloatVal = v
		} else {
			return nil, fmt.Errorf("invalid value type for FLOAT: %T", rawValue)
		}
	case DOUBLE:
		if v, ok := rawValue.(float64); ok {
			val.DoubleVal = v
		} else {
			return nil, fmt.Errorf("invalid value type for DOUBLE: %T", rawValue)
		}
	case TEXT:
		if v, ok := rawValue.(string); ok {
			val.StringVal = v
		} else {
			return nil, fmt.Errorf("invalid value type for TEXT: %T", rawValue)
		}
	default:
		return nil, errors.New("unsupported data type")
	}

	return val, nil
}

// GetValue returns the value stored in the Value instance as an interface{}
func (v *Value) GetValue() interface{} {
	switch v.Type {
	case BOOLEAN:
		return v.BoolVal
	case INT32:
		return v.Int32Val
	case INT64:
		return v.Int64Val
	case FLOAT:
		return v.FloatVal
	case DOUBLE:
		return v.DoubleVal
	case TEXT:
		return v.StringVal
	default:
		return nil
	}
}

// Free clears the content of the Value instance (used when dealing with TEXT types)
func (v *Value) Free() {
	if v.Type == TEXT {
		v.StringVal = ""
	}
}
