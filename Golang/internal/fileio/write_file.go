package fileio

import (
	"fmt"
	"os"
	_ "syscall"
)

/*
WriteBuf from buffer to file
(from io writer to here)
Used to handle high level file operations
Create file, sync file, write to file, close file
*/

// WriteFile represents a wrapper around system file operations for handling high-level file functionality
type WriteFile struct {
	fd   *os.File // A reference to the open file
	path string   // Path of the file
}

// Create opens a new file or ensures the file doesn’t already exist
// - Parameters:
//   - filePath: Absolute or relative path of the file to create
//   - flags: Options passed for file opening (O_RDWR, O_CREATE, etc.)
//     (Use constants from the `os` or `syscall` package)
//   - mode: File permissions (e.g., 0644)
//
// - Returns: Error if any, or nil if creation succeeds
func (wf *WriteFile) Create(filePath string, flags int, mode os.FileMode) error {
	// First, check if the file descriptor is already in use
	if wf.fd != nil {
		return fmt.Errorf("file already open: fd=%v, path=%s", wf.fd, wf.path)
	}

	// Assign file path to internal state
	wf.path = filePath

	// Open the file with the provided flags and mode
	file, err := os.OpenFile(filePath, flags, mode)
	if err != nil {
		// Handle errors specific to file opening
		return fmt.Errorf("failed to open file: path=%s, error=%v", filePath, err)
	}

	// Set file descriptor on success
	wf.fd = file
	return nil
}

// Internal function to handle the actual creation of the file
func (wf *WriteFile) doCreate(flags int, mode os.FileMode) error {
	file, err := os.OpenFile(wf.path, flags, mode)
	if err != nil {
		return fmt.Errorf("failed to open file: path=%s, error=%v", wf.path, err)
	}

	wf.fd = file
	return nil
}

// Write writes the provided buffer to the file
func (wf *WriteFile) Write(buf []byte, len uint32) error {
	if wf.fd == nil {
		return fmt.Errorf("file is not open: path=%s", wf.path)
	}

	var totalWritten = uint32(0)

	for totalWritten < len {
		written, err := wf.fd.Write(buf[totalWritten:])
		if err != nil {
			return fmt.Errorf("failed to write to file: path=%s, error=%v", wf.path, err)
		}
		totalWritten += uint32(written)
	}

	return nil
}

// Sync flushes any written data to disk
func (wf *WriteFile) Sync() error {
	if wf.fd == nil {
		return fmt.Errorf("file is not open: path=%s", wf.path)
	}

	// Call filesystem sync operation to ensure data is persisted
	err := wf.fd.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file: path=%s, error=%v", wf.path, err)
	}

	return nil
}

// Close closes the file descriptor and releases resources
func (wf *WriteFile) Close() error {
	if wf.fd == nil {
		return fmt.Errorf("file is not open: path=%s", wf.path)
	}

	// Close the file descriptor
	err := wf.fd.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: path=%s, error=%v", wf.path, err)
	}

	// Reset internal state
	wf.fd = nil
	wf.path = ""

	return nil
}

// IsFileOpened checks whether the file is currently open
func (wf *WriteFile) IsFileOpened() bool {
	return wf.fd != nil
}

// GetFilePath returns the currently loaded file path
func (wf *WriteFile) GetFilePath() string {
	return wf.path
}
