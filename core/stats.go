package core

import "log"

var KeyspaceStat [4]map[string]int

func Update(num int, metric string, value int) {
	log.Println("Called KeyStat")
	KeyspaceStat[num][metric] = value
}