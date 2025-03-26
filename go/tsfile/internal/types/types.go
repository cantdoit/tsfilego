package types

import "unsafe"

type ResultSet struct {
	cSet       unsafe.Pointer // C ResultSet
	columns    []ColumnMeta
	currentRow int
}

type ColumnMeta struct {
	Name string
	Type TSDataType
}

func ErrorFromCode(code int) error {
	// Maps C error codes to Go errors
}

func NewResultSet(cSet unsafe.Pointer) *ResultSet {
	// Initializes with metadata from C ResultSetMetaData
}
