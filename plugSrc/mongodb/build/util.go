package build

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io"
	"time"
)

func GetNowStr(isClient bool) string {
	var msg string
	if isClient {
		msg = time.Now().Format("2006-01-02 15:04:05")+"| cli -> ser |"
	}else{
		msg = time.Now().Format("2006-01-02 15:04:05")+"| ser -> cli |"
	}
	return msg
}

func ReadInt32(r io.Reader) (n int32) {
	binary.Read(r, binary.LittleEndian, &n)
	return
}

func ReadInt64(r io.Reader) int64 {
	var n int64
	binary.Read(r, binary.LittleEndian, &n)
	return n
}

func ReadString(r io.Reader) string {

	var result []byte
	var b = make([]byte, 1)
	for {

		_, err := r.Read(b)

		if err != nil {
			panic(err)
		}

		if b[0] == '\x00' {
			break
		}

		result = append(result, b[0])
	}

	return string(result)
}

func ReadBson2Json(r io.Reader) (string) {

	//read len
	docLen := ReadInt32(r)
	if docLen == 0 {
		return ""
	}

	//document []byte
	docBytes := make([]byte, int(docLen))
	binary.LittleEndian.PutUint32(docBytes, uint32(docLen))
	if _, err := io.ReadFull(r, docBytes[4:]); err != nil {
		panic(err)
	}

	//resolve document
	var bsn bson.M
	err := bson.Unmarshal(docBytes, &bsn)
	if err != nil {
		panic(err)
	}

	//format to Json
	jsonStr, err := json.Marshal(bsn)
	if err != nil {
		return fmt.Sprintf("{\"error\":%s}", err.Error())
	}
	return string(jsonStr)
}

