package commontest

import (
	"Golang/internal/common"
	"testing"
)

func TestNewFieldDefaultConstructor(t *testing.T) {
	field := common.NewField(common.NULL_TYPE)

	if field.Type != common.NULL_TYPE {
		t.Errorf("Expected Type: %s, Got: %s", common.NULL_TYPE, field.Type)
	}
}

func TestNewFieldTypeConstructor(t *testing.T) {
	field := common.NewField(common.BOOLEAN)

	if field.Type != common.BOOLEAN {
		t.Errorf("Expected Type: %s, Got: %s", common.BOOLEAN, field.Type)
	}
}

func TestFieldSetValueAndGetValue(t *testing.T) {
	// Test INT32
	field := common.NewField(common.INT32)
	err := field.SetValue(int32(123))
	if err != nil {
		t.Fatalf("Failed to set field value: %v", err)
	}
	if field.Type != common.INT32 {
		t.Errorf("Expected Type: %s, Got: %s", common.INT32, field.Type)
	}
	if field.Int32Val != 123 {
		t.Errorf("Expected Value: %d, Got: %d", 123, field.Int32Val)
	}

	// Test DOUBLE
	field = common.NewField(common.DOUBLE)
	err = field.SetValue(3.14)
	if err != nil {
		t.Fatalf("Failed to set field value: %v", err)
	}
	if field.Type != common.DOUBLE {
		t.Errorf("Expected Type: %s, Got: %s", common.DOUBLE, field.Type)
	}
	if field.DoubleVal != 3.14 {
		t.Errorf("Expected Value: %f, Got: %f", 3.14, field.DoubleVal)
	}

	// Test TEXT
	field = common.NewField(common.TEXT)
	err = field.SetValue("test")
	if err != nil {
		t.Fatalf("Failed to set field value: %v", err)
	}
	if field.Type != common.TEXT {
		t.Errorf("Expected Type: %s, Got: %s", common.TEXT, field.Type)
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
	field := common.NewField(common.TEXT)
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
	row := common.NewRowRecord(0, 5)

	if len(row.Fields) != 5 {
		t.Fatalf("Expected 5 fields, Got: %d", len(row.Fields))
	}

	for _, field := range row.Fields {
		if field.Type != common.NULL_TYPE {
			t.Errorf("Expected Field Type: %s, Got: %s", common.NULL_TYPE, field.Type)
		}
	}
}

func TestNewRowRecordWithTimestamp(t *testing.T) {
	row := common.NewRowRecord(1625140800, 3)

	if row.Timestamp != 1625140800 {
		t.Errorf("Expected Timestamp: %d, Got: %d", 1625140800, row.Timestamp)
	}

	if len(row.Fields) != 3 {
		t.Fatalf("Expected 3 fields, Got: %d", len(row.Fields))
	}

	for _, field := range row.Fields {
		if field.Type != common.NULL_TYPE {
			t.Errorf("Expected Field Type: %s, Got: %s", common.NULL_TYPE, field.Type)
		}
	}
}

func TestAddFieldToRowRecord(t *testing.T) {
	row := common.NewRowRecord(0, 2)
	field := common.NewField(common.INT64)
	err := field.SetValue(int64(12345))
	if err != nil {
		return
	}

	row.AddField(field)
	if len(row.Fields) != 3 {
		t.Fatalf("Expected 3 fields after addition, Got: %d", len(row.Fields))
	}

	lastField := row.Fields[2]
	if lastField.Type != common.INT64 {
		t.Errorf("Expected Type: %s, Got: %s", common.INT64, lastField.Type)
	}
	if lastField.Int64Val != 12345 {
		t.Errorf("Expected Value: %d, Got: %d", 12345, lastField.Int64Val)
	}
}

func TestSetFieldValue(t *testing.T) {
	row := common.NewRowRecord(0, 2)

	// Set first field to INT64
	err := row.SetFieldValue(0, int64(12345), common.INT64)
	if err != nil {
		t.Fatalf("Failed to set field value: %v", err)
	}

	field := row.Fields[0]
	if field.Type != common.INT64 {
		t.Errorf("Expected Type: %s, Got: %s", common.INT64, field.Type)
	}
	if field.Int64Val != 12345 {
		t.Errorf("Expected Value: %d, Got: %d", 12345, field.Int64Val)
	}

	// Test out-of-bound index
	err = row.SetFieldValue(10, int64(67890), common.INT64)
	if err == nil {
		t.Errorf("Expected error for out-of-bound index, got nil")
	}
}

func TestGetField(t *testing.T) {
	row := common.NewRowRecord(0, 2)

	// Modify the field at index 1
	row.Fields[1].Type = common.TEXT
	err := row.Fields[1].SetValue("test_value")
	if err != nil {
		return
	}

	retrievedField, err := row.GetField(1)
	if err != nil {
		t.Fatalf("Failed to get field: %v", err)
	}

	if retrievedField.Type != common.TEXT {
		t.Errorf("Expected Type: %s, Got: %s", common.TEXT, retrievedField.Type)
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
