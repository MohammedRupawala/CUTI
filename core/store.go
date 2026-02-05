package core

import (
	"echoserver/mod/config"
	"log"
	// "log"
	"time"
)

type value struct {
	val          interface{}
	TypeEncoding uint8
	lastAccessed uint32
	// lru_bits     uint32
}

func getLruClock() uint32 {
	return uint32(time.Now().UnixMilli()) & 0x00FFFFFF
}

var storage map[string]*value
var expires map[*value]uint64

func init() {
	storage = make(map[string]*value)
	expires = make(map[*value]uint64)
}

func setExpire(val *value, expireMs int64) {
	log.Println("Setting expiry to ", expireMs)
	expires[val] = (uint64(time.Now().UnixMilli()) + uint64(expireMs))
}

func CreateObj(val interface{}, expiration int64, t uint8, encoding uint8) *value {

	obj := &value{
		val:          val,
		TypeEncoding: t | encoding,
		lastAccessed: getLruClock(),
	}
	if expiration > 0 {
		setExpire(obj, expiration)
	}
	return obj
}

func GET(key string) *value {
	val := storage[key]
	log.Print("Value is " ,val)
	if val != nil {
		if hasexpired(val) {
			DELETE([]string{key})
			return nil
		}
		val.lastAccessed = getLruClock()
		// log.Println("")	
	}

	return val
}

func PUT(key string, val *value) {
	if len(storage) >= config.MaxKeys {
		RemoveKey()
	}
	val.lastAccessed = getLruClock()
	if KeyspaceStat[0] == nil {
		KeyspaceStat[0] = make(map[string]int)
	}

	// if GET(key) != nil {
	KeyspaceStat[0]["keys"]++
	// }
	storage[key] = val
}

func DELETE(keys []string) int64 {
	var total int64 = 0
	log.Println("Key to be deleted ", keys[0])
	for _, key := range keys {
		if _, exists := storage[key]; exists {
			KeyspaceStat[0]["keys"]--
			delete(expires, storage[key])
			delete(storage, key)
			total++
		}
	}

	return total
}

func TTL(val *value) uint64{
	return expires[val]
}
func Expire(key string, number int64) int64 {
	val := GET(key)

	if val == nil {
		return int64(0)
	} else {
		expires[val] = uint64(time.Now().UnixMilli()) + uint64(number)
		return int64(1)
	}

}
