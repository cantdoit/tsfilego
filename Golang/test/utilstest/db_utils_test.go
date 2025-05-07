package utils_test

import (
	"Golang/internal/utils"
	"testing"
)

func TestTsIDConstructor(t *testing.T) {
	tsID := utils.TsID{}
	if tsID.IsValid() {
		t.Errorf("Expected TsID to be initialized to zero values, got %+v", tsID)
	}
}

func TestTsIDParameterizedConstructor(t *testing.T) {
	tsID := utils.TsID{DbNID: 1, DeviceNID: 2, MeasurementNID: 3}
	if tsID.DbNID != 1 || tsID.DeviceNID != 2 || tsID.MeasurementNID != 3 {
		t.Errorf("Expected TsID(1, 2, 3), got %+v", tsID)
	}
}

func TestTsIDReset(t *testing.T) {
	tsID := utils.TsID{DbNID: 1, DeviceNID: 2, MeasurementNID: 3}
	tsID.Reset()
	if tsID.DbNID != 0 || tsID.DeviceNID != 0 || tsID.MeasurementNID != 0 {
		t.Errorf("Expected TsID to be reset to zero values, got %+v", tsID)
	}
}

func TestTsIDEquals(t *testing.T) {
	tsID1 := utils.TsID{DbNID: 1, DeviceNID: 2, MeasurementNID: 3}
	tsID2 := utils.TsID{DbNID: 1, DeviceNID: 2, MeasurementNID: 3}
	if !tsID1.Equals(tsID2) {
		t.Errorf("Expected TsID1 to equal TsID2, got TsID1 %+v and TsID2 %+v", tsID1, tsID2)
	}

	tsID2.DbNID = 4
	if tsID1.Equals(tsID2) {
		t.Errorf("Expected TsID1 to not equal TsID2, got TsID1 %+v and TsID2 %+v", tsID1, tsID2)
	}
}

func TestTsIDNotEquals(t *testing.T) {
	tsID1 := utils.TsID{DbNID: 1, DeviceNID: 2, MeasurementNID: 3}
	tsID2 := utils.TsID{DbNID: 1, DeviceNID: 2, MeasurementNID: 3}
	if !tsID1.Equals(tsID2) {
		t.Errorf("Expected TsID1 to equal TsID2, got TsID1 %+v and TsID2 %+v", tsID1, tsID2)
	}

	tsID2.DbNID = 4
	if tsID1.Equals(tsID2) {
		t.Errorf("Expected TsID1 to not equal TsID2, got TsID1 %+v and TsID2 %+v", tsID1, tsID2)
	}
}

func TestTsIDToInt64(t *testing.T) {
	tsID := utils.TsID{DbNID: 1, DeviceNID: 2, MeasurementNID: 3}
	expected := (1 << 32) | (2 << 16) | 3
	if tsID.ToInt64() != int64(expected) {
		t.Errorf("Expected ToInt64() to return %d, got %d", expected, tsID.ToInt64())
	}
}

func TestTsIDLessOperator(t *testing.T) {
	tsID1 := utils.TsID{DbNID: 1, DeviceNID: 2, MeasurementNID: 3}
	tsID2 := utils.TsID{DbNID: 1, DeviceNID: 2, MeasurementNID: 4}

	// Simulating < operator with custom comparison
	isLess := func(a, b utils.TsID) bool {
		if a.DbNID != b.DbNID {
			return a.DbNID < b.DbNID
		}
		if a.DeviceNID != b.DeviceNID {
			return a.DeviceNID < b.DeviceNID
		}
		return a.MeasurementNID < b.MeasurementNID
	}

	if !isLess(tsID1, tsID2) {
		t.Errorf("Expected tsID1 to be less than tsID2, got %+v and %+v", tsID1, tsID2)
	}

	if isLess(tsID2, tsID1) {
		t.Errorf("Expected tsID2 to not be less than tsID1, got %+v and %+v", tsID2, tsID1)
	}
}
