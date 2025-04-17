package writer

import (
	"Golang/internal/common"  // Schema-related utilities/constants
	"Golang/internal/file"    // Path to the file package
	_ "Golang/internal/utils" // For error codes or handling
	_ "errors"
	"fmt"
	"os"
)

/*
top level handing of writing tsfile with the level as such
tsfile_writer -> tsfile_io_writer -> tsfile_writer
*/

// TsFileWriter represents the high-level writer managing TSFiles
type TsFileWriter struct {
	WriteFile      *file.WriteFile                 // Instance of WriteFile used to manage underlying file
	DeviceSchemas  map[string]*common.DeviceSchema // Map of device names to their respective schemas
	WriteFileReady bool                            // Tracks if the write file is initialized and ready

}

// NewTsFileWriter creates a new instance of TsFileWriter
func NewTsFileWriter() *TsFileWriter {
	return &TsFileWriter{
		DeviceSchemas: make(map[string]*common.DeviceSchema),
	}
}

// MeasurementSchema defines the schema for a single measurement
type MeasurementSchema struct {
	Name       string // Measurement name (e.g., temperature)
	DataType   string // Data type (e.g., INT32, FLOAT)
	Encoding   string // Encoding method (e.g., RLE, TS_2DIFF)
	Compressor string // Compression algorithm (e.g., LZ4, SNAPPY)
}

// DeviceSchema holds the schema for a specific device
type DeviceSchema struct {
	Measurements map[string]*MeasurementSchema // Map of measurement names to measurement schemas
}

// Open handles opening or creating a TSFile
// - Parameters:
//   - filePath: Path to the TSFile
//   - flags: Flags for file opening (e.g., O_RDWR | O_CREATE)
//   - mode: Permissions to be assigned to the file if created
func (tf *TsFileWriter) Open(filePath string, flags int, mode os.FileMode) error {
	// Initialize a new WriteFile instance
	tf.WriteFile = &file.WriteFile{}

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists: path=%s", filePath)
	}

	// Call the Create function from WriteFile
	if err := tf.WriteFile.Create(filePath, flags, mode); err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	tf.WriteFileReady = true
	// If successful, file is ready for writing
	return nil
}

// RegisterTimeseries registers a timeseries (measurement) for a specific device
// Parameters:
// - deviceName: Name of the device (e.g., "sensor1")
// - measurementName: Name of the measurement (e.g., "temperature")
// - dataType: Type of the measurement data (e.g., INT32, FLOAT)
// - encoding: Encoding method for the measurement data
// - compressor: Compression used for the measurement data
// Returns:
// - Error if the registration fails, or nil if it succeeds
func (tf *TsFileWriter) RegisterTimeseries(deviceName, measurementName, dataType, encoding, compressor string, defaultValue interface{}) error {
	// Create a TSDataType-compatible value object
	var value *common.Value
	if defaultValue != nil {
		var err error
		value, err = common.NewValue(common.TSDataType(dataType), defaultValue)
		if err != nil {
			return err
		}
	}

	// Update the schema using the common package logic
	err := common.RegisterOrUpdateSchema(tf.DeviceSchemas, deviceName, measurementName, dataType, encoding, compressor)
	if err != nil {
		return err
	}

	// Set the default value in the schema
	tf.DeviceSchemas[deviceName].Measurements[measurementName].DefaultValue = value
	return nil
}

// Close Closes the file and finalizes the TSFile writing process
func (tf *TsFileWriter) Close() error {
	if tf.WriteFile == nil || !tf.WriteFileReady {
		return fmt.Errorf("file is not open or ready")
	}

	// Finalize writing to the TSFile and close it
	if err := tf.WriteFile.Close(); err != nil {
		return fmt.Errorf("failed to close the file: %v", err)
	}

	tf.WriteFileReady = false
	return nil
}
