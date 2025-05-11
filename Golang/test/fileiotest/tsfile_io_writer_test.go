package fileiotest

import (
	"Golang/internal/common/base"
	"Golang/internal/common/core"
	"Golang/internal/fileio"
	"Golang/internal/utils"
	"reflect"
	"testing"
)

func TestFileIndexWritingMemManager_AddIndexNode(t *testing.T) {
	type fields struct {
		ByteStream    *base.ByteStream
		AllIndexNodes []*core.MetaIndexNode
	}
	type args struct {
		node *core.MetaIndexNode
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &fileio.FileIndexWritingMemManager{
				ByteStream:    tt.fields.ByteStream,
				AllIndexNodes: tt.fields.AllIndexNodes,
			}
			m.AddIndexNode(tt.args.node)
		})
	}
}

func TestFileIndexWritingMemManager_Free(t *testing.T) {
	type fields struct {
		ByteStream    *base.ByteStream
		AllIndexNodes []*core.MetaIndexNode
	}
	var tests []struct {
		name   string
		fields fields
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &fileio.FileIndexWritingMemManager{
				ByteStream:    tt.fields.ByteStream,
				AllIndexNodes: tt.fields.AllIndexNodes,
			}
			m.Free()
		})
	}
}

func TestNewFileIndexWritingMemManager(t *testing.T) {
	type args struct {
		pageSize uint32
	}
	var tests []struct {
		name    string
		args    args
		want    *fileio.FileIndexWritingMemManager
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fileio.NewFileIndexWritingMemManager(tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileIndexWritingMemManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileIndexWritingMemManager() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTsFileIOWriter(t *testing.T) {
	var tests []struct {
		name string
		want *fileio.TsFileIOWriter
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileio.NewTsFileIOWriter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTsFileIOWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsFileIOWriter_AddCurIndexNodeToQueue(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		node  *core.MetaIndexNode
		queue *[]*core.MetaIndexNode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.AddCurIndexNodeToQueue(tt.args.node, tt.args.queue); (err != nil) != tt.wantErr {
				t.Errorf("AddCurIndexNodeToQueue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_AddDeviceNode(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		deviceMap                 map[string]*core.MetaIndexNode
		deviceName                string
		measurementIndexNodeQueue []*core.MetaIndexNode
		wmm                       *fileio.FileIndexWritingMemManager
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.AddDeviceNode(tt.args.deviceMap, tt.args.deviceName, tt.args.measurementIndexNodeQueue, tt.args.wmm); (err != nil) != tt.wantErr {
				t.Errorf("AddDeviceNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_AllocAndInitMetaIndexEntry(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		wmm  *fileio.FileIndexWritingMemManager
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *core.MetaIndexEntry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			got, err := io.AllocAndInitMetaIndexEntry(tt.args.wmm, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllocAndInitMetaIndexEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllocAndInitMetaIndexEntry() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsFileIOWriter_AllocAndInitMetaIndexNode(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		wmm      *fileio.FileIndexWritingMemManager
		nodeType core.MetaIndexNodeType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *core.MetaIndexNode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			got, err := io.AllocAndInitMetaIndexNode(tt.args.wmm, tt.args.nodeType)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllocAndInitMetaIndexNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllocAndInitMetaIndexNode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsFileIOWriter_AllocMetaIndexNodeQueue(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		wmm *fileio.FileIndexWritingMemManager
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*core.MetaIndexNode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			got, err := io.AllocMetaIndexNodeQueue(tt.args.wmm)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllocMetaIndexNodeQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllocMetaIndexNodeQueue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsFileIOWriter_BuildDeviceLevel(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		deviceMap map[string]*core.MetaIndexNode
		retRoot   **core.MetaIndexNode
		wmm       *fileio.FileIndexWritingMemManager
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.BuildDeviceLevel(tt.args.deviceMap, tt.args.retRoot, tt.args.wmm); (err != nil) != tt.wantErr {
				t.Errorf("BuildDeviceLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_CloneNodeList(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		src []*core.MetaIndexNode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*core.MetaIndexNode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			got, err := io.CloneNodeList(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloneNodeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloneNodeList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsFileIOWriter_CurrentFilePosition(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if got := io.CurrentFilePosition(); got != tt.want {
				t.Errorf("CurrentFilePosition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsFileIOWriter_Destroy(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			io.Destroy()
		})
	}
}

func TestTsFileIOWriter_EndFile(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.EndFile(); (err != nil) != tt.wantErr {
				t.Errorf("EndFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_EndFlushChunk(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.EndFlushChunk(); (err != nil) != tt.wantErr {
				t.Errorf("EndFlushChunk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_EndFlushChunkGroup(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		isAligned bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.EndFlushChunkGroup(tt.args.isAligned); (err != nil) != tt.wantErr {
				t.Errorf("EndFlushChunkGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_FlushChunk(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		chunkData *base.ByteStream
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.FlushChunk(tt.args.chunkData); (err != nil) != tt.wantErr {
				t.Errorf("FlushChunk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_FlushStreamToFile(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.FlushStreamToFile(); (err != nil) != tt.wantErr {
				t.Errorf("FlushStreamToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_GenerateRoot(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		nodeQueue []*core.MetaIndexNode
		nodeType  core.MetaIndexNodeType
		wmm       *fileio.FileIndexWritingMemManager
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *core.MetaIndexNode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			got, err := io.GenerateRoot(tt.args.nodeQueue, tt.args.nodeType, tt.args.wmm)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateRoot() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsFileIOWriter_GetFilePath(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if got := io.GetFilePath(); got != tt.want {
				t.Errorf("GetFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsFileIOWriter_GetPathCount(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		cgmList []*core.ChunkGroupMeta
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if got := io.GetPathCount(tt.args.cgmList); got != tt.want {
				t.Errorf("GetPathCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsFileIOWriter_Init(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		writeFile *fileio.WriteFile
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.Init(tt.args.writeFile); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_StartFile(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.StartFile(); (err != nil) != tt.wantErr {
				t.Errorf("StartFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_StartFlushChunk(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		chunkData       *base.ByteStream
		measurementName string
		dataType        base.TSDataType
		encoding        base.TSEncoding
		compression     base.CompressionType
		numOfPages      int32
		TsID            utils.TsID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.StartFlushChunk(tt.args.chunkData, tt.args.measurementName, tt.args.dataType, tt.args.encoding, tt.args.compression, tt.args.numOfPages, tt.args.TsID); (err != nil) != tt.wantErr {
				t.Errorf("StartFlushChunk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_StartFlushChunkGroup(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		deviceName string
		isAligned  bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.StartFlushChunkGroup(tt.args.deviceName, tt.args.isAligned); (err != nil) != tt.wantErr {
				t.Errorf("StartFlushChunkGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_WriteBuf(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		tsfile    string
		tsfileLen int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.WriteBuf(tt.args.tsfile, tt.args.tsfileLen); (err != nil) != tt.wantErr {
				t.Errorf("WriteBuf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_WriteByte(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		written byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.WriteByte(tt.args.written); (err != nil) != tt.wantErr {
				t.Errorf("WriteByte() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_WriteChunkData(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		chunkData *base.ByteStream
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.WriteChunkData(tt.args.chunkData); (err != nil) != tt.wantErr {
				t.Errorf("WriteChunkData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_WriteFileFooter(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.WriteFileFooter(); (err != nil) != tt.wantErr {
				t.Errorf("WriteFileFooter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_WriteFileIndex(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.WriteFileIndex(); (err != nil) != tt.wantErr {
				t.Errorf("WriteFileIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_WriteLogIndexRange(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.WriteLogIndexRange(); (err != nil) != tt.wantErr {
				t.Errorf("WriteLogIndexRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_WriteSeperatorMarker(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		tsfIleMeta *core.TsFileMeta
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.WriteSeperatorMarker(tt.args.tsfIleMeta); (err != nil) != tt.wantErr {
				t.Errorf("WriteSeperatorMarker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_WriteString(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if err := io.WriteString(tt.args.str); (err != nil) != tt.wantErr {
				t.Errorf("WriteString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTsFileIOWriter_curFilePosition(t *testing.T) {
	type fields struct {
		writeStream         *base.ByteStream
		writeStreamConsumer *base.ByteStreamConsumer
		curChunkMeta        *core.ChunkMeta
		curChunkGroupMeta   *core.ChunkGroupMeta
		chunkMetaCount      int
		chunkGroupMetaList  []*core.ChunkGroupMeta
		usePrevAllocCgm     bool
		curDeviceName       string
		file                *fileio.WriteFile
		tsTimeIndexVector   []core.TimeseriesTimeIndexEntry
		writeFileCreated    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io := &fileio.TsFileIOWriter{
				writeStream:         tt.fields.writeStream,
				writeStreamConsumer: tt.fields.writeStreamConsumer,
				curChunkMeta:        tt.fields.curChunkMeta,
				curChunkGroupMeta:   tt.fields.curChunkGroupMeta,
				chunkMetaCount:      tt.fields.chunkMetaCount,
				chunkGroupMetaList:  tt.fields.chunkGroupMetaList,
				usePrevAllocCgm:     tt.fields.usePrevAllocCgm,
				curDeviceName:       tt.fields.curDeviceName,
				file:                tt.fields.file,
				tsTimeIndexVector:   tt.fields.tsTimeIndexVector,
				writeFileCreated:    tt.fields.writeFileCreated,
			}
			if got := io.curFilePosition(); got != tt.want {
				t.Errorf("curFilePosition() = %v, want %v", got, tt.want)
			}
		})
	}
}
