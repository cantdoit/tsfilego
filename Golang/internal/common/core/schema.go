package core

import (
	"Golang/internal/common/base"
	"errors"
)

// DeviceSchema holds the schema for a specific device
type DeviceSchema struct {
	Measurements map[string]*MeasurementSchema // Map of measurement names to measurement schemas
	IsAligned    bool                          // Optional: device alignment
}

// MeasurementSchema defines the schema for a single measurement
type MeasurementSchema struct {
	Name         string               // Measurement name (e.g., temperature)
	DataType     base.TSDataType      // Data type (e.g., INT32, BOOLEAN, TEXT)
	Encoding     base.TSEncoding      // Encoding method (e.g., RLE, PLAIN)
	Compressor   base.CompressionType // Compression algorithm (e.g., LZ4, SNAPPY)
	DefaultValue *base.Value          // Optional: Default or pre-defined value for this measurement
}

// RegisterOrUpdateSchema registers or updates a schema for a given device and measurement
func RegisterOrUpdateSchema(deviceSchemas map[string]*DeviceSchema, deviceName, measurementName string,
	dataType base.TSDataType, encoding base.TSEncoding, compressor base.CompressionType) error {

	// Validate data type
	if !base.IsValidDataType(dataType) {
		return errors.New("invalid data type: " + string(dataType))
	}

	// Validate encoding and apply default if necessary
	if encoding == "" {
		encoding = base.GetDefaultEncoding(dataType)
	} else if !base.IsValidEncoding(encoding) {
		return errors.New("invalid encoding: " + string(encoding))
	}

	// Validate compression and apply default if necessary
	if compressor == "" {
		compressor = base.GetDefaultCompression(dataType)
	} else if !base.IsValidCompression(compressor) {
		return errors.New("invalid compression: " + string(compressor))
	}

	// Fetch or create device schema
	deviceSchema, exists := deviceSchemas[deviceName]
	if !exists {
		deviceSchema = &DeviceSchema{
			Measurements: make(map[string]*MeasurementSchema),
			IsAligned:    false, // Default to non-aligned
		}
		deviceSchemas[deviceName] = deviceSchema
	}

	// Check if the measurement already exists
	if _, exists := deviceSchema.Measurements[measurementName]; exists {
		return errors.New("measurement already exists: " + measurementName)
	}

	// Register the new measurement
	deviceSchema.Measurements[measurementName] = &MeasurementSchema{
		Name:       measurementName,
		DataType:   dataType,
		Encoding:   encoding,
		Compressor: compressor,
	}

	return nil
}
