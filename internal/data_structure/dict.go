package data_structure

import (
	"Go-Redis/internal/config"
	"log"
	"time"
)

type Obj struct {
	Value interface{}
	LastAccessTime uint32
}

type Dict struct {
	dicStore map[string]*Obj
	expiredDictStore map[string]uint64
}

func CreateDict() *Dict {
	res := &Dict{
		dicStore: make(map[string]*Obj),
		expiredDictStore:  make(map[string]uint64),
	}
	return res
}

func(d *Dict) GetDictStore() map[string]*Obj{
	return d.dicStore
}

func now() uint32 {
	return uint32(time.Now().Unix())
}

func (d *Dict) GetExpireDictStore() map[string]uint64 {
	return d.expiredDictStore
}

func (d *Dict) NewObj(key string , value interface{} , ttlMs int64) *Obj {
	obj := &Obj{
		Value: value,
		LastAccessTime: now(),
	}

	// add key in the expiredDictStore
	if ttlMs > 0 {
		d.SetExpiry(key , ttlMs)
	}
	return obj
}

func (d *Dict) GetExpiry(key string) (uint64 , bool) {
	exp , exist := d.expiredDictStore[key]
	return exp , exist
}

func (d *Dict) SetExpiry(key string , ttlMs int64) {
	d.expiredDictStore[key] = uint64(time.Now().UnixMilli()) + uint64(ttlMs)
}

func (d *Dict) HasExpired(key string) bool {
	exp , exist := d.expiredDictStore[key]
	if !exist {
		return false
	}
	return exp <= uint64(time.Now().UnixMilli())
}

func (d *Dict) Get(k string) *Obj {
	v := d.dicStore[k]
	if v != nil {
		v.LastAccessTime = now()
		if d.HasExpired(k) {
			d.Del(k)
			return nil
		}
	}
	return v
}


func (d *Dict) Set(k string , obj *Obj) {
	if len(d.dicStore) == config.MaxKeyNumber {
		d.evict()
	}
	if _ , exist := d.dicStore[k] ; !exist {
		HashKeySpaceStat.Key++
	}
	d.dicStore[k] = obj
}

func (d *Dict) Del(k string) bool {
	log.Printf("Delete key %s" , k)
	if _ , exist := d.dicStore[k] ; exist {
		delete(d.dicStore , k)
		delete(d.expiredDictStore , k)
		HashKeySpaceStat.Key--
		return true
	}
	return false
}

func (d *Dict) populateEpool() {
	remain := config.EpoolLruSampleSize
	for k := range d.dicStore {
		ePool.Push(k , d.dicStore[k].LastAccessTime)
		remain--
		if remain == 0 {
			break
		}
	}
	log.Println("Epool:")
	for _ , item := range ePool.pool {
		log.Println(item.key , item.lastAccessTime)
	}
}

func (d *Dict) evictLru() {
	d.populateEpool()
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))
	log.Print("trigger LRU eviction")
	for i := 0 ; i < int(evictCount) && len(ePool.pool) > 0 ; i++ {
		item := ePool.Pop()
		if item != nil {
			d.Del(item.key)
		}
	}
}

func (d *Dict) evictRandom() {
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))
	log.Printf("trigger random eviction")
	for k := range d.dicStore {
		d.Del(k)
		evictCount--
		if evictCount == 0 {
			break
		}
	}
}

func (d *Dict) evict() {
	switch config.EvictionPolicy {
	case "allkeys-random":
		d.evictRandom()
	case "allkeys-lru":
		d.evictLru()
	}
}
