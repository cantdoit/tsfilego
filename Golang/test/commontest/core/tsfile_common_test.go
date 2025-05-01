package core

import (
	"Golang/internal/common/base"
	"Golang/internal/common/core"
	"Golang/internal/common/statistic"
	"Golang/internal/utils"
	"testing"
)

/////////////////
// page header //
/////////////////

// TestPageHeaderDefaultConstructor verifies the default initialization of PageHeader.
func TestPageHeaderDefaultConstructor(t *testing.T) {
	header := core.NewPageHeader() // Create a new PageHeader instance.

	if header.UncompressedSize != 0 {
		t.Errorf("UncompressedSize mismatch. Expected: 0, Got: %d", header.UncompressedSize)
	}

	if header.CompressedSize != 0 {
		t.Errorf("CompressedSize mismatch. Expected: 0, Got: %d", header.CompressedSize)
	}

	if header.Statistic != nil {
		t.Errorf("Statistic mismatch. Expected: nil, Got: %v", header.Statistic)
	}
}

// TestPageHeaderReset verifies that calling Reset clears all fields of PageHeader.
func TestPageHeaderReset(t *testing.T) {
	// Create a new PageHeader instance.
	header := core.NewPageHeader()

	// Simulate values being set.
	header.UncompressedSize = 100
	header.CompressedSize = 50

	// Allocate a statistic object.
	factory := statistic.Factory{}
	stat, err := factory.AllocStatistic(base.BOOLEAN)
	if err != nil {
		t.Fatalf("Failed to allocate Statistic: %v", err)
	}
	header.Statistic = stat

	// Call Reset and verify fields are cleared.
	header.Reset(factory)

	if header.UncompressedSize != 0 {
		t.Errorf("UncompressedSize mismatch after Reset. Expected: 0, Got: %d", header.UncompressedSize)
	}

	if header.CompressedSize != 0 {
		t.Errorf("CompressedSize mismatch after Reset. Expected: 0, Got: %d", header.CompressedSize)
	}

	if header.Statistic != nil {
		t.Errorf("Statistic mismatch after Reset. Expected: nil, Got: %v", header.Statistic)
	}
}

//////////////////
// chunk header //
//////////////////

// TestChunkHeaderDefaultConstructor verifies the default initialization of ChunkHeader.
func TestChunkHeaderDefaultConstructor(t *testing.T) {
	header := core.NewChunkHeader() // Create a new ChunkHeader instance.

	if header.MeasurementName != "" {
		t.Errorf("MeasurementName mismatch. Expected: \"\", Got: %s", header.MeasurementName)
	}

	if header.DataSize != 0 {
		t.Errorf("DataSize mismatch. Expected: 0, Got: %d", header.DataSize)
	}

	if header.DataType != base.INVALID_TS {
		t.Errorf("DataType mismatch. Expected: INVALID_TS, Got: %v", header.DataType)
	}

	if header.CompressionType != base.INVALID_C {
		t.Errorf("CompressionType mismatch. Expected: INVALID_C, Got: %v", header.CompressionType)
	}

	if header.EncodingType != base.INVALID_E {
		t.Errorf("EncodingType mismatch. Expected: INVALID_E, Got: %v", header.EncodingType)
	}

	if header.NumOfPages != 0 {
		t.Errorf("NumOfPages mismatch. Expected: 0, Got: %d", header.NumOfPages)
	}

	if header.SerializedSize != 0 {
		t.Errorf("SerializedSize mismatch. Expected: 0, Got: %d", header.SerializedSize)
	}

	if header.ChunkType != 0 {
		t.Errorf("ChunkType mismatch. Expected: 0, Got: %d", header.ChunkType)
	}
}

// TestChunkHeaderReset verifies that calling Reset on ChunkHeader clears all fields.
func TestChunkHeaderReset(t *testing.T) {
	// Create a new ChunkHeader instance.
	header := core.NewChunkHeader()

	// Simulate values being set.
	header.MeasurementName = "test"
	header.DataSize = 100
	header.DataType = base.INT32
	header.CompressionType = base.SNAPPY
	header.EncodingType = base.PLAIN
	header.NumOfPages = 5
	header.SerializedSize = 50
	header.ChunkType = 1

	// Call Reset method.
	header.Reset()

	if header.MeasurementName != "" {
		t.Errorf("MeasurementName mismatch after Reset. Expected: \"\", Got: %s", header.MeasurementName)
	}

	if header.DataSize != 0 {
		t.Errorf("DataSize mismatch after Reset. Expected: 0, Got: %d", header.DataSize)
	}

	if header.DataType != base.INVALID_TS {
		t.Errorf("DataType mismatch after Reset. Expected: INVALID_TS, Got: %v", header.DataType)
	}

	if header.CompressionType != base.INVALID_C {
		t.Errorf("CompressionType mismatch after Reset. Expected: INVALID_C, Got: %v", header.CompressionType)
	}

	if header.EncodingType != base.INVALID_E {
		t.Errorf("EncodingType mismatch after Reset. Expected: INVALID_E, Got: %v", header.EncodingType)
	}

	if header.NumOfPages != 0 {
		t.Errorf("NumOfPages mismatch after Reset. Expected: 0, Got: %d", header.NumOfPages)
	}

	if header.SerializedSize != 0 {
		t.Errorf("SerializedSize mismatch after Reset. Expected: 0, Got: %d", header.SerializedSize)
	}

	if header.ChunkType != 0 {
		t.Errorf("ChunkType mismatch after Reset. Expected: 0, Got: %d", header.ChunkType)
	}
}

////////////////
// chunk meta //
////////////////

// TestChunkMetaDefaultConstructor verifies that ChunkMeta is initialized with default values.
func TestChunkMetaDefaultConstructor(t *testing.T) {
	meta := core.ChunkMeta{} // Create a new ChunkMeta instance.

	if meta.OffsetOfHeader != 0 {
		t.Errorf("OffsetOfChunkHeader mismatch. Expected: 0, Got: %d", meta.OffsetOfHeader)
	}

	if meta.Statistic != nil {
		t.Errorf("Statistic mismatch. Expected: nil, Got: %v", meta.Statistic)
	}

	if meta.Mask != 0 {
		t.Errorf("Mask mismatch. Expected: 0, Got: %d", meta.Mask)
	}
}

// TestChunkMetaInit verifies the initialization of ChunkMeta with valid inputs.
func TestChunkMetaInit(t *testing.T) {
	meta := core.ChunkMeta{} // Create a new ChunkMeta instance.
	name := "test"
	measurementName := name // Assuming base.NewString abstracts string creation.
	factory := statistic.Factory{}
	stat, err := factory.AllocStatistic(base.INT32)
	tsID := utils.TsID{}
	mask := 1

	// Initialize the ChunkMeta instance.

	err = meta.Initialize(measurementName, base.INT32, 100, stat, tsID, byte(mask))
	if err != nil {
		return
	}

	// Validate the initialized fields.
	if meta.DataType != base.INT32 {
		t.Errorf("DataType mismatch. Expected: INT32, Got: %v", meta.DataType)
	}

	if meta.OffsetOfHeader != 100 {
		t.Errorf("OffsetOfChunkHeader mismatch. Expected: 100, Got: %d", meta.OffsetOfHeader)
	}

	if meta.Statistic != stat {
		t.Errorf("Statistic mismatch. Expected: %v, Got: %v", stat, meta.Statistic)
	}

	if meta.TsID != tsID {
		t.Errorf("TsID mismatch. Expected: %v, Got: %v", tsID, meta.TsID)
	}

	if meta.Mask != 1 {
		t.Errorf("Mask mismatch. Expected: 1, Got: %d", meta.Mask)
	}
}
