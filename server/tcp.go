package server

import (
	"echoserver/mod/config"
	"echoserver/mod/core"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func readCommand(conn io.ReadWriter) (*core.RedisCmd, error) {
	var read [512]byte

	n, err := conn.Read(read[:])
	if err != nil {
		return nil, err
	}

	data, error := core.DecodeArrayString(read[:n])
	if error != nil {
		return nil, error
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	data[0] = strings.ToLower(data[0])
	return &core.RedisCmd{
		Cmd:  data[0],
		Args: data[1:], 
	}, nil

	// return string(read[:n]), nil
}

func writeCommand(conn io.ReadWriter, command *core.RedisCmd) {
	// mess := strings.ToLower(command.Cmd)

	// log.Println("Write Command")
	lowerCommand := strings.ToLower(command.Cmd)
	if lowerCommand == "ping" {
		err := core.EvalAndRespond(command, conn)
		if err != nil {
			core.ErrorResponse(err,conn)
		}
	} else {
		conn.Write([]byte("+OK\r\n"))
	}
}

func TcpSyncServer() {
	log.Println("You are In Server And The Host is ", config.Host+"Running on Port ", config.Port)

	var num_client int = 0

	listen, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {
		connect, err := listen.Accept()

		if err != nil {
			panic(err)
		}

		num_client += 1
		// log.Println("Connected with Client with Host ", connect.RemoteAddr(), " Connected Client ", num_client)

		for {
			cmd, err := readCommand(connect)

			if err != nil {

				num_client -= 1
				log.Println("Connection Disconnected With host ", connect.RemoteAddr(), "Connected Clients ", num_client)

				if err == io.EOF {
					break
				}

				log.Println("err ", err)

			}

			// log.Println("Command Recieved From Client ", cmd)

			// if err = writeCommand(connect, cmd); err != nil {
			// 	log.Println("Error while Writing ", err)
			// }

			writeCommand(connect, cmd)
		}

	}

}
