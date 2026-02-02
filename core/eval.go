package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"
)

func evalPing(args []string) []byte {
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
		return Encode(errors.New("ERR wrong number of arguments for 'PING' command"), "simpleString")
	}
	return b

}

func evalSet(args []string) []byte {
	n := len(args)
	var b []byte
	if n < 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'set' command"), "simpleString")
	}

	key := args[0]
	val := args[1]

	vType, vEnc := deduceTypeEncoding(val)
	var exDurationMs int64
	var expFound bool
	for i := 2; i < n; i++ {
		switch args[i] {
		case "ex":
			if expFound {
				return Encode(errors.New("ERR syntax error"), "simpleString")
			}
			i++
			expFound = true
			if i == n {
				return Encode(errors.New("ERR syntax error"), "simpleString")
			} else {

				exDurationSecs, parseErr := strconv.ParseInt(args[i], 10, 64)
				if parseErr != nil {
					return Encode(errors.New("ERR value is not an integer or out of range"), "simpleString")
				}

				exDurationMs = exDurationSecs * 1000
			}
		case "px":
			if expFound {
				return Encode(errors.New("ERR syntax error"), "simpleString")
			}
			i++
			expFound = true
			if i == n {
				return Encode(errors.New("ERR syntax error"), "simpleString")
			} else {
				exDurationSecs, parseErr := strconv.ParseInt(args[i], 10, 64)
				if parseErr != nil {
					return Encode(errors.New("ERR value is not an integer or out of range"), "simpleString")
				}

				exDurationMs = exDurationSecs

			}
		default:
			return Encode(errors.New("ERR syntax error"), "simpleString")
		}
	}

	if exDurationMs > 0 {
		PUT(key, CreateObj(val, exDurationMs, vType, vEnc))
	} else {
		PUT(key, CreateObj(val, -1, vType, vEnc))
	}
	b = Encode("OK", "simpleString")
	return b

}

func evalGet(args []string) []byte {
	n := len(args)
	if n > 1 || n < 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'get' command"), "simpleString")
	} else {
		value := GET(args[0])
		if value == nil || value.expiresAt != -1 && value.expiresAt < time.Now().UnixMilli(){
			delete(storage, args[0])
			return Encode("(nil)", "simpleString")
		} else {
			return Encode(value.val, "bulkString")
		}
	}
}

func evalTTL(args []string) []byte {
	n := len(args)
	if n > 1 || n < 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ttl' command"), "simpleString")
	} else {
		value := GET(args[0])
		if value == nil {
			return Encode(int64(-2), "number")
		} else if value.expiresAt == -1 {
			return Encode(int64(-1), "number")
		} else if value.expiresAt < time.Now().UnixMilli() {
			return Encode(int64(-2), "number")
		} else {
			remainingSeconds := (value.expiresAt - time.Now().UnixMilli()) / 1000
			log.Printf("%d remaining Seconds\n", remainingSeconds)
			fmt.Printf("type: %T\n", remainingSeconds)
			return Encode(remainingSeconds, "number")
		}
	}
}

func evalDel(args []string) []byte {
	n := len(args)

	if n < 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'del' command"), "simpleString")
	} else {
		res := DELETE(args)

		return Encode(res, "number")
	}
}

func evalExpire(args []string) []byte {
	n := len(args)

	if n < 2 || n > 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'expire' command"), "simpleString")
	} else {
		expireSeconds, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return Encode(errors.New("ERR value is not an integer or out of range"), "simpleString")
		}
		res := Expire(args[0], expireSeconds*int64(1000))
		return Encode(res, "number")
	}
}

func evalAOF(args []string) []byte {
	if len(args) > 0 {
		return Encode(errors.New("ERR wrong number of arguments for 'Bbgrewriteaof' command"), "simpleString")
	}

	err := BgWrite(storage)
	if err != nil {
		return Encode(string(err.Error()), "simpleString")
	} else {
		return Encode("OK", "simpleString")

	}

}

func evalIncr(args []string) []byte {

	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'incr' command"), "simpleString")
	}

	key := args[0]
	obj := GET(args[0])

	if obj == nil {
		obj = CreateObj("0", -1, OBJ_TYPE_STRING, OBJ_ENCODING_INT)
		PUT(key, obj)
	} else {
		if err := chckType(obj.TypeEncoding, OBJ_TYPE_STRING); err != nil {
			return Encode(err, "simpleString")
		}

		if err := checkEncoding(obj.TypeEncoding, OBJ_ENCODING_INT); err != nil {
			return Encode(err, "simpleString")
		}

	}
	count, _ := strconv.ParseInt(obj.val.(string), 10, 64)
	count++

	obj.val = strconv.FormatInt(count, 10)
	return Encode(count, "number")
}

func evalInfo(args []string) []byte {

	var buff = bytes.NewBuffer([]byte{})

	buff.Write(Encode("# Keyspace", "bulkString"))
	for i := range KeyspaceStat {
		buff.WriteString(fmt.Sprintf("db%d:keys=%d,expires=0,avg_ttl=0\r\n", i, KeyspaceStat[i]["keys"]))
	}
	return (Encode(buff.String(),"bulkString"))

}
func EvalAndRespond(cmd RedisCmds, conn io.ReadWriter) {
	b := []byte{}
	buff := bytes.NewBuffer(b)
	for _, val := range cmd {
		command := val.Cmd
		args := val.Args

		switch command {
		case "ping":
			buff.Write(evalPing(args))
		case "set":
			// log.Println("Set Command Found")
			buff.Write(evalSet(args))
		case "get":
			buff.Write(evalGet(args))
		case "ttl":
			buff.Write(evalTTL(args))
		case "del":
			buff.Write(evalDel(args))
		case "expire":
			buff.Write(evalExpire(args))
		case "bgrewriteaof":
			buff.Write(evalAOF(args))
		case "incr":
			buff.Write(evalIncr(args))
		case "info":
			buff.Write(evalInfo(args))
		default:
			errMsg := fmt.Sprintf("+(error) ERR unknown command '%s', with args beginning with:\r\n", command)
			buff.Write([]byte(errMsg))
		}
	}

	conn.Write(buff.Bytes())
}

func ErrorResponse(err error, conn io.ReadWriter) {
	conn.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}
