package core

import "sort"

type evictionPoolElement struct {
	key            string
	lastAccessedAt uint32
}

type EvictionPool struct {
	pool   []*evictionPoolElement
	keyset map[string]*evictionPoolElement
}

type idleTime []*evictionPoolElement

func (arr idleTime) Len() int {
	return len(arr)
}

func (arr idleTime) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (arr idleTime) Less(i, j int) bool {
	return getIdleTime(arr[i].lastAccessedAt) > getIdleTime(arr[j].lastAccessedAt)
}

func (pq *EvictionPool) Push(key string, lastAccessedAt uint32) {

	// arr := pq.pool
	_, ok := pq.keyset[key]
	if ok {
		return
	}

	item := &evictionPoolElement{
		key:            key,
		lastAccessedAt: lastAccessedAt,
	}
	if len(pq.pool) < 16 {

		pq.pool = append(pq.pool, item)
		pq.keyset[key] = item

		sort.Sort(idleTime(pq.pool))
	}else if(lastAccessedAt > pq.pool[len(pq.pool) - 1].lastAccessedAt){
		delete(pq.keyset,pq.pool[len(pq.pool)-1].key)
		pq.pool = pq.pool[1:]
		pq.pool = append(pq.pool, item)
		pq.keyset[key] = item
	}
}


func (pq *EvictionPool) Pop() *evictionPoolElement{


	if len(pq.pool) == 0{
		return nil
	}
	
	item := pq.pool[0]
	delete(pq.keyset,item.key)
	pq.pool = pq.pool[1:]

	return item
}


func createEvictionPool(size int) *EvictionPool{
	return &EvictionPool{
		pool : make([]*evictionPoolElement,size),
		keyset : make(map[string]*evictionPoolElement),
	}
}

var ePoolSizeMax int = 16
var ePool *EvictionPool = createEvictionPool(0)