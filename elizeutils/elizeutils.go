package elizeutils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func Errchk(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func ToBytes(data interface{}) []byte {
	var blockBuffer bytes.Buffer
	err := gob.NewEncoder(&blockBuffer).Encode(data)
	Errchk(err)
	return blockBuffer.Bytes()
}

func FromBytes(emptyStruct interface{}, data []byte) {
	Errchk(gob.NewDecoder(bytes.NewReader(data)).Decode(emptyStruct))
}
