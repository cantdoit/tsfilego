package reader

import (
	"Golang/internal/common/base"
	"errors"
	_ "errors"
	"hash/fnv"
	"math"
)

////////////
// Bitset //
////////////

// BitSet represents a data structure for managing a set of bits efficiently.
type BitSet struct {
	Words     []uint64 // Words is a slice of 64-bit integers storing the bits.
	WordCount int32    // WordCount tracks the number of words in the BitSet.
}

// Init initializes the BitSet with a specific size in bits.
func (bs *BitSet) Init(size int32) {
	if size <= 1 {
		panic("Size must be greater than 1")
	}
	bs.WordCount = (size-1)/64 + 1
	bs.Words = make([]uint64, bs.WordCount)
}

// Destroy clears the BitSet.
func (bs *BitSet) Destroy() {
	bs.Words = nil
	bs.WordCount = 0
}

// Set sets a specific bit in the BitSet.
func (bs *BitSet) Set(pos int32) {
	wordIdx := pos / 64
	wordOffset := pos % 64
	bs.Words[wordIdx] |= 1 << wordOffset
}

// GetWordsInUse returns the number of 64-bit words actually used by the BitSet.
func (bs *BitSet) GetWordsInUse() int32 {
	for i := bs.WordCount - 1; i >= 0; i-- {
		if bs.Words[i] != 0 {
			return i + 1
		}
	}
	return 0
}

// SerializeTo writes the BitSet into a ByteStream.
func (bs *BitSet) SerializeTo(stream *base.ByteStream) error {
	serial := base.SerializationUtil{}
	wordsInUse := bs.GetWordsInUse()
	if err := serial.WriteUint32(uint32(wordsInUse), stream); err != nil {
		return err
	}
	for i := int32(0); i < wordsInUse; i++ {
		if err := serial.WriteUint64(bs.Words[i], stream); err != nil {
			return err
		}
	}
	return nil
}

// DeserializeFrom reads the BitSet from a ByteStream.
func (bs *BitSet) DeserializeFrom(stream *base.ByteStream) error {
	serial := base.SerializationUtil{}
	wordsInUse, err := serial.ReadUint32(stream)
	if err != nil {
		return err
	}

	bs.WordCount = int32(wordsInUse)
	bs.Words = make([]uint64, wordsInUse)
	for i := int32(0); i < int32(wordsInUse); i++ {
		bs.Words[i], err = serial.ReadUint64(stream)
		if err != nil {
			return err
		}
	}
	return nil
}

//////////////////
// hashfunction //
//////////////////

// HashFunction represents a structure for hashing with a fixed capacity and an optional seed for initialization.
// It allows the computation of consistent hashes within a bounded range defined by the capacity.
// The seed provides variability to the hash values for different instances of HashFunction.
type HashFunction struct {
	Cap  int32
	Seed int32
}

// Init initializes the hash function with capacity and seed.
func (hf *HashFunction) Init(cap, seed int32) {
	if cap <= 1 {
		panic("Capacity must be greater than 1")
	}
	hf.Cap = cap
	hf.Seed = seed
}

// Hash computes a hash for the given buffer.
func (hf *HashFunction) Hash(buf string) int32 {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(buf))
	hashValue := int32(hasher.Sum32()) + hf.Seed

	// Ensure hash is positive and bounded by capacity
	if hashValue == math.MinInt32 {
		hashValue = 0
	}
	hashValue = int32(math.Abs(float64(hashValue))) % hf.Cap
	return hashValue
}

/////////////////
// bloomfilter //
/////////////////

// BloomFilter represents a probabilistic data structure for efficient membership testing with a configurable error rate.
type BloomFilter struct {
	Size           int32
	HashFuncCount  int32
	HashFunctions  []HashFunction
	BitSetInstance BitSet
}

const (
	MaxHashFuncCount int32 = 8
	MinSize          int32 = 256
	MinBFErrorRate         = 0.01
	MaxBFErrorRate         = 0.1
)

// Init initializes the BloomFilter with error percentage and entry count.
func (bf *BloomFilter) Init(errorPercent float64, entryCount int) error {
	if errorPercent < MinBFErrorRate || errorPercent > MaxBFErrorRate {
		return errors.New("errorPercent is out of valid range")
	}
	if entryCount <= 0 {
		return errors.New("entryCount must be greater than 0")
	}

	bf.Size = int32(-1 * float64(entryCount) * math.Log(errorPercent) / math.Pow(math.Log(2), 2))
	if bf.Size < MinSize {
		bf.Size = MinSize
	}

	bf.HashFuncCount = int32(math.Ceil(math.Log(2) * float64(bf.Size) / float64(entryCount)))
	if bf.HashFuncCount > MaxHashFuncCount {
		bf.HashFuncCount = MaxHashFuncCount
	}

	bf.HashFunctions = make([]HashFunction, bf.HashFuncCount)
	for i, size := 0, bf.Size; i < int(bf.HashFuncCount); i++ {
		bf.HashFunctions[i].Init(size, int32(i))
	}

	bf.BitSetInstance.Init(bf.Size)
	return nil
}

// AddPathEntry adds an entry to the BloomFilter.
func (bf *BloomFilter) AddPathEntry(deviceName, measurementName string) {
	entryString := deviceName + ":" + measurementName
	for i := int32(0); i < bf.HashFuncCount; i++ {
		hash := bf.HashFunctions[i].Hash(entryString)
		bf.BitSetInstance.Set(hash)
	}
}

// SerializeTo serializes the BloomFilter into a ByteStream.
func (bf *BloomFilter) SerializeTo(stream *base.ByteStream) error {
	serial := base.SerializationUtil{}
	if err := serial.WriteUint32(uint32(bf.Size), stream); err != nil {
		return err
	}
	if err := serial.WriteUint32(uint32(bf.HashFuncCount), stream); err != nil {
		return err
	}
	for _, hashFunc := range bf.HashFunctions {
		if err := serial.WriteUint32(uint32(hashFunc.Seed), stream); err != nil {
			return err
		}
	}
	if err := bf.BitSetInstance.SerializeTo(stream); err != nil {
		return err
	}
	return nil
}

// DeserializeFrom deserializes a BloomFilter from a ByteStream.
func (bf *BloomFilter) DeserializeFrom(stream *base.ByteStream) error {
	serial := base.SerializationUtil{}
	size, err := serial.ReadUint32(stream)
	if err != nil {
		return err
	}
	bf.Size = int32(size)

	hashFuncCount, err := serial.ReadUint32(stream)
	if err != nil {
		return err
	}
	bf.HashFuncCount = int32(hashFuncCount)

	bf.HashFunctions = make([]HashFunction, hashFuncCount)
	for i := int32(0); i < int32(hashFuncCount); i++ {
		seed, err := serial.ReadUint32(stream)
		if err != nil {
			return err
		}
		bf.HashFunctions[i].Init(bf.Size, int32(seed))
	}

	err = bf.BitSetInstance.DeserializeFrom(stream)
	if err != nil {
		return err
	}
	return nil
}
