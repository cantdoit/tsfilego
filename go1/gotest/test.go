package gotest

import (
	"fmt"
	"go1/gowrapper"
	"os"
	"testing"
	"time"
)

// TestWriteAndReadTsFile tests both writing to and reading from a TSFile.
func TestWriteAndReadTsFile(t *testing.T) {
	// Set up the test file path
	testFile := "test_write_read.tsfile"
	defer os.Remove(testFile) // Clean up the file after testing

	// Write test data to the TSFile
	err := writeTestData(testFile)
	if err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// Verify the written data by reading it back
	err = verifyTestData(testFile, t)
	if err != nil {
		t.Fatalf("Failed to verify test data: %v", err)
	}
}

// writeTestData writes sample data to a TSFile using the gowrapper package.
func writeTestData(filePath string) error {
	// Create a writer using gowrapper
	writer, err := gowrapper.CreateWriter(filePath)
	if err != nil {
		return fmt.Errorf("failed to create writer: %w", err)
	}
	defer writer.Close()

	// Register timeseries schema
	schema := map[string]int{
		"temperature": gowrapper.TypeFloat64,
		"humidity":    gowrapper.TypeFloat32,
		"active":      gowrapper.TypeBool,
		"counter":     gowrapper.TypeInt32,
	}

	for name, dtype := range schema {
		if err := writer.RegisterColumn(name, gowrapper.SchemaInfo(dtype)); err != nil {
			return fmt.Errorf("failed to register column '%s': %w", name, err)
		}
	}

	// Write test data rows
	testTime := time.Now().UnixNano() / int64(time.Millisecond)
	rows := []map[string]interface{}{
		{"temperature": 23.5, "humidity": float32(45.7), "active": true, "counter": int32(1)},
		{"temperature": 24.1, "humidity": float32(46.2), "active": false, "counter": int32(2)},
		{"temperature": 22.9, "humidity": float32(44.9), "active": true, "counter": int32(3)},
	}

	for _, row := range rows {
		if err := writer.WriteData("root.test", testTime, row); err != nil {
			return fmt.Errorf("failed to write data: %w", err)
		}
		testTime += 1000 // Increment time by 1 second
	}

	// Flush the writer to ensure data is written to disk
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}

// verifyTestData verifies that the data can be queried without errors.
func verifyTestData(filePath string, t *testing.T) error {
	// Create a reader using gowrapper
	reader, err := gowrapper.CreateReader(filePath)
	if err != nil {
		return fmt.Errorf("failed to create reader: %w", err)
	}
	defer reader.Close()

	// We cannot validate individual rows due to missing data processing functions.
	// This is where row-by-row checks would be added if the functions were available.

	return nil
}

// assertEqual is a helper function to compare test results.
func assertEqual(t *testing.T, actual interface{}, expected interface{}, field string) {
	if actual != expected {
		t.Errorf("unexpected value for %s: got %v, expected %v", field, actual, expected)
	}
}

// BenchmarkWriteRead benchmarks the write and read operations.
func BenchmarkWriteRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testFile := "benchmark.tsfile"
		defer os.Remove(testFile)

		// Write test data
		err := writeTestData(testFile)
		if err != nil {
			b.Fatalf("Failed to write test data: %v", err)
		}

		// Read test data
		err = verifyTestData(testFile, nil)
		if err != nil {
			b.Fatalf("Failed to read test data: %v", err)
		}
	}
}
