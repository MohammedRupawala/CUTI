package core

import "echoserver/mod/config"

func randomDelete() {
	for k := range storage {
		delete(storage, k)
		return
	}
}


func allKeyRandom(){
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

	}

}