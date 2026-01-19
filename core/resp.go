package core

import (
	"errors"
	"fmt"
	// "log"
)

func readLength(data []byte) (int, int) {
	val := 0
	pos := 1 // skip prefix: $, *

	for ; data[pos] != '\r'; pos++ {
		val = val*10 + int(data[pos]-'0')
	}

	return val, pos + 2
}

func decodeSimpleString(data []byte) (interface{}, error, int) {
	pos := 1
	for ; data[pos] != '\r'; pos++ {

	}

	return string(data[1:pos]), nil, pos + 2
}

func decodeError(data []byte) (interface{}, error, int) {
	return decodeSimpleString(data)
}

func decodeBulkString(data []byte) (interface{}, error, int) {
	length, delta := readLength(data)
	pos := delta

	return string(data[pos : pos+length]), nil, pos + length + 2
}

func decodeInt(data []byte) (interface{}, error, int) {
	var num int64 = 0
	pos := 1

	for ; data[pos] != '\r'; pos++ {
		num = num*10 + int64(data[pos]-'0')
	}

	return num, nil, pos + 2
}

func decodeArrays(data []byte) ([]interface{}, error, int) {
	len, delta := readLength(data)

	var ele []interface{} = make([]interface{}, len)

	pos := delta

	for i := range ele {
		val, err, delta := Parse(data[pos:])

		if err != nil {
			return nil, err, 0
		}

		ele[i] = val

		pos += delta
	}

	// var str []string

	return ele, nil, pos
}

func DecodeArrayString(data []byte) ([]string, error) {
	res, err := Decode(data)

	if err != nil {
		return []string{}, nil
	}

	val := res.([]interface{})
	var arr = make([]string, len(val))

	for i := range val {
		arr[i] = val[i].(string)
	}

	return arr, nil
}
func Parse(data []byte) (interface{}, error, int) {
	switch data[0] {
	case '-':
		return decodeError(data)
	case ':':
		return decodeInt(data)

	case '*':
		return decodeArrays(data)
	case '$':
		return decodeBulkString(data)

	case '+':
		return decodeSimpleString(data)
	}

	return nil, nil, 0

}

func Decode(data []byte) (interface{}, error) {

	// log.Println("Recived Command " , string(data))
	if len(data) <= 0 {
		return nil, errors.New("No Data")
	}

	val, err, _ := Parse(data)

	return val, err
}

func Encode(val interface{}, isSimple bool) []byte {

	switch v := val.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n",v))
		}

		return []byte(fmt.Sprintf("$%d\r\n%s\r\n",len(v),v))
	}

	return []byte{}

}
