package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte { // gob package 사용
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) { // 인자는 포인터와 복원할 데이터
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(decoder.Decode(i))
	// 포인터로 접근해서 값 자체를 변경해주었기 때문에
	// 이렇게 간단하게 코드가 마무리 되었다
}

func Hash(i interface{}) string { // 인터페이스라는 것은 모든 형태의 변수를 뜻함
	s := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func Spliter(arr, deli string, i int) string {
	r := strings.Split(arr, deli)
	if len(r)-1 < i {
		return ""
	}
	return r[i]
}

func ToJSON(i interface{}) []byte {
	b, err := json.Marshal(i)
	HandleErr(err)
	return b
}
