package commontest

import (
	"Golang/internal/common"
	"testing"
)

// RecordTestSuite manages setup and teardown for tests related to TsRecord.
type RecordTestSuite struct {
	tsRecord *common.TsRecord
}

// Setup initializes a new TsRecord instance.
func (suite *RecordTestSuite) Setup(deviceID string, timestamp int64, capacity int) {
	suite.tsRecord = &common.TsRecord{
		Timestamp: timestamp,
		DeviceID:  deviceID,
		Points:    make([]common.DataPoint, 0, capacity),
	}
}

// TearDown clears the TsRecord instance.
func (suite *RecordTestSuite) TearDown() {
	suite.tsRecord = nil
}

// Test for boolean DataPoint constructor
func TestDataPoint_BoolConstructor(t *testing.T) {
	dp := common.NewDataPointBool("touch_sensor", true)
	if dp.MeasurementName != "touch_sensor" {
		t.Errorf("Expected MeasurementName to be 'touch_sensor', got '%s'", dp.MeasurementName)
	}
	if dp.DataType != common.BOOLEAN {
		t.Errorf("Expected DataType to be BOOLEAN, got '%v'", dp.DataType)
	}
	if dp.BoolVal == nil || *dp.BoolVal != true {
		t.Errorf("Expected BoolVal to be true, got '%v'", dp.BoolVal)
	}
}

// Test for int32 DataPoint constructor
func TestDataPoint_Int32Constructor(t *testing.T) {
	dp := common.NewDataPointInt32("temperature", 100)
	if dp.MeasurementName != "temperature" {
		t.Errorf("Expected MeasurementName to be 'temperature', got '%s'", dp.MeasurementName)
	}
	if dp.DataType != common.INT32 {
		t.Errorf("Expected DataType to be INT32, got '%v'", dp.DataType)
	}
	if dp.Int32Val == nil || *dp.Int32Val != 100 {
		t.Errorf("Expected Int32Val to be 100, got '%v'", dp.Int32Val)
	}
}

// Test for setting int32 value in DataPoint
func TestDataPoint_SetInt32(t *testing.T) {
	dp := common.DataPoint{MeasurementName: "temperature"}
	dp.SetInt32(42)

	if dp.DataType != common.INT32 {
		t.Errorf("Expected DataType to be INT32, got '%v'", dp.DataType)
	}
	if dp.Int32Val == nil || *dp.Int32Val != 42 {
		t.Errorf("Expected Int32Val to be 42, got '%v'", dp.Int32Val)
	}
}

// Test for a TsRecord constructed with a device name
func TestTsRecord_ConstructorWithDeviceName(t *testing.T) {
	tsr := common.NewTsRecord("device1")

	if tsr.DeviceID != "device1" {
		t.Errorf("Expected DeviceID to be 'device1', got '%s'", tsr.DeviceID)
	}
	if len(tsr.Points) != 0 {
		t.Errorf("Expected Points to be empty, got %d points", len(tsr.Points))
	}
}

// Test for adding a DataPoint to TsRecord
func TestTsRecord_AddDataPoint(t *testing.T) {
	suite := &RecordTestSuite{}
	suite.Setup("device1", 0, 1)
	defer suite.TearDown()

	dp := common.NewDataPointDouble("temperature", 36.6)
	if err := suite.tsRecord.AddDataPoint(dp); err != nil {
		t.Errorf("Failed to add DataPoint: %v", err)
	}

	if len(suite.tsRecord.Points) != 1 {
		t.Errorf("Expected 1 DataPoint, got %d", len(suite.tsRecord.Points))
	}
	if suite.tsRecord.Points[0].MeasurementName != "temperature" {
		t.Errorf("Expected first DataPoint's MeasurementName to be 'temperature', got '%s'", suite.tsRecord.Points[0].MeasurementName)
	}
	if suite.tsRecord.Points[0].DataType != common.DOUBLE {
		t.Errorf("Expected first DataPoint's DataType to be DOUBLE, got '%v'", suite.tsRecord.Points[0].DataType)
	}
	if suite.tsRecord.Points[0].DoubleVal == nil || *suite.tsRecord.Points[0].DoubleVal != 36.6 {
		t.Errorf("Expected first DataPoint's DoubleVal to be 36.6, got '%v'", suite.tsRecord.Points[0].DoubleVal)
	}
}

// Test for adding a large number of DataPoints to TsRecord
func TestTsRecord_LargeQuantities(t *testing.T) {
	suite := &RecordTestSuite{}
	suite.Setup("device1", 0, 10000)
	defer suite.TearDown()

	for i := 0; i < 10000; i++ {
		dp := common.NewDataPointInt64("measurement_"+string(rune(i)), int64(i))
		if err := suite.tsRecord.AddDataPoint(dp); err != nil {
			t.Fatalf("Failed to add DataPoint %d: %v", i, err)
		}
	}

	if len(suite.tsRecord.Points) != 10000 {
		t.Errorf("Expected 10000 DataPoints, got %d", len(suite.tsRecord.Points))
	}
	for i := 0; i < 10000; i++ {
		if suite.tsRecord.Points[i].MeasurementName != "measurement_"+string(rune(i)) {
			t.Errorf("Expected MeasurementName to be 'measurement_%d', got '%s'", i, suite.tsRecord.Points[i].MeasurementName)
		}
		if suite.tsRecord.Points[i].Int64Val == nil || *suite.tsRecord.Points[i].Int64Val != int64(i) {
			t.Errorf("Expected Int64Val to be %d, got '%v'", i, suite.tsRecord.Points[i].Int64Val)
		}
	}
}
