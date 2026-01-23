package server

import (
	"echoserver/mod/config"
	"echoserver/mod/core"
	"log"
	"net"
	"strconv"
	"syscall"
)


var con_clients int = 0
func AsyncTCPServer() error {
	log.Println("Welcome To Async Server " + config.Host + " Running at Port " + strconv.Itoa(config.Port))

	var max_client int = 20000

	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_client)

	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM|syscall.O_NONBLOCK, 0)

	if err != nil {
		return err
	}
	
	defer syscall.Close(serverFD)

	err = syscall.SetNonblock(serverFD, true)
	if err != nil {
		return err
	}

	ip4 := net.ParseIP(config.Host)
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		return err
	}

	if err = syscall.Listen(serverFD, max_client); err != nil {
		return err
	}

	epollFD, err := syscall.EpollCreate1(0)
	if err != nil{
		return err
	}

	defer syscall.Close(epollFD)
	var socketServerEvent syscall.EpollEvent = syscall.EpollEvent{
		Events : syscall.EPOLLIN,
		Fd : int32(serverFD),
	}

	if err = syscall.EpollCtl(epollFD,syscall.EPOLL_CTL_ADD,serverFD,&socketServerEvent); err != nil{
		return err
	}

	for {
		nevents , e := syscall.EpollWait(epollFD,events[:],-1)
		if e != nil{
			continue;
		}


		for i := 0;i<nevents;i++{
				if(int(events[i].Fd) == serverFD){
					fd , _ , err := syscall.Accept(serverFD)
					if err != nil{
						continue
					}

					con_clients++;
					syscall.SetNonblock(serverFD,true)
					log.Println("Client Connected with Fd Number " , fd )

					var socketClient syscall.EpollEvent = syscall.EpollEvent{
						Events: syscall.EPOLLIN,
						Fd : int32(fd),
					}

					if err = syscall.EpollCtl(epollFD,syscall.EPOLL_CTL_ADD,fd,&socketClient); err != nil{
						log.Fatal(err)
					}


					// log.Println("Connected With Client " + sockaddrToString(addr))
					
				}else{
					comm := core.FDComm{Fd : int(events[i].Fd)}
					cmd , err := readCommand(comm)

					if(err != nil){
						syscall.Close(int(events[i].Fd))
						con_clients--;
						continue
					}

					writeCommand(comm,cmd)
					
				}
		}
	}
}
