package tsfile

import (
	"tsfile/go/tsfile/internal/cwrapper"
	"tsfile/go/tsfile/internal/types"
)

type Reader struct {
	cReader *cwrapper.CReader
}

func NewReader(path string) (*Reader, error) {
	cReader, err := cwrapper.NewReader(path)
	if err != nil {
		return nil, err
	}
	return &Reader{cReader: cReader}, nil
}

func (r *Reader) Query(table string, columns []string, start, end int64) (*types.ResultSet, error) {
	return cwrapper.Query(r.cReader, table, columns, start, end)
}

func (r *Reader) Close() error {
	// Calls C.tsfile_reader_close through cwrapper
}
