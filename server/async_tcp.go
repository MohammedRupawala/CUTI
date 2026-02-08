package server

import (
	"echoserver/mod/config"
	"echoserver/mod/core"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var connectedClients map[int]*core.Client

const SO_REUSEPORT = 0x0F

var con_clients int = 0
var cronFrequency time.Duration = 1 * time.Second
var lastCronExecTime time.Time = time.Now()

const EngineStatus_WAITING int32 = 1 << 1
const EngineStatus_BUSY int32 = 1 << 2
const EngineStatus_SHUTTING_DOWN int32 = 1 << 3

var eStatus int32 = EngineStatus_WAITING

func init() {
	connectedClients = make(map[int]*core.Client)
}

func WaitForSignal(signal chan os.Signal, wg *sync.WaitGroup) {
	<-signal
	log.Println("Interrupt Signal Detected")

	for atomic.LoadInt32(&eStatus) == EngineStatus_BUSY {

	}

	log.Println("Server Trying To Shutdown Intialized")

	atomic.StoreInt32(&eStatus, EngineStatus_SHUTTING_DOWN)
	core.Shutdown()
	log.Println("Server Ready to Exit")

	os.Exit(0)
}
func AsyncTCPServer(wg *sync.WaitGroup) error {

	defer wg.Done()
	defer func() {
		atomic.StoreInt32(&eStatus, EngineStatus_SHUTTING_DOWN)
	}()

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

	for atomic.LoadInt32(&eStatus) != EngineStatus_SHUTTING_DOWN {

		if time.Now().After(lastCronExecTime.Add(cronFrequency)) {
			core.ActiveDelete()
			lastCronExecTime = time.Now()
		}

		nevents, e := syscall.EpollWait(epollFD, events[:], -1)
		if e != nil {
			continue
		}

		if !atomic.CompareAndSwapInt32(&eStatus, EngineStatus_WAITING, EngineStatus_BUSY) {
			switch eStatus {
			case EngineStatus_SHUTTING_DOWN:
				return nil
			}
		}
		log.Print("Server Is Not in Shutting Down Stage\n")

		for i := 0; i < nevents; i++ {
			if int(events[i].Fd) == serverFD {
				fd, _, err := syscall.Accept(serverFD)
				log.Printf("New Client Accepted with fd %d\n", fd)
				if err != nil {
					continue
				}
				connectedClients[fd] = core.NewClient(fd)
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
				comm := connectedClients[int(events[i].Fd)]
				cmds, err := ReadMultipleCommands(comm)
				if err != nil {
					syscall.Close(int(events[i].Fd))
					// con_clients--
					delete(connectedClients, int(events[i].Fd))

					continue
				}
				log.Println("Command Recieved: ", cmds)

				writeCommand(comm, cmds)

			}
		}

		atomic.StoreInt32(&eStatus, EngineStatus_WAITING)
	}

	return nil
}
