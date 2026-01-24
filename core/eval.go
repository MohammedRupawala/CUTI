package core

import (
	"errors"
	// "log"
	"io"
	"fmt"
)

func EvalAndRespond(cmd *RedisCmd, conn io.ReadWriter)(error){
	command := cmd.Cmd
	args := cmd.Args
	var b []byte

	// log.Println("Evaluating")
	if(command == "ping" && len(args) > 1){
		return errors.New("ERR wrong number of arguments for 'PING' command")
	}else if(command == "ping"){
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

func ErrorResponse(err error,conn io.ReadWriter) {
		conn.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}