package data_structure

import (
	"Go-Redis/internal/config"
	"sort"
)

type EvictionCandidate struct {
	key string
	lastAccessTime uint32
}

type EvictionPool struct {
	pool []*EvictionCandidate
}

type ByLastAccessTime []*EvictionCandidate

func (a ByLastAccessTime) Len() int {return len(a)}

func (a ByLastAccessTime) Swap(i , j int) {
	a[i] , a[j] = a[j] , a[i]
}

func (a ByLastAccessTime) Less(i , j int) bool {
	return a[i].lastAccessTime < a[j].lastAccessTime
}

//Push a new item to the pool , maintains the accesslastTime acsending order (old items are on the left)
//pool size > EpoolMaxSize , removes the newest item

func (p * EvictionPool) Push (key string , lastAccessTime uint32) {
	newItem := &EvictionCandidate{
		key: key,
		lastAccessTime: lastAccessTime,
	}

	//ref : https://github.com/redis/redis/blob/unstable/src/evict.c#L126
	exist := false
	for i := range len(p.pool) {
		if p.pool[i].key == key {
			exist = true
			p.pool[i] = newItem
		}
	}
	if !exist {
		p.pool = append(p.pool, newItem)
	}
	sort.Sort(ByLastAccessTime(p.pool))
	if len(p.pool) > config.EpoolMaxSize {
		lastIndex := len(p.pool) - 1
		key = p.pool[lastIndex].key
		p.pool = p.pool[:lastIndex]
	}
}

func (p *EvictionPool) Pop() *EvictionCandidate {
	if len(p.pool) == 0 {
		return  nil
	}
	oldestItem := p.pool[0]
	p.pool = p.pool[1:]
	return oldestItem
}

func newPool(sz int) *EvictionPool {
	return &EvictionPool{
		pool: make([]*EvictionCandidate, sz),
	}
}

var ePool *EvictionPool = newPool(0)