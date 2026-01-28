package core

import (
	"errors"
	"strconv"
)

func EleType(t uint8) uint8 {
	return ((t >> 4) << 4)
}

func EleEncode(t uint8) uint8 {
	return ((t << 4) >> 4)
}

func chckType(qt uint8, ot uint8) error {
	if EleType(qt) != ot {
		return errors.New("the operation is not permmited on this type")
	}

	return nil
}


func checkEncoding(qt uint8, ot uint8) error {
	if EleEncode(qt) != ot {
		return errors.New("the operation is not permmited on this type")
	}

	return nil
}


func deduceTypeEncoding(value string)(uint8, uint8){
	dType := OBJ_TYPE_STRING

	if _,ok := strconv.ParseInt(value,10,64); ok == nil{
		return dType,OBJ_ENCODING_INT
	}

	if(len(value) < 44){
		return dType,OBJ_ENCODING_EMDEB
	}

	return OBJ_TYPE_STRING, OBJ_ENCODING_RAW
}
