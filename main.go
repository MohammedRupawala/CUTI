package main

import (
	"echoserver/mod/config"
	"echoserver/mod/server"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func setUp() {
	flag.StringVar(&config.Host,"Host","0.0.0.0","Host from Cuti")
	flag.IntVar(&config.Port,"Port",7379,"Port for Cuti")
	flag.Parse()
}
func main() {
	log.Println("Lets Win THis")
	setUp()


	var signalPipe chan os.Signal = make(chan os.Signal,1)
	signal.Notify(signalPipe,syscall.SIGTERM,syscall.SIGINT)


	var wg sync.WaitGroup 
	wg.Add(2)



	go server.AsyncTCPServer(&wg)
	go server.WaitForSignal(signalPipe,&wg)

	if err := server.AsyncTCPServer(&wg); err != nil {
		log.Fatal("Server failed:", err)
	}

	// core_test.TestArrayDecode


	wg.Wait()

}