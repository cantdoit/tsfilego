package common

import (
	"errors"
)

// ByteStream provides a mechanism for buffered reading and writing of binary data.
// Data is stored in Pages (fixed-size buffers), and the ByteStream dynamically allocates Pages as needed.
// It also provides wrapping functionality for handling external buffers.
type ByteStream struct {
	PageSize       uint32  // Size of each page in bytes
	TotalSize      uint32  // Total size of the byte stream (written bytes)
	ReadPos        uint32  // Current read position
	MarkedReadPos  uint32  // Marked read position for calculating offsets
	Pages          []*Page // Pages representing the stream's data
	CurrentPage    *Page   // The page currently being written to
	CurrentPageIdx uint32  // Index of the current page during reading
	WrappedBuffer  []byte  // External wrapped buffer (if any)
	IsWrapped      bool    // Whether the stream is wrapping an external buffer
}

// Page represents a fixed-size buffer used within the ByteStream.
type Page struct {
	data []byte // The data buffer of the page
}

// NewByteStream initializes a new ByteStream with a given page size.
// The page size must be greater than zero, or it returns an error.
func NewByteStream(pageSize uint32) (*ByteStream, error) {
	if pageSize == 0 {
		return nil, errors.New("page size must be greater than zero")
	}
	return &ByteStream{
		PageSize:       pageSize,
		TotalSize:      0,
		ReadPos:        0,
		MarkedReadPos:  0,
		Pages:          []*Page{},
		CurrentPage:    nil,
		CurrentPageIdx: 0,
		IsWrapped:      false,
	}, nil
}

///////////////////////////////////////
// Part 0: Wrapping External Buffers //
///////////////////////////////////////

// WrapFrom wraps the ByteStream around an external buffer.
// This corresponds to the C++ function `wrap_from`.
func (bs *ByteStream) WrapFrom(buffer []byte, size uint32) {
	bs.WrappedBuffer = buffer
	bs.TotalSize = size
	bs.ReadPos = 0
	bs.MarkedReadPos = 0
	bs.IsWrapped = true
	bs.Pages = nil // Clear any internal Pages since this is an external buffer.
}

// IsWrappedBuffer returns whether the ByteStream is currently wrapping an external buffer.
func (bs *ByteStream) IsWrappedBuffer() bool {
	return bs.IsWrapped
}

// ClearWrappedBuffer clears the ByteStream's reference to the external buffer.
func (bs *ByteStream) ClearWrappedBuffer() {
	bs.WrappedBuffer = nil
	bs.IsWrapped = false
}

//////////////////////////////////
// Part 1: Basic Functionality //
//////////////////////////////////

// RemainingSize returns the number of remaining bytes to be read.
func (bs *ByteStream) RemainingSize() uint32 {
	if bs.TotalSize < bs.ReadPos {
		panic("TotalSize cannot be smaller than ReadPos")
	}
	return bs.TotalSize - bs.ReadPos
}

// HasRemaining returns whether there are any remaining bytes to be read.
func (bs *ByteStream) hasRemaining() bool {
	return bs.RemainingSize() > 0
}

// MarkReadPos records the current read position for offset calculations.
func (bs *ByteStream) MarkReadPos() {
	bs.MarkedReadPos = bs.ReadPos
}

// GetMarkLen calculates the number of bytes read since the marked position.
func (bs *ByteStream) GetMarkLen() uint32 {
	if bs.MarkedReadPos > bs.ReadPos {
		panic("MarkedReadPos cannot be greater than ReadPos")
	}
	return bs.ReadPos - bs.MarkedReadPos
}

// Reset clears the ByteStream and reallocates all memory.
func (bs *ByteStream) Reset() {
	if bs.IsWrapped {
		// Do not free the wrapped buffer, just clear its reference
		bs.ClearWrappedBuffer()
	}

	// Free all internal Pages
	bs.Pages = nil
	bs.CurrentPage = nil

	// Reset metadata
	bs.TotalSize = 0
	bs.ReadPos = 0
	bs.MarkedReadPos = 0
	bs.CurrentPageIdx = 0
}

//////////////////////////////////////
// Part 2: Writing and Reading Data //
//////////////////////////////////////

// WriteBuf writes a buffer of uint8 data to the ByteStream, dynamically allocating memory as needed.
func (bs *ByteStream) WriteBuf(buf []uint8, bufLen uint32) error {
	if buf == nil || bufLen == 0 {
		return nil // No-op for empty buffers
	}

	writeLen := uint32(0) // Tracks the amount of bytes written in this operation

	// Write loop: ensures all data from `buf` is written, splitting across Pages if needed
	for writeLen < bufLen {
		// Prepare the current page or allocate a new one if the current is full
		if bs.CurrentPage == nil || uint32(len(bs.CurrentPage.data)) == bs.PageSize {
			bs.AllocatePage()
		}

		// Calculate the remaining space in the current page
		remainingPageSpace := bs.PageSize - uint32(len(bs.CurrentPage.data))

		// Determine how many bytes to write in this iteration
		copyLen := minUint32(bufLen-writeLen, remainingPageSpace)

		// Append the data to the current page
		bs.CurrentPage.data = append(bs.CurrentPage.data, buf[writeLen:writeLen+copyLen]...)

		// Update counters
		writeLen += copyLen
		bs.TotalSize += copyLen
	}

	return nil
}

// ReadBuf reads up to `wantLen` bytes from the ByteStream into the provided buffer.
// Returns the actual number of bytes read.
// This corresponds to the C++ function `read_buf`.
func (bs *ByteStream) ReadBuf(buf []uint8, wantLen uint32) (uint32, error) {
	if bs.ReadPos >= bs.TotalSize {
		return 0, errors.New("read position exceeds total size")
	}

	readLen := uint32(0) // Tracks the number of bytes read into `buf`

	for readLen < wantLen && bs.ReadPos < bs.TotalSize {
		// Get the current page using the read position
		page := bs.Pages[bs.CurrentPageIdx]
		pageOffset := bs.ReadPos % bs.PageSize

		// Calculate the remaining data in the current page
		remainingPageData := uint32(len(page.data)) - pageOffset

		// Determine how many bytes to copy in this iteration
		copyLen := minUint32(wantLen-readLen, remainingPageData)

		// Copy the data from the page to the buffer
		copy(buf[readLen:], page.data[pageOffset:pageOffset+copyLen])

		// Update counters
		readLen += copyLen
		bs.ReadPos += copyLen

		// Advance to the next page if needed
		if bs.ReadPos%bs.PageSize == 0 && bs.CurrentPageIdx < uint32(len(bs.Pages)-1) {
			bs.CurrentPageIdx++
		}
	}

	return readLen, nil
}

///////////////////////////////////////
// Part 3: Buffer Management Helpers //
///////////////////////////////////////

// AcquireBuffer returns a writable buffer for the next write operation.
// It does not span multiple Pages.
func (bs *ByteStream) AcquireBuffer() ([]byte, uint32) {
	if bs.CurrentPage == nil || uint32(len(bs.CurrentPage.data)) == bs.PageSize {
		bs.AllocatePage()
	}
	remainingSpace := bs.PageSize - uint32(len(bs.CurrentPage.data))
	return bs.CurrentPage.data, remainingSpace
}

// BufferUsed updates the current page state after writing data to it.
func (bs *ByteStream) BufferUsed(usedBytes uint32) {
	if usedBytes == 0 {
		panic("usedBytes must be greater than 0")
	}
	bs.TotalSize += usedBytes
}

//////////////////////////////////////////
// Part 4: Buffer Iterators and Readers //
//////////////////////////////////////////

// BufferIterator provides a way to iterate over all buffers in the ByteStream.
func (bs *ByteStream) BufferIterator() func() ([]byte, uint32, bool) {
	currentIndex := 0
	return func() ([]byte, uint32, bool) {
		if currentIndex >= len(bs.Pages) {
			return nil, 0, false
		}
		page := bs.Pages[currentIndex]
		currentIndex++
		return page.data, uint32(len(page.data)), true
	}
}

//////////////////////////////////////
// Utility Functions and Allocators //
//////////////////////////////////////

// AllocatePage creates and appends a new page to the ByteStream.
func (bs *ByteStream) AllocatePage() {
	page := &Page{
		data: make([]byte, 0, bs.PageSize),
	}
	bs.Pages = append(bs.Pages, page)
	bs.CurrentPage = page
}

// minUint32 returns the smaller of two uint32 values.
func minUint32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
