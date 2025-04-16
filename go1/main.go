package main

/*
#cgo CFLAGS: -IC:/Users/User/Documents/GitHub/tsfilego/cpp/src
#cgo LDFLAGS: -LC:/Users/User/Documents/GitHub/tsfilego/cpp/build/lib -ltsfile
#include "cwrapper/TsFile-cwrapper.h"
*/
import "C"
import (
	"fmt"
)

func main() {
	fmt.Println("Testing TsFile C wrapper library...")

	// Test opening a TSFile reader
	var errCode C.ErrorCode
	testFile := C.CString("test.tsfile") // Replace with an actual TsFile path if you have one

	reader := C.ts_reader_open(testFile, &errCode)
	if reader == nil {
		fmt.Printf("Failed to open TsFile: error code %d\n", errCode)
	} else {
		fmt.Println("Successfully opened TsFile!")
		C.ts_reader_close(reader)
	}
}
