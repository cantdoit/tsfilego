package common

import "errors"

// DeviceSchema holds the schema for a specific device
type DeviceSchema struct {
	Measurements map[string]*MeasurementSchema // Map of measurement names to measurement schemas
	IsAligned    bool                          // Optional: device alignment
}

// MeasurementSchema defines the schema for a single measurement
type MeasurementSchema struct {
	Name         string          // Measurement name (e.g., temperature)
	DataType     TSDataType      // Data type (e.g., INT32, BOOLEAN, TEXT)
	Encoding     TSEncoding      // Encoding method (e.g., RLE, PLAIN)
	Compressor   CompressionType // Compression algorithm (e.g., LZ4, SNAPPY)
	DefaultValue *Value          // Optional: Default or pre-defined value for this measurement
}

// RegisterOrUpdateSchema registers or updates a schema for a given device and measurement
func RegisterOrUpdateSchema(deviceSchemas map[string]*DeviceSchema, deviceName, measurementName string,
	dataType TSDataType, encoding TSEncoding, compressor CompressionType) error {

	// Validate data type
	if !IsValidDataType(dataType) {
		return errors.New("invalid data type: " + string(dataType))
	}

	// Validate encoding and apply default if necessary
	if encoding == "" {
		encoding = GetDefaultEncoding(dataType)
	} else if !IsValidEncoding(encoding) {
		return errors.New("invalid encoding: " + string(encoding))
	}

	// Validate compression and apply default if necessary
	if compressor == "" {
		compressor = GetDefaultCompression(dataType)
	} else if !IsValidCompression(compressor) {
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
