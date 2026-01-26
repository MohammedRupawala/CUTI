package server

import (
	"echoserver/mod/config"
	"echoserver/mod/core"
	"log"
	"net"
	"strconv"
	"syscall"
	"time"
)

const SO_REUSEPORT = 0x0F

var con_clients int = 0
var cronFrequency time.Duration = 1 * time.Second
var lastCronExecTime time.Time = time.Now()
func AsyncTCPServer() error {


	log.Println("Welcome To Async Server " + config.Host + " Running at Port " + strconv.Itoa(config.Port))

	var max_client int = 20000

	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_client)

	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM|syscall.O_NONBLOCK, 0)

	if err != nil {
		return err
	}

	log.Println("Server is Running")
	defer syscall.Close(serverFD)

	err = syscall.SetNonblock(serverFD, true)
	if err != nil {
		return err
	}

	syscall.SetsockoptInt(
		serverFD,
		syscall.SOL_SOCKET,
		SO_REUSEPORT,
		1,
	)

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
	if err != nil {
		return err
	}

	defer syscall.Close(epollFD)
	var socketServerEvent syscall.EpollEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFD),
	}

	if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &socketServerEvent); err != nil {
		return err
	}

	for {

		if(time.Now().After(lastCronExecTime.Add(cronFrequency))){
			core.ActiveDelete()
			lastCronExecTime = time.Now()
		}
		nevents, e := syscall.EpollWait(epollFD, events[:], -1)
		if e != nil {
			continue
		}

		for i := 0; i < nevents; i++ {
			if int(events[i].Fd) == serverFD {
				fd, _, err := syscall.Accept(serverFD)
				if err != nil {
					continue
				}

				con_clients++
				syscall.SetNonblock(fd, true)
				// log.Println("Client Connected with Fd Number ", fd)

				var socketClient syscall.EpollEvent = syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(fd),
				}

				if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketClient); err != nil {
					log.Fatal(err)
				}

				// log.Println("Connected With Client " + sockaddrToString(addr))

			} else {
				comm := core.FDComm{Fd: int(events[i].Fd)}
				cmd, err := readCommand(comm)
				// log.Println("Command Recieved: " + cmd.Cmd)
				if err != nil {
					syscall.Close(int(events[i].Fd))
					con_clients--
					continue
				}

				writeCommand(comm, cmd)

			}
		}
	}
}
