package writertest

import (
	"Golang/internal/common/base"
	"Golang/internal/common/statistic"
	"Golang/internal/writer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPageWriter_WriteBooleanSuccess(t *testing.T) {
	pageWriter := &writer.PageWriter{}
	err := pageWriter.Initialize(base.BOOLEAN, base.PLAIN, base.UNCOMPRESSED)
	assert.NoError(t, err)

	err = pageWriter.Write(1234567890, true)
	assert.NoError(t, err)
}

func TestPageWriter_WriteInt32Success(t *testing.T) {
	pageWriter := &writer.PageWriter{}
	err := pageWriter.Initialize(base.INT32, base.PLAIN, base.UNCOMPRESSED)
	assert.NoError(t, err)

	err = pageWriter.Write(1234567890, int32(42))
	assert.NoError(t, err)
}

func TestPageWriter_WriteInt64Success(t *testing.T) {
	pageWriter := &writer.PageWriter{}
	err := pageWriter.Initialize(base.INT64, base.PLAIN, base.UNCOMPRESSED)
	assert.NoError(t, err)

	err = pageWriter.Write(1234567890, int64(42))
	assert.NoError(t, err)
}

func TestPageWriter_WriteFloatSuccess(t *testing.T) {
	pageWriter := &writer.PageWriter{}
	err := pageWriter.Initialize(base.FLOAT, base.PLAIN, base.UNCOMPRESSED)
	assert.NoError(t, err)

	err = pageWriter.Write(1234567890, float32(42.0))
	assert.NoError(t, err)
}

func TestPageWriter_WriteDoubleSuccess(t *testing.T) {
	pageWriter := &writer.PageWriter{}
	err := pageWriter.Initialize(base.DOUBLE, base.PLAIN, base.UNCOMPRESSED)
	assert.NoError(t, err)

	err = pageWriter.Write(1234567890, float64(42.0))
	assert.NoError(t, err)
}

func TestPageWriter_WriteLargeDataSet(t *testing.T) {
	pageWriter := &writer.PageWriter{}
	err := pageWriter.Initialize(base.DOUBLE, base.PLAIN, base.UNCOMPRESSED)
	assert.NoError(t, err)

	for i := int64(0); i < 10000; i++ {
		err = pageWriter.Write(i, float64(i)*0.1)
		if err != nil {
			t.Fatalf("error writing data at index %d: %v", i, err)
		}
	}

	assert.Equal(t, 10000, pageWriter.GetPointCount())
}

func TestPageWriter_ResetPageWriter(t *testing.T) {

	pageWriter := &writer.PageWriter{}
	err := pageWriter.Initialize(base.INT64, base.PLAIN, base.UNCOMPRESSED)
	assert.NoError(t, err)

	err = pageWriter.Write(1234567890, int64(42))
	assert.NoError(t, err)
	// stat := pageWriter.Statistic.(*statistic.Int64Statistic)
	// t.Logf("stat: %v", stat)

	err = pageWriter.Write(1234567891, int64(40))
	assert.NoError(t, err)
	// t.Logf("stat: %v", stat)
	pageWriter.Reset()
	// t.Logf("stat: %v", stat)

	assert.Equal(t, 0, pageWriter.GetPointCount())
	assert.Equal(t, uint32(0), pageWriter.TimeOutStream.TotalSize)
}

func TestPageWriter_DestroyPageWriter(t *testing.T) {
	pageWriter := writer.PageWriter{}
	err := pageWriter.Initialize(base.INT64, base.PLAIN, base.UNCOMPRESSED)
	assert.NoError(t, err)

	err = pageWriter.Write(1234567890, int64(42))
	assert.NoError(t, err)

	stat := pageWriter.Statistic.(*statistic.Int64Statistic)
	t.Logf("stat: %v", stat)

	assert.NotNil(t, stat)
	assert.Equal(t, int64(1), stat.Count)
}
