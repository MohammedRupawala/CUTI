package core

import (
	"time"
)

type value struct {
	val   interface {}
	expiresAt int64	
}

var storage map[string]*value
func init(){
 	storage = make(map[string]*value)
}
// type Time = time.Time

func CreateObj(val interface{}, expiration int64) *value {

	var expiresAt int64 = -1
	if expiration > 0{
		expiresAt = time.Now().UnixMilli() + expiration
	}
	return &value{
		val: val,
		expiresAt: expiresAt,
	}
}

func PUT(key string, val *value) {
	storage[key] = val
}

func GET(key string) *value {
	return storage[key]
}

// func GetTTL(key string) int {
// 	v, ok := storage[key]
// 	if ok {
// 		getDiff := int(v.expiresAt.Sub(time.Now()).Milliseconds())

// 		if()
// 		return getDiff
// 	}
// }
