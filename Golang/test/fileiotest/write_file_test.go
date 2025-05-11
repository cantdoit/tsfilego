package fileiotest

import (
	"Golang/internal/fileio" // Import your module for the WriteFile implementation
	"os"
	"testing"
)

// TestWriteFileCreateFile tests the creation of a file.
func TestWriteFileCreateFile(t *testing.T) {
	fileName := "write_file_test.dat"

	// Remove the file if it already exists (cleanup)
	_ = os.Remove(fileName)

	writeFile := &fileio.WriteFile{}

	// Create the file
	err := writeFile.Create(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Ensure the file is opened
	if !writeFile.IsFileOpened() {
		t.Fatal("Expected the file to be opened, but it is not.")
	}

	// Ensure the file path is correct
	if writeFile.GetFilePath() != fileName {
		t.Fatalf("Expected file path to be %s, but got %s", fileName, writeFile.GetFilePath())
	}

	// Close the file
	if err := writeFile.CloseFile(); err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}

	// Cleanup
	_ = os.Remove(fileName)
}

// TestWriteFileWriteToFile tests writing content to a file.
func TestWriteFileWriteToFile(t *testing.T) {
	fileName := "test_file_write.dat"

	// Remove the file if it already exists (cleanup)
	_ = os.Remove(fileName)

	writeFile := &fileio.WriteFile{}

	// Create the file
	err := writeFile.Create(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Ensure the file is opened
	if !writeFile.IsFileOpened() {
		t.Fatal("Expected the file to be opened, but it is not.")
	}

	// Write some content to the file
	content := []byte("Hello, World!")
	err = writeFile.Write(content, uint32(len(content)))
	if err != nil {
		t.Fatalf("Failed to write content to file: %v", err)
	}

	// Close the file
	if err := writeFile.CloseFile(); err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}

	// Verify that the written content matches
	data, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Failed to read file content: %v", err)
	}

	if string(data) != string(content) {
		t.Fatalf("Expected file content to be %q, but got %q", string(content), string(data))
	}

	// Cleanup
	_ = os.Remove(fileName)
}

// TestWriteFileSyncFile tests syncing data to disk.
func TestWriteFileSyncFile(t *testing.T) {
	fileName := "test_file_sync.dat"

	// Remove the file if it already exists (cleanup)
	_ = os.Remove(fileName)

	writeFile := &fileio.WriteFile{}

	// Create the file
	err := writeFile.Create(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Ensure the file is opened
	if !writeFile.IsFileOpened() {
		t.Fatal("Expected the file to be opened, but it is not.")
	}

	// Write some content to the file
	content := []byte("Hello, Sync!")
	err = writeFile.Write(content, uint32(len(content)))
	if err != nil {
		t.Fatalf("Failed to write content to file: %v", err)
	}

	// Sync the file to disk
	if err := writeFile.SyncFile(); err != nil {
		t.Fatalf("Failed to sync file: %v", err)
	}

	// Close the file
	if err := writeFile.CloseFile(); err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}

	// Cleanup
	_ = os.Remove(fileName)
}

// TestWriteFileCloseFile tests closing an open file.
func TestWriteFileCloseFile(t *testing.T) {
	fileName := "test_file_close.dat"

	// Remove the file if it already exists (cleanup)
	_ = os.Remove(fileName)

	writeFile := &fileio.WriteFile{}

	// Create the file
	err := writeFile.Create(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Ensure the file is opened
	if !writeFile.IsFileOpened() {
		t.Fatal("Expected the file to be opened, but it is not.")
	}

	// Write some content to the file
	content := []byte("Closing file.")
	err = writeFile.Write(content, uint32(len(content)))
	if err != nil {
		t.Fatalf("Failed to write content to file: %v", err)
	}

	// Close the file
	if err := writeFile.CloseFile(); err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}

	// Cleanup
	_ = os.Remove(fileName)
}
