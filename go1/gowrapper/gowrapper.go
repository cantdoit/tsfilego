package gowrapper

// #cgo CFLAGS: -IC:/Users/User/Documents/GitHub/tsfilego/cpp/src
// #include "cwrapper/TsFile-cwrapper.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// TsFileReader provides an interface for reading data from a TSFile.
type TsFileReader struct {
	reader          C.CTsFileReader
	queryResult     C.QueryDataRet
	isReaderOpened  bool
	tableName       string
	selectedColumns []string
}

const (
	TypeInt32   = 1 << 8
	TypeBool    = 1 << 9
	TypeFloat32 = 1 << 10
	TypeFloat64 = 1 << 11
	TypeInt64   = 1 << 12
)

func SchemaInfo(dataType int) int {
	return dataType
}

// Open initializes the TSFile reader.
func (r *TsFileReader) Open(filepath string) error {
	cFilePath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cFilePath))

	var errCode C.ErrorCode
	r.reader = C.ts_reader_open(cFilePath, &errCode)

	if int(errCode) != 0 {
		return fmt.Errorf("failed to open TSFile '%s', error code: %v", filepath, int(errCode))
	}

	r.isReaderOpened = true
	return nil
}

// CreateWriter creates a new TsFileWriter instance for the specified file path.
func CreateWriter(filePath string) (*TsFileWriter, error) {
	writer := &TsFileWriter{}
	err := writer.Open(filePath)
	if err != nil {
		return nil, err
	}
	return writer, nil
}

func (w *TsFileWriter) RegisterColumn(columnName string, schemaType int) error {
	if !w.isWriterOpened {
		return errors.New("writer is not opened")
	}

	cColumnName := C.CString(columnName)
	defer C.free(unsafe.Pointer(cColumnName))

	columnSchema := C.ColumnSchema{
		name:       cColumnName,
		column_def: C.SchemaInfo(schemaType),
	}

	errCode := C.tsfile_register_table_column(w.writer, cColumnName, &columnSchema)
	if int(errCode) != 0 {
		return fmt.Errorf("failed to register column '%s', error code: %d", columnName, int(errCode))
	}

	return nil
}

// CreateReader creates a new TsFileReader instance for the specified file path.
func CreateReader(filePath string) (*TsFileReader, error) {
	reader := &TsFileReader{}
	err := reader.Open(filePath)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

// QueryData queries data from the TSFile based on the table name and columns.
func (r *TsFileReader) QueryData(tableName string, columns []string, startTime, endTime *int64) error {
	if !r.isReaderOpened {
		return errors.New("reader is not opened")
	}

	cTableName := C.CString(tableName)
	defer C.free(unsafe.Pointer(cTableName))

	cColumns := make([]*C.char, len(columns))
	for i, col := range columns {
		cColumns[i] = C.CString(col)
		defer C.free(unsafe.Pointer(cColumns[i]))
	}
	cColumnsPtr := (**C.char)(unsafe.Pointer(&cColumns[0]))

	if startTime != nil || endTime != nil {
		start := C.timestamp(C.LLONG_MIN)
		end := C.timestamp(C.LLONG_MAX)
		if startTime != nil {
			start = C.timestamp(*startTime)
		}
		if endTime != nil {
			end = C.timestamp(*endTime)
		}

		r.queryResult = C.ts_reader_begin_end(r.reader, cTableName, cColumnsPtr, C.int(len(columns)), start, end)
	} else {
		r.queryResult = C.ts_reader_read(r.reader, cTableName, cColumnsPtr, C.int(len(columns)))
	}

	if r.queryResult == nil {
		return errors.New("failed to query data from TSFile")
	}

	r.tableName = tableName
	r.selectedColumns = columns
	return nil
}

// GetNextBatch retrieves the next batch of data from the queried result.
func (r *TsFileReader) GetNextBatch(expectLineCount int) (map[string]interface{}, error) {
	if r.queryResult == nil {
		return nil, errors.New("no query result available")
	}

	batch := C.ts_next(r.queryResult, C.int(expectLineCount))
	if batch == nil {
		return nil, errors.New("no more data")
	}

	// Result parsing is application-specific. Implement specific processing logic here.

	return map[string]interface{}{}, nil
}

// Close closes the reader and releases resources.
func (r *TsFileReader) Close() error {
	if !r.isReaderOpened {
		return errors.New("reader is not opened")
	}

	var errCode C.ErrorCode
	errCode = C.ts_reader_close(r.reader)
	if int(errCode) != 0 {
		return fmt.Errorf("failed to close TSFile reader, error code: %d", int(errCode))
	}

	r.isReaderOpened = false
	return nil
}

// TsFileWriter provides an interface for writing data to a TSFile.
type TsFileWriter struct {
	writer         C.CTsFileWriter
	isWriterOpened bool
}

// Open initializes the TSFile writer.
func (w *TsFileWriter) Open(filepath string) error {
	cFilePath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cFilePath))

	var errCode C.ErrorCode
	w.writer = C.ts_writer_open(cFilePath, &errCode)

	if int(errCode) != 0 {
		return fmt.Errorf("failed to open TSFile for writing '%s', error code: %d", filepath, int(errCode))
	}

	w.isWriterOpened = true
	return nil
}

// RegisterTable registers a new table schema in the TSFile.
func (w *TsFileWriter) RegisterTable(tableName string, columns []string, schemaInfo []int64) error {
	if !w.isWriterOpened {
		return errors.New("writer is not opened")
	}

	if len(columns) != len(schemaInfo) {
		return errors.New("columns and schemaInfo length do not match")
	}

	var columnSchemas []*C.ColumnSchema
	for i := range columns {
		cColumnName := C.CString(columns[i])
		defer C.free(unsafe.Pointer(cColumnName))

		columnSchemas = append(columnSchemas, &C.ColumnSchema{
			name:       cColumnName,
			column_def: C.SchemaInfo(schemaInfo[i]),
		})
	}

	cTableSchema := &C.TableSchema{
		table_name:    C.CString(tableName),
		column_schema: (**C.ColumnSchema)(unsafe.Pointer(&columnSchemas[0])),
		column_num:    C.int(len(columns)),
	}
	defer C.free(unsafe.Pointer(cTableSchema.table_name))

	errCode := C.tsfile_register_table(w.writer, cTableSchema)
	if int(errCode) != 0 {
		return fmt.Errorf("failed to register table '%s', error code: %d", tableName, int(errCode))
	}

	return nil
}

// WriteRow writes a data row to the TSFile.
func (w *TsFileWriter) WriteRow(tableName string, timestamp int64, data map[string]interface{}) error {
	if !w.isWriterOpened {
		return errors.New("writer is not opened")
	}

	cTableName := C.CString(tableName)
	defer C.free(unsafe.Pointer(cTableName))

	rowData := C.create_tsfile_row(cTableName, C.timestamp(timestamp), C.int(len(data)))
	if rowData == nil {
		return errors.New("failed to create TSFile row")
	}
	defer C.destory_tsfile_row(rowData)

	for column, value := range data {
		cColumnName := C.CString(column)
		defer C.free(unsafe.Pointer(cColumnName))

		switch v := value.(type) {
		case int32:
			C.insert_data_into_tsfile_row_int32(rowData, cColumnName, C.int(v))
		case int64:
			C.insert_data_into_tsfile_row_int64(rowData, cColumnName, C.longlong(v))
		case float32:
			C.insert_data_into_tsfile_row_float(rowData, cColumnName, C.float(v))
		case float64:
			C.insert_data_into_tsfile_row_double(rowData, cColumnName, C.double(v))
		case bool:
			C.insert_data_into_tsfile_row_boolean(rowData, cColumnName, C.bool(v))
		default:
			return fmt.Errorf("unsupported data type for column '%s'", column)
		}
	}

	errCode := C.tsfile_write_row_data(w.writer, rowData)
	if int(errCode) != 0 {
		return fmt.Errorf("failed to write row data into TSFile, error code: %d", int(errCode))
	}

	return nil
}

func (w *TsFileWriter) WriteData(tableName string, timestamp int64, data map[string]interface{}) error {
	if !w.isWriterOpened {
		return errors.New("writer is not opened")
	}

	cTableName := C.CString(tableName)
	defer C.free(unsafe.Pointer(cTableName))

	rowData := C.create_tsfile_row(cTableName, C.timestamp(timestamp), C.int(len(data)))
	if rowData == nil {
		return errors.New("failed to create row data")
	}
	defer C.destory_tsfile_row(rowData)

	for columnName, value := range data {
		cColumnName := C.CString(columnName)
		defer C.free(unsafe.Pointer(cColumnName))

		var errCode C.ErrorCode
		switch v := value.(type) {
		case int32:
			errCode = C.insert_data_into_tsfile_row_int32(rowData, cColumnName, C.int(v))
		case float64:
			errCode = C.insert_data_into_tsfile_row_double(rowData, cColumnName, C.double(v))
		case float32:
			errCode = C.insert_data_into_tsfile_row_float(rowData, cColumnName, C.float(v))
		case bool:
			errCode = C.insert_data_into_tsfile_row_boolean(rowData, cColumnName, C.bool(v))
		default:
			return fmt.Errorf("unsupported data type for column '%s'", columnName)
		}

		if int(errCode) != 0 {
			return fmt.Errorf("failed to insert data into column '%s', error code: %d", columnName, int(errCode))
		}
	}

	errCode := C.tsfile_write_row_data(w.writer, rowData)
	if int(errCode) != 0 {
		return fmt.Errorf("failed to write row data, error code: %d", int(errCode))
	}

	return nil
}

// Close closes the TSFile writer and flushes all data.
func (w *TsFileWriter) Close() error {
	if !w.isWriterOpened {
		return errors.New("writer is not opened")
	}

	errCode := C.ts_writer_close(w.writer)
	if int(errCode) != 0 {
		return fmt.Errorf("failed to close TSFile writer, error code: %d", int(errCode))
	}

	w.isWriterOpened = false
	return nil
}
