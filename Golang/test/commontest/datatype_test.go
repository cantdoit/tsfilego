package commontest

import (
	"Golang/internal/common"
	"testing"
)

func TestValueConstructor(t *testing.T) {
	v, err := common.NewValue(common.INT64, int64(123456789))
	if err != nil {
		t.Fatalf("Failed to create Value instance: %v", err)
	}

	if v.Type != common.INT64 {
		t.Errorf("Expected Type: %s, Got: %s", common.INT64, v.Type)
	}

	if v.GetValue().(int64) != 123456789 {
		t.Errorf("Expected Value: %d, Got: %d", 123456789, v.GetValue().(int64))
	}
}

func TestValueSetAndGetValue(t *testing.T) {
	// Test INT64
	v, err := common.NewValue(common.INT64, int64(123456789))
	if err != nil {
		t.Fatalf("Failed to create Value instance: %v", err)
	}
	if v.Type != common.INT64 {
		t.Errorf("Expected Type: %s, Got: %s", common.INT64, v.Type)
	}
	if v.GetValue().(int64) != 123456789 {
		t.Errorf("Expected Value: %d, Got: %d", 123456789, v.GetValue().(int64))
	}

	// Test DOUBLE
	v, err = common.NewValue(common.DOUBLE, float64(123.456))
	if err != nil {
		t.Fatalf("Failed to create Value instance: %v", err)
	}
	if v.Type != common.DOUBLE {
		t.Errorf("Expected Type: %s, Got: %s", common.DOUBLE, v.Type)
	}
	if v.GetValue().(float64) != 123.456 {
		t.Errorf("Expected Value: %f, Got: %f", 123.456, v.GetValue().(float64))
	}

	// Test TEXT
	v, err = common.NewValue(common.TEXT, "hello")
	if err != nil {
		t.Fatalf("Failed to create Value instance: %v", err)
	}
	if v.Type != common.TEXT {
		t.Errorf("Expected Type: %s, Got: %s", common.TEXT, v.Type)
	}
	if v.GetValue().(string) != "hello" {
		t.Errorf("Expected Value: %s, Got: %s", "hello", v.GetValue().(string))
	}
}

func TestValueFree(t *testing.T) {
	// Test Free method for TEXT data type
	v, err := common.NewValue(common.TEXT, "hello")
	if err != nil {
		t.Fatalf("Failed to create Value instance: %v", err)
	}
	v.Free()
	if v.GetValue().(string) != "" {
		t.Errorf("Expected Value after Free: empty string, Got: %s", v.GetValue().(string))
	}
}

func TestInvalidValueCreation(t *testing.T) {
	_, err := common.NewValue(common.INT64, "invalid_type")
	if err == nil {
		t.Errorf("Expected error when creating Value with invalid type, but got none")
	}
}
