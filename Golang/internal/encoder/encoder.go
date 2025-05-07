package encoder

import (
	"Golang/internal/common/base"
	"bytes"
	"errors"
)

// all the different types of encoders

// Encoder is a generic interface for encoding data into a ByteStream.
type Encoder interface {
	Encode(value interface{}, stream *base.ByteStream) error
	Destroy()
}

// NewEncoder is a factory function that returns an appropriate encoder instance based on the provided type and encoding.
func NewEncoder(dataType base.TSDataType, encoding base.TSEncoding) (Encoder, error) {
	if encoding == "PLAIN" {
		// For plain encoding, return a PlainEncoder
		return NewPlainEncoder(dataType), nil
	}

	// Add other encoder types here as needed
	// Example: if encoding == "dictionary" { return NewDictionaryEncoder(dataType), nil }

	return nil, errors.New(string("unknown encoding type: " + encoding))
}

type PlainEncoder struct {
	Datatype string
}

// NewPlainEncoder initializes the PlainEncoder with serialization utilities.
func NewPlainEncoder(datatype base.TSDataType) *PlainEncoder {
	return &PlainEncoder{
		Datatype: "plain",
	}
}

// Encode encodes a value (of supported types) and writes it into the ByteStream.
func (pe *PlainEncoder) Encode(value interface{}, stream *base.ByteStream) error {
	su := base.SerializationUtil{}
	switch v := value.(type) {
	case bool:
		return su.WriteBool(v, stream)
	case int32:
		return su.WriteVarUint(uint32(v), stream)
	case int64:
		return su.WriteUint64(uint64(v), stream)
	case float32:
		return su.WriteFloat(v, stream)
	case float64:
		return su.WriteDouble(v, stream)
	default:
		return errors.New("unsupported type for encoding")
	}
}

// Destroy releases resources for the encoder (if any).
func (pe *PlainEncoder) Destroy() {
	// No resources to release for PlainEncoder
}

// Decode reads data from the buffer and converts it back to the original value.
func (pe *PlainEncoder) Decode(dataType string, buffer *bytes.Buffer) (interface{}, error) {
	// TODO: implement decoder
	/*
		switch dataType {
		case "uint8":
			return pe.Util.ReadUI8(buffer)
		case "uint16":
			return pe.Util.ReadUI16(buffer)
		case "uint32":
			return pe.Util.ReadUI32(buffer)
		case "uint64":
			return pe.Util.ReadUI64(buffer)
		default:
			return nil, errors.New("unsupported data type for decoding")
		}
	*/
	return nil, errors.New("unsupported data type for decoding")
}
