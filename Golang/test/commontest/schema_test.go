package commontest

import (
	"Golang/internal/common"
	"testing"
)

func TestDefaultMeasurementSchema(t *testing.T) {
	schema := common.MeasurementSchema{}

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
	defaultValue, _ := common.NewValue(common.INT64, int64(0))
	schema := common.MeasurementSchema{
		Name:         "temperature",
		DataType:     common.INT64,
		Encoding:     common.RLE,
		Compressor:   common.SNAPPY,
		DefaultValue: defaultValue,
	}

	if schema.Name != "temperature" {
		t.Errorf("Expected Name: 'temperature', Got: %s", schema.Name)
	}
	if schema.DataType != common.INT64 {
		t.Errorf("Expected DataType: %s, Got: %s", common.INT64, schema.DataType)
	}
	if schema.Encoding != common.RLE {
		t.Errorf("Expected Encoding: %s, Got: %s", common.RLE, schema.Encoding)
	}
	if schema.Compressor != common.SNAPPY {
		t.Errorf("Expected Compressor: %s, Got: %s", common.SNAPPY, schema.Compressor)
	}
	if schema.DefaultValue.GetValue().(int64) != int64(0) {
		t.Errorf("Expected DefaultValue: %d, Got: %d", int64(0), schema.DefaultValue.GetValue().(int64))
	}
}

func TestDefaultDeviceSchema(t *testing.T) {
	deviceSchema := common.DeviceSchema{
		Measurements: make(map[string]*common.MeasurementSchema),
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
	deviceSchemas := make(map[string]*common.DeviceSchema)

	err := common.RegisterOrUpdateSchema(deviceSchemas, "device1", "temperature", common.INT32, common.RLE, common.SNAPPY)
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

	if measurementSchema.DataType != common.INT32 {
		t.Errorf("Expected DataType: %s, Got: %s", common.INT32, measurementSchema.DataType)
	}
	if measurementSchema.Encoding != common.RLE {
		t.Errorf("Expected Encoding: %s, Got: %s", common.RLE, measurementSchema.Encoding)
	}
	if measurementSchema.Compressor != common.SNAPPY {
		t.Errorf("Expected Compressor: %s, Got: %s", common.SNAPPY, measurementSchema.Compressor)
	}
}

func TestRegisterOrUpdateSchemaDuplicate(t *testing.T) {
	deviceSchemas := make(map[string]*common.DeviceSchema)

	err := common.RegisterOrUpdateSchema(deviceSchemas, "device1", "temperature", common.INT32, common.RLE, common.SNAPPY)
	if err != nil {
		t.Fatalf("Failed to register schema: %v", err)
	}

	// Try to register the same measurement again
	err = common.RegisterOrUpdateSchema(deviceSchemas, "device1", "temperature", common.INT32, common.RLE, common.SNAPPY)
	if err == nil {
		t.Fatalf("Expected error when registering duplicate measurement, but got none")
	}
}
