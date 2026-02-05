package core

import (
	"echoserver/mod/config"
	"log"
)

func randomDelete() {
	for k := range storage {
		delete(storage, k)
		return
	}
}

func getIdleTime(lastAccessed uint32) uint32 {
	c := getLruClock()
	if c >= lastAccessed {
		return c - lastAccessed
	} else {
		return (0x00FFFFFF - lastAccessed) + c
	}
}

func populatePool() {
	sample := 5
	for k := range storage {
		ePool.Push(k, storage[k].lastAccessed)
		sample--
		if sample == 0 {
			break
		}
	}
}

func approxLru() {
	populatePool()
	log.Println("Populated Pool", len(ePool.pool))
	evictCount := int16(config.EvictionFactor * float64(config.MaxKeys))

	for i := 0; i < int(evictCount) && len(ePool.pool) > 0; i++ {
		item := ePool.Pop()

		if item == nil {
			return
		}

		DELETE([]string{item.key})
	}
}

func allKeyRandom() {
	count := config.EvictionFactor * float64(config.MaxKeys)
	var keysList []string
	for k := range storage {
		keysList = append(keysList, k)
		count--
		if count <= 0 {
			break
		}
	}

	DELETE(keysList)
}

func RemoveKey() {
	switch config.EvictionStrategy {
	case "simple":
		randomDelete()
	case "all-key-random":
		allKeyRandom()
	case "approx-lru":

		log.Println("Approx-Lru Called")
		approxLru()

	}

}
