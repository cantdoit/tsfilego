package common

import (
	"errors"
	"fmt"
)

// TSDataType defines the supported data types for timeseries values
type TSDataType string

const (
	BOOLEAN TSDataType = "BOOLEAN"
	INT32   TSDataType = "INT32"
	INT64   TSDataType = "INT64"
	FLOAT   TSDataType = "FLOAT"
	DOUBLE  TSDataType = "DOUBLE"
	TEXT    TSDataType = "TEXT"
)

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
