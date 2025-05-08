package writertest

import (
	"Golang/internal/common/base"
	"Golang/internal/utils"
	"Golang/internal/writer"
	"testing"
)

// TestChunkWriterInitWithParameters tests initialization with direct parameters.
func TestChunkWriterInitWithParameters(t *testing.T) {
	chunkWriter := writer.ChunkWriter{}
	// t.Log(chunkWriter)
	// Initialize ChunkWriter with test parameters
	err := chunkWriter.Initialize("test_measurement", base.DOUBLE, base.PLAIN, base.UNCOMPRESSED)
	if err != nil {
		t.Errorf("ChunkWriter initialization failed: %v", err)
	}
	// t.Log(chunkWriter)

	// Destroy to clean up resources
	chunkWriter.Destroy()
}

// TestChunkWriterWriteBoolean tests writing a boolean value to a double-only chunk.
func TestChunkWriterWriteBoolean(t *testing.T) {
	chunkWriter := writer.ChunkWriter{}
	err := chunkWriter.Initialize("test_measurement", base.DOUBLE, base.PLAIN, base.UNCOMPRESSED)
	if err != nil {
		t.Fatalf("ChunkWriter initialization failed: %v", err)
	}
	// t.Log(chunkWriter.PageWriter)
	// Writing a boolean should result in a type mismatch error
	err = chunkWriter.Write(1234567890, true)
	if err == nil || 27 != utils.ErrTypeNotMatch {
		t.Errorf("Expected type mismatch error when writing boolean, but got: %v", err)
	}
	// t.Log(chunkWriter.PageWriter)
	chunkWriter.Destroy()
}

// TestChunkWriterWriteInt32 tests writing an int32 value to a double-only chunk.
func TestChunkWriterWriteInt32(t *testing.T) {
	chunkWriter := writer.ChunkWriter{}
	err := chunkWriter.Initialize("test_measurement", base.DOUBLE, base.PLAIN, base.UNCOMPRESSED)
	if err != nil {
		t.Fatalf("ChunkWriter initialization failed: %v", err)
	}

	// Writing an int32 should result in a type mismatch error
	err = chunkWriter.Write(1234567890, int32(42))
	if err == nil || 27 != utils.ErrTypeNotMatch {
		t.Errorf("Expected type mismatch error when writing int32, but got: %v", err)
	}

	chunkWriter.Destroy()
}

// TestChunkWriterWriteDouble tests writing a double value to a double-only chunk.
func TestChunkWriterWriteDouble(t *testing.T) {
	chunkWriter := writer.ChunkWriter{}
	err := chunkWriter.Initialize("test_measurement", base.DOUBLE, base.PLAIN, base.UNCOMPRESSED)
	if err != nil {
		t.Fatalf("ChunkWriter initialization failed: %v", err)
	}

	// Writing a double should succeed
	err = chunkWriter.Write(1234567890, float64(42.0))
	if err != nil {
		t.Errorf("Expected successful write of double value, but got error: %v", err)
	}

	chunkWriter.Destroy()
}

// TestChunkWriterWriteLargeDataSet tests writing a large set of data points to a chunk.
func TestChunkWriterWriteLargeDataSet(t *testing.T) {
	chunkWriter := writer.ChunkWriter{}
	err := chunkWriter.Initialize("test_measurement", base.DOUBLE, base.PLAIN, base.UNCOMPRESSED)
	if err != nil {
		t.Fatalf("ChunkWriter initialization failed: %v", err)
	}

	// Write a large dataset (10,000 points)
	for i := 0; i < 10000; i++ {
		err = chunkWriter.Write(int64(i), float64(i)*0.1)
		if err != nil {
			t.Fatalf("Failed to write data at index %d: %v", i, err)
		}
	}

	// Verify the count in the chunk's statistics
	if chunkWriter.PageWriter.PointCount != 10000 {
		t.Errorf("Expected 10000 data points, but got %d", chunkWriter.PageWriter.PointCount)
	}

	chunkWriter.Destroy()
}

// TestChunkWriterEndEncodeChunk tests finalizing the chunk encoding.
func TestChunkWriterEndEncodeChunk(t *testing.T) {
	chunkWriter := writer.ChunkWriter{}
	err := chunkWriter.Initialize("test_measurement", base.DOUBLE, base.PLAIN, base.UNCOMPRESSED)
	if err != nil {
		t.Fatalf("ChunkWriter initialization failed: %v", err)
	}
	t.Log(chunkWriter.PageWriter)
	// Write some data
	err = chunkWriter.Write(1234567890, float64(42.0))
	if err != nil {
		t.Fatalf("Failed to write data point: %v", err)
	}
	t.Log(chunkWriter.PageWriter.Statistic)
	// End the encoding of the chunk
	err = chunkWriter.EndEncodeChunk()
	if err != nil {
		t.Errorf("Failed to end chunk encoding: %v", err)
	}

	// Verify the chunk data size is greater than zero
	if chunkWriter.ChunkData.TotalSize <= 0 {
		t.Errorf("Expected non-zero chunk data size, but got %d", chunkWriter.ChunkData.TotalSize)
	}

	chunkWriter.Destroy()
}

// TestChunkWriterDestroy tests proper cleanup and resource release via Destroy.
func TestChunkWriterDestroy(t *testing.T) {
	chunkWriter := writer.ChunkWriter{}
	err := chunkWriter.Initialize("test_measurement", base.DOUBLE, base.PLAIN, base.UNCOMPRESSED)
	if err != nil {
		t.Fatalf("ChunkWriter initialization failed: %v", err)
	}
	t.Log(chunkWriter.ChunkData)
	// Write some data
	err = chunkWriter.Write(1234567890, float64(42.0))
	if err != nil {
		t.Fatalf("Failed to write data point: %v", err)
	}
	t.Log(chunkWriter.ChunkData)
	// Destroy the ChunkWriter
	chunkWriter.Destroy()
	// Ensure the statistics are reset
	if chunkWriter.ChunkStatistic != nil {
		t.Errorf("Expected statistics to be nil after destroy, but got %+v", chunkWriter.ChunkStatistic)
	}
	t.Log(chunkWriter.ChunkData)
	// Ensure the chunk data is cleared
	if chunkWriter.ChunkData.TotalSize != 0 {
		t.Errorf("Expected chunk data size to be zero after destroy, but got %d", chunkWriter.ChunkData.TotalSize)
	}
}
