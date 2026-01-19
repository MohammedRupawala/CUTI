package core

import (
	"errors"
	// "log"
	"net"
)

func EvalAndRespond(cmd *RedisCmd, conn net.Conn)(error){
	command := cmd.Cmd
	args := cmd.Args
	var b []byte

	// log.Println("Evaluating")
	if(command == "PING" && len(args) > 1){
		return errors.New("ERR wrong number of arguments for 'ping' command")
	}else if(command == "PING"){
		if len(args) == 1{
			b = Encode(args[0],false)
		}else{
			b = Encode("PONG",true)
		}
	}

	_,err := conn.Write(b)
	// if(err != nil){
	// 	retirn err
	// }
	return err
}

func ErrorResponse(conn net.Conn) error {
	_,err := conn.Write([]byte("ERR wrong number of arguments for 'ping' command"))
	return err
}