package common

// DataPoint represents a data value of one measurement of some device.
type DataPoint struct {
	IsNull          bool       // Indicates if the value is null
	MeasurementName string     // Measurement name
	DataType        TSDataType // Data type of the value
	BoolVal         *bool      // Boolean value (optional)
	Int32Val        *int32     // Int32 value (optional)
	Int64Val        *int64     // Int64 value (optional)
	FloatVal        *float32   // Float value (optional)
	DoubleVal       *float64   // Double value (optional)
	TextVal         *TextType  // Text value (optional)
}

// NewDataPointBool initializes a boolean DataPoint.
func NewDataPointBool(measurementName string, value bool) DataPoint {
	return DataPoint{
		IsNull:          false,
		MeasurementName: measurementName,
		DataType:        BOOLEAN,
		BoolVal:         &value,
	}
}

// NewDataPointInt32 initializes an int32 DataPoint.
func NewDataPointInt32(measurementName string, value int32) DataPoint {
	return DataPoint{
		IsNull:          false,
		MeasurementName: measurementName,
		DataType:        INT32,
		Int32Val:        &value,
	}
}

// NewDataPointInt64 initializes an int64 DataPoint.
func NewDataPointInt64(measurementName string, value int64) DataPoint {
	return DataPoint{
		IsNull:          false,
		MeasurementName: measurementName,
		DataType:        INT64,
		Int64Val:        &value,
	}
}

// NewDataPointFloat initializes a float DataPoint.
func NewDataPointFloat(measurementName string, value float32) DataPoint {
	return DataPoint{
		IsNull:          false,
		MeasurementName: measurementName,
		DataType:        FLOAT,
		FloatVal:        &value,
	}
}

// NewDataPointDouble initializes a double DataPoint.
func NewDataPointDouble(measurementName string, value float64) DataPoint {
	return DataPoint{
		IsNull:          false,
		MeasurementName: measurementName,
		DataType:        DOUBLE,
		DoubleVal:       &value,
	}
}

// SetInt32 updates the DataPoint with a new int32 value.
func (dp *DataPoint) SetInt32(value int32) {
	dp.DataType = INT32
	dp.Int32Val = &value
	dp.IsNull = false
}

// SetInt64 updates the DataPoint with a new int64 value.
func (dp *DataPoint) SetInt64(value int64) {
	dp.DataType = INT64
	dp.Int64Val = &value
	dp.IsNull = false
}

// SetFloat updates the DataPoint with a new float value.
func (dp *DataPoint) SetFloat(value float32) {
	dp.DataType = FLOAT
	dp.FloatVal = &value
	dp.IsNull = false
}

// SetDouble updates the DataPoint with a new double value.
func (dp *DataPoint) SetDouble(value float64) {
	dp.DataType = DOUBLE
	dp.DoubleVal = &value
	dp.IsNull = false
}

// TextType represents textual data with a buffer and a length.
type TextType struct {
	Buffer []byte
	Length int32
}

// TsRecord represents a record containing a timestamp, device ID, and associated DataPoints.
type TsRecord struct {
	Timestamp int64       // Timestamp of the record
	DeviceID  string      // Device ID associated with the record
	Points    []DataPoint // List of DataPoints
}

// NewTsRecord initializes a TsRecord with just a device name.
func NewTsRecord(deviceName string) TsRecord {
	return TsRecord{
		DeviceID: deviceName,
	}
}

// NewTimestampedTsRecord initializes a TsRecord with a timestamp, device name, and an optional initial capacity.
func NewTimestampedTsRecord(timestamp int64, deviceName string, pointCountInRow int) TsRecord {
	record := TsRecord{
		Timestamp: timestamp,
		DeviceID:  deviceName,
	}
	// Preallocate capacity for points
	if pointCountInRow > 0 {
		record.Points = make([]DataPoint, 0, pointCountInRow)
	}
	return record
}

// AddDataPoint adds a new DataPoint to the record.
func (tsr *TsRecord) AddDataPoint(dp DataPoint) error {
	tsr.Points = append(tsr.Points, dp)
	return nil
}
