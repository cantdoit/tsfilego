package base

import (
	"Golang/internal/utils"
	"errors"
	"fmt"
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
func (bs *ByteStream) ReadBuf(buf []uint8, wantLen uint32) (uint32, error) {
	if bs.ReadPos >= bs.TotalSize {
		return 0, utils.GetError(utils.ErrNoMoreData)
	}

	readLen := uint32(0) // Tracks the number of bytes read into `buf`

	for readLen < wantLen && bs.ReadPos < bs.TotalSize {
		if bs.CurrentPageIdx >= uint32(len(bs.Pages)) {
			return 0, utils.GetError(utils.ErrNoMoreData)
		}
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

	// If fewer bytes were read than requested, this is considered a partial read.
	if readLen < wantLen {
		return readLen, utils.GetError(utils.ErrPartialRead)
	}

	return readLen, nil
}

// PurgePrevPages removes a specified number of pages from the beginning of the ByteStream.
// If purgePageCount is greater than the current number of pages, all pages except the last one will be purged.
func (bs *ByteStream) PurgePrevPages(purgePageCount int) {
	if len(bs.Pages) == 0 || purgePageCount <= 0 {
		return // Nothing to purge or invalid input
	}

	if purgePageCount >= len(bs.Pages) {
		// Retain only the last page
		lastPage := bs.Pages[len(bs.Pages)-1]
		bs.Pages = []*Page{lastPage}
		bs.TotalSize = uint32(len(lastPage.data))
		bs.ReadPos = 0
		bs.CurrentPageIdx = 0
		return
	}

	// Remove pages and adjust metadata
	bs.Pages = bs.Pages[purgePageCount:]
	bs.TotalSize -= uint32(purgePageCount) * bs.PageSize

	if bs.ReadPos >= bs.PageSize*uint32(purgePageCount) {
		bs.ReadPos -= bs.PageSize * uint32(purgePageCount)
	} else {
		bs.ReadPos = 0
	}
	bs.CurrentPageIdx = 0
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
			return nil, 0, false // No more pages
		}
		page := bs.Pages[currentIndex]
		currentIndex++
		// Return the actual portion of page data that was written, not full capacity
		dataLen := uint32(len(page.data))
		return page.data[:dataLen], dataLen, true

	}
}

// InitBufferIterator resets the internal buffer iterator to the beginning of the ByteStream.
// It should be called before starting a new iteration.
func (bs *ByteStream) InitBufferIterator() {
	bs.CurrentPageIdx = 0 // Reset the current page index to the beginning
}

// GetNextBuffer retrieves the buffer and its length from the current page in the ByteStream.
// Advances the internal iterator to the next page. Returns nil if there are no more pages.
func (bs *ByteStream) GetNextBuffer() ([]byte, uint32, error) {
	// Check if there are any remaining pages to process
	if int(bs.CurrentPageIdx) >= len(bs.Pages) {
		return nil, 0, errors.New("no more buffers available")
	}

	// Get the current page
	page := bs.Pages[bs.CurrentPageIdx]

	// Advance to the next page
	bs.CurrentPageIdx++

	// Return the buffer and its length
	return page.data, uint32(len(page.data)), nil
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

// MergeByteStream merges the contents of the "river" ByteStream into the "sea" ByteStream.
// If purgeRiver is true, it clears the read sections of the river after they're merged.
func (bs *ByteStream) MergeByteStream(sea *ByteStream, river *ByteStream, purgeRiver bool) error {
	// Initialize the buffer iterator for the river
	river.InitBufferIterator()

	// Iterate through all buffers in the river
	for {
		buf, length, err := river.GetNextBuffer()
		if err != nil {
			// No more buffers to process
			break
		}

		// Write the buffer from the river into the sea
		if err := sea.WriteBuf(buf, length); err != nil {
			return err // Abort if write to the sea fails
		}

		// Optionally purge the old page in the river
		if purgeRiver {
			river.PurgePrevPages(1)
		}
	}

	return nil
}

// CopyBSToBuffer copies the contents of the ByteStream into the provided buffer (`destBuf`) up to the buffer's length.
// Returns an error if the buffer is not large enough.
func (bs *ByteStream) CopyBSToBuffer(stream *ByteStream, destBuf []byte, destBufLen uint32) error {
	it := stream.BufferIterator() // Initialize the buffer iterator
	var copiedLen uint32 = 0

	for {
		// Get the next buffer from the iterator
		buf, length, ok := it()
		if !ok {
			// No more data in the byte stream, finish copying
			break
		}
		// Check if there's enough space in the destination buffer
		if destBufLen-copiedLen < length {
			return utils.GetError(utils.ErrBufNotEnough)
		}

		// Copy the data from the ByteStream buffer to the destination buffer
		copy(destBuf[copiedLen:], buf[:length])

		// Update the copied length
		copiedLen += length
	}

	return utils.GetError(utils.ErrOk)
}

// GetBytesFromByteStream reads all data from the internal ByteStream and returns it as a single byte array.
// If there is no data, or if memory allocation fails, it returns an error.
func (bs *ByteStream) GetBytesFromByteStream() ([]byte, error) {
	// Check if the stream has data
	if bs.TotalSize == 0 {
		return nil, nil // No data in the stream
	}

	// Allocate a buffer to hold the entire byte stream
	retBuf := make([]byte, bs.TotalSize)
	if retBuf == nil {
		return nil, errors.New("memory allocation failed for retBuf")
	}

	// Iterate through each page in the ByteStream and copy data to the retBuf
	offset := 0
	for _, page := range bs.Pages {
		if len(page.data) == 0 {
			break // No more data in this page
		}
		copy(retBuf[offset:], page.data)
		offset += len(page.data)
	}

	// Ensure the entire TotalSize is copied
	if offset != int(bs.TotalSize) {
		return nil, fmt.Errorf("unexpected offset mismatch: expected %d, got %d", bs.TotalSize, offset)
	}

	return retBuf, nil
}

// deserializeBufNotEnough checks if the buffer state indicates an out-of-range or partial-read error.
func deserializeBufNotEnough(ret int) bool {
	return nil != utils.GetError(utils.ErrOutOfRange) || nil != utils.GetError(utils.ErrPartialRead)
}

func GetVarUintSize(ui32 uint32) uint32 {
	var bytes uint32 = 0
	for (ui32 & 0xFFFFFF80) != 0 { // While more than 7 bits are set
		bytes++
		ui32 = ui32 >> 7 // Right-shift by 7 bits
	}
	return bytes + 1 // Add 1 for the last byte

}
