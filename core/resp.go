package core

import (
	// "errors"
	"fmt"
	"log"
	"strings"
	// "log"
)

func readLength(data []byte) (int, int) {
	val := 0
	pos := 1

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
	arrLen, delta := readLength(data)
	log.Printf("Number of Character %d and curr posi is %d\n", arrLen, delta)
	pos := delta

	return string(data[pos : pos+arrLen]), nil, pos + arrLen + 2
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

	arrLen, delta := readLength(data)
	
	
	var ele []interface{} = make([]interface{}, arrLen)
	
	// log.Printf("Number of Elements %d and curr posi is %d\n", arrLen, delta)
	// log.Printf("Last Char %q ", string(data[delta]))

	pos := delta
	for i := range arrLen{
		val , err ,delta := Parse(data[pos:])
		if(err != nil){
			return nil,err,0
		}

		ele[i] = val
		pos += delta
	}

	return ele , nil , pos
}

func DecodeArrayStrings(data []byte) ([][]string, error) {

	var input [][]string
	n := len(data)

	i := 0
	for i < n {
		res, idx, err := Decode(data[i:])
		// log.Printf("Result of %d ele %v", i, res)
		i += idx
		if err != nil {
			return nil, err
		}

		if it, ok := res.([]interface{}); ok {
			var arr = make([]string, len(it))
			for j := range it {
				if s, ok := it[j].(string); ok {
					arr[j] = strings.ToLower(s)
				} else {
					// log.Println("Not String")
					log.Println(it[j])
				}
			}
			input = append(input, arr)
		} else {
			// log.Println("Error Ocurred")
			return nil, err
		}
	}
	return input, nil
}
func Parse(data []byte) (interface{}, error, int) {
	log.Printf("Parse called with data: %q", string(data))

	if data[0] == '\r' || data[0] == '\n' {
        val, err, delta := Parse(data[1:])
        return val, err, delta + 1
    }
	switch data[0] {
	case '-':
		log.Println("Decoding Error type")
		return decodeError(data)
	case ':':
		log.Println("Decoding Integer type")
		return decodeInt(data)
	case '*':
		log.Println("Decoding Array type")
		return decodeArrays(data)
	case '$':
		log.Println("Decoding Bulk String type")
		return decodeBulkString(data)
	case '+':
		log.Println("Decoding Simple String type")
		return decodeSimpleString(data)
	}

	log.Println("Unknown type encountered")
	return nil, nil, 0
}

func Decode(data []byte) (interface{}, int, error) {

	// log.Println("Recived Command " , string(data))
	// if len(data) <= 0 {
	// 	return []byte{}, -1, errors.New("No Data")
	// }

	val, err, idx := Parse(data)
	return val, idx, err
}

func Encode(val interface{}, resType string) []byte {

	switch resType {
	case "simpleString":
		return []byte(fmt.Sprintf("+%s\r\n", val))
	case "bulkString":
		if s, ok := val.(string); ok {
			return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
		}
		return []byte{}
	case "number":
		if n, ok := val.(int64); ok {
			log.Println("Found Number")
			return []byte(fmt.Sprintf(":%d\r\n", n))
		}
		return []byte{}
	default:
		return []byte{}
	}

}
