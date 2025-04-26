package writertest

import (
	"Golang/internal/writer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageWriter_WriteBooleanSuccess(t *testing.T) {
	pageWriter := NewPageWriter()
	pageWriter.Init(BOOLEAN, PLAIN, UNCOMPRESSED)

	result := pageWriter.Write(1234567890, true)

	assert.Equal(t, E_OK, result)
}

func TestPageWriter_WriteInt32Success(t *testing.T) {
	pageWriter := NewPageWriter()
	pageWriter.Init(INT32, PLAIN, UNCOMPRESSED)

	result := pageWriter.Write(1234567890, int32(42))

	assert.Equal(t, E_OK, result)
}

func TestPageWriter_WriteInt64Success(t *testing.T) {
	pageWriter := NewPageWriter()
	pageWriter.Init(INT64, PLAIN, UNCOMPRESSED)

	result := pageWriter.Write(1234567890, int64(42))

	assert.Equal(t, E_OK, result)
}

func TestPageWriter_WriteFloatSuccess(t *testing.T) {
	pageWriter := NewPageWriter()
	pageWriter.Init(FLOAT, PLAIN, UNCOMPRESSED)

	result := pageWriter.Write(1234567890, float32(42.0))

	assert.Equal(t, E_OK, result)
}

func TestPageWriter_WriteDoubleSuccess(t *testing.T) {
	pageWriter := NewPageWriter()
	pageWriter.Init(DOUBLE, PLAIN, UNCOMPRESSED)

	result := pageWriter.Write(1234567890, float64(42.0))

	assert.Equal(t, E_OK, result)
}

func TestPageWriter_WriteLargeDataSet(t *testing.T) {
	pageWriter := NewPageWriter()
	pageWriter.Init(DOUBLE, PLAIN, UNCOMPRESSED)

	for i := 0; i < 10000; i++ {
		pageWriter.Write(i, float64(i)*0.1)
	}

	assert.Equal(t, 10000, pageWriter.GetPointNumber())
}

func TestPageWriter_ResetPageWriter(t *testing.T) {
	pageWriter := NewPageWriter()
	pageWriter.Init(INT64, PLAIN, UNCOMPRESSED)
	pageWriter.Write(1234567890, int64(42))

	pageWriter.Reset()

	assert.Equal(t, 0, pageWriter.GetPointNumber())
	assert.Equal(t, 0, pageWriter.GetTimeOutStreamSize())
}

func TestPageWriter_DestroyPageWriter(t *testing.T) {
	pageWriter := NewPageWriter()
	pageWriter.Init(INT64, PLAIN, UNCOMPRESSED)
	pageWriter.Write(1234567890, int64(42))

	stat := pageWriter.GetStatistic().(*Int64Statistic)

	assert.NotNil(t, stat)
	assert.Equal(t, int64(1), stat.Count)
}
