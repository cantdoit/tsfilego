package statistic

import (
	"Golang/internal/common/base"
	"errors"
	"fmt"
)

// Statistic defines the base statistic structure
type Statistic struct {
	Count     int32           // Number of entries in the statistic
	StartTime int64           // Start time of the statistic
	EndTime   int64           // End time of the statistic
	DataType  base.TSDataType // The data type of this statistic (e.g., BOOLEAN, INT32)
}

func (s *Statistic) GetType() base.TSDataType {
	//TODO implement me
	panic("implement me")
}

func (s *Statistic) Clone() Interface {
	//TODO implement me
	panic("implement me")
}

func (s *Statistic) DeserializeTypedStat(stream *base.ByteStream) error {
	//TODO implement me
	panic("implement me")
}

func (s *Statistic) SerializeTypedStat(stream *base.ByteStream) error {
	//TODO implement me
	panic("implement me")
}

// NewStatistic creates a new base statistic
func NewStatistic(dataType base.TSDataType) (*Statistic, error) {
	if !base.IsValidDataType(dataType) {
		return nil, errors.New("invalid data type for statistic")
	}
	return &Statistic{
		Count:     0,
		StartTime: 0,
		EndTime:   0,
		DataType:  dataType,
	}, nil
}

// Reset clears the data in the statistic
func (s *Statistic) Reset() {
	s.Count = 0
	s.StartTime = 0
	s.EndTime = 0
}

func (s *Statistic) cloneFrom(src *Statistic) {
	s.Count = src.Count
	s.StartTime = src.StartTime
	s.EndTime = src.EndTime
}

// Update updates the statistic with a time and value (int64)
func (s *Statistic) Update(time int64, value int64) error {
	if s.Count == 0 {
		// Initialize first values
		s.StartTime = time
		s.EndTime = time
	} else {
		if time < s.StartTime {
			s.StartTime = time
		}
		if time > s.EndTime {
			s.EndTime = time
		}
	}
	// s.Count++
	return nil
}

/*
// UpdateBoolean is for updating boolean statistics
func (s *Statistic) UpdateBoolean(time int64, value bool) error {
	return s.Update(time, toInt64(value))
}
*/

// Serialize serializes the statistic to a ByteStream
func (s *Statistic) Serialize(out *base.ByteStream) error {
	// Create an instance of the utility
	util := base.SerializationUtil{}

	if err := util.WriteVarUint(uint32(s.Count), out); err != nil {
		return err
	}
	if err := util.WriteUint64(uint64(s.StartTime), out); err != nil {
		return err
	}
	if err := util.WriteUint64(uint64(s.EndTime), out); err != nil {
		return err
	}
	return nil
}

// Deserialize deserializes the statistic from a ByteStream
func (s *Statistic) Deserialize(in *base.ByteStream) error {
	// Create an instance of the utility
	util := base.SerializationUtil{}

	count, err := util.ReadVarUint(in)
	if err != nil {
		return err
	}
	s.Count = int32(count)

	startTime, err := util.ReadUint64(in)
	if err != nil {
		return err
	}
	s.StartTime = int64(startTime)

	endTime, err := util.ReadUint64(in)
	if err != nil {
		return err
	}
	s.EndTime = int64(endTime)

	return nil
}

// ToString produces a string representation of the statistic
func (s *Statistic) ToString() string {
	return fmt.Sprintf("Statistic{Count: %d, StartTime: %d, EndTime: %d, DataType: %s}",
		s.Count, s.StartTime, s.EndTime, s.DataType)
}

// Helper: Converts a boolean to int64 (used in statistics)
func toInt64(value bool) int64 {
	if value {
		return 1
	}
	return 0
}

//////////////////////
/// Bool statistic ///
//////////////////////

// BooleanStatistic represents a statistic for boolean data.
type BooleanStatistic struct {
	Statistic
	SumValue   int64 // Sum of all boolean values (as integers)
	FirstValue bool  // First boolean value observed
	LastValue  bool  // Last boolean value observed
}

// NewBooleanStatistic creates a new BooleanStatistic.
func NewBooleanStatistic() *BooleanStatistic {
	stat, _ := NewStatistic(base.BOOLEAN)
	return &BooleanStatistic{
		Statistic:  *stat,
		SumValue:   0,
		FirstValue: false,
		LastValue:  false,
	}
}

// GetType returns the data type of the statistic
func (bs *BooleanStatistic) GetType() base.TSDataType {
	return base.BOOLEAN
}

func (bs *BooleanStatistic) cloneFrom(src *BooleanStatistic) {
	bs.Statistic.cloneFrom(&src.Statistic) // Copy common statistic fields
	bs.SumValue = src.SumValue
	bs.FirstValue = src.FirstValue
	bs.LastValue = src.LastValue
}

func (bs *BooleanStatistic) Clone() Interface {
	clone := NewBooleanStatistic()
	clone.cloneFrom(bs)
	return clone
}

// Update updates the statistic with a boolean value and its associated timestamp.
func (bs *BooleanStatistic) Update(time int64, value interface{}) error {
	if err := bs.Statistic.Update(time, 0); err != nil {
		return err
	}

	boolValue, ok := value.(bool)
	if !ok {
		return errors.New("invalid value type, expected bool")
	}

	if bs.Count == 0 {
		bs.FirstValue = boolValue
	}
	bs.SumValue += toInt64(boolValue)
	bs.LastValue = boolValue
	bs.Count++
	return nil
}

// SerializeTypedStat serializes the BooleanStatistic-specific fields.
func (bs *BooleanStatistic) SerializeTypedStat(out *base.ByteStream) error {
	util := base.SerializationUtil{}
	if err := util.WriteBool(bs.FirstValue, out); err != nil {
		return err
	}
	if err := util.WriteBool(bs.LastValue, out); err != nil {
		return err
	}
	if err := util.WriteUint64(uint64(bs.SumValue), out); err != nil {
		return err
	}
	return nil
}

// DeserializeTypedStat deserializes the BooleanStatistic-specific fields.
func (bs *BooleanStatistic) DeserializeTypedStat(in *base.ByteStream) error {
	util := base.SerializationUtil{}
	firstValue, err := util.ReadBool(in)
	if err != nil {
		return err
	}
	bs.FirstValue = firstValue

	lastValue, err := util.ReadBool(in)
	if err != nil {
		return err
	}
	bs.LastValue = lastValue

	sumValue, err := util.ReadUint64(in)
	if err != nil {
		return err
	}
	bs.SumValue = int64(sumValue)

	return nil
}

// MergeWith merges another BooleanStatistic into the current one.
func (bs *BooleanStatistic) MergeWith(other Interface) error {
	// Ensure the other statistic is of the BooleanStatistic type.
	otherStat, ok := other.(*BooleanStatistic)
	if !ok {
		return errors.New("statistic type does not match for merge")
	}

	// If the other statistic has no data, there's nothing to merge.
	if otherStat.Count == 0 {
		return nil
	}

	// If the current statistic has no data, clone from the other.
	if bs.Count == 0 {
		bs.Count = otherStat.Count
		bs.StartTime = otherStat.StartTime
		bs.EndTime = otherStat.EndTime
		bs.SumValue = otherStat.SumValue
		bs.FirstValue = otherStat.FirstValue
		bs.LastValue = otherStat.LastValue
	} else {
		// Merge values.
		bs.Count += otherStat.Count
		if otherStat.StartTime < bs.StartTime {
			bs.StartTime = otherStat.StartTime
			bs.FirstValue = otherStat.FirstValue
		}
		if otherStat.EndTime > bs.EndTime {
			bs.EndTime = otherStat.EndTime
			bs.LastValue = otherStat.LastValue
		}
		bs.SumValue += otherStat.SumValue
	}

	return nil
}

// Reset clears the data in the statistic
func (bs *BooleanStatistic) Reset() {
	bs.SumValue = 0
	bs.Count = 0
	bs.StartTime = 0
	bs.EndTime = 0
	bs.FirstValue = false
	bs.LastValue = false
}

///////////////////////
/// Int32 statistic ///
///////////////////////

// Int32Statistic represents a statistic for int32 data.
type Int32Statistic struct {
	Statistic
	SumValue   int64 // Sum of all observed values
	MinValue   int32 // Minimum value observed
	MaxValue   int32 // Maximum value observed
	FirstValue int32 // First observed value
	LastValue  int32 // Last observed value
}

// NewInt32Statistic creates a new Int32Statistic.
func NewInt32Statistic() *Int32Statistic {
	return &Int32Statistic{
		Statistic: Statistic{DataType: base.INT32},
		SumValue:  0,
	}
}

// GetType returns the data type of the statistic
func (is *Int32Statistic) GetType() base.TSDataType {
	return base.INT32
}

func (is *Int32Statistic) cloneFrom(src *Int32Statistic) {
	is.Statistic.cloneFrom(&src.Statistic) // Copy common statistic fields
	is.SumValue = src.SumValue
	is.MinValue = src.MinValue
	is.MaxValue = src.MaxValue
	is.FirstValue = src.FirstValue
	is.LastValue = src.LastValue
}

func (is *Int32Statistic) Clone() Interface {
	clone := NewInt32Statistic()
	clone.cloneFrom(is)
	return clone
}

// Update updates the statistic with an int32 value and its associated timestamp.
func (is *Int32Statistic) Update(time int64, value interface{}) error {
	if err := is.Statistic.Update(time, 0); err != nil {
		return err
	}
	int32Value, ok := value.(int32)
	if !ok {
		return errors.New("invalid value type, expected int32")
	}

	if is.Count == 0 {
		is.MinValue = int32Value
		is.MaxValue = int32Value
		is.FirstValue = int32Value
	} else {
		if int32Value < is.MinValue {
			is.MinValue = int32Value
		}
		if int32Value > is.MaxValue {
			is.MaxValue = int32Value
		}
	}
	is.SumValue += int64(int32Value)
	is.LastValue = int32Value
	is.Count++
	return nil
}

// SerializeTypedStat serializes the Int32Statistic-specific fields.
func (is *Int32Statistic) SerializeTypedStat(out *base.ByteStream) error {
	if err := is.Statistic.Serialize(out); err != nil {
		return err
	}

	util := base.SerializationUtil{}
	if err := util.WriteUint32(uint32(is.MinValue), out); err != nil {
		return err
	}
	if err := util.WriteUint32(uint32(is.MaxValue), out); err != nil {
		return err
	}
	if err := util.WriteUint32(uint32(is.FirstValue), out); err != nil {
		return err
	}
	if err := util.WriteUint32(uint32(is.LastValue), out); err != nil {
		return err
	}
	if err := util.WriteUint64(uint64(is.SumValue), out); err != nil {
		return err
	}
	return nil
}

// DeserializeTypedStat deserializes the Int32Statistic-specific fields.
func (is *Int32Statistic) DeserializeTypedStat(in *base.ByteStream) error {
	if err := is.Statistic.Deserialize(in); err != nil {
		return err
	}

	util := base.SerializationUtil{}
	minValue, err := util.ReadUint32(in)
	if err != nil {
		return err
	}
	is.MinValue = int32(minValue)

	maxValue, err := util.ReadUint32(in)
	if err != nil {
		return err
	}
	is.MaxValue = int32(maxValue)

	firstValue, err := util.ReadUint32(in)
	if err != nil {
		return err
	}
	is.FirstValue = int32(firstValue)

	lastValue, err := util.ReadUint32(in)
	if err != nil {
		return err
	}
	is.LastValue = int32(lastValue)

	sumValue, err := util.ReadUint64(in)
	if err != nil {
		return err
	}
	is.SumValue = int64(sumValue)

	return nil
}

// MergeWith merges another Int32Statistic into the current one.
func (is *Int32Statistic) MergeWith(other Interface) error {
	// Ensure the other statistic is of the Int32Statistic type.
	otherStat, ok := other.(*Int32Statistic)
	if !ok {
		return errors.New("statistic type does not match for merge")
	}

	// If the other statistic has no data, there's nothing to merge.
	if otherStat.Count == 0 {
		return nil
	}

	// If the current statistic has no data, clone from the other.
	if is.Count == 0 {
		is.Count = otherStat.Count
		is.StartTime = otherStat.StartTime
		is.EndTime = otherStat.EndTime
		is.SumValue = otherStat.SumValue
		is.FirstValue = otherStat.FirstValue
		is.LastValue = otherStat.LastValue
		is.MinValue = otherStat.MinValue
		is.MaxValue = otherStat.MaxValue
	} else {
		// Merge values.
		is.Count += otherStat.Count
		if otherStat.StartTime < is.StartTime {
			is.StartTime = otherStat.StartTime
			is.FirstValue = otherStat.FirstValue
		}
		if otherStat.EndTime > is.EndTime {
			is.EndTime = otherStat.EndTime
			is.LastValue = otherStat.LastValue
		}
		is.SumValue += otherStat.SumValue
		if otherStat.MinValue < is.MinValue {
			is.MinValue = otherStat.MinValue
		}
		if otherStat.MaxValue > is.MaxValue {
			is.MaxValue = otherStat.MaxValue
		}
	}

	return nil
}

// Reset clears the data in the statistic
func (is *Int32Statistic) Reset() {
	is.MinValue = 0
	is.MaxValue = 0
	is.FirstValue = 0
	is.LastValue = 0
	is.SumValue = 0
	is.Count = 0
	is.StartTime = 0
	is.EndTime = 0
}

///////////////////////
/// Int64 statistic ///
///////////////////////

// Int64Statistic represents a statistic for int64 data.
type Int64Statistic struct {
	Statistic
	SumValue   float64 // Sum of all observed values as a floating-point
	MinValue   int64   // Minimum value observed
	MaxValue   int64   // Maximum value observed
	FirstValue int64   // First observed value
	LastValue  int64   // Last observed value
}

// NewInt64Statistic creates a new Int64Statistic.
func NewInt64Statistic() *Int64Statistic {
	return &Int64Statistic{
		Statistic: Statistic{DataType: base.INT64},
	}
}

func (is *Int64Statistic) cloneFrom(src *Int64Statistic) {
	is.Statistic.cloneFrom(&src.Statistic) // Copy common statistic fields
	is.SumValue = src.SumValue
	is.MinValue = src.MinValue
	is.MaxValue = src.MaxValue
	is.FirstValue = src.FirstValue
	is.LastValue = src.LastValue
}

func (is *Int64Statistic) Clone() Interface {
	clone := NewInt64Statistic()
	clone.cloneFrom(is)
	return clone
}

// GetType returns the data type of the statistic
func (is *Int64Statistic) GetType() base.TSDataType {
	return base.INT64
}

// Update updates the statistic with an int64 value and its associated timestamp.
func (is *Int64Statistic) Update(time int64, value interface{}) error {
	if err := is.Statistic.Update(time, 0); err != nil {
		return err
	}
	int64Value, ok := value.(int64)
	if !ok {
		return errors.New("invalid value type, expected int64")
	}
	if is.Count == 0 {
		is.MinValue = int64Value
		is.MaxValue = int64Value
		is.FirstValue = int64Value
	} else {
		if int64Value < is.MinValue {
			is.MinValue = int64Value
		}
		if int64Value > is.MaxValue {
			is.MaxValue = int64Value
		}
	}
	is.SumValue += float64(int64Value)
	is.LastValue = int64Value
	is.Count++
	return nil
}

// SerializeTypedStat serializes the Int64Statistic-specific fields.
func (is *Int64Statistic) SerializeTypedStat(out *base.ByteStream) error {
	util := base.SerializationUtil{}
	if err := util.WriteUint64(uint64(is.MinValue), out); err != nil {
		return err
	}
	if err := util.WriteUint64(uint64(is.MaxValue), out); err != nil {
		return err
	}
	if err := util.WriteUint64(uint64(is.FirstValue), out); err != nil {
		return err
	}
	if err := util.WriteUint64(uint64(is.LastValue), out); err != nil {
		return err
	}
	if err := util.WriteDouble(is.SumValue, out); err != nil {
		return err
	}
	return nil
}

// DeserializeTypedStat deserializes the Int64Statistic-specific fields.
func (is *Int64Statistic) DeserializeTypedStat(in *base.ByteStream) error {
	util := base.SerializationUtil{}
	minValue, err := util.ReadUint64(in)
	if err != nil {
		return err
	}
	is.MinValue = int64(minValue)

	maxValue, err := util.ReadUint64(in)
	if err != nil {
		return err
	}
	is.MaxValue = int64(maxValue)

	firstValue, err := util.ReadUint64(in)
	if err != nil {
		return err
	}
	is.FirstValue = int64(firstValue)

	lastValue, err := util.ReadUint64(in)
	if err != nil {
		return err
	}
	is.LastValue = int64(lastValue)

	sumValue, err := util.ReadDouble(in)
	if err != nil {
		return err
	}
	is.SumValue = float64(sumValue)

	return nil
}

// MergeWith merges another Int64Statistic into the current one.
func (is *Int64Statistic) MergeWith(other Interface) error {
	// Ensure the other statistic is of the Int64Statistic type.
	otherStat, ok := other.(*Int64Statistic)
	if !ok {
		return errors.New("statistic type does not match for merge")
	}

	// If the other statistic has no data, there's nothing to merge.
	if otherStat.Count == 0 {
		return nil
	}

	// If the current statistic has no data, clone from the other.
	if is.Count == 0 {
		is.Count = otherStat.Count
		is.StartTime = otherStat.StartTime
		is.EndTime = otherStat.EndTime
		is.SumValue = otherStat.SumValue
		is.FirstValue = otherStat.FirstValue
		is.LastValue = otherStat.LastValue
		is.MinValue = otherStat.MinValue
		is.MaxValue = otherStat.MaxValue
	} else {
		// Merge values.
		is.Count += otherStat.Count
		if otherStat.StartTime < is.StartTime {
			is.StartTime = otherStat.StartTime
			is.FirstValue = otherStat.FirstValue
		}
		if otherStat.EndTime > is.EndTime {
			is.EndTime = otherStat.EndTime
			is.LastValue = otherStat.LastValue
		}
		is.SumValue += otherStat.SumValue
		if otherStat.MinValue < is.MinValue {
			is.MinValue = otherStat.MinValue
		}
		if otherStat.MaxValue > is.MaxValue {
			is.MaxValue = otherStat.MaxValue
		}
	}

	return nil
}

// Reset clears the data in the statistic
func (is *Int64Statistic) Reset() {
	is.MinValue = 0
	is.MaxValue = 0
	is.FirstValue = 0
	is.LastValue = 0
	is.SumValue = 0
	is.Count = 0
	is.StartTime = 0
	is.EndTime = 0
}

///////////////////////
/// Float statistic ///
///////////////////////

// FloatStatistic represents statistics for float data.
type FloatStatistic struct {
	Statistic
	SumValue   float64 // Sum of all observed values
	MinValue   float32 // Minimum value observed
	MaxValue   float32 // Maximum value observed
	FirstValue float32 // First observed value
	LastValue  float32 // Last observed value
}

// NewFloatStatistic creates a new FloatStatistic.
func NewFloatStatistic() *FloatStatistic {
	return &FloatStatistic{
		Statistic: Statistic{DataType: base.FLOAT},
		SumValue:  0,
	}
}

// GetType returns the data type of the statistic
func (fs *FloatStatistic) GetType() base.TSDataType {
	return base.FLOAT
}

func (fs *FloatStatistic) cloneFrom(src *FloatStatistic) {
	fs.Statistic.cloneFrom(&src.Statistic)
	fs.SumValue = src.SumValue
	fs.MinValue = src.MinValue
	fs.MaxValue = src.MaxValue
	fs.FirstValue = src.FirstValue
	fs.LastValue = src.LastValue
}

func (fs *FloatStatistic) Clone() Interface {
	clone := NewFloatStatistic()
	clone.cloneFrom(fs)
	return clone
}

// Update updates the statistic with a float value and its associated timestamp.
func (fs *FloatStatistic) Update(time int64, value interface{}) error {
	// Validate the value type as float32
	floatValue, ok := value.(float32)
	if !ok {
		return errors.New("invalid value type, expected float32")
	}

	// Update the base Statistic fields
	if err := fs.Statistic.Update(time, int64(floatValue)); err != nil {
		return err
	}

	// Handle the first update differently
	if fs.Count == 0 {
		fs.MinValue = floatValue
		fs.MaxValue = floatValue
		fs.FirstValue = floatValue
	} else {
		if floatValue < fs.MinValue {
			fs.MinValue = floatValue
		}
		if floatValue > fs.MaxValue {
			fs.MaxValue = floatValue
		}
	}

	// Accumulate the sum and update the last value
	fs.SumValue += float64(floatValue)
	fs.LastValue = floatValue

	// Increment the count
	fs.Count++
	return nil
}

// SerializeTypedStat serializes the FloatStatistic-specific fields.
func (fs *FloatStatistic) SerializeTypedStat(out *base.ByteStream) error {
	util := base.SerializationUtil{}
	if err := util.WriteFloat(fs.MinValue, out); err != nil {
		return err
	}
	if err := util.WriteFloat(fs.MaxValue, out); err != nil {
		return err
	}
	if err := util.WriteFloat(fs.FirstValue, out); err != nil {
		return err
	}
	if err := util.WriteFloat(fs.LastValue, out); err != nil {
		return err
	}
	if err := util.WriteDouble(fs.SumValue, out); err != nil {
		return err
	}
	return nil
}

// DeserializeTypedStat deserializes the FloatStatistic-specific fields.
func (fs *FloatStatistic) DeserializeTypedStat(in *base.ByteStream) error {
	util := base.SerializationUtil{}
	minValue, err := util.ReadFloat(in)
	if err != nil {
		return err
	}
	fs.MinValue = minValue

	maxValue, err := util.ReadFloat(in)
	if err != nil {
		return err
	}
	fs.MaxValue = maxValue

	firstValue, err := util.ReadFloat(in)
	if err != nil {
		return err
	}
	fs.FirstValue = firstValue

	lastValue, err := util.ReadFloat(in)
	if err != nil {
		return err
	}
	fs.LastValue = lastValue

	sumValue, err := util.ReadDouble(in)
	if err != nil {
		return err
	}
	fs.SumValue = sumValue

	return nil
}

// MergeWith merges another FloatStatistic into the current one.
func (fs *FloatStatistic) MergeWith(other Interface) error {
	// Ensure the other statistic is of the FloatStatistic type.
	otherStat, ok := other.(*FloatStatistic)
	if !ok {
		return errors.New("statistic type does not match for merge")
	}

	// If the other statistic has no data, there's nothing to merge.
	if otherStat.Count == 0 {
		return nil
	}

	// If the current statistic has no data, clone from the other.
	if fs.Count == 0 {
		fs.Count = otherStat.Count
		fs.StartTime = otherStat.StartTime
		fs.EndTime = otherStat.EndTime
		fs.SumValue = otherStat.SumValue
		fs.FirstValue = otherStat.FirstValue
		fs.LastValue = otherStat.LastValue
		fs.MinValue = otherStat.MinValue
		fs.MaxValue = otherStat.MaxValue
	} else {
		// Merge values.
		fs.Count += otherStat.Count
		if otherStat.StartTime < fs.StartTime {
			fs.StartTime = otherStat.StartTime
			fs.FirstValue = otherStat.FirstValue
		}
		if otherStat.EndTime > fs.EndTime {
			fs.EndTime = otherStat.EndTime
			fs.LastValue = otherStat.LastValue
		}
		fs.SumValue += otherStat.SumValue
		if otherStat.MinValue < fs.MinValue {
			fs.MinValue = otherStat.MinValue
		}
		if otherStat.MaxValue > fs.MaxValue {
			fs.MaxValue = otherStat.MaxValue
		}
	}

	return nil
}

////////////////////////
/// Double statistic ///
////////////////////////

// DoubleStatistic represents statistics for double data.
type DoubleStatistic struct {
	Statistic
	SumValue   float64 // Sum of all observed values
	MinValue   float64 // Minimum value observed
	MaxValue   float64 // Maximum value observed
	FirstValue float64 // First observed value
	LastValue  float64 // Last observed value
}

// NewDoubleStatistic creates a new DoubleStatistic.
func NewDoubleStatistic() *DoubleStatistic {
	return &DoubleStatistic{
		Statistic: Statistic{DataType: base.DOUBLE},
		SumValue:  0,
	}
}

// GetType returns the data type of the statistic
func (ds *DoubleStatistic) GetType() base.TSDataType {
	return base.DOUBLE
}

func (ds *DoubleStatistic) cloneFrom(src *DoubleStatistic) {
	ds.Statistic.cloneFrom(&src.Statistic)
	ds.SumValue = src.SumValue
	ds.MinValue = src.MinValue
	ds.MaxValue = src.MaxValue
	ds.FirstValue = src.FirstValue
	ds.LastValue = src.LastValue
}

func (ds *DoubleStatistic) Clone() Interface {
	clone := NewDoubleStatistic()
	clone.cloneFrom(ds)
	return clone
}

// Update updates the statistic with a double value and its associated timestamp.
func (ds *DoubleStatistic) Update(time int64, value interface{}) error {

	doubleValue, ok := value.(float64)
	if !ok {
		return errors.New("invalid value type, expected float64")
	}
	if err := ds.Statistic.Update(time, int64(doubleValue)); err != nil {
		return err
	}
	if ds.Count == 0 {
		ds.MinValue = doubleValue
		ds.MaxValue = doubleValue
		ds.FirstValue = doubleValue
	} else {
		if doubleValue < ds.MinValue {
			ds.MinValue = doubleValue
		}
		if doubleValue > ds.MaxValue {
			ds.MaxValue = doubleValue
		}
	}
	ds.SumValue += doubleValue
	ds.LastValue = doubleValue
	ds.Count++
	return nil
}

// SerializeTypedStat serializes the DoubleStatistic-specific fields.
func (ds *DoubleStatistic) SerializeTypedStat(out *base.ByteStream) error {
	util := base.SerializationUtil{}
	if err := util.WriteDouble(ds.MinValue, out); err != nil {
		return err
	}
	if err := util.WriteDouble(ds.MaxValue, out); err != nil {
		return err
	}
	if err := util.WriteDouble(ds.FirstValue, out); err != nil {
		return err
	}
	if err := util.WriteDouble(ds.LastValue, out); err != nil {
		return err
	}
	if err := util.WriteDouble(ds.SumValue, out); err != nil {
		return err
	}
	return nil
}

// DeserializeTypedStat deserializes the DoubleStatistic-specific fields.
func (ds *DoubleStatistic) DeserializeTypedStat(in *base.ByteStream) error {
	util := base.SerializationUtil{}
	minValue, err := util.ReadDouble(in)
	if err != nil {
		return err
	}
	ds.MinValue = minValue

	maxValue, err := util.ReadDouble(in)
	if err != nil {
		return err
	}
	ds.MaxValue = maxValue

	firstValue, err := util.ReadDouble(in)
	if err != nil {
		return err
	}
	ds.FirstValue = firstValue

	lastValue, err := util.ReadDouble(in)
	if err != nil {
		return err
	}
	ds.LastValue = lastValue

	sumValue, err := util.ReadDouble(in)
	if err != nil {
		return err
	}
	ds.SumValue = sumValue

	return nil
}

// MergeWith merges another DoubleStatistic into the current one.
func (ds *DoubleStatistic) MergeWith(other Interface) error {
	// Ensure the other statistic is of the DoubleStatistic type.
	otherStat, ok := other.(*DoubleStatistic)
	if !ok {
		return errors.New("statistic type does not match for merge")
	}

	// If the other statistic has no data, there's nothing to merge.
	if otherStat.Count == 0 {
		return nil
	}

	// Merge the statistics.
	if ds.Count == 0 {
		// If the current statistic has no data, clone from the other.
		ds.cloneFrom(otherStat)
	} else {
		// Combine the fields.
		ds.Count += otherStat.Count
		if otherStat.StartTime < ds.StartTime {
			ds.StartTime = otherStat.StartTime
			ds.FirstValue = otherStat.FirstValue
		}
		if otherStat.EndTime > ds.EndTime {
			ds.EndTime = otherStat.EndTime
			ds.LastValue = otherStat.LastValue
		}
		ds.SumValue += otherStat.SumValue
		if otherStat.MinValue < ds.MinValue {
			ds.MinValue = otherStat.MinValue
		}
		if otherStat.MaxValue > ds.MaxValue {
			ds.MaxValue = otherStat.MaxValue
		}
	}

	return nil
}

// ToString produces a string representation of the FloatStatistic.
func (fs *FloatStatistic) ToString() string {
	return fmt.Sprintf("FloatStatistic{Count: %d, StartTime: %d, EndTime: %d, FirstValue: %f, LastValue: %f, SumValue: %f, MinValue: %f, MaxValue: %f}",
		fs.Count, fs.StartTime, fs.EndTime, fs.FirstValue, fs.LastValue, fs.SumValue, fs.MinValue, fs.MaxValue)
}

/////////////////////////////
/// Vector/Time statistic ///
/////////////////////////////

// TimeStatistic represents statistics for time data.
type TimeStatistic struct {
	Statistic
}

// NewTimeStatistic creates a new TimeStatistic.
func NewTimeStatistic() *TimeStatistic {
	return &TimeStatistic{
		Statistic: Statistic{DataType: base.VECTOR},
	}
}

// GetType returns the data type of the statistic.
func (ts *TimeStatistic) GetType() base.TSDataType {
	return base.VECTOR
}

func (ts *TimeStatistic) cloneFrom(src *TimeStatistic) {
	ts.Statistic.cloneFrom(&src.Statistic)
}

// Clone creates a deep copy of the current TimeStatistic.
func (ts *TimeStatistic) Clone() Interface {
	clone := NewTimeStatistic()
	clone.cloneFrom(ts)
	return clone
}

// Update updates the statistic with an associated timestamp.
func (ts *TimeStatistic) Update(time int64, value interface{}) error {
	// Since TimeStatistic only needs time updates, ignore value and just use time.
	err := ts.Statistic.Update(time, 0) // Pass 0 as there is no "value".
	if err != nil {
		return err
	}
	ts.Count++
	return nil
}

// SerializeTypedStat provides a no-op serialize for TimeStatistic-specific fields.
func (ts *TimeStatistic) SerializeTypedStat(out *base.ByteStream) error {
	// No additional fields in TimeStatistic to serialize.
	return nil
}

// DeserializeTypedStat provides a no-op deserialize for TimeStatistic-specific fields.
func (ts *TimeStatistic) DeserializeTypedStat(in *base.ByteStream) error {
	// No additional fields in TimeStatistic to deserialize.
	return nil
}

// ToString produces a string representation of the TimeStatistic.
func (ts *TimeStatistic) ToString() string {
	return fmt.Sprintf("TimeStatistic{Count: %d, StartTime: %d, EndTime: %d}",
		ts.Count, ts.StartTime, ts.EndTime)
}

// MergeWith merges another TimeStatistic into the current one.
func (ts *TimeStatistic) MergeWith(other Interface) error {
	// Ensure the other statistic is of the TimeStatistic type.
	otherStat, ok := other.(*TimeStatistic)
	if !ok {
		return errors.New("statistic type does not match for merge")
	}

	// If the other statistic has no data, there's nothing to merge.
	if otherStat.Count == 0 {
		return nil
	}

	// If the current statistic has no data, clone from the other.
	if ts.Count == 0 {
		ts.Count = otherStat.Count
		ts.StartTime = otherStat.StartTime
		ts.EndTime = otherStat.EndTime
	} else {
		// Merge values.
		ts.Count += otherStat.Count
		if otherStat.StartTime < ts.StartTime {
			ts.StartTime = otherStat.StartTime
		}
		if otherStat.EndTime > ts.EndTime {
			ts.EndTime = otherStat.EndTime
		}
	}

	return nil
}

/////////////////////////
/// Factory statistic ///
/////////////////////////

// Factory provides methods to create instances of different types of statistics.
type Factory struct{}

// AllocStatistic dynamically allocates a new statistic instance based on the provided data type.
func (sf *Factory) AllocStatistic(dataType base.TSDataType) (Interface, error) {
	switch dataType {
	case base.BOOLEAN:
		return NewBooleanStatistic(), nil
	case base.INT32:
		return NewInt32Statistic(), nil
	case base.INT64:
		return NewInt64Statistic(), nil
	case base.FLOAT:
		return NewFloatStatistic(), nil
	case base.DOUBLE:
		return NewDoubleStatistic(), nil
	case base.VECTOR:
		return NewTimeStatistic(), nil
	case base.TEXT:
		return nil, errors.New("STATISTIC_FACTORY: unsupported data type for allocation")
	default:
		return nil, fmt.Errorf("STATISTIC_FACTORY: unknown data type '%v'", dataType)
	}
}

// CloneStatistic clones an existing statistic object to a new one, maintaining its type.
func (sf *Factory) CloneStatistic(from Interface, to Interface, dataType base.TSDataType) error {
	if from == nil || to == nil {
		return errors.New("CloneStatistic: input statistics cannot be nil")
	}

	switch dataType {
	case base.BOOLEAN:
		if src, ok := from.(*BooleanStatistic); ok {
			if dest, ok := to.(*BooleanStatistic); ok {
				dest.cloneFrom(src)
				return nil
			}
		}
	case base.INT32:
		if src, ok := from.(*Int32Statistic); ok {
			if dest, ok := to.(*Int32Statistic); ok {
				dest.cloneFrom(src)
				return nil
			}
		}
	case base.INT64:
		if src, ok := from.(*Int64Statistic); ok {
			if dest, ok := to.(*Int64Statistic); ok {
				dest.cloneFrom(src)
				return nil
			}
		}
	case base.FLOAT:
		if src, ok := from.(*FloatStatistic); ok {
			if dest, ok := to.(*FloatStatistic); ok {
				dest.cloneFrom(src)
				return nil
			}
		}
	case base.DOUBLE:
		if src, ok := from.(*DoubleStatistic); ok {
			if dest, ok := to.(*DoubleStatistic); ok {
				dest.cloneFrom(src)
				return nil
			}
		}
	default:
		return fmt.Errorf("CloneStatistic: unsupported type '%v'", dataType)
	}

	return errors.New("CloneStatistic: type mismatch while cloning")
}

// FreeStatistic "frees" the statistic in Go terms (simply sets the reference to nil).
func (sf *Factory) FreeStatistic(stat Interface) {
	stat = nil // In Go, this simply removes the reference for GC
}

// Interface Example interface to unify different typed statistics
type Interface interface {
	// Update common method signatures that every statistic needs to implement
	Update(time int64, value interface{}) error
	GetType() base.TSDataType
	ToString() string
	Clone() Interface
	DeserializeTypedStat(stream *base.ByteStream) error
	SerializeTypedStat(stream *base.ByteStream) error
	Reset()
	MergeWith(other Interface) error
}

// CloneStatistic dynamically clones a statistic from one object to another, based on the type.
func CloneStatistic(from Interface, to Interface, dataType base.TSDataType) error {
	if from == nil || to == nil {
		return errors.New("clone_statistic: 'from' or 'to' statistic cannot be nil")
	}

	// Cast and clone based on type
	switch dataType {
	case base.BOOLEAN:
		if src, ok := from.(*BooleanStatistic); ok {
			if dest, okTo := to.(*BooleanStatistic); okTo {
				dest.cloneFrom(src)
				return nil
			}
		}
	case base.INT32:
		if src, ok := from.(*Int32Statistic); ok {
			if dest, okTo := to.(*Int32Statistic); okTo {
				dest.cloneFrom(src)
				return nil
			}
		}
	case base.INT64:
		if src, ok := from.(*Int64Statistic); ok {
			if dest, okTo := to.(*Int64Statistic); okTo {
				dest.cloneFrom(src)
				return nil
			}
		}
	case base.FLOAT:
		if src, ok := from.(*FloatStatistic); ok {
			if dest, okTo := to.(*FloatStatistic); okTo {
				dest.cloneFrom(src)
				return nil
			}
		}
	case base.DOUBLE:
		if src, ok := from.(*DoubleStatistic); ok {
			if dest, okTo := to.(*DoubleStatistic); okTo {
				dest.cloneFrom(src)
				return nil
			}
		}
	case base.TEXT:
		return errors.New("clone_statistic: TEXT type not supported")
	case base.VECTOR:
		// If you have a TimeStatistic type, implement cloning logic here
		return errors.New("clone_statistic: VECTOR type not supported")
	default:
		return fmt.Errorf("clone_statistic: unsupported data type '%v'", dataType)
	}

	return errors.New("clone_statistic: type mismatch between 'from' and 'to'")
}
