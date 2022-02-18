package elizeutils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

func Hash(anything interface{}) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v", anything))))
}

func Splitter(s, sep string, i int) string {
	r := strings.Split(s, sep)
	if len(r) < i {
		return ""
	}
	return r[i]
}

func ToJSON(i interface{}) []byte {
	b, err := json.Marshal(i)
	Errchk(err)
	return b
}
