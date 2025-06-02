package core

import (
	"Golang/internal/common/base"
	"Golang/internal/writer"
	"testing"
)

func TestDefaultMeasurementSchema(t *testing.T) {
	schema := writer.MeasurementSchema{}

	if schema.Name != "" {
		t.Errorf("Expected Name: '', Got: %s", schema.Name)
	}
	if schema.DataType != "" {
		t.Errorf("Expected DataType: '', Got: %s", schema.DataType)
	}
	if schema.Encoding != "" {
		t.Errorf("Expected Encoding: '', Got: %s", schema.Encoding)
	}
	if schema.Compressor != "" {
		t.Errorf("Expected Compressor: '', Got: %s", schema.Compressor)
	}
	if schema.DefaultValue != nil {
		t.Errorf("Expected DefaultValue: nil, Got: %v", schema.DefaultValue)
	}
}

func TestParameterizedMeasurementSchema(t *testing.T) {
	defaultValue, _ := base.NewValue(base.INT64, int64(0))
	schema := writer.MeasurementSchema{
		Name:         "temperature",
		DataType:     base.INT64,
		Encoding:     base.RLE,
		Compressor:   base.SNAPPY,
		DefaultValue: defaultValue,
	}

	if schema.Name != "temperature" {
		t.Errorf("Expected Name: 'temperature', Got: %s", schema.Name)
	}
	if schema.DataType != base.INT64 {
		t.Errorf("Expected DataType: %s, Got: %s", base.INT64, schema.DataType)
	}
	if schema.Encoding != base.RLE {
		t.Errorf("Expected Encoding: %s, Got: %s", base.RLE, schema.Encoding)
	}
	if schema.Compressor != base.SNAPPY {
		t.Errorf("Expected Compressor: %s, Got: %s", base.SNAPPY, schema.Compressor)
	}
	if schema.DefaultValue.GetValue().(int64) != int64(0) {
		t.Errorf("Expected DefaultValue: %d, Got: %d", int64(0), schema.DefaultValue.GetValue().(int64))
	}
}

func TestDefaultDeviceSchema(t *testing.T) {
	deviceSchema := writer.DeviceSchema{
		Measurements: make(map[string]*writer.MeasurementSchema),
		IsAligned:    false,
	}

	if len(deviceSchema.Measurements) != 0 {
		t.Errorf("Expected empty Measurements map, Got: %d", len(deviceSchema.Measurements))
	}

	if deviceSchema.IsAligned != false {
		t.Errorf("Expected IsAligned: false, Got: %v", deviceSchema.IsAligned)
	}
}

func TestRegisterOrUpdateSchema(t *testing.T) {
	deviceSchemas := make(map[string]*writer.DeviceSchema)

	err := writer.RegisterOrUpdateSchema(deviceSchemas, "device1", "temperature", base.INT32, base.RLE, base.SNAPPY)
	if err != nil {
		t.Fatalf("Failed to register schema: %v", err)
	}

	if len(deviceSchemas) != 1 {
		t.Fatalf("Expected 1 device schema, Got: %d", len(deviceSchemas))
	}

	deviceSchema := deviceSchemas["device1"]
	measurementSchema, exists := deviceSchema.Measurements["temperature"]
	if !exists {
		t.Fatalf("Measurement 'temperature' not found in device schema")
	}

	if measurementSchema.DataType != base.INT32 {
		t.Errorf("Expected DataType: %s, Got: %s", base.INT32, measurementSchema.DataType)
	}
	if measurementSchema.Encoding != base.RLE {
		t.Errorf("Expected Encoding: %s, Got: %s", base.RLE, measurementSchema.Encoding)
	}
	if measurementSchema.Compressor != base.SNAPPY {
		t.Errorf("Expected Compressor: %s, Got: %s", base.SNAPPY, measurementSchema.Compressor)
	}
}

func TestRegisterOrUpdateSchemaDuplicate(t *testing.T) {
	deviceSchemas := make(map[string]*writer.DeviceSchema)

	err := writer.RegisterOrUpdateSchema(deviceSchemas, "device1", "temperature", base.INT32, base.RLE, base.SNAPPY)
	if err != nil {
		t.Fatalf("Failed to register schema: %v", err)
	}

	// Try to register the same measurement again
	err = writer.RegisterOrUpdateSchema(deviceSchemas, "device1", "temperature", base.INT32, base.RLE, base.SNAPPY)
	if err == nil {
		t.Fatalf("Expected error when registering duplicate measurement, but got none")
	}
}
