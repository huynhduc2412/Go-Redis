package data_structure

import (
	"math"
	"github.com/spaolacci/murmur3"
)

//log10(0.5)
//precomputed valued
const Log10PointFive = -0.30102999566

//CMS is Count-Min Sketch data structure.
//this is matrix 2D with width x depth
type CMS struct {
	width uint32
	depth uint32
	
	//counter is matrix 2D with row(depth) x col(width)
	counter [][]uint32
} 

func CreateCMS(w uint32 , d uint32) *CMS {
	cms := &CMS{
		width: w,
		depth: d,
	}
	cms.counter = make([][]uint32, d)
	for i := range d {
		cms.counter[i] = make([]uint32, w)
	}
	return cms
}

//CalcCMSDim calculates the dimensions (width and depth) of the CMS
//based on the desired error rate and probability
func CalcCMSDim(errRate float64 , errProb float64) (uint32 , uint32) {
	w := uint32(math.Ceil(2.0 / errRate))
	d := uint32(math.Ceil(math.Log10(errProb) / Log10PointFive))
	return w , d
}

//calcHash calculates a 32-bit hash for the given item and seed 
func (c *CMS) calcHash(item string , seed uint32) uint32 {
	hasher := murmur3.New32WithSeed(seed)
	hasher.Write([]byte(item))
	return hasher.Sum32()
}

//incre the count for an item by a specific value
// return estimate count for the item after increasing
func (c *CMS) IncrBy(item string , value uint32) uint32 {
	var minCount uint32 = math.MaxUint32

	for i := range c.depth {
		//calculate hash
		hash := c.calcHash(item , i)
		//use the hash to get the column index within the row
		j := hash % c.width

		//avoid overflow 32-bit
		if math.MaxUint32 - c.counter[i][j] < value {
			c.counter[i][j] = math.MaxUint32
		} else {
			c.counter[i][j] += value
		}

		//keep track of value min in position hash value calculated respective index col
		minCount = min(minCount , c.counter[i][j])
	}
	return minCount
}

//return the estimated count for an item.
func (c *CMS) Count(item string) uint32 {
	var minCount uint32 = math.MaxUint32
	for i := range c.depth {
		hash := c.calcHash(item , i)
		j := hash % c.width
		minCount = min(minCount , c.counter[i][j])
	}
	return minCount
}

