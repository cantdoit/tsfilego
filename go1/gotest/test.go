package main

import (
	"fmt"
	_ "log"
	"os"
	"testing"
	"time"
)

func TestWriteAndReadTsFile(t *testing.T) {
	// Setup test file path
	testFile := "test_write_read.tsfile"
	defer os.Remove(testFile) // Clean up after test

	// 1. First write test data to the file
	err := gowarpper.NewTsFile()
	if err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// 2. Then read back the data we just wrote
	err = verifyTestData(testFile, t)
	if err != nil {
		t.Fatalf("Failed to verify test data: %v", err)
	}
}

// writeTestData writes sample data to a TsFile
func writeTestData(filePath string) error {
	tf, err := OpenForWriting(filePath)
	if err != nil {
		return fmt.Errorf("failed to open writer: %w", err)
	}
	defer tf.Close()

	// Register our timeseries schema
	schema := map[string]string{
		"temperature": "float64",
		"humidity":    "float32",
		"active":      "bool",
		"counter":     "int32",
	}

	for name, dtype := range schema {
		if err := tf.RegisterTimeseries("root.test", name, dtype); err != nil {
			return fmt.Errorf("failed to register %s: %w", name, err)
		}
	}

	// Write some test data
	testTime := time.Now()
	testData := []map[string]interface{}{
		{
			"temperature": 23.5,
			"humidity":    float32(45.7),
			"active":      true,
			"counter":     int32(1),
		},
		{
			"temperature": 24.1,
			"humidity":    float32(46.2),
			"active":      false,
			"counter":     int32(2),
		},
		{
			"temperature": 22.9,
			"humidity":    float32(44.9),
			"active":      true,
			"counter":     int32(3),
		},
	}

	for _, data := range testData {
		if err := tf.WriteRow("root.test", testTime, data); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
		testTime = testTime.Add(time.Second)
	}

	// Ensure data is flushed to disk
	if err := tf.Flush(); err != nil {
		return fmt.Errorf("failed to flush data: %w", err)
	}

	return nil
}

// verifyTestData reads back the test data and verifies it
func verifyTestData(filePath string, t *testing.T) error {
	tf, err := OpenForReading(
		filePath,
		"root.test",
		[]string{"temperature", "humidity", "active", "counter"},
		nil, nil, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to open reader: %w", err)
	}
	defer tf.Close()

	// Read all data
	data, err := tf.Read()
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	// Verify we got 3 rows as we wrote
	if len(data) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(data))
	}

	// Check the values in the first row
	row := data[0]
	if temp, ok := row["temperature"].(float64); !ok || temp != 23.5 {
		t.Errorf("Unexpected temperature value: %v", row["temperature"])
	}
	if humid, ok := row["humidity"].(float32); !ok || humid != 45.7 {
		t.Errorf("Unexpected humidity value: %v", row["humidity"])
	}
	if active, ok := row["active"].(bool); !ok || !active {
		t.Errorf("Unexpected active value: %v", row["active"])
	}
	if counter, ok := row["counter"].(int32); !ok || counter != 1 {
		t.Errorf("Unexpected counter value: %v", row["counter"])
	}

	// Check timestamp exists and is reasonable
	if _, ok := row["Time"].(time.Time); !ok {
		t.Error("Missing timestamp in row")
	}

	return nil
}

// BenchmarkWriteRead measures performance of write+read operations
func BenchmarkWriteRead(b *testing.B) {
	testFile := "benchmark.ts"
	defer os.Remove(testFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Write test data
		err := writeTestData(testFile)
		if err != nil {
			b.Fatalf("Write failed: %v", err)
		}

		// Read test data
		tf, err := OpenForReading(testFile, "root.test",
			[]string{"temperature", "humidity", "active", "counter"}, nil, nil, nil)
		if err != nil {
			b.Fatalf("Open for read failed: %v", err)
		}

		_, err = tf.Read()
		if err != nil {
			b.Fatalf("Read failed: %v", err)
		}
		tf.Close()
	}
}

func main() {
	TestWriteAndReadTsFile()
}
