package base

import (
	"errors"
	"fmt"
	"math"
	"unsafe"
)

// SerializationUtil provides utility methods for serializing and deserializing binary data
// in a ByteStream (big-endian encoding is used).
// TODO consider Stack Allocation due to GO garbage collector impact
type SerializationUtil struct{}

// WriteBool writes a boolean value as uint8 (1 for true, 0 for false) to the ByteStream
func (su *SerializationUtil) WriteBool(value bool, out *ByteStream) error {
	uint8Value := uint8(0)
	if value {
		uint8Value = 1
	}
	return su.WriteUint8(uint8Value, out)
}

// ReadBool reads a boolean value encoded as a uint8 (1 for true, 0 for false) from the ByteStream.
func (su *SerializationUtil) ReadBool(in *ByteStream) (bool, error) {
	// Create a buffer to read 1 byte
	buf := make([]byte, 1)

	// Use `ReadBuf` to read into the buffer
	_, err := in.ReadBuf(buf, 1) // Pass buffer and length
	if err != nil {
		return false, err
	}

	// Extract the first byte and interpret as boolean
	switch buf[0] {
	case 1:
		return true, nil
	case 0:
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %v", buf[0])
	}
}

// WriteUint8 writes a single uint8 value to the ByteStream.
func (su *SerializationUtil) WriteUint8(ui8 uint8, out *ByteStream) error {
	return out.WriteBuf([]byte{ui8}, 1)
}

// WriteUint16 writes a single uint16 value (big-endian) to the ByteStream.
func (su *SerializationUtil) WriteUint16(ui16 uint16, out *ByteStream) error {
	buf := []byte{
		uint8((ui16 >> 8) & 0xFF),
		uint8(ui16 & 0xFF),
	}
	return out.WriteBuf(buf, 2)
}

// WriteUint32 writes a single uint32 value (big-endian) to the ByteStream.
func (su *SerializationUtil) WriteUint32(ui32 uint32, out *ByteStream) error {
	buf := []byte{
		uint8((ui32 >> 24) & 0xFF),
		uint8((ui32 >> 16) & 0xFF),
		uint8((ui32 >> 8) & 0xFF),
		uint8(ui32 & 0xFF),
	}
	return out.WriteBuf(buf, 4)
}

// WriteUint64 writes a single uint64 value (big-endian) to the ByteStream.
func (su *SerializationUtil) WriteUint64(ui64 uint64, out *ByteStream) error {
	buf := []byte{
		uint8((ui64 >> 56) & 0xFF),
		uint8((ui64 >> 48) & 0xFF),
		uint8((ui64 >> 40) & 0xFF),
		uint8((ui64 >> 32) & 0xFF),
		uint8((ui64 >> 24) & 0xFF),
		uint8((ui64 >> 16) & 0xFF),
		uint8((ui64 >> 8) & 0xFF),
		uint8(ui64 & 0xFF),
	}
	return out.WriteBuf(buf, 8)
}

// ReadUint8 reads a single uint8 value from the ByteStream.
func (su *SerializationUtil) ReadUint8(in *ByteStream) (uint8, error) {
	buf := make([]byte, 1)
	_, err := in.ReadBuf(buf, 1)
	return buf[0], err
}

// ReadUint16 reads a single uint16 value (big-endian) from the ByteStream.
func (su *SerializationUtil) ReadUint16(in *ByteStream) (uint16, error) {
	buf := make([]byte, 2)
	_, err := in.ReadBuf(buf, 2)
	if err != nil {
		return 0, err
	}
	return uint16(buf[0])<<8 | uint16(buf[1]), nil
}

// ReadUint32 reads a single uint32 value (big-endian) from the ByteStream.
func (su *SerializationUtil) ReadUint32(in *ByteStream) (uint32, error) {
	buf := make([]byte, 4)
	_, err := in.ReadBuf(buf, 4)
	if err != nil {
		return 0, err
	}
	return uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3]), nil
}

// ReadUint64 reads a single uint64 value (big-endian) from the ByteStream.
func (su *SerializationUtil) ReadUint64(in *ByteStream) (uint64, error) {
	buf := make([]byte, 8)
	_, err := in.ReadBuf(buf, 8)
	if err != nil {
		return 0, err
	}
	return uint64(buf[0])<<56 | uint64(buf[1])<<48 | uint64(buf[2])<<40 | uint64(buf[3])<<32 |
		uint64(buf[4])<<24 | uint64(buf[5])<<16 | uint64(buf[6])<<8 | uint64(buf[7]), nil
}

// WriteFloat writes a single float32 value to the ByteStream.
func (su *SerializationUtil) WriteFloat(f float32, out *ByteStream) error {
	bytes := su.FloatToBytes(f)
	return out.WriteBuf(bytes[:], 4)
}

// ReadFloat reads a single float32 value from the ByteStream.
func (su *SerializationUtil) ReadFloat(in *ByteStream) (float32, error) {
	buf := make([]byte, 4)
	_, err := in.ReadBuf(buf, 4)
	if err != nil {
		return 0, err
	}
	var fixedBuf [4]byte                  // Create a fixed-size array
	copy(fixedBuf[:], buf)                // Copy slice data into the fixed-size array
	return su.BytesToFloat(fixedBuf), nil // Pass the fixed-size array

}

// WriteDouble writes a single float64 (double) value to the ByteStream.
func (su *SerializationUtil) WriteDouble(d float64, out *ByteStream) error {
	bytes := su.DoubleToBytes(d)
	return out.WriteBuf(bytes[:], 8)
}

// ReadDouble reads a single float64 (double) value from the ByteStream.
func (su *SerializationUtil) ReadDouble(in *ByteStream) (float64, error) {
	buf := make([]byte, 8)
	_, err := in.ReadBuf(buf, 8)
	if err != nil {
		return 0, err
	}
	var fixedBuf [8]byte                   // Create a fixed-size array
	copy(fixedBuf[:], buf)                 // Copy slice data into the fixed-size array
	return su.BytesToDouble(fixedBuf), nil // Pass the fixed-size array

}

// WriteVarUint writes a variable-length uint32 value to the ByteStream.
func (su *SerializationUtil) WriteVarUint(ui32 uint32, out *ByteStream) error {
	for (ui32 & 0xFFFFFF80) != 0 {
		err := su.WriteUint8(uint8((ui32&0x7F)|0x80), out)
		if err != nil {
			return err
		}
		ui32 >>= 7
	}
	return su.WriteUint8(uint8(ui32&0x7F), out)
}

// ReadVarUint reads a variable-length uint32 value from the ByteStream.
func (su *SerializationUtil) ReadVarUint(in *ByteStream) (uint32, error) {
	ui32 := uint32(0)
	i := 0
	for {
		byteVal, err := su.ReadUint8(in)
		if err != nil {
			return 0, errors.New("readUint8 failed: " + err.Error())
		}
		ui32 |= uint32(byteVal&0x7F) << (i * 7)
		if (byteVal & 0x80) == 0 {
			break
		}
		i++
	}
	return ui32, nil
}

// WriteString writes a variable-length string to the ByteStream.
func (su *SerializationUtil) WriteString(str string, out *ByteStream) error {
	err := su.WriteVarUint(uint32(len(str)), out)
	if err != nil {
		return err
	}
	return out.WriteBuf([]byte(str), uint32(len(str)))
}

// ReadString reads a string from the ByteStream.
func (su *SerializationUtil) ReadString(in *ByteStream) (string, error) {
	length, err := su.ReadVarUint(in)
	if err != nil {
		return "", errors.New("readVarUint failed: " + err.Error())
	}
	buf := make([]byte, length)
	_, err = in.ReadBuf(buf, length)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

/////////////////////////////////////////
// Conversion from one type to another //
/////////////////////////////////////////

// FloatToInt converts a float32 to its binary representation as int32.
func (su *SerializationUtil) FloatToInt(f float32) int32 {
	return *(*int32)(unsafe.Pointer(&f))
}

// IntToFloat converts an int32 to its binary representation as float32.
func (su *SerializationUtil) IntToFloat(i int32) float32 {
	return *(*float32)(unsafe.Pointer(&i))
}

// FloatToBytes converts a float32 and returns its IEEE 754 big-endian byte representation.
func (su *SerializationUtil) FloatToBytes(f float32) [4]byte {
	if math.IsNaN(float64(f)) {
		// NaN representation: IEEE 754 0x7FC00000
		return [4]byte{0x7F, 0xC0, 0x00, 0x00}
	}

	intBits := su.FloatToInt(f) // Use low-level float-to-int conversion
	return [4]byte{
		byte(intBits >> 24),
		byte(intBits >> 16),
		byte(intBits >> 8),
		byte(intBits),
	}
}

// BytesToFloat converts a big-endian byte array to a float32.
func (su *SerializationUtil) BytesToFloat(bytes [4]byte) float32 {
	intBits := int32(bytes[0])<<24 | int32(bytes[1])<<16 | int32(bytes[2])<<8 | int32(bytes[3])
	return su.IntToFloat(intBits)
}

// DoubleToInt converts a float64 (double) to its binary representation as int64.
func (su *SerializationUtil) DoubleToInt(d float64) int64 {
	return *(*int64)(unsafe.Pointer(&d))
}

// IntToDouble converts an int64 to its binary representation as float64 (double).
func (su *SerializationUtil) IntToDouble(i int64) float64 {
	return *(*float64)(unsafe.Pointer(&i))
}

// DoubleToBytes converts a float64 (double) to its IEEE 754 big-endian byte representation.
func (su *SerializationUtil) DoubleToBytes(d float64) [8]byte {
	if math.IsNaN(d) {
		// NaN representation: IEEE 754 0x7FF8000000000000L
		return [8]byte{0x7F, 0xF8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	}

	intBits := su.DoubleToInt(d) // Use low-level double-to-int conversion
	return [8]byte{
		byte(intBits >> 56),
		byte(intBits >> 48),
		byte(intBits >> 40),
		byte(intBits >> 32),
		byte(intBits >> 24),
		byte(intBits >> 16),
		byte(intBits >> 8),
		byte(intBits),
	}
}

// BytesToDouble converts a big-endian byte array to a float64 (double).
func (su *SerializationUtil) BytesToDouble(bytes [8]byte) float64 {
	intBits := int64(bytes[0])<<56 | int64(bytes[1])<<48 | int64(bytes[2])<<40 | int64(bytes[3])<<32 |
		int64(bytes[4])<<24 | int64(bytes[5])<<16 | int64(bytes[6])<<8 | int64(bytes[7])
	return su.IntToDouble(intBits)
}
