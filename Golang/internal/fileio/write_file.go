package fileio

import (
	"fmt"
	"os"
	_ "syscall"
)

/*
WriteBuf from buffer to File
(from io writer to here)
Used to handle high level File operations
Create File, sync File, write to File, close File
*/

// WriteFile represents a wrapper around system File operations for handling high-level File functionality
type WriteFile struct {
	fd   *os.File // A reference to the open File
	path string   // Path of the File
}

// Create opens a new File or ensures the File doesn’t already exist
func (wf *WriteFile) Create(filePath string, flags int, mode os.FileMode) error {
	// First, check if the File descriptor is already in use
	// fmt.Println("Creating File...")
	if wf.fd != nil {
		return fmt.Errorf("File already open: fd=%v, path=%s", wf.fd, wf.path)
	}

	// Assign File path to internal state
	wf.path = filePath

	// Open the File with the provided flags and mode
	file, err := os.OpenFile(filePath, flags, mode)
	if err != nil {
		// Handle errors specific to File opening
		return fmt.Errorf("failed to open File: path=%s, error=%v", filePath, err)
	}

	// Set File descriptor on success
	wf.fd = file
	// fmt.Println("File created successfully.", wf.fd)
	return nil
}

// Internal function to handle the actual creation of the File
func (wf *WriteFile) doCreate(flags int, mode os.FileMode) error {
	file, err := os.OpenFile(wf.path, flags, mode)
	if err != nil {
		return fmt.Errorf("failed to open File: path=%s, error=%v", wf.path, err)
	}

	wf.fd = file
	return nil
}

// Write writes the provided buffer to the File
func (wf *WriteFile) Write(buf []byte, len uint32) error {

	// fmt.Println("Writing to File...")
	if wf.fd == nil {
		return fmt.Errorf("file is not open: path=%s", wf.path)
	}

	var totalWritten = uint32(0)

	for totalWritten < len {
		written, err := wf.fd.Write(buf[totalWritten:])
		if err != nil {
			return fmt.Errorf("failed to write to File: path=%s, error=%v", wf.path, err)
		}
		totalWritten += uint32(written)
	}

	return nil
}

// SyncFile flushes any written data to disk
func (wf *WriteFile) SyncFile() error {
	if wf.fd == nil {
		return fmt.Errorf("File is not open: path=%s", wf.path)
	}

	// Call filesystem sync operation to ensure data is persisted
	err := wf.fd.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync File: path=%s, error=%v", wf.path, err)
	}

	return nil
}

// CloseFile closes the File descriptor and releases resources
func (wf *WriteFile) CloseFile() error {
	if wf.fd == nil {
		return fmt.Errorf("File is not open: path=%s", wf.path)
	}

	// Close the File descriptor
	err := wf.fd.Close()
	if err != nil {
		return fmt.Errorf("failed to close File: path=%s, error=%v", wf.path, err)
	}

	// Reset internal state
	wf.fd = nil
	wf.path = ""

	return nil
}

// IsFileOpened checks whether the File is currently open
func (wf *WriteFile) IsFileOpened() bool {
	return wf.fd != nil
}

// GetFilePath returns the currently loaded File path
func (wf *WriteFile) GetFilePath() string {
	return wf.path
}
