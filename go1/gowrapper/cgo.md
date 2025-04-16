/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.


package wrapper


#cgo CFLAGS: -I../../../cpp/src/cwrapper
#include "tsfile_cwrapper.h"

import "C"
import (
	"errors"
	"fmt"
	"time"
	"unsafe"
)

const (
	TimeColumn = "Time"
)

// Type mappings
var typeMapping = map[string]C.TSDataType{
	"int32":   C.TS_DATATYPE_INT32,
	"int64":   C.TS_DATATYPE_INT64,
	"float32": C.TS_DATATYPE_FLOAT,
	"float64": C.TS_DATATYPE_DOUBLE,
	"bool":    C.TS_DATATYPE_BOOLEAN,
}

// Reader represents a TsFile reader
type Reader struct {
	cReader       C.TsFileReader
	cRet          C.QueryDataRet
	batchSize     int
	readAllAtOnce bool
}

// NewReader creates a new TsFile reader
func NewReader(path string, tableName string, columns []string, startTime, endTime *time.Time, batchSize *int) (*Reader, error) {
	r := &Reader{
		batchSize:     1024,
		readAllAtOnce: true,
	}

	if batchSize != nil {
		r.batchSize = *batchSize
		r.readAllAtOnce = false
	}

	// Open the reader
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var errCode C.ErrorCode
	r.cReader = C.ts_reader_open(cPath, &errCode)
	if errCode != 0 {
		return nil, fmt.Errorf("failed to open reader (code %v)", errCode)
	}

	// Query the data
	if err := r.queryData(tableName, columns, startTime, endTime); err != nil {
		C.ts_reader_close(r.cReader)
		return nil, err
	}

	return r, nil
}

func (r *Reader) queryData(tableName string, columns []string, startTime, endTime *time.Time) error {
	cTable := C.CString(tableName)
	defer C.free(unsafe.Pointer(cTable))

	cCols := make([]*C.char, len(columns))
	for i, col := range columns {
		cCols[i] = C.CString(col)
		defer C.free(unsafe.Pointer(cCols[i]))
	}

	var cRet C.QueryDataRet
	if startTime != nil || endTime != nil {
		var start, end C.longlong
		if startTime == nil {
			start = C.LLONG_MIN
		} else {
			start = C.longlong(startTime.UnixNano() / int64(time.Millisecond))
		}
		if endTime == nil {
			end = C.LLONG_MAX
		} else {
			end = C.longlong(endTime.UnixNano() / int64(time.Millisecond))
		}

		cRet = C.ts_reader_begin_end(
			r.cReader,
			cTable,
			(**C.char)(unsafe.Pointer(&cCols[0])),
			C.int(len(columns)),
			start,
			end,
		)
	} else {
		cRet = C.ts_reader_read(
			r.cReader,
			cTable,
			(**C.char)(unsafe.Pointer(&cCols[0])),
			C.int(len(columns)),
		)
	}

	if cRet.data == nil {
		return errors.New("query returned no results")
	}

	r.cRet = cRet
	return nil
}

// Read reads all data from the TsFile
func (r *Reader) Read() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	for {
		data, err := r.nextBatch()
		if err != nil {
			return nil, err
		}
		if data == nil {
			break
		}
		results = append(results, data...)
	}

	return results, nil
}

// Next returns the next batch of data
func (r *Reader) Next() ([]map[string]interface{}, error) {
	return r.nextBatch()
}

func (r *Reader) nextBatch() ([]map[string]interface{}, error) {
	result := C.ts_next(r.cRet, C.int(r.batchSize))
	if result == nil || result.column_schema == nil {
		return nil, nil
	}
	defer C.destory_tablet(result)

	// Get column names
	colNames := make([]string, r.cRet.column_num)
	for i := 0; i < int(r.cRet.column_num); i++ {
		colNames[i] = C.GoString(*(**C.char)(unsafe.Pointer(
			uintptr(unsafe.Pointer(r.cRet.column_names)) + uintptr(i)*unsafe.Sizeof(*r.cRet.column_names),
		)))
	}

	// Process rows
	var rows []map[string]interface{}
	rowCount := int(result.cur_num)

	for row := 0; row < rowCount; row++ {
		rowData := make(map[string]interface{})

		// Add timestamp
		ts := *(*C.longlong)(unsafe.Pointer(
			uintptr(unsafe.Pointer(result.times)) + uintptr(row)*unsafe.Sizeof(*result.times),
		))
		rowData[TimeColumn] = time.Unix(0, int64(ts)*int64(time.Millisecond))

		// Add column values
		for col := 0; col < int(r.cRet.column_num); col++ {
			colName := colNames[col]
			schema := *(**C.ColumnSchema)(unsafe.Pointer(
				uintptr(unsafe.Pointer(result.column_schema)) + uintptr(col)*unsafe.Sizeof(*result.column_schema),
			))

			// Check if value is NULL
			isNull := *(*C.bint)(unsafe.Pointer(
				uintptr(unsafe.Pointer(result.bitmap[col])) + uintptr(row)*unsafe.Sizeof(*result.bitmap[col]),
			)) == 0
			if isNull {
				rowData[colName] = nil
				continue
			}

			// Get value based on type
			valPtr := unsafe.Pointer(
				uintptr(unsafe.Pointer(result.value[col])) + uintptr(row)*unsafe.Sizeof(*result.value[col]),
			)

			switch schema.column_def {
			case C.TS_DATATYPE_INT32:
				rowData[colName] = int32(*(*C.int)(valPtr))
			case C.TS_DATATYPE_INT64:
				rowData[colName] = int64(*(*C.longlong)(valPtr))
			case C.TS_DATATYPE_FLOAT:
				rowData[colName] = float32(*(*C.float)(valPtr))
			case C.TS_DATATYPE_DOUBLE:
				rowData[colName] = float64(*(*C.double)(valPtr))
			case C.TS_DATATYPE_BOOLEAN:
				rowData[colName] = *(*C.bint)(valPtr) != 0
			default:
				rowData[colName] = nil
			}
		}

		rows = append(rows, rowData)
	}

	return rows, nil
}

// Close closes the reader
func (r *Reader) Close() error {
	if r.cRet != nil {
		if C.destory_query_dataret(r.cRet) != 0 {
			return errors.New("failed to free query data")
		}
	}
	if r.cReader != nil {
		if C.ts_reader_close(r.cReader) != 0 {
			return errors.New("failed to close reader")
		}
	}
	return nil
}

// Writer represents a TsFile writer
type Writer struct {
	cWriter C.TsFileWriter
}

// NewWriter creates a new TsFile writer
func NewWriter(path string) (*Writer, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var errCode C.ErrorCode
	cWriter := C.ts_writer_open(cPath, &errCode)
	if errCode != 0 {
		return nil, fmt.Errorf("failed to open writer (code %v)", errCode)
	}

	return &Writer{cWriter: cWriter}, nil
}

// RegisterTimeseries registers a new timeseries
func (w *Writer) RegisterTimeseries(tableName, columnName, dataType string) error {
	cTable := C.CString(tableName)
	defer C.free(unsafe.Pointer(cTable))

	cCol := C.CString(columnName)
	defer C.free(unsafe.Pointer(cCol))

	typ, ok := typeMapping[dataType]
	if !ok {
		return fmt.Errorf("unsupported data type: %s", dataType)
	}

	schema := C.ColumnSchema{
		name:      cCol,
		data_type: typ,
	}

	if C.tsfile_register_table_column(w.cWriter, cTable, &schema) != 0 {
		return errors.New("failed to register timeseries")
	}

	return nil
}

// WriteRow writes a row of data
func (w *Writer) WriteRow(tableName string, timestamp time.Time, data map[string]interface{}) error {
	cTable := C.CString(tableName)
	defer C.free(unsafe.Pointer(cTable))

	ts := C.longlong(timestamp.UnixNano() / int64(time.Millisecond))
	row := C.create_tsfile_row(cTable, ts, C.int(len(data)))
	defer C.destory_tsfile_row(row)

	for col, val := range data {
		cCol := C.CString(col)
		defer C.free(unsafe.Pointer(cCol))

		switch v := val.(type) {
		case int32:
			C.insert_data_into_tsfile_row_int32(row, cCol, C.int(v))
		case int64:
			C.insert_data_into_tsfile_row_int64(row, cCol, C.longlong(v))
		case float32:
			C.insert_data_into_tsfile_row_float(row, cCol, C.float(v))
		case float64:
			C.insert_data_into_tsfile_row_double(row, cCol, C.double(v))
		case bool:
			var b C.bint
			if v {
				b = 1
			}
			C.insert_data_into_tsfile_row_boolean(row, cCol, b)
		default:
			return fmt.Errorf("unsupported data type for column %s", col)
		}
	}

	if C.tsfile_write_row_data(w.cWriter, row) != 0 {
		return errors.New("failed to write row")
	}

	return nil
}

// Flush flushes the data to disk
func (w *Writer) Flush() error {
	if C.tsfile_flush_data(w.cWriter) != 0 {
		return errors.New("failed to flush data")
	}
	return nil
}

// Close closes the writer
func (w *Writer) Close() error {
	if C.ts_writer_close(w.cWriter) != 0 {
		return errors.New("failed to close writer")
	}
	return nil
}

*/