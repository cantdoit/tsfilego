package cwrapper

/*
#cgo CFLAGS: -I${SRCDIR}/../../../../cpp/src/cwrapper
#cgo LDFLAGS: -L${SRCDIR}/../../../../cpp/src/cwrapper -ltsfile
#include "TsFile-cwrapper.h"
*/
import "C"
import (
	"tsfile/go/tsfile/internal/types"
	"unsafe"
)

type (
	CReader = C.TsFileReader
	CWriter = C.TsFileWriter
	CTablet = C.Tablet
)

const (
	TypeBoolean = C.TS_DATATYPE_BOOLEAN
	TypeInt32   = C.TS_DATATYPE_INT32
	TypeInt64   = C.TS_DATATYPE_INT64
	TypeFloat   = C.TS_DATATYPE_FLOAT
	TypeDouble  = C.TS_DATATYPE_DOUBLE
)

func NewReader(path string) (*CReader, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	var errCode C.ERRNO
	reader := C.tsfile_reader_new(cpath, &errCode)
	if errCode != 0 {
		return nil, types.ErrorFromCode(int(errCode))
	}
	return (*CReader)(reader), nil
}

func Query(reader *CReader, table string, columns []string, start, end int64) (*types.ResultSet, error) {
	// Implementation using C.tsfile_reader_query_table
	// Handles conversion between Go and C types
}
