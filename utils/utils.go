// Package utils contains functions to be used across the apllication.

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

var logFn = log.Panic

// HandleErr handle a error
func HandleErr(err error) {
	if err != nil {
		logFn(err)
	}
}

// ToBytes takes an interface and then will return []byte made by using gob Encoder
func ToBytes(i interface{}) []byte { // gob package 사용
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()
}

// FromBytes takes an interface and data  and then will encode the data to the interface
func FromBytes(i interface{}, data []byte) { // 인자는 포인터와 복원할 데이터
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(decoder.Decode(i))
	// 포인터로 접근해서 값 자체를 변경해주었기 때문에
	// 이렇게 간단하게 코드가 마무리 되었다
}

// Hash takes an interface and returns hex string made by using sha256
func Hash(i interface{}) string { // 인터페이스라는 것은 모든 형태의 변수를 뜻함
	s := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

// Spliter takes string that is going to split, delimiter, index and then will return string that is the i-th element of the array seperated from string received as parameter
func Spliter(arr, deli string, i int) string {
	// i는 나눠진 배열에서 원하는 데이터의 인덱스
	r := strings.Split(arr, deli)
	if len(r)-1 < i {
		return ""
	}
	return r[i]
}

// ToJSON takes an interface and then returns []bytes that is encoded by method of JSON
func ToJSON(i interface{}) []byte {
	b, err := json.Marshal(i)
	HandleErr(err)
	return b
}
