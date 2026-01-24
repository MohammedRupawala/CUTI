package main

import (
	"echoserver/mod/config"
	"echoserver/mod/server"
	"flag"
	"log"
)

func setUp() {
	flag.StringVar(&config.Host,"Host","0.0.0.0","Host from Cuti")
	flag.IntVar(&config.Port,"Port",7379,"Port for Cuti")
	flag.Parse()
}
func main() {
	log.Println("Lets Win THis")
	setUp()
	server.AsyncTCPServer()

	if err := server.AsyncTCPServer(); err != nil {
		log.Fatal("Server failed:", err)
	}

	// core_test.TestArrayDecode


}