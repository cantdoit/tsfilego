package core

import (
	"Golang/internal/common/base"
	"errors"
)

// Field represents a single value with its associated type and optional column name.
type Field struct {
	Type       base.TSDataType // Data type of the field (e.g., INT64, BOOLEAN)
	ColumnName string          // Optional: Column name associated with the field
	BoolVal    bool
	Int64Val   int64
	Int32Val   int32
	FloatVal   float32
	DoubleVal  float64
	StringVal  string
}

// NewField creates a new Field instance with a specific type.
func NewField(dataType base.TSDataType) *Field {
	return &Field{
		Type:       dataType,
		ColumnName: "",
	}
}

// SetValue sets the value of the Field based on its data type.
func (f *Field) SetValue(value interface{}) error {
	switch f.Type {
	case base.BOOLEAN:
		if v, ok := value.(bool); ok {
			f.BoolVal = v
		} else {
			return errors.New("invalid value type for BOOLEAN")
		}
	case base.INT32:
		if v, ok := value.(int32); ok {
			f.Int32Val = v
		} else {
			return errors.New("invalid value type for INT32")
		}
	case base.INT64:
		if v, ok := value.(int64); ok {
			f.Int64Val = v
		} else {
			return errors.New("invalid value type for INT64")
		}
	case base.FLOAT:
		if v, ok := value.(float32); ok {
			f.FloatVal = v
		} else {
			return errors.New("invalid value type for FLOAT")
		}
	case base.DOUBLE:
		if v, ok := value.(float64); ok {
			f.DoubleVal = v
		} else {
			return errors.New("invalid value type for DOUBLE")
		}
	case base.TEXT:
		if v, ok := value.(string); ok {
			f.StringVal = v
		} else {
			return errors.New("invalid value type for TEXT")
		}
	default:
		return errors.New("unsupported data type")
	}
	return nil
}

// Free clears the StringVal memory for TEXT type.
func (f *Field) Free() {
	if f.Type == base.TEXT {
		f.StringVal = ""
	}
}

// RowRecord represents a row of data with associated fields and a timestamp.
type RowRecord struct {
	Timestamp int64    // Optional timestamp for the row
	Fields    []*Field // Fields representing the row's values
}

// NewRowRecord constructs a RowRecord with a given column count and optional timestamp.
func NewRowRecord(timestamp int64, columnCount int) *RowRecord {
	record := &RowRecord{
		Timestamp: timestamp,
		Fields:    make([]*Field, columnCount),
	}

	// Initialize fields with NULL_TYPE by default
	for i := 0; i < columnCount; i++ {
		record.Fields[i] = NewField(base.NULL_TYPE)
	}
	return record
}

// AddField adds a Field to the RowRecord.
func (r *RowRecord) AddField(field *Field) {
	r.Fields = append(r.Fields, field)
}

// GetField returns the Field present at a specific index.
func (r *RowRecord) GetField(index int) (*Field, error) {
	if index < 0 || index >= len(r.Fields) {
		return nil, errors.New("field index out of range")
	}
	return r.Fields[index], nil
}

// SetFieldValue modifies the value of a field at a specific index.
func (r *RowRecord) SetFieldValue(index int, value interface{}, dataType base.TSDataType) error {
	if index < 0 || index >= len(r.Fields) {
		return errors.New("field index out of range")
	}

	field := r.Fields[index]
	field.Type = dataType
	return field.SetValue(value)
}
