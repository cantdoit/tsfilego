package base

import (
	"Golang/internal/common/base"
	"testing"
)

type SerializationUtilTestSuite struct {
	byteStream        *base.ByteStream
	serializationUtil *base.SerializationUtil
}

// Setup initializes the test suite with a ByteStream and SerializationUtil instance
func (suite *SerializationUtilTestSuite) Setup(t *testing.T, pageSize int) {
	var err error
	suite.byteStream, err = base.NewByteStream(uint32(pageSize))
	if err != nil {
		t.Fatalf("Failed to set up ByteStream: %v", err)
	}
	suite.serializationUtil = &base.SerializationUtil{}
}

// TearDown cleans up the test suite
func (suite *SerializationUtilTestSuite) TearDown() {
	suite.byteStream = nil // Allow garbage collection
	suite.serializationUtil = nil
}

// TestWriteReadUint8 verifies writing and reading a uint8 value.
func TestWriteReadUint8(t *testing.T) {
	suite := &SerializationUtilTestSuite{}
	suite.Setup(t, 16) // Page size = 16 bytes
	defer suite.TearDown()

	valueToWrite := uint8(0x12)
	if err := suite.serializationUtil.WriteUint8(valueToWrite, suite.byteStream); err != nil {
		t.Fatalf("Failed to write uint8: %v", err)
	}

	valueRead, err := suite.serializationUtil.ReadUint8(suite.byteStream)
	if err != nil {
		t.Fatalf("Failed to read uint8: %v", err)
	}

	if valueToWrite != valueRead {
		t.Errorf("uint8 mismatch. Expected: %d, Got: %d", valueToWrite, valueRead)
	}
}

// TestWriteReadUint16 verifies writing and reading a uint16 value.
func TestWriteReadUint16(t *testing.T) {
	suite := &SerializationUtilTestSuite{}
	suite.Setup(t, 16) // Page size = 16 bytes
	defer suite.TearDown()

	valueToWrite := uint16(0x1234)
	if err := suite.serializationUtil.WriteUint16(valueToWrite, suite.byteStream); err != nil {
		t.Fatalf("Failed to write uint16: %v", err)
	}

	valueRead, err := suite.serializationUtil.ReadUint16(suite.byteStream)
	if err != nil {
		t.Fatalf("Failed to read uint16: %v", err)
	}

	if valueToWrite != valueRead {
		t.Errorf("uint16 mismatch. Expected: %d, Got: %d", valueToWrite, valueRead)
	}
}

// TestWriteReadUint32 verifies writing and reading a uint32 value.
func TestWriteReadUint32(t *testing.T) {
	suite := &SerializationUtilTestSuite{}
	suite.Setup(t, 16) // Page size = 16 bytes
	defer suite.TearDown()

	valueToWrite := uint32(0x12345678)
	if err := suite.serializationUtil.WriteUint32(valueToWrite, suite.byteStream); err != nil {
		t.Fatalf("Failed to write uint32: %v", err)
	}

	valueRead, err := suite.serializationUtil.ReadUint32(suite.byteStream)
	if err != nil {
		t.Fatalf("Failed to read uint32: %v", err)
	}

	if valueToWrite != valueRead {
		t.Errorf("uint32 mismatch. Expected: %d, Got: %d", valueToWrite, valueRead)
	}
}

// TestWriteReadUint64 verifies writing and reading a uint64 value.
func TestWriteReadUint64(t *testing.T) {
	suite := &SerializationUtilTestSuite{}
	suite.Setup(t, 16) // Page size = 16 bytes
	defer suite.TearDown()

	valueToWrite := uint64(0x123456789ABCDEF0)
	if err := suite.serializationUtil.WriteUint64(valueToWrite, suite.byteStream); err != nil {
		t.Fatalf("Failed to write uint64: %v", err)
	}

	valueRead, err := suite.serializationUtil.ReadUint64(suite.byteStream)
	if err != nil {
		t.Fatalf("Failed to read uint64: %v", err)
	}

	if valueToWrite != valueRead {
		t.Errorf("uint64 mismatch. Expected: %d, Got: %d", valueToWrite, valueRead)
	}
}

// TestWriteReadFloat32 verifies writing and reading a float32 value.
func TestWriteReadFloat32(t *testing.T) {
	suite := &SerializationUtilTestSuite{}
	suite.Setup(t, 16) // Page size = 16 bytes
	defer suite.TearDown()

	valueToWrite := float32(3.14)
	if err := suite.serializationUtil.WriteFloat(valueToWrite, suite.byteStream); err != nil {
		t.Fatalf("Failed to write float32: %v", err)
	}

	valueRead, err := suite.serializationUtil.ReadFloat(suite.byteStream)
	if err != nil {
		t.Fatalf("Failed to read float32: %v", err)
	}

	if valueToWrite != valueRead {
		t.Errorf("float32 mismatch. Expected: %f, Got: %f", valueToWrite, valueRead)
	}
}

// TestWriteReadFloat64 verifies writing and reading a float64 value.
func TestWriteReadFloat64(t *testing.T) {
	suite := &SerializationUtilTestSuite{}
	suite.Setup(t, 16) // Page size = 16 bytes
	defer suite.TearDown()

	valueToWrite := float64(3.141592653589793)
	if err := suite.serializationUtil.WriteDouble(valueToWrite, suite.byteStream); err != nil {
		t.Fatalf("Failed to write float64: %v", err)
	}

	valueRead, err := suite.serializationUtil.ReadDouble(suite.byteStream)
	if err != nil {
		t.Fatalf("Failed to read float64: %v", err)
	}

	if valueToWrite != valueRead {
		t.Errorf("float64 mismatch. Expected: %f, Got: %f", valueToWrite, valueRead)
	}
}

// TestWriteReadString verifies writing and reading a string.
func TestWriteReadString(t *testing.T) {
	suite := &SerializationUtilTestSuite{}
	suite.Setup(t, 16) // Page size = 16 bytes
	defer suite.TearDown()

	valueToWrite := "Hello, World!"
	if err := suite.serializationUtil.WriteString(valueToWrite, suite.byteStream); err != nil {
		t.Fatalf("Failed to write string: %v", err)
	}

	valueRead, err := suite.serializationUtil.ReadString(suite.byteStream)
	if err != nil {
		t.Fatalf("Failed to read string: %v", err)
	}

	if valueToWrite != valueRead {
		t.Errorf("String mismatch. Expected: %s, Got: %s", valueToWrite, valueRead)
	}
}
