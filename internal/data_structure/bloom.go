package data_structure

import (
	"math"

	"github.com/spaolacci/murmur3"
)


const Ln2 float64 = 0.693147180559945
const Ln2Square float64 = 0.480453013918201
const ABigSeed uint32 = 0x9747b28c

type Bloom struct {
	Hashes int
	Entries uint64
	Error float64
	bitPerEntry float64
	bf []bool
	bytes uint64 // size of bf in byte
}

type HashValue struct {
	a , b uint64
}

func calcBpe(err float64) float64 {
	num := math.Log(err)
	return math.Abs(-(num / Ln2Square))
}

/*
http://en.wikipedia.org/wiki/Bloom_filter
- Optimal number of bits is: bits = (entries * ln(error)) / ln(2)^2
- bitPerEntry = bits/entries
- Optimal number of hash functions is: hashes = bitPerEntry * ln(2)
*/
func CreateBloomFilter(entries uint64 , errRate float64) *Bloom {
	bloom := &Bloom{ 
		Entries: entries,
		Error: errRate,
	}
	bloom.bitPerEntry = calcBpe(errRate)
	bits := uint64(float64(entries) * bloom.bitPerEntry)
	if bits % 64 != 0 {
		bloom.bytes = ((bits / 64) + 1) * 8
	} else {
		bloom.bytes = bits / 8
	}
	bloom.Hashes = int(math.Ceil(Ln2 * bloom.bitPerEntry))
	bloom.bf = make([]bool, bloom.bytes)
	return bloom
}

func (b *Bloom) CalcHash(entry string) HashValue {
	hasher := murmur3.New128WithSeed(ABigSeed)
	hasher.Write([]byte(entry))
	x , y := hasher.Sum128()
	return HashValue{
		a: x,
		b: y,
	}
}

func (b *Bloom) Add(entry string) {
	var bytePos uint64
	initHash := b.CalcHash(entry)
	for i := range b.Hashes {
		bytePos = (initHash.a + initHash.b * uint64(i)) % b.bytes
		b.bf[bytePos] = true
	}
}

func (b *Bloom) Exist(entry string) bool {
	var bytePos uint64
	initHash := b.CalcHash(entry)
	for i := range b.Hashes {
		bytePos = (initHash.a + initHash.b * uint64(i)) % b.bytes
		if !b.bf[bytePos] {
			return false
		}
	}
	return true
}

func (b *Bloom) AddHash(initHash HashValue) {
	var bytePos uint64
	for i := range b.Hashes {
		bytePos = (initHash.a + initHash.b * uint64(i)) % b.bytes
		b.bf[bytePos] = true
	}
}

func (b *Bloom) ExistHash(initHash HashValue) bool {
	var bytePos uint64
	for i := range b.Hashes {
		bytePos = (initHash.a + initHash.b * uint64(i)) % b.bytes
		if !b.bf[bytePos] {
			return false
		}
	}
	return true
}
