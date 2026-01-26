package core

import "echoserver/mod/config"

func randomDelete() {
	for k := range storage {
		delete(storage, k)
		return
	}
}

func RemoveKey() {
	switch config.EvictionStrategy {
	case "simple":
		randomDelete()
	}
}