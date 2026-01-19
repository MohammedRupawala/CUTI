package server

import (
	"echoserver/mod/config"
	"echoserver/mod/core"
	"io"
	"log"
	"net"
	"strconv"
)


func readCommand(conn net.Conn) (string,error) {
	var read [512]byte

	n , err := conn.Read(read[:])
	if err != nil {
		return "",err
	}

	//  := core.Decode(read[:n]);

	return string(read[:n]),nil
}

func writeCommand(conn net.Conn,command string) error{
	_ , err := conn.Write([]byte(command))

	if err != nil{
		return err
	}
	return nil
}




func TcpSyncServer() {
	log.Println("You are In Server And The Host is ", config.Host + "Running on Port ", config.Port)

	var num_client int = 0

	listen, err := net.Listen("tcp",config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {
		connect,err := listen.Accept()

		if err != nil{
			panic(err)
		}

		num_client += 1;
		log.Println("Connected with Client with Host ",connect.RemoteAddr(), " Connected Client ",num_client)
		
		for{
			cmd,err := readCommand(connect)
			
			if err != nil{
				
				num_client -= 1;
				log.Println("Connection Disconnected With host ", connect.RemoteAddr() , "Connected Clients ", num_client )

				if err == io.EOF{
					break
				}

				log.Println("err ",err)


			}

			log.Println("Command Recieved From Client " , cmd)

			if err = writeCommand(connect,cmd);err != nil{
				log.Println("Error while Writing " , err)
			}
		}

	}

}