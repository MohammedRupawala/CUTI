package server

import (
	"echoserver/mod/config"
	"echoserver/mod/core"
	"io"
	"log"
	"net"
	"strconv"
)

func ReadMultipleCommands(conn io.ReadWriter) (core.RedisCmds, error) {
	var read [512]byte
	n, err := conn.Read(read[:])
	if err != nil {
		return nil, err
	}

	data, err := core.DecodeArrayStrings(read[:n])
	if err != nil {
		return nil, err
	}

	var input core.RedisCmds

	for _, arr := range data {
		if len(arr) > 0 {
			cmd := &core.RedisCmd{
				Cmd: arr[0],
			}
			if len(arr) > 1 {
				cmd.Args = arr[1:]
			}
			input = append(input, cmd)
		}
	}
	return input, nil

}

func writeCommand(conn io.ReadWriter, commands core.RedisCmds) {
	core.EvalAndRespond(commands,conn)
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
			cmds, err := ReadMultipleCommands(connect)

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

			writeCommand(connect, cmds)
		}

	}

}
