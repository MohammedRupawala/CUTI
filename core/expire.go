package core

import (
	"time"
	"log"
)

func DeleteExpiredKeys() float32 {
	var limit int = 20
	var expiredCount int = 0

	for key, val := range storage {

		if val.expiresAt != -1 {
			limit--
			if val.expiresAt <= time.Now().UnixMilli() {
				expiredCount++;
				delete(storage,key)
			}
		}
		if(limit == 0){
			break
		}
	}

	return float32(expiredCount) / float32(20.0)
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