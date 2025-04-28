package utils

import "fmt"

/*
Package for defining error numbers, their meanings, and mappings for Go error handling
*/

// GetError translates an integer error code to its corresponding `error`.
func GetError(code int) error {
	if err, exists := errorMap[code]; exists {
		return err
	}
	return fmt.Errorf("unknown error code: %d", code)
}

// Integer error codes
const (
	ErrOk                       = 0  // E_OK
	ErrOOM                      = 1  // E_OOM
	ErrNotExist                 = 2  // E_NOT_EXIST
	ErrAlreadyExist             = 3  // E_ALREADY_EXIST
	ErrInvalidArg               = 4  // E_INVALID_ARG
	ErrOutOfRange               = 5  // E_OUT_OF_RANGE
	ErrPartialRead              = 6  // E_PARTIAL_READ
	ErrNetBindErr               = 7  // E_NET_BIND_ERR
	ErrNetSocketErr             = 8  // E_NET_SOCKET_ERR
	ErrNetEpollErr              = 9  // E_NET_EPOLL_ERR
	ErrNetEpollWaitErr          = 10 // E_NET_EPOLL_WAIT_ERR
	ErrNetRecvErr               = 11 // E_NET_RECV_ERR
	ErrNetAcceptErr             = 12 // E_NET_ACCEPT_ERR
	ErrNetFcntlErr              = 13 // E_NET_FCNTL_ERR
	ErrNetListenErr             = 14 // E_NET_LISTEN_ERR
	ErrNetSendErr               = 15 // E_NET_SEND_ERR
	ErrPipeErr                  = 16 // E_PIPE_ERR
	ErrThreadCreateErr          = 17 // E_THREAD_CREATE_ERR
	ErrMutexErr                 = 18 // E_MUTEX_ERR
	ErrCondErr                  = 19 // E_COND_ERR
	ErrOverflow                 = 20 // E_OVERFLOW
	ErrNoMoreData               = 21 // E_NO_MORE_DATA
	ErrOutOfOrder               = 22 // E_OUT_OF_ORDER
	ErrTsBlockTypeNotSupported  = 23 // E_TSBLOCK_TYPE_NOT_SUPPORTED
	ErrTsBlockDataInconsistency = 24 // E_TSBLOCK_DATA_INCONSISTENCY
	ErrDDLUnknownType           = 25 // E_DDL_UNKNOWN_TYPE
	ErrTypeNotSupported         = 26 // E_TYPE_NOT_SUPPORTED
	ErrTypeNotMatch             = 27 // E_TYPE_NOT_MATCH
	ErrFileOpenErr              = 28 // E_FILE_OPEN_ERR
	ErrFileCloseErr             = 29 // E_FILE_CLOSE_ERR
	ErrFileWriteErr             = 30 // E_FILE_WRITE_ERR
	ErrFileReadErr              = 31 // E_FILE_READ_ERR
	ErrFileSyncErr              = 32 // E_FILE_SYNC_ERR
	ErrTsFileWriterMetaErr      = 33 // E_TSFILE_WRITER_META_ERR
	ErrFileStatErr              = 34 // E_FILE_STAT_ERR
	ErrTsFileCorrupted          = 35 // E_TSFILE_CORRUPTED
	ErrBufNotEnough             = 36 // E_BUF_NOT_ENOUGH
	ErrInvalidPath              = 37 // E_INVALID_PATH
	ErrNotMatch                 = 38 // E_NOT_MATCH
	ErrJsonInvalid              = 39 // E_JSON_INVALID
	ErrNotSupport               = 40 // E_NOT_SUPPORT
	ErrParserErr                = 41 // E_PARSER_ERR
	ErrAnalyzeErr               = 42 // E_ANALYZE_ERR
	ErrInvalidDataPoint         = 43 // E_INVALID_DATA_POINT
	ErrDeviceNotExist           = 44 // E_DEVICE_NOT_EXIST
	ErrMeasurementNotExist      = 45 // E_MEASUREMENT_NOT_EXIST
	ErrInvalidQuery             = 46 // E_INVALID_QUERY
	ErrSdkQueryOptimizeErr      = 47 // E_SDK_QUERY_OPTIMIZE_ERR
	ErrCompressErr              = 48 // E_COMPRESS_ERR
	ErrBufferNotEnough          = 49 // E_BUFFER_NOT_ENOUGH
)

// Map integer error codes to Go `error` objects
var errorMap = map[int]error{
	ErrOk:                       nil,
	ErrOOM:                      fmt.Errorf("out of memory"),
	ErrNotExist:                 fmt.Errorf("not exist"),
	ErrAlreadyExist:             fmt.Errorf("already exist"),
	ErrInvalidArg:               fmt.Errorf("invalid argument"),
	ErrOutOfRange:               fmt.Errorf("out of range"),
	ErrPartialRead:              fmt.Errorf("partial read"),
	ErrNetBindErr:               fmt.Errorf("network bind error"),
	ErrNetSocketErr:             fmt.Errorf("network socket error"),
	ErrNetEpollErr:              fmt.Errorf("network epoll error"),
	ErrNetEpollWaitErr:          fmt.Errorf("network epoll wait error"),
	ErrNetRecvErr:               fmt.Errorf("network receive error"),
	ErrNetAcceptErr:             fmt.Errorf("network accept error"),
	ErrNetFcntlErr:              fmt.Errorf("network fcntl error"),
	ErrNetListenErr:             fmt.Errorf("network listen error"),
	ErrNetSendErr:               fmt.Errorf("network send error"),
	ErrPipeErr:                  fmt.Errorf("pipe error"),
	ErrThreadCreateErr:          fmt.Errorf("thread creation error"),
	ErrMutexErr:                 fmt.Errorf("mutex error"),
	ErrCondErr:                  fmt.Errorf("condition variable error"),
	ErrOverflow:                 fmt.Errorf("overflow error"),
	ErrNoMoreData:               fmt.Errorf("no more data"),
	ErrOutOfOrder:               fmt.Errorf("out of order"),
	ErrTsBlockTypeNotSupported:  fmt.Errorf("unsupported TS block type"),
	ErrTsBlockDataInconsistency: fmt.Errorf("TS block data inconsistency"),
	ErrDDLUnknownType:           fmt.Errorf("unknown DDL type"),
	ErrTypeNotSupported:         fmt.Errorf("type not supported"),
	ErrTypeNotMatch:             fmt.Errorf("type not match"),
	ErrFileOpenErr:              fmt.Errorf("file open error"),
	ErrFileCloseErr:             fmt.Errorf("file close error"),
	ErrFileWriteErr:             fmt.Errorf("file write error"),
	ErrFileReadErr:              fmt.Errorf("file read error"),
	ErrFileSyncErr:              fmt.Errorf("file sync error"),
	ErrTsFileWriterMetaErr:      fmt.Errorf("TS file writer metadata error"),
	ErrFileStatErr:              fmt.Errorf("file stat error"),
	ErrTsFileCorrupted:          fmt.Errorf("TS file corrupted"),
	ErrBufNotEnough:             fmt.Errorf("buffer not enough"),
	// Other errors to be added
}
