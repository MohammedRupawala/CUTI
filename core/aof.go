package core

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)


func dump(key string, val *value, f *os.File){
	if(expires[val] != 0){
		ttl := int64(expires[val]) - time.Now().UnixMilli()

		if(ttl > 0){
			cmd := fmt.Sprintf("SET %s %s EX %d", key,val.val,ttl)
			tokens := strings.Split(cmd," ")
			f.Write(Encode(tokens,"arrayString"))
		}
	}else{
		cmd := fmt.Sprintf("SET %s %s", key,val.val)
		tokens := strings.Split(cmd," ")
		f.Write(Encode(tokens,"arrayString"))
	}
}
func BgWrite(data map[string] *value) error {
	tmpFile, err := os.CreateTemp(".", "temp-rewrite.aof")
	if err != nil {
		log.Println("Error Ocurred While Creating new File ", err)
		return err
	}

	defer tmpFile.Close()

	for k, v := range data {
		dump(k,v,tmpFile)
	}

	tmpFile.Sync()

	if err := os.Rename(tmpFile.Name(), "database.aof"); err != nil{
		log.Println("Rewrite Error: Rename failed", err)
        return err
	}


	log.Println("AOF File Written Successfully")
	return nil
}
