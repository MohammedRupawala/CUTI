package core

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"
)

func evalPing(args []string) ([]byte, error) {
	n := len(args)
	log.Print("PING Command Found")
	var b []byte
	switch n {
	case 0:
		b = Encode("PONG", "simpleString")
	case 1:
		log.Println(args[0] + " This is Argument")
		b = Encode(args[0], "bulkString")
	default:
		return []byte{}, errors.New("ERR wrong number of arguments for 'PING' command")
	}
	return b, nil

}

func evalSet(args []string) ([]byte, error) {
	n := len(args)
	var b []byte
	if n < 2 {
		return []byte{}, errors.New("ERR wrong number of arguments for 'set' command")
	}

	key := args[0]
	val := args[1]
	var exDurationMs int64
	var expFound bool
	for i := 2; i < n; i++ {
		switch args[i] {
		case "ex":
			if expFound {
				return []byte{}, errors.New("ERR syntax error")
			}
			i++
			expFound = true
			if i == n {
				return []byte{}, errors.New("ERR syntax error")
			} else {

				exDurationSecs, parseErr := strconv.ParseInt(args[i], 10, 64)
				if parseErr != nil {
					return []byte{}, errors.New("ERR value is not an integer or out of range")
				}

				exDurationMs = exDurationSecs * 1000
			}
		case "px":
			if expFound {
				return []byte{}, errors.New("ERR syntax error")
			}
			i++
			expFound = true
			if i == n {
				return []byte{}, errors.New("ERR syntax error")
			} else {
				exDurationSecs, parseErr := strconv.ParseInt(args[i], 10, 64)
				if parseErr != nil {
					return []byte{}, errors.New("ERR value is not an integer or out of range")
				}

				exDurationMs = exDurationSecs

			}
		default:
			return []byte{}, errors.New("Err syntax error")
		}
	}

	if exDurationMs > 0 {
		PUT(key, CreateObj(val, exDurationMs))
	} else {
		PUT(key, CreateObj(val, -1))
	}
	b = Encode("OK", "simpleString")
	return b, nil

}

func evalGet(args []string) ([]byte, error){
	n := len(args)
	if(n > 1 || n < 1){
		return []byte{},errors.New("ERR wrong number of arguments for 'get' command")
	}else{
		value := GET(args[0])
		if(value.expiresAt != -1 && value.expiresAt < time.Now().UnixMilli() || value == nil ){
			return Encode("(nil)","simpleString"),nil
		}else{
			return Encode(value.val,"bulkString"),nil
		}
	}
}
func EvalAndRespond(cmd *RedisCmd, conn io.ReadWriter) error {
	command := cmd.Cmd
	args := cmd.Args
	var b []byte
	var e error

	switch command {
	case "ping":
		b, e = evalPing(args)
	case "set":
		log.Println("Set Command Found")
		b, e = evalSet(args)
	case "get":
		b,e = evalGet(args)
	}

	if e != nil {
		return e
	}

	_, err := conn.Write(b)
	return err
}

func ErrorResponse(err error, conn io.ReadWriter) {
	conn.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}
