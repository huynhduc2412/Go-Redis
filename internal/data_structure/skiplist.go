package data_structure

import (
	"math"
	"math/rand"
	"strings"
)

const SkiplistMaxLevel = 32
const minScore = float64(math.MinInt64)

type SkiplistLevel struct {
	forward *SkiplistNode
	span    uint32
}

type SkiplistNode struct {
	ele      string
	score    float64
	backward *SkiplistNode
	levels   []SkiplistLevel
}

type Skiplist struct {
	head   *SkiplistNode
	tail   *SkiplistNode
	length uint32
	level  int
}

func (sl *Skiplist) randomLevel() int {
	level := 1
	for rand.Intn(2) == 1 {
		level++
	}
	if level > SkiplistMaxLevel {
		return SkiplistMaxLevel
	}
	return level
}

func (sl *Skiplist) CreateNode(level int, score float64, ele string) *SkiplistNode {
	res := &SkiplistNode{
		ele:      ele,
		score:    score,
		backward: nil,
	}
	res.levels = make([]SkiplistLevel, level)
	return res
}

func CreateSkipList() *Skiplist {
	sl := &Skiplist{
		length: 0,
		level:  1,
	}
	sl.head = sl.CreateNode(SkiplistMaxLevel, minScore, "")
	sl.head.backward = nil
	sl.tail = nil
	return sl
}

/*
find the rank for ele by score and key
return 0 when ele not found , otherwise return rank
*/
func (sl *Skiplist) GetRank(score float64, ele string) uint32 {
	iter := sl.head
	rank := uint32(0)
	for i := sl.length - 1; i >= uint32(0); i-- {
		for iter.levels[i].forward != nil && (iter.levels[i].forward.score < score ||
			(iter.levels[i].forward.score == score && strings.Compare(iter.levels[i].forward.ele, ele) <= 0)) {
			rank += iter.levels[i].span
			iter = iter.levels[i].forward
		}
		if iter.score == score && strings.Compare(iter.ele, ele) == 0 {
			return rank
		}
	}
	return 0
}

/*
Insert new ele to Skiplist , allow duplicated scores
following step by step:

	1)find and store position need to insert new ele at each level skiplist `update`
	while we need to store all span traverse skiplist to recalculate span after insert `rank`
	2)roll coin in skiplist idea (need insert new ele at each level is determined probability)
	3)insert and recalculate span after add new ele
	4)increase span at level skiplist not update new ele
	5)Update edge case to complete skiplist insert new node
*/
func (sl *Skiplist) Insert(score float64, ele string) *SkiplistNode {
	update := [SkiplistMaxLevel]*SkiplistNode{}
	rank := [SkiplistMaxLevel]uint32{}
	iter := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] += rank[i+1]
		}
		//find position need update at each level skiplist and store update[i]
		for iter.levels[i].forward != nil && (iter.levels[i].forward.score < score ||
			(iter.levels[i].forward.score == score && strings.Compare(iter.levels[i].forward.ele, ele) == -1)) {
			rank[i] += iter.levels[i].span
			iter = iter.levels[i].forward
		}
		//store the last node < score of the ele to update insert after that
		update[i] = iter
	}

	//roll coin idea to determine the number level need insert new ele
	level := sl.randomLevel()

	//if the new node's level is higher than the highest current level of the skiplist
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.head
			// The span for new levels from the head to the end is the entire list length
			update[i].levels[i].span = sl.length
		}

		// update the highest level skiplist
		sl.level = level
	}
	//create new node
	iter = sl.CreateNode(level, score, ele)
	
	//insert new node and calculate span at each level skiplist
	for i := 0; i < level; i++ {
		//update foward pointer to insert new node
		iter.levels[i].forward = update[i].levels[i].forward
		update[i].levels[i].forward = iter
		//recalculate span
		iter.levels[i].span = update[i].levels[i].span - (rank[0] - rank[i])
		//Update the span
		update[i].levels[i].span = rank[0] - rank[i] + 1
	}

	// increase span for untouched level because we have a new node
	for i := level; i < sl.level; i++ {
		update[i].levels[i].span++
	}
	
	// Update the backward pointer for the new node, which is at the bottom level (0).
	if update[0] == sl.head {
		iter.backward = nil
	} else {
		iter.backward = update[0]
	}

	// Update the backward pointer of the node that comes after the new node.
	if iter.levels[0].forward != nil {
		iter.levels[0].forward.backward = iter
	} else {
		// If the new node is the last one in the list, update the tail.
		sl.tail = iter
	}

	// Increment the total length of the skiplist.
	sl.length++
	// Return the newly inserted node.
	return iter
}

//This function assumes that the element must exist and must match 'score'
func (sl *Skiplist) UpdateScore(curScore float64 , ele string , newScore float64) *SkiplistNode {
	update := [SkiplistMaxLevel]*SkiplistNode{}
	iter := sl.head
	for i := sl.level - 1 ; i >= 0 ; i-- {
		//find position need update at each level skiplist and store update[i]
		for iter.levels[i].forward != nil && (iter.levels[i].forward.score < curScore ||
			(iter.levels[i].forward.score == curScore && strings.Compare(iter.levels[i].forward.ele, ele) == -1)) {
			iter = iter.levels[i].forward
		}
		//store the last node < score of the ele to update insert after that
		update[i] = iter
	}
	iter = iter.levels[0].forward
	if (iter.backward == nil || iter.backward.score < newScore) &&
		(iter.levels[0].forward == nil || iter.levels[0].forward.score > newScore) {
			iter.score = newScore
			return iter
	}
	sl.DeleteNode(iter , update)
	newNode := sl.Insert(newScore , ele)	
	return newNode
}

func (sl *Skiplist) DeleteNode(x *SkiplistNode , update [SkiplistMaxLevel]*SkiplistNode) {
	for i := 0 ; i < sl.level ; i++ {
		if update[i].levels[i].forward == x {
			update[i].levels[i].span += x.levels[i].span - 1
			update[i].levels[i].forward = x.levels[i].forward
		}else{
			update[i].levels[i].span--
		}
	}
	if x.levels[0].forward != nil {
		x.levels[0].forward.backward = x.backward
	}else{
		sl.tail = x.backward
	}
	for sl.level > 1 && sl.head.levels[sl.level - 1].forward == nil {
		sl.level--
	}
	sl.length--
}

func (sl *Skiplist) Delete(score float64 , ele string) int {
	update := [SkiplistMaxLevel]*SkiplistNode{}
	iter := sl.head
	for i := sl.level - 1 ; i >= 0 ; i-- {
		//find position need update at each level skiplist and store update[i]
		for iter.levels[i].forward != nil && (iter.levels[i].forward.score < score ||
			(iter.levels[i].forward.score == score && strings.Compare(iter.levels[i].forward.ele, ele) == -1)) {
			iter = iter.levels[i].forward
		}
		//store the last node < score of the ele to update insert after that
		update[i] = iter
	}
	iter = iter.levels[0].forward
	if iter != nil && iter.score == score && strings.Compare(iter.ele , ele) == 0 {
		sl.DeleteNode(iter , update)
		return 1
	}
	return 0
}