package writertest

import (
	"Golang/internal/common/base"
	"Golang/internal/common/core"
	"Golang/internal/writer"
	"os"
	"testing"
)

func Test_TsFileWriter_WriteDataPointsToFile(t *testing.T) {
	// Set up the TSFileWriter
	tsFileWriter := writer.NewTsFileWriter()
	fileName := "test_tsfile_writer.tsfile"

	// Cleanup: ensure the file is deleted after the test
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(fileName)

	// Open the file
	err := tsFileWriter.Open(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// Define device and measurements
	deviceName := "device1"
	measurementSchemas := []*core.MeasurementSchema{
		{Name: "temperature", DataType: base.FLOAT, Encoding: base.PLAIN, Compressor: base.UNCOMPRESSED},
		{Name: "humidity", DataType: base.FLOAT, Encoding: base.PLAIN, Compressor: base.UNCOMPRESSED},
	}

	// Register the timeseries
	for _, schema := range measurementSchemas {
		err = tsFileWriter.RegisterTimeseries(deviceName, schema)
		if err != nil {
			t.Fatalf("Failed to register timeseries '%s': %v", schema.Name, err)
		}
	}
	t.Log(tsFileWriter)
	// Create and write records
	rows := []struct {
		Timestamp int64
		Values    map[string]interface{}
	}{
		{Timestamp: 1622505600000, Values: map[string]interface{}{"temperature": 22.5, "humidity": 60.0}},
		{Timestamp: 1622505601000, Values: map[string]interface{}{"temperature": 23.1, "humidity": 65.2}},
		{Timestamp: 1622505602000, Values: map[string]interface{}{"temperature": 24.0, "humidity": 70.0}},
	}

	for _, row := range rows {
		record := &core.TsRecord{
			Timestamp: row.Timestamp,
			DeviceID:  deviceName,
			Points:    []core.DataPoint{},
		}

		// Populate points based on the row's values
		for key, value := range row.Values {
			var pointValue *base.Value
			var err error
			switch v := value.(type) {
			case float32:
				pointValue, err = base.NewValue(base.FLOAT, v)
			case float64:
				pointValue, err = base.NewValue(base.FLOAT, float32(v))
			default:
				t.Fatalf("Unsupported value type for '%s': %T", key, v)
			}

			if err != nil {
				t.Fatalf("Failed to create value for point '%s': %v", key, err)
			}

			record.Points = append(record.Points, core.DataPoint{
				MeasurementName: key,
				DataType:        measurementSchemas[0].DataType,
				Value:           &pointValue,
			})
		}

		// Write the record
		err := tsFileWriter.WriteRecord(record)
		if err != nil {
			t.Fatalf("Failed to write record: %v", err)
		}
	}

	// Flush written data
	err = tsFileWriter.Flush()
	if err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	// Close the TSFile
	err = tsFileWriter.Close()
	if err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}

	// Confirm file creation (basic verification)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Fatalf("TSFile does not exist upon test completion: %s", fileName)
	}
}
