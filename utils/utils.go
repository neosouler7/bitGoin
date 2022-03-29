package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr((encoder.Encode(i)))
	return aBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(decoder.Decode(i))
}

func Hash(i interface{}) string {
	hash := sha256.Sum256([]byte(fmt.Sprint(i)))
	return fmt.Sprintf("%x", hash)
}
