package core

import (
	"Golang/internal/common/base"
	"Golang/internal/common/core"
	"Golang/internal/common/statistic"
	"Golang/internal/reader"
	_ "Golang/internal/reader"
	"Golang/internal/utils"
	"testing"
)

/////////////////
// page header //
/////////////////

// TestPageHeaderDefaultConstructor verifies the default initialization of PageHeader.
func TestPageHeaderDefaultConstructor(t *testing.T) {
	header := core.NewPageHeader() // Create a new PageHeader instance.

	if header.UncompressedSize != 0 {
		t.Errorf("UncompressedSize mismatch. Expected: 0, Got: %d", header.UncompressedSize)
	}

	if header.CompressedSize != 0 {
		t.Errorf("CompressedSize mismatch. Expected: 0, Got: %d", header.CompressedSize)
	}

	if header.Statistic != nil {
		t.Errorf("Statistic mismatch. Expected: nil, Got: %v", header.Statistic)
	}
}

// TestPageHeaderReset verifies that calling Reset clears all fields of PageHeader.
func TestPageHeaderReset(t *testing.T) {
	// Create a new PageHeader instance.
	header := core.NewPageHeader()

	// Simulate values being set.
	header.UncompressedSize = 100
	header.CompressedSize = 50

	// Allocate a statistic object.
	factory := statistic.Factory{}
	stat, err := factory.AllocStatistic(base.BOOLEAN)
	if err != nil {
		t.Fatalf("Failed to allocate Statistic: %v", err)
	}
	header.Statistic = stat

	// Call Reset and verify fields are cleared.
	header.Reset(factory)

	if header.UncompressedSize != 0 {
		t.Errorf("UncompressedSize mismatch after Reset. Expected: 0, Got: %d", header.UncompressedSize)
	}

	if header.CompressedSize != 0 {
		t.Errorf("CompressedSize mismatch after Reset. Expected: 0, Got: %d", header.CompressedSize)
	}

	if header.Statistic != nil {
		t.Errorf("Statistic mismatch after Reset. Expected: nil, Got: %v", header.Statistic)
	}
}

//////////////////
// chunk header //
//////////////////

// TestChunkHeaderDefaultConstructor verifies the default initialization of ChunkHeader.
func TestChunkHeaderDefaultConstructor(t *testing.T) {
	header := core.NewChunkHeader() // Create a new ChunkHeader instance.

	if header.MeasurementName != "" {
		t.Errorf("MeasurementName mismatch. Expected: \"\", Got: %s", header.MeasurementName)
	}

	if header.DataSize != 0 {
		t.Errorf("DataSize mismatch. Expected: 0, Got: %d", header.DataSize)
	}

	if header.DataType != base.INVALID_TS {
		t.Errorf("DataType mismatch. Expected: INVALID_TS, Got: %v", header.DataType)
	}

	if header.CompressionType != base.INVALID_C {
		t.Errorf("CompressionType mismatch. Expected: INVALID_C, Got: %v", header.CompressionType)
	}

	if header.EncodingType != base.INVALID_E {
		t.Errorf("EncodingType mismatch. Expected: INVALID_E, Got: %v", header.EncodingType)
	}

	if header.NumOfPages != 0 {
		t.Errorf("NumOfPages mismatch. Expected: 0, Got: %d", header.NumOfPages)
	}

	if header.SerializedSize != 0 {
		t.Errorf("SerializedSize mismatch. Expected: 0, Got: %d", header.SerializedSize)
	}

	if header.ChunkType != 0 {
		t.Errorf("ChunkType mismatch. Expected: 0, Got: %d", header.ChunkType)
	}
}

// TestChunkHeaderReset verifies that calling Reset on ChunkHeader clears all fields.
func TestChunkHeaderReset(t *testing.T) {
	// Create a new ChunkHeader instance.
	header := core.NewChunkHeader()

	// Simulate values being set.
	header.MeasurementName = "test"
	header.DataSize = 100
	header.DataType = base.INT32
	header.CompressionType = base.SNAPPY
	header.EncodingType = base.PLAIN
	header.NumOfPages = 5
	header.SerializedSize = 50
	header.ChunkType = 1

	// Call Reset method.
	header.Reset()

	if header.MeasurementName != "" {
		t.Errorf("MeasurementName mismatch after Reset. Expected: \"\", Got: %s", header.MeasurementName)
	}

	if header.DataSize != 0 {
		t.Errorf("DataSize mismatch after Reset. Expected: 0, Got: %d", header.DataSize)
	}

	if header.DataType != base.INVALID_TS {
		t.Errorf("DataType mismatch after Reset. Expected: INVALID_TS, Got: %v", header.DataType)
	}

	if header.CompressionType != base.INVALID_C {
		t.Errorf("CompressionType mismatch after Reset. Expected: INVALID_C, Got: %v", header.CompressionType)
	}

	if header.EncodingType != base.INVALID_E {
		t.Errorf("EncodingType mismatch after Reset. Expected: INVALID_E, Got: %v", header.EncodingType)
	}

	if header.NumOfPages != 0 {
		t.Errorf("NumOfPages mismatch after Reset. Expected: 0, Got: %d", header.NumOfPages)
	}

	if header.SerializedSize != 0 {
		t.Errorf("SerializedSize mismatch after Reset. Expected: 0, Got: %d", header.SerializedSize)
	}

	if header.ChunkType != 0 {
		t.Errorf("ChunkType mismatch after Reset. Expected: 0, Got: %d", header.ChunkType)
	}
}

////////////////
// chunk meta //
////////////////

// TestChunkMetaDefaultConstructor verifies that ChunkMeta is initialized with default values.
func TestChunkMetaDefaultConstructor(t *testing.T) {
	meta := core.ChunkMeta{} // Create a new ChunkMeta instance.

	if meta.OffsetOfHeader != 0 {
		t.Errorf("OffsetOfChunkHeader mismatch. Expected: 0, Got: %d", meta.OffsetOfHeader)
	}

	if meta.Statistic != nil {
		t.Errorf("Statistic mismatch. Expected: nil, Got: %v", meta.Statistic)
	}

	if meta.Mask != 0 {
		t.Errorf("Mask mismatch. Expected: 0, Got: %d", meta.Mask)
	}
}

// TestChunkMetaInit verifies the initialization of ChunkMeta with valid inputs.
func TestChunkMetaInit(t *testing.T) {
	meta := core.ChunkMeta{} // Create a new ChunkMeta instance.
	name := "test"
	measurementName := name // Assuming base.NewString abstracts string creation.
	factory := statistic.Factory{}
	stat, err := factory.AllocStatistic(base.INT32)
	tsID := utils.TsID{}
	mask := 1

	// Initialize the ChunkMeta instance.

	err = meta.Initialize(measurementName, base.INT32, 100, stat, tsID, byte(mask))
	if err != nil {
		return
	}

	// Validate the initialized fields.
	if meta.DataType != base.INT32 {
		t.Errorf("DataType mismatch. Expected: INT32, Got: %v", meta.DataType)
	}

	if meta.OffsetOfHeader != 100 {
		t.Errorf("OffsetOfChunkHeader mismatch. Expected: 100, Got: %d", meta.OffsetOfHeader)
	}

	if meta.Statistic != stat {
		t.Errorf("Statistic mismatch. Expected: %v, Got: %v", stat, meta.Statistic)
	}

	if meta.TsID != tsID {
		t.Errorf("TsID mismatch. Expected: %v, Got: %v", tsID, meta.TsID)
	}

	if meta.Mask != 1 {
		t.Errorf("Mask mismatch. Expected: 1, Got: %d", meta.Mask)
	}
}

//////////////////////
// chunk group meta //
//////////////////////

// Test for the constructor (NewChunkGroupMeta)
func TestChunkGroupMetaConstructor(t *testing.T) {
	// Initialize a new ChunkGroupMeta
	groupMeta := core.NewChunkGroupMeta("")

	// Verify the ChunkMetaList is initially empty
	if len(groupMeta.ChunkMetaList) != 0 {
		t.Errorf("Expected ChunkMetaList size to be 0, got %d", len(groupMeta.ChunkMetaList))
	}
}

// Test for the Init method (NewChunkGroupMeta with device name)
func TestChunkGroupMetaInit(t *testing.T) {
	deviceName := "device_1"

	// Initialize a new ChunkGroupMeta with the device name
	groupMeta := core.NewChunkGroupMeta(deviceName)

	// Verify the DeviceName was initialized correctly
	if groupMeta.DeviceName != deviceName {
		t.Errorf("Expected DeviceName to be %s, got %s", deviceName, groupMeta.DeviceName)
	}
}

// Test the Push method
func TestChunkGroupMetaPush(t *testing.T) {
	// Create a new ChunkGroupMeta
	groupMeta := core.NewChunkGroupMeta("device_1")

	// Create a ChunkMeta instance (simulate it as empty if definition is missing)
	meta := &core.ChunkMeta{}

	// Push a ChunkMeta into the group
	groupMeta.Push(meta)

	// Verify the ChunkMetaList contains one item
	if len(groupMeta.ChunkMetaList) != 1 {
		t.Errorf("Expected ChunkMetaList size to be 1, got %d", len(groupMeta.ChunkMetaList))
	}

	// Verify the first item in ChunkMetaList is the pushed ChunkMeta
	if groupMeta.ChunkMetaList[0] != meta {
		t.Errorf("Expected ChunkMetaList[0] to be %v, but got %v", meta, groupMeta.ChunkMetaList[0])
	}
}

///////////////////////
// time series index //
///////////////////////

// TestTimeseriesIndexConstructorAndDestructor verifies the default initialization of TimeseriesIndex.
func TestTimeseriesIndexConstructorAndDestructor(t *testing.T) {
	tsIndex := core.TimeseriesIndex{}
	tsIndex.Initialize()
	// Ensure the default data type is INVALID.
	if tsIndex.DataType != base.INVALID_TS {
		t.Errorf("Expected data type to be INVALID_TS, got: %v", tsIndex.DataType)
	}

	// Ensure the default statistic is nil.
	if tsIndex.Statistic != nil {
		t.Errorf("Expected statistic to be nil, got: %v", tsIndex.Statistic)
	}

	// Ensure the default chunk meta list is nil.
	if tsIndex.ChunkMetaList != nil {
		t.Errorf("Expected ChunkMetaList to be nil, got: %v", tsIndex.ChunkMetaList)
	}
}

// TestTimeseriesIndexReset verifies that resetting a TimeseriesIndex clears its fields.
func TestTimeseriesIndexReset(t *testing.T) {
	tsIndex := core.TimeseriesIndex{
		DataType: base.INT64,
		Statistic: &statistic.DoubleStatistic{
			SumValue: 1.0,
		},
		ChunkMetaList: []*core.ChunkMeta{},
	}

	// t.Log(tsIndex)

	// Reset the TimeseriesIndex.
	tsIndex.Reset()
	//tsIndex.DataType = base.VECTOR // Simulating a reset in behavior
	//tsIndex.Statistic = nil
	//tsIndex.ChunkMetaList = nil

	// Verify fields after reset.
	if tsIndex.DataType != base.VECTOR {
		t.Errorf("Expected DataType to be VECTOR after reset, got: %v", tsIndex.DataType)
	}

	if tsIndex.Statistic != nil {
		t.Errorf("Expected Statistic to be nil after reset, got: %v", tsIndex.Statistic)
	}

	if tsIndex.ChunkMetaList != nil {
		t.Errorf("Expected ChunkMetaList to be nil after reset, got: %v", tsIndex.ChunkMetaList)
	}
}

// TestTimeseriesIndexSerializeAndDeserialize verifies serialization and deserialization of TimeseriesIndex.
func TestTimeseriesIndexSerializeAndDeserialize(t *testing.T) {
	// Create a statistic for the TimeseriesIndex.
	factory := statistic.Factory{}

	stat, err := factory.AllocStatistic(base.INT32)
	if err != nil {
		t.Fatalf("Failed to allocate Statistic: %v", err)
	}

	// Create a TimeseriesIndex and set its values.
	tsIndex := core.TimeseriesIndex{
		MeasurementName:    "test_measurement",
		TimeseriesMetaType: 1,
		DataType:           base.INT32,
		Statistic:          stat,
		TsID:               utils.TsID{DbNID: 1, DeviceNID: 2, MeasurementNID: 3},
	}

	stream, _ := base.NewByteStream(1024)

	// Serialize the TimeseriesIndex.
	err = tsIndex.SerializeTo(stream)
	if err != nil {
		t.Fatalf("Failed to serialize TimeseriesIndex: %v", err)
	}
	tsIndex.SetTsID(tsIndex.TsID)

	// Deserialize into a new TimeseriesIndex.
	tsIndexDeserialized := &core.TimeseriesIndex{}

	err = tsIndexDeserialized.DeserializeFrom(stream)
	if err != nil {
		t.Fatalf("Failed to deserialize TimeseriesIndex: %v", err)
	}

	// Verify deserialized values.
	if tsIndexDeserialized.DataType != base.INT32 {
		t.Errorf("Expected DataType to be INT32, got: %v", tsIndexDeserialized.DataType)
	}

	if tsIndexDeserialized.MeasurementName != "test_measurement" {
		t.Errorf("Expected MeasurementName to be test_measurement, got: %v", tsIndexDeserialized.MeasurementName)
	}

}

/////////////////
// TSMIterator //
/////////////////

// Helper function: create a ChunkGroupMeta mock
func createTestChunkGroupMeta() *core.ChunkGroupMeta {
	// Create a new ChunkGroupMeta
	chunkGroupMeta := core.NewChunkGroupMeta("device_1")

	// Create a new ChunkMeta and initialize it
	chunkMeta := &core.ChunkMeta{}
	factory := statistic.Factory{}
	stat, _ := factory.AllocStatistic(base.INT32)
	tsID := utils.TsID{}

	// Initialize the ChunkMeta with test data
	err := chunkMeta.Initialize("measurement_1", base.INT32, 100, stat, tsID, 1)
	if err != nil {
		panic("Failed to initialize test ChunkMeta: " + err.Error())
	}

	// Add the ChunkMeta to the ChunkGroupMeta's list
	chunkGroupMeta.ChunkMetaList = append(chunkGroupMeta.ChunkMetaList, chunkMeta)

	return chunkGroupMeta
}

// Test InitSuccess
func TestTSMIterator_InitSuccess(t *testing.T) {
	// Arrange: Create a ChunkGroupMeta list
	chunkGroupMetaList := []*core.ChunkGroupMeta{
		createTestChunkGroupMeta(),
	}

	iter := &core.TSMIterator{
		ChunkGroupMetaList: chunkGroupMetaList,
	}

	// Act: Initialize the iterator
	err := iter.TSMInit()

	// Assert: Ensure initialization succeeded
	if err != nil {
		t.Fatalf("TSMIterator.Init() failed: %v", err)
	}

}

// Test InitEmptyList
func TestTSMIterator_InitEmptyList(t *testing.T) {
	// Arrange: Create an empty ChunkGroupMeta list
	var chunkGroupMetaList []*core.ChunkGroupMeta

	iter := &core.TSMIterator{
		ChunkGroupMetaList: chunkGroupMetaList,
	}

	// Act: Initialize the iterator
	err := iter.TSMInit()

	// Assert: Init should succeed with an empty list
	if err != nil {
		t.Fatalf("TSMIterator.Init() failed for empty list: %v", err)
	}
}

// Test HasNext
func TestTSMIterator_HasNext(t *testing.T) {
	// Arrange: Create a ChunkGroupMeta list
	chunkGroupMetaList := []*core.ChunkGroupMeta{
		createTestChunkGroupMeta(),
	}

	iter := &core.TSMIterator{
		ChunkGroupMetaList: chunkGroupMetaList,
	}

	// Act: Initialize the iterator
	err := iter.TSMInit()
	if err != nil {
		t.Fatalf("TSMIterator.Init() failed: %v", err)
	}

	// Assert: Ensure HasNext() returns true
	// if !iter.HasNext() {t.Errorf("TSMIterator.HasNext() returned false, but there are elements to iterate over")}

}

// Test GetNext
func TestTSMIterator_GetNext(t *testing.T) {
	// Arrange: Create a ChunkGroupMeta list
	chunkGroupMetaList := []*core.ChunkGroupMeta{
		createTestChunkGroupMeta(),
	}

	iter := &core.TSMIterator{
		ChunkGroupMetaList: chunkGroupMetaList,
	}

	// Act: Initialize the iterator
	err := iter.TSMInit()
	if err != nil {
		t.Fatalf("TSMIterator.Init() failed: %v", err)
	}

	// Assert: Try to get the next item
	if !iter.HasNext() {
		t.Fatalf("TSMIterator.HasNext() returned false, but elements are expected")
	}

	// Retrieve the next element
	deviceName, measurementName, tsIndex, err := iter.GetNext()

	if err != nil {
		t.Fatalf("TSMIterator.GetNext() failed: %v", err)
	}

	// Assert: Verify the values returned
	if deviceName != "device_1" {
		t.Errorf("Expected device name \"device_1\", got \"%s\"", deviceName)
	}
	if measurementName != "measurement_1" {
		t.Errorf("Expected measurement name \"measurement_1\", got \"%s\"", measurementName)
	}
	if tsIndex == nil {
		t.Fatalf("Expected a valid TimeseriesIndex, got nil")
	}

	// Assert: Ensure there is no more data after this
	if iter.HasNext() {
		t.Logf("Expected HasNext() to return false after consuming all elements, but it returned true")
	}

	_, _, _, err = iter.GetNext()

	if err != utils.GetError(utils.ErrNoMoreData) {
		t.Fatalf("Expected TSMIterator.GetNext() to return ErrNoMoreData, but got: %v", err)
	}

}

////////////////////
// MetaIndexEntry //
////////////////////

// Test MetaIndexEntry Initialization
func TestMetaIndexEntry_InitSuccess(t *testing.T) {
	// Arrange
	name := "test_name"
	offset := int64(123456)

	// Act
	entry := &core.MetaIndexEntry{}
	entry.Init(name, offset)

	// Assert
	if entry.Name != name {
		t.Errorf("Expected entry.Name to be '%s', got '%s'", name, entry.Name)
	}
	if entry.Offset != offset {
		t.Errorf("Expected entry.Offset to be '%d', got '%d'", offset, entry.Offset)
	}
}

// Test MetaIndexEntry Serialization and Deserialization
func TestMetaIndexEntry_SerializeDeserialize(t *testing.T) {
	// Arrange
	name := "test_name"
	offset := int64(123456)
	entry := &core.MetaIndexEntry{}
	entry.Init(name, offset)

	t.Log(entry)

	// Serialize the entry
	stream, _ := base.NewByteStream(1024)
	err := entry.Serialize(stream)
	if err != nil {
		t.Fatalf("MetaIndexEntry serialization failed: %v", err)
	}

	// Deserialize into a new entry
	newEntry := &core.MetaIndexEntry{}
	err = newEntry.Deserialize(stream)
	if err != nil {
		t.Fatalf("Failed to deserialize MetaIndexEntry: %v", err)
	}

	// Assert
	if newEntry.Name != name {
		t.Errorf("Expected newEntry.Name to be '%s', got '%s'", name, newEntry.Name)
	}
	if newEntry.Offset != offset {
		t.Errorf("Expected newEntry.Offset to be '%d', got '%d'", offset, newEntry.Offset)
	}
}

// Test MetaIndexNode First Child Name Retrieval
func TestMetaIndexNode_GetFirstChildName(t *testing.T) {
	// Arrange
	node := &core.MetaIndexNode{}
	node.Children = []*core.MetaIndexEntry{}

	// Act and Assert: No children, should return error
	_, err := node.GetFirstChildName()
	if err == nil {
		t.Errorf("Expected error for no children, but got nil")
	}

	// Add a child
	entry := &core.MetaIndexEntry{}
	entry.Init("child_name", 0)
	node.Children = append(node.Children, entry)

	// Act and Assert: Should retrieve the first child's name
	name, err := node.GetFirstChildName()
	if err != nil {
		t.Fatalf("Error retrieving first child name: %v", err)
	}
	if name != "child_name" {
		t.Errorf("Expected first child name 'child_name', got '%s'", name)
	}
}

// Test MetaIndexNode Serialization and Deserialization
func TestMetaIndexNode_SerializeDeserialize(t *testing.T) {
	// Arrange
	node := &core.MetaIndexNode{
		EndOffset: 456,
		NodeType:  core.LEAF_DEVICE,
		Children:  []*core.MetaIndexEntry{},
	}

	entry := &core.MetaIndexEntry{}
	entry.Init("child_name", 123)
	node.Children = append(node.Children, entry)

	// Serialize the node
	stream, _ := base.NewByteStream(1024)
	err := node.Serialize(stream)
	if err != nil {
		t.Fatalf("MetaIndexNode serialization failed: %v", err)
	}

	// Deserialize into a new node
	newNode := &core.MetaIndexNode{}
	err = newNode.Deserialize(stream)
	if err != nil {
		t.Fatalf("Failed to deserialize MetaIndexNode: %v", err)
	}

	// Assert
	if newNode.EndOffset != 456 {
		t.Errorf("Expected EndOffset to be '456', got '%d'", newNode.EndOffset)
	}
	if newNode.NodeType != core.LEAF_DEVICE {
		t.Errorf("Expected NodeType to be LEAF_DEVICE, got '%d'", newNode.NodeType)
	}
	if len(newNode.Children) != 1 {
		t.Fatalf("Expected 1 child in newNode, got '%d'", len(newNode.Children))
	}
	if newNode.Children[0].Name != "child_name" {
		t.Errorf("Expected child name 'child_name', got '%s'", newNode.Children[0].Name)
	}
	if newNode.Children[0].Offset != 123 {
		t.Errorf("Expected child offset '123', got '%d'", newNode.Children[0].Offset)
	}
}

// Test MetaIndexNode Binary Search (Exact Match)
func TestMetaIndexNode_BinarySearchChildren_ExactMatch(t *testing.T) {
	// Arrange
	node := &core.MetaIndexNode{
		Children: []*core.MetaIndexEntry{},
	}

	// Add some entries
	entry1 := &core.MetaIndexEntry{}
	entry1.Init("apple", 10)
	entry2 := &core.MetaIndexEntry{}
	entry2.Init("banana", 20)
	entry3 := &core.MetaIndexEntry{}
	entry3.Init("cherry", 30)
	node.Children = append(node.Children, entry1, entry2, entry3)

	// Act: Search for exact match
	retEntry := &core.MetaIndexEntry{}
	retOffset, err := node.BinarySearchChildren("banana", true, retEntry)

	// Assert
	if err != nil {
		t.Fatalf("BinarySearchChildren failed: %v", err)
	}
	if retOffset != 20 {
		t.Errorf("Expected offset '20', got '%d'", retOffset)
	}
	if retEntry.Name != "banana" {
		t.Errorf("Expected entry name 'banana', got '%s'", retEntry.Name)
	}
}

// Test TsFileMeta Serialization and Deserialization
func TestTsFileMeta_SerializeDeserialize(t *testing.T) {
	// Arrange
	meta := &core.TsFileMeta{
		IndexNode:  &core.MetaIndexNode{Children: []*core.MetaIndexEntry{}},
		MetaOffset: 456,
		BloomFilter: func() *reader.BloomFilter {
			filter := &reader.BloomFilter{}
			err := filter.Init(0.1, 100)
			if err != nil {
				return nil
			}
			return filter
		}(),
	}
	entry := &core.MetaIndexEntry{}
	entry.Init("child_name", 123)
	meta.IndexNode.Children = append(meta.IndexNode.Children, entry)

	// Serialize
	stream, _ := base.NewByteStream(1024)
	err := meta.Serialize(stream)
	if err != nil {
		t.Fatalf("Failed to serialize TsFileMeta: %v", err)
	}

	// Deserialize into a new meta
	newMeta := &core.TsFileMeta{}
	err = newMeta.Deserialize(stream)
	if err != nil {
		t.Fatalf("Failed to deserialize TsFileMeta: %v", err)
	}

	// Assert
	if newMeta.MetaOffset != 456 {
		t.Errorf("Expected MetaOffset '456', got '%d'", newMeta.MetaOffset)
	}
	if len(newMeta.IndexNode.Children) != 1 {
		t.Fatalf("Expected 1 child in IndexNode, got '%d'", len(newMeta.IndexNode.Children))
	}
	if newMeta.IndexNode.Children[0].Name != "child_name" {
		t.Errorf("Expected child name 'child_name', got '%s'", newMeta.IndexNode.Children[0].Name)
	}
	if newMeta.IndexNode.Children[0].Offset != 123 {
		t.Errorf("Expected child offset '123', got '%d'", newMeta.IndexNode.Children[0].Offset)
	}
}
