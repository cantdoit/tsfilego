package gowrapper

import (
	"fmt"
)

type TsFileReader struct {
	cReader CTsFileReader
}

func NewReader(path string) (*TsFileReader, error) {
	reader, errCode := TsReaderOpen(path)
	if errCode != 0 {
		return nil, fmt.Errorf("failed to open TsFile (code %d)", errCode)
	}
	return &TsFileReader{cReader: reader}, nil
}

func (r *TsFileReader) Close() error {
	if errCode := TsReaderClose(r.cReader); errCode != 0 {
		return fmt.Errorf("failed to close reader (code %d)", errCode)
	}
	return nil
}
