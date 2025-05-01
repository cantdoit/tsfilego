package core

import (
	"Golang/internal/common/base"
	"Golang/internal/common/core"
	"testing"
)

func TestNewFieldDefaultConstructor(t *testing.T) {
	field := core.NewField(base.NULL_TYPE)

	if field.Type != base.NULL_TYPE {
		t.Errorf("Expected Type: %s, Got: %s", base.NULL_TYPE, field.Type)
	}
}

func TestNewFieldTypeConstructor(t *testing.T) {
	field := core.NewField(base.BOOLEAN)

	if field.Type != base.BOOLEAN {
		t.Errorf("Expected Type: %s, Got: %s", base.BOOLEAN, field.Type)
	}
}

func TestFieldSetValueAndGetValue(t *testing.T) {
	// Test INT32
	field := core.NewField(base.INT32)
	err := field.SetValue(int32(123))
	if err != nil {
		t.Fatalf("Failed to set field value: %v", err)
	}
	if field.Type != base.INT32 {
		t.Errorf("Expected Type: %s, Got: %s", base.INT32, field.Type)
	}
	if field.Int32Val != 123 {
		t.Errorf("Expected Value: %d, Got: %d", 123, field.Int32Val)
	}

	// Test DOUBLE
	field = core.NewField(base.DOUBLE)
	err = field.SetValue(3.14)
	if err != nil {
		t.Fatalf("Failed to set field value: %v", err)
	}
	if field.Type != base.DOUBLE {
		t.Errorf("Expected Type: %s, Got: %s", base.DOUBLE, field.Type)
	}
	if field.DoubleVal != 3.14 {
		t.Errorf("Expected Value: %f, Got: %f", 3.14, field.DoubleVal)
	}

	// Test TEXT
	field = core.NewField(base.TEXT)
	err = field.SetValue("test")
	if err != nil {
		t.Fatalf("Failed to set field value: %v", err)
	}
	if field.Type != base.TEXT {
		t.Errorf("Expected Type: %s, Got: %s", base.TEXT, field.Type)
	}
	if field.StringVal != "test" {
		t.Errorf("Expected Value: 'test', Got: '%s'", field.StringVal)
	}

	// Test invalid type
	err = field.SetValue(123) // Attempt to set integer while type is TEXT
	if err == nil {
		t.Errorf("Expected error when setting invalid type for field, got nil")
	}
}

func TestFieldFree(t *testing.T) {
	field := core.NewField(base.TEXT)
	err := field.SetValue("test")
	if err != nil {
		t.Fatalf("Failed to set field value: %v", err)
	}

	field.Free()
	if field.StringVal != "" {
		t.Errorf("Expected StringVal after Free: '', Got: '%s'", field.StringVal)
	}
}

func TestNewRowRecordWithColumnCount(t *testing.T) {
	row := core.NewRowRecord(0, 5)

	if len(row.Fields) != 5 {
		t.Fatalf("Expected 5 fields, Got: %d", len(row.Fields))
	}

	for _, field := range row.Fields {
		if field.Type != base.NULL_TYPE {
			t.Errorf("Expected Field Type: %s, Got: %s", base.NULL_TYPE, field.Type)
		}
	}
}

func TestNewRowRecordWithTimestamp(t *testing.T) {
	row := core.NewRowRecord(1625140800, 3)

	if row.Timestamp != 1625140800 {
		t.Errorf("Expected Timestamp: %d, Got: %d", 1625140800, row.Timestamp)
	}

	if len(row.Fields) != 3 {
		t.Fatalf("Expected 3 fields, Got: %d", len(row.Fields))
	}

	for _, field := range row.Fields {
		if field.Type != base.NULL_TYPE {
			t.Errorf("Expected Field Type: %s, Got: %s", base.NULL_TYPE, field.Type)
		}
	}
}

func TestAddFieldToRowRecord(t *testing.T) {
	row := core.NewRowRecord(0, 2)
	field := core.NewField(base.INT64)
	err := field.SetValue(int64(12345))
	if err != nil {
		return
	}

	row.AddField(field)
	if len(row.Fields) != 3 {
		t.Fatalf("Expected 3 fields after addition, Got: %d", len(row.Fields))
	}

	lastField := row.Fields[2]
	if lastField.Type != base.INT64 {
		t.Errorf("Expected Type: %s, Got: %s", base.INT64, lastField.Type)
	}
	if lastField.Int64Val != 12345 {
		t.Errorf("Expected Value: %d, Got: %d", 12345, lastField.Int64Val)
	}
}

func TestSetFieldValue(t *testing.T) {
	row := core.NewRowRecord(0, 2)

	// Set first field to INT64
	err := row.SetFieldValue(0, int64(12345), base.INT64)
	if err != nil {
		t.Fatalf("Failed to set field value: %v", err)
	}

	field := row.Fields[0]
	if field.Type != base.INT64 {
		t.Errorf("Expected Type: %s, Got: %s", base.INT64, field.Type)
	}
	if field.Int64Val != 12345 {
		t.Errorf("Expected Value: %d, Got: %d", 12345, field.Int64Val)
	}

	// Test out-of-bound index
	err = row.SetFieldValue(10, int64(67890), base.INT64)
	if err == nil {
		t.Errorf("Expected error for out-of-bound index, got nil")
	}
}

func TestGetField(t *testing.T) {
	row := core.NewRowRecord(0, 2)

	// Modify the field at index 1
	row.Fields[1].Type = base.TEXT
	err := row.Fields[1].SetValue("test_value")
	if err != nil {
		return
	}

	retrievedField, err := row.GetField(1)
	if err != nil {
		t.Fatalf("Failed to get field: %v", err)
	}

	if retrievedField.Type != base.TEXT {
		t.Errorf("Expected Type: %s, Got: %s", base.TEXT, retrievedField.Type)
	}
	if retrievedField.StringVal != "test_value" {
		t.Errorf("Expected Value: 'test_value', Got: '%s'", retrievedField.StringVal)
	}

	// Test out-of-bound index
	_, err = row.GetField(5)
	if err == nil {
		t.Errorf("Expected error for out-of-bound index, got nil")
	}
}
