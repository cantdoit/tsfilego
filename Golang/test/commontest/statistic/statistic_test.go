package statistic

import (
	"Golang/internal/common/base"
	"Golang/internal/common/statistic"
	"testing"
)

func TestBooleanStatistic_BasicFunctionality(t *testing.T) {
	// Create a new BooleanStatistic
	stat := &statistic.BooleanStatistic{}

	// Verify initial state
	if stat.Count != 0 {
		t.Errorf("expected Count to be 0, got %d", stat.Count)
	}
	if stat.StartTime != 0 {
		t.Errorf("expected StartTime to be 0, got %d", stat.StartTime)
	}
	if stat.EndTime != 0 {
		t.Errorf("expected EndTime to be 0, got %d", stat.EndTime)
	}
	if stat.SumValue != 0 {
		t.Errorf("expected SumValue to be 0, got %d", stat.SumValue)
	}
	if stat.FirstValue != false {
		t.Errorf("expected FirstValue to be false, got %v", stat.FirstValue)
	}
	if stat.LastValue != false {
		t.Errorf("expected LastValue to be false, got %v", stat.LastValue)
	}

	// Update the statistic
	err := stat.Update(1000, true)
	// t.Logf("err: %v", stat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = stat.Update(2000, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify updated state
	if stat.Count != 2 {
		t.Errorf("expected Count to be 2, got %d", stat.Count)
	}
	if stat.StartTime != 1000 {
		t.Errorf("expected StartTime to be 1000, got %d", stat.StartTime)
	}
	if stat.EndTime != 2000 {
		t.Errorf("expected EndTime to be 2000, got %d", stat.EndTime)
	}
	if stat.SumValue != 1 {
		t.Errorf("expected SumValue to be 1, got %d", stat.SumValue)
	}
	if stat.FirstValue != true {
		t.Errorf("expected FirstValue to be true, got %v", stat.FirstValue)
	}
	if stat.LastValue != false {
		t.Errorf("expected LastValue to be false, got %v", stat.LastValue)
	}

	// Test serialization and deserialization using ByteStream
	byteStream, err := base.NewByteStream(1024) // Initialize a ByteStream with a page size of 1024 bytes.
	if err != nil {
		t.Fatalf("failed to create ByteStream: %v", err)
	}

	err = stat.SerializeTypedStat(byteStream) // Serialize the statistic into the ByteStream.
	if err != nil {
		t.Fatalf("serialization failed: %v", err)
	}

	// Rewind ByteStream for reading (reset read position).
	byteStream.ReadPos = 0

	// Deserialize the statistic from the ByteStream.
	deserializedStat := &statistic.BooleanStatistic{}
	err = deserializedStat.DeserializeTypedStat(byteStream)
	if err != nil {
		t.Fatalf("deserialization failed: %v", err)
	}

	// Verify that deserialized data matches the original
	if deserializedStat.SumValue != stat.SumValue {
		t.Errorf("expected deserialized SumValue to match %v  got %d", stat.SumValue, deserializedStat.SumValue)
	}
	if deserializedStat.FirstValue != stat.FirstValue {
		t.Errorf("expected deserialized FirstValue to match %v , got %v", stat.FirstValue, deserializedStat.FirstValue)
	}
	if deserializedStat.LastValue != stat.LastValue {
		t.Errorf("expected deserialized LastValue to match %v , got %v", stat.LastValue, deserializedStat.LastValue)
	}
}

func TestInt32Statistic_BasicFunctionality(t *testing.T) {
	// Create a new Int32Statistic
	stat := &statistic.Int32Statistic{}

	// Verify initial state
	if stat.Count != 0 {
		t.Errorf("expected Count to be 0, got %d", stat.Count)
	}
	if stat.StartTime != 0 {
		t.Errorf("expected StartTime to be 0, got %d", stat.StartTime)
	}
	if stat.EndTime != 0 {
		t.Errorf("expected EndTime to be 0, got %d", stat.EndTime)
	}
	if stat.SumValue != 0 {
		t.Errorf("expected SumValue to be 0, got %d", stat.SumValue)
	}
	if stat.MinValue != 0 {
		t.Errorf("expected MinValue to be 0, got %d", stat.MinValue)
	}
	if stat.MaxValue != 0 {
		t.Errorf("expected MaxValue to be 0, got %d", stat.MaxValue)
	}
	if stat.FirstValue != 0 {
		t.Errorf("expected FirstValue to be 0, got %d", stat.FirstValue)
	}
	if stat.LastValue != 0 {
		t.Errorf("expected LastValue to be 0, got %d", stat.LastValue)
	}

	// Update the statistic
	err := stat.Update(1000, int32(10))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = stat.Update(2000, int32(20))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify updated state
	if stat.Count != 2 {
		t.Errorf("expected Count to be 2, got %d", stat.Count)
	}
	if stat.StartTime != 1000 {
		t.Errorf("expected StartTime to be 1000, got %d", stat.StartTime)
	}
	if stat.EndTime != 2000 {
		t.Errorf("expected EndTime to be 2000, got %d", stat.EndTime)
	}
	if stat.SumValue != 30 {
		t.Errorf("expected SumValue to be 30, got %d", stat.SumValue)
	}
	if stat.MinValue != 10 {
		t.Errorf("expected MinValue to be 10, got %d", stat.MinValue)
	}
	if stat.MaxValue != 20 {
		t.Errorf("expected MaxValue to be 20, got %d", stat.MaxValue)
	}
	if stat.FirstValue != 10 {
		t.Errorf("expected FirstValue to be 10, got %d", stat.FirstValue)
	}
	if stat.LastValue != 20 {
		t.Errorf("expected LastValue to be 20, got %d", stat.LastValue)
	}

	// Test serialization and deserialization using ByteStream
	byteStream, err := base.NewByteStream(1024) // Initialize ByteStream.
	if err != nil {
		t.Fatalf("failed to create ByteStream: %v", err)
	}
	err = stat.SerializeTypedStat(byteStream) // Serialize the statistic to the ByteStream.
	if err != nil {
		t.Fatalf("serialization failed: %v", err)
	}

	// Rewind ByteStream for reading.
	byteStream.ReadPos = 0

	// Deserialize the statistic from the ByteStream.
	deserializedStat := &statistic.Int32Statistic{}
	err = deserializedStat.DeserializeTypedStat(byteStream)
	if err != nil {
		t.Fatalf("deserialization failed: %v", err)
	}

	// Verify that deserialized data matches the original
	if deserializedStat.SumValue != stat.SumValue {
		t.Errorf("expected deserialized SumValue to match %v , got %d", stat.SumValue, deserializedStat.SumValue)
	}
	if deserializedStat.MinValue != stat.MinValue {
		t.Errorf("expected deserialized MinValue to match %v, got %d", stat.MinValue, deserializedStat.MinValue)
	}
	if deserializedStat.MaxValue != stat.MaxValue {
		t.Errorf("expected deserialized MaxValue to match %v, got %d", stat.MaxValue, deserializedStat.MaxValue)
	}
	if deserializedStat.FirstValue != stat.FirstValue {
		t.Errorf("expected deserialized FirstValue to match %v, got %d", stat.FirstValue, deserializedStat.FirstValue)
	}
	if deserializedStat.LastValue != stat.LastValue {
		t.Errorf("expected deserialized LastValue to match %v, got %d", stat.LastValue, deserializedStat.LastValue)
	}
}

func TestInt64Statistic_BasicFunctionality(t *testing.T) {
	// Create a new Int64Statistic
	stat := &statistic.Int64Statistic{}

	// Verify initial state
	if stat.Count != 0 {
		t.Errorf("expected Count to be 0, got %d", stat.Count)
	}
	if stat.StartTime != 0 {
		t.Errorf("expected StartTime to be 0, got %d", stat.StartTime)
	}
	if stat.EndTime != 0 {
		t.Errorf("expected EndTime to be 0, got %d", stat.EndTime)
	}
	if stat.SumValue != 0 {
		t.Errorf("expected SumValue to be 0, got %v", stat.SumValue)
	}
	if stat.MinValue != 0 {
		t.Errorf("expected MinValue to be 0, got %v", stat.MinValue)
	}
	if stat.MaxValue != 0 {
		t.Errorf("expected MaxValue to be 0, got %v", stat.MaxValue)
	}
	if stat.FirstValue != 0 {
		t.Errorf("expected FirstValue to be 0, got %v", stat.FirstValue)
	}
	if stat.LastValue != 0 {
		t.Errorf("expected LastValue to be 0, got %v", stat.LastValue)
	}

	// Update the statistic
	err := stat.Update(1000, int64(100))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = stat.Update(2000, int64(200))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify updated state
	if stat.Count != 2 {
		t.Errorf("expected Count to be 2, got %d", stat.Count)
	}
	if stat.StartTime != 1000 {
		t.Errorf("expected StartTime to be 1000, got %d", stat.StartTime)
	}
	if stat.EndTime != 2000 {
		t.Errorf("expected EndTime to be 2000, got %d", stat.EndTime)
	}
	if stat.SumValue != 300 {
		t.Errorf("expected SumValue to be 300, got %v", stat.SumValue)
	}
	if stat.MinValue != 100 {
		t.Errorf("expected MinValue to be 100, got %d", stat.MinValue)
	}
	if stat.MaxValue != 200 {
		t.Errorf("expected MaxValue to be 200, got %d", stat.MaxValue)
	}
	if stat.FirstValue != 100 {
		t.Errorf("expected FirstValue to be 100, got %d", stat.FirstValue)
	}
	if stat.LastValue != 200 {
		t.Errorf("expected LastValue to be 200, got %d", stat.LastValue)
	}

	// Test serialization and deserialization using ByteStream
	byteStream, err := base.NewByteStream(1024) // Initialize a ByteStream with a page size of 1024 bytes
	if err != nil {
		t.Fatalf("Failed to create ByteStream: %v", err)
	}

	err = stat.SerializeTypedStat(byteStream)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Rewind ByteStream for reading
	byteStream.ReadPos = 0

	// Deserialize the statistic from the ByteStream
	deserializedStat := &statistic.Int64Statistic{}
	err = deserializedStat.DeserializeTypedStat(byteStream)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	// Verify the deserialized data matches the original
	if deserializedStat.SumValue != stat.SumValue {
		t.Errorf("expected deserialized SumValue to match %v, got %v", stat.SumValue, deserializedStat.SumValue)
	}
	if deserializedStat.MinValue != stat.MinValue {
		t.Errorf("expected deserialized MinValue to match %d, got %d", stat.MinValue, deserializedStat.MinValue)
	}
	if deserializedStat.MaxValue != stat.MaxValue {
		t.Errorf("expected deserialized MaxValue to match %d, got %d", stat.MaxValue, deserializedStat.MaxValue)
	}
	if deserializedStat.FirstValue != stat.FirstValue {
		t.Errorf("expected deserialized FirstValue to match %d, got %d", stat.FirstValue, deserializedStat.FirstValue)
	}
	if deserializedStat.LastValue != stat.LastValue {
		t.Errorf("expected deserialized LastValue to match %d, got %d", stat.LastValue, deserializedStat.LastValue)
	}
}

func TestFloatStatistic_BasicFunctionality(t *testing.T) {
	// Create a new FloatStatistic
	stat := &statistic.FloatStatistic{}

	// Verify initial state
	if stat.Count != 0 {
		t.Errorf("expected Count to be 0, got %d", stat.Count)
	}
	if stat.StartTime != 0 {
		t.Errorf("expected StartTime to be 0, got %d", stat.StartTime)
	}
	if stat.EndTime != 0 {
		t.Errorf("expected EndTime to be 0, got %d", stat.EndTime)
	}
	if stat.SumValue != 0 {
		t.Errorf("expected SumValue to be 0, got %f", stat.SumValue)
	}
	if stat.MinValue != 0 {
		t.Errorf("expected MinValue to be 0, got %f", stat.MinValue)
	}
	if stat.MaxValue != 0 {
		t.Errorf("expected MaxValue to be 0, got %f", stat.MaxValue)
	}
	if stat.FirstValue != 0 {
		t.Errorf("expected FirstValue to be 0, got %f", stat.FirstValue)
	}
	if stat.LastValue != 0 {
		t.Errorf("expected LastValue to be 0, got %f", stat.LastValue)
	}

	// Update the statistic
	err := stat.Update(1000, float32(1.5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = stat.Update(2000, float32(2.5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Serialize and Deserialize Typed Statistic
	byteStream, err := base.NewByteStream(1024)
	if err != nil {
		t.Fatalf("Failed to create ByteStream: %v", err)
	}

	err = stat.SerializeTypedStat(byteStream)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Rewind ByteStream for reading
	byteStream.ReadPos = 0

	deserializedStat := &statistic.FloatStatistic{}
	err = deserializedStat.DeserializeTypedStat(byteStream)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	// Verify deserialization values
	if deserializedStat.SumValue != stat.SumValue {
		t.Errorf("expected deserialized SumValue to match %f, got %f", stat.SumValue, deserializedStat.SumValue)
	}
	if deserializedStat.MaxValue != stat.MaxValue {
		t.Errorf("expected deserialized MaxValue to match %f, got %f", stat.MaxValue, deserializedStat.MaxValue)
	}
}

func TestDoubleStatistic_BasicFunctionality(t *testing.T) {
	// Create a new DoubleStatistic
	stat := &statistic.DoubleStatistic{}

	// Verify initial state
	if stat.Count != 0 {
		t.Errorf("expected Count to be 0, got %d", stat.Count)
	}

	// Update the statistic
	err := stat.Update(1000, float64(100.5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = stat.Update(2000, float64(200.7))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Serialize and deserialize
	byteStream, err := base.NewByteStream(1024)
	if err != nil {
		t.Fatalf("Failed to create ByteStream: %v", err)
	}

	err = stat.SerializeTypedStat(byteStream)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	byteStream.ReadPos = 0

	deserializedStat := &statistic.DoubleStatistic{}
	err = deserializedStat.DeserializeTypedStat(byteStream)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	// Verify data matches
	if deserializedStat.SumValue != stat.SumValue {
		t.Errorf("expected SumValue to match %f, got %f", stat.SumValue, deserializedStat.SumValue)
	}
}
