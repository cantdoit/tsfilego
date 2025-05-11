package writer

import (
	_ "Golang/internal/common/base"
	// "Golang/internal/common/bitmap"
	"Golang/internal/common/core"
	"Golang/internal/fileio"
	"fmt"
	"os"
)

// TsFileWriter represents the high-level writer managing TSFiles
type TsFileWriter struct {
	WriteFile                 *fileio.WriteFile                  // Managing the underlying TSFile
	DeviceSchemas             map[string]*core.DeviceSchema      // Device to schema map
	ChunkWriters              map[string]*ChunkWriter            // Device to ChunkWriter map
	Schemas                   map[string]*MeasurementSchemaGroup // Device schemas group for timeseries
	StartFileDone             bool                               // Indication if file writing has been initialized
	RecordCountSinceLastFlush int64                              // Count of records since the last flush
	RecordCountForNextCheck   int64                              // Count for the next memory boundary check
	WriteFileCreated          bool                               // Indicates if the WriteFile has been created
}

// MeasurementSchemaGroup represents a group of measurements for a device
type MeasurementSchemaGroup struct {
	MeasurementSchemas map[string]*core.MeasurementSchema
	IsAligned          bool
}

// Function to create a new instance of TsFileWriter
func NewTsFileWriter() *TsFileWriter {
	return &TsFileWriter{
		DeviceSchemas: make(map[string]*core.DeviceSchema),
		ChunkWriters:  make(map[string]*ChunkWriter),
		Schemas:       make(map[string]*MeasurementSchemaGroup),
	}
}

// Open handles opening or creating a TSFile
func (tf *TsFileWriter) Open(filePath string, flags int, mode os.FileMode) error {
	// Initialize a new WriteFile instance
	tf.WriteFile = &fileio.WriteFile{}

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists: path=%s", filePath)
	}

	// Create the TSFile
	if err := tf.WriteFile.Create(filePath, flags, mode); err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	tf.WriteFileCreated = true
	// Indicate the file is prepared for writing
	return nil
}

// Flush ensures any buffered data is written to the TSFile
func (tf *TsFileWriter) Flush() error {
	if tf.WriteFile == nil || !tf.WriteFileCreated {
		return fmt.Errorf("no file to write to")
	}

	// Add logic to flush buffers and ensure all writes are saved
	return nil
}

// Close finalizes the TSFile and writes index/data footer
func (tf *TsFileWriter) Close() error {
	if tf.WriteFile == nil || !tf.WriteFileCreated {
		return fmt.Errorf("file is not open or ready")
	}

	// Add logic for cleaning up, writing indexes, etc.
	if err := tf.WriteFile.CloseFile(); err != nil {
		return fmt.Errorf("failed to close the file: %v", err)
	}

	tf.WriteFileCreated = false
	return nil
}

// RegisterTimeseries registers a timeseries (a measurement) for a specific device
func (tf *TsFileWriter) RegisterTimeseries(deviceID string, measurementSchema *core.MeasurementSchema) error {
	group, exists := tf.Schemas[deviceID]
	if !exists {
		group = &MeasurementSchemaGroup{
			MeasurementSchemas: make(map[string]*core.MeasurementSchema),
			IsAligned:          false,
		}
		tf.Schemas[deviceID] = group
	}

	// Add or update the measurement schema in the group
	group.MeasurementSchemas[measurementSchema.Name] = measurementSchema
	return nil
}

// RegisterAlignedTimeseries registers aligned timeseries for a particular device
func (tf *TsFileWriter) RegisterAlignedTimeseries(deviceID string, measurementSchemas []*core.MeasurementSchema) error {
	group, exists := tf.Schemas[deviceID]
	if !exists {
		group = &MeasurementSchemaGroup{
			MeasurementSchemas: make(map[string]*core.MeasurementSchema),
			IsAligned:          true,
		}
		tf.Schemas[deviceID] = group
	}

	// Add or update aligned schemas
	for _, schema := range measurementSchemas {
		group.MeasurementSchemas[schema.Name] = schema
	}
	return nil
}

// WriteRecord writes a single record to the TSFile
func (tf *TsFileWriter) WriteRecord(record *core.TsRecord) error {
	// Ensure the file is initialized and ready
	if tf.WriteFile == nil || !tf.WriteFileCreated {
		return fmt.Errorf("TSFile is not ready for writing")
	}

	deviceName := record.DeviceID
	timestamp := record.Timestamp

	// Check and prepare schema and chunk writers
	chunkWriters, err := tf.checkAndPrepareSchema(deviceName, record)
	if err != nil {
		return err
	}

	// Ensure the chunk writers match the size of the record's points
	if len(chunkWriters) != len(record.Points) {
		return fmt.Errorf("mismatch between chunk writers and points in the record")
	}

	// Iterate through the points and write to the respective ChunkWriter
	for index, point := range record.Points {
		chunkWriter := chunkWriters[index]
		if chunkWriter == nil {
			// If no valid chunk writer exists for the measurement, skip
			continue
		}
		if err := chunkWriter.Write(timestamp, point); err != nil {
			// Log or handle individual failure, but allow other points to continue
			return fmt.Errorf("failed to write point at index %d: %v", index, err)
		}
	}

	// Increment the record count and check memory thresholds for flushing
	tf.RecordCountSinceLastFlush++
	if err := tf.checkMemoryAndFlushChunks(); err != nil {
		return fmt.Errorf("memory check or flush failed: %v", err)
	}
	return nil
}

// checkAndPrepareSchema validates the schema for the given device and ensures chunk writers are ready.
func (tf *TsFileWriter) checkAndPrepareSchema(deviceName string, record *core.TsRecord) ([]*ChunkWriter, error) {
	// Locate the device's schema
	deviceSchemaGroup, exists := tf.Schemas[deviceName]
	if !exists || deviceSchemaGroup == nil {
		return nil, fmt.Errorf("device '%s' does not exist or has no schema", deviceName)
	}

	measurementSchemaMap := deviceSchemaGroup.MeasurementSchemas
	measurementCount := len(record.Points)
	chunkWriters := make([]*ChunkWriter, measurementCount)

	// Iterate through measurement names in the record
	for idx, point := range record.Points {
		measurementName := point.MeasurementName
		schema, exists := measurementSchemaMap[measurementName]

		if !exists {
			// Measurement does not exist in the schema, mark as nil
			chunkWriters[idx] = nil
			continue
		}
		chunkWriter := ChunkWriter{}
		err := chunkWriter.Initialize(schema.Name, schema.DataType, schema.Encoding, schema.Compressor)
		// If the chunk writer does not exist, initialize it
		if tf.ChunkWriters == nil {

			if err != nil {
				// Cleanup in case of initialization failure
				for _, writer := range chunkWriters {
					if writer != nil {
						writer.Destroy() // Assume Close safely cleans up resources
					}
				}
				return nil, fmt.Errorf("failed to initialize chunk writer for measurement '%s': %v", measurementName, err)
			}
		}

		// Add the chunk writer to the list
		chunkWriters[idx] = &chunkWriter
	}

	return chunkWriters, nil
}

// checkMemoryAndFlushChunks checks if memory size exceeds threshold and flushes chunks if necessary.
func (tf *TsFileWriter) checkMemoryAndFlushChunks() error {
	// Logic to determine if memory usage requires flushing
	// This is just a placeholder for more advanced memory management
	if tf.RecordCountSinceLastFlush >= tf.RecordCountForNextCheck {
		// Perform flushing operation (implementation depends on other layers)
		for _, schemaGroup := range tf.Schemas {
			err := tf.flushChunkGroup(schemaGroup, false)
			if err != nil {
				return fmt.Errorf("failed to flush chunks for device schema: %v", err)
			}
		}
		// Reset the record count after flushing
		tf.RecordCountSinceLastFlush = 0
	}
	return nil
}

// WriteTablet writes a tablet of records (batch) to the TSFile
func (tf *TsFileWriter) WriteTablet(tablet *core.Tablet) error {
	// Validate the tablet schema with registered schemas
	// Write data in batch using an underlying chunk writer for each column
	return nil
}

// CheckSchema validates schema compatibility and initializes chunk writers if necessary
func (tf *TsFileWriter) checkSchema(deviceID string, measurementNames []string) error {
	if tf.WriteFile == nil || !tf.WriteFileCreated {
		return fmt.Errorf("file is not ready")
	}

	deviceSchema, exists := tf.DeviceSchemas[deviceID]
	if !exists {
		return fmt.Errorf("device '%s' schema not registered", deviceID)
	}

	// Check if all measurement names exist in the schema
	for _, measurement := range measurementNames {
		if _, ok := deviceSchema.Measurements[measurement]; !ok {
			return fmt.Errorf("measurement '%s' not found for device '%s'", measurement, deviceID)
		}
	}
	return nil
}

// Private function to flush a specific chunk group (device)
func (tf *TsFileWriter) flushChunkGroup(group *MeasurementSchemaGroup, isAligned bool) error {
	// Logic to flush a chunk group
	return nil
}
