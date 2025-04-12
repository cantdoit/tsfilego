package gowrapper

/*
#cgo CFLAGS: -I${SRCDIR}/../../cwrapper
#cgo LDFLAGS: -L${SRCDIR}/../../cwrapper -ltsfile

#include "TsFile-cwrapper.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type CTsFileReader C.CTsFileReader

func TsReaderOpen(path string) (*CTsFileReader, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var errCode C.ErrorCode
	cReader := C.ts_reader_open(cPath, &errCode)
	if errCode != 0 {
		return nil, fmt.Errorf("C error code: %d", errCode)
	}
	return (*CTsFileReader)(cReader), nil
}

func TsReaderClose(reader *CTsFileReader) error {
	if C.ts_reader_close((*C.CTsFileReader)(reader)) != 0 {
		return fmt.Errorf("failed to close reader")
	}
	return nil
}
