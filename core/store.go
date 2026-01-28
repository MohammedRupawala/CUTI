package core

import (
	"echoserver/mod/config"
	"time"
)

type value struct {
	val       interface{}
	TypeEncoding uint8
	expiresAt int64
}



var storage map[string]*value

func init() {
	storage = make(map[string]*value)
}


func CreateObj(val interface{}, expiration int64, t uint8 , encoding uint8) *value {

	var expiresAt int64 = -1
	if expiration > 0 {
		expiresAt = time.Now().UnixMilli() + expiration
	}
	return &value{
		val:  val,
		TypeEncoding: t | encoding,
		expiresAt: expiresAt,
	}
}

func PUT(key string, val *value) {
	if(len(storage) == config.MaxKeys){
		RemoveKey()
	}
	storage[key] = val
}

func GET(key string) *value {
	return storage[key]
}

func DELETE(keys []string) int64 {
    var total int64 = 0

    for _, key := range keys {
        if _, exists := storage[key]; exists {
            delete(storage, key)
            total++ 
        }
    }
    
    return total
}

func Expire(key string , number int64) int64{
	val := GET(key)

	if(val == nil){
		return int64(0)
	}else{
		val.expiresAt = time.Now().UnixMilli() + number
		return int64(1)
	}

}
