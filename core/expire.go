package core

import (
	"time"
	"log"
)

func DeleteExpiredKeys() float32 {
	var limit int = 20
	var expiredCount int = 0

	for key, val := range storage {

		exp, ok := expires[val]
		if ok && exp <= uint64(time.Now().UnixMilli()){
			delete(storage,key)
			delete(expires,val)
			limit--;
		}
		if(limit == 0){
			break
		}
	}

	return float32(expiredCount) / float32(20.0)
}


func hasexpired(val *value) bool{
	exp, ok := expires[val]
	// log.Print("Expires is thier or not ", exp)
	if !ok{
		return false
	}
	log.Println("Has Expired?")
	return exp <= uint64(time.Now().UnixMilli())

}


func ActiveDelete(){
	for{
		num := DeleteExpiredKeys()

		if(num < 0.25){
			break;
		}
	}

	log.Println("deleted the expired keys. total keys", len(storage))
}