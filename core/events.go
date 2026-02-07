package core

import "log"

func Shutdown() {

	log.Println("STARTING WRITING AOF FILE")
	evalAOF([]string{})
	log.Print("FINISHED WRITING AOF FILE")

}
