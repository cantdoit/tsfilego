package common

import "errors"

// DeviceSchema holds the schema for a specific device
type DeviceSchema struct {
	Measurements map[string]*MeasurementSchema // Map of measurement names to measurement schemas
}

// MeasurementSchema defines the schema for a single measurement
type MeasurementSchema struct {
	Name         string     // Measurement name (e.g., temperature)
	DataType     TSDataType // Data type (e.g., INT32, BOOLEAN, TEXT)
	Encoding     string     // Encoding method (e.g., RLE)
	Compressor   string     // Compression algorithm (e.g., LZ4, SNAPPY)
	DefaultValue *Value     // Optional: Default or pre-defined value for this measurement
}

// RegisterOrUpdateSchema registers or updates a schema for a given device and measurement
func RegisterOrUpdateSchema(deviceSchemas map[string]*DeviceSchema, deviceName, measurementName, dataType string, encoding, compressor string) error {
	// Validate data type
	validType := false
	for _, t := range []TSDataType{BOOLEAN, INT32, INT64, FLOAT, DOUBLE, TEXT} {
		if TSDataType(dataType) == t {
			validType = true
			break
		}
	}
	if !validType {
		return errors.New("invalid data type: " + dataType)
	}

	// Check the encoding (example validation can be expanded)
	validEncodings := []string{"RLE", "PLAIN", "TS_2DIFF"}
	isValidEncoding := false
	for _, e := range validEncodings {
		if encoding == e {
			isValidEncoding = true
			break
		}
	}
	if !isValidEncoding {
		return errors.New("invalid encoding: " + encoding)
	}

	deviceSchema, exists := deviceSchemas[deviceName]
	if !exists {
		// Create a new device schema if it doesn't exist
		deviceSchema = &DeviceSchema{
			Measurements: make(map[string]*MeasurementSchema),
		}
		deviceSchemas[deviceName] = deviceSchema
	}

	// Check if the measurement already exists for this device
	if _, exists := deviceSchema.Measurements[measurementName]; exists {
		return errors.New("measurement already exists: " + measurementName)
	}

	// Register a new measurement
	deviceSchema.Measurements[measurementName] = &MeasurementSchema{
		Name:       measurementName,
		DataType:   TSDataType(dataType),
		Encoding:   encoding,
		Compressor: compressor,
	}

	return nil
}
