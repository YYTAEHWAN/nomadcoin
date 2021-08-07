package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"learngo/github.com/nomadcoders/utils"
	"math/big"
	"os"
)

const (
	fileName string = "nomadcoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat("nomadcoin.wallet")
	return !os.IsNotExist(err)
}

func createKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key) // key의 parsing 과 marshaling을 책임지는 x509를 사용
	//privateKey를 []byte로 만들어주고
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, 0644) // 이 함수는 파일을 열어주고, 데이터를 써주고, 파일을 자동으로 닫아준다 0644는 umask 번호
	utils.HandleErr(err)

}

func restoreKey() (key *ecdsa.PrivateKey) { // naked return 을 보여주기 위해 일부러 사용해봄
	keyAsBytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)

	return
}

func encodeBigInts(a, b []byte) string {
	slice := append(a, b...)
	return fmt.Sprintf("%x", slice)
}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func Sign(payload string, w wallet) string { // 여기서 w wallet을 포인터로 줄 필요는 없음 복사해서 값만 넣으면 되기 때문
	payloadAsB, err := hex.DecodeString(payload) // 또한 이 과정을 안거치고 바로 []byte()함수를 이용해 바꿔줄 수 도 있지만
	utils.HandleErr(err)                         // 16진수인 string인지를 확인하여 정확성을 더하기 위해선 hex.DecodeString()함수를 쓰는게 맞음
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsB)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(payload string) (*big.Int, *big.Int) { // 컴퓨팅에 있어서 페이로드는 전송되는 데이터를 뜻한다
	Bytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	firstB := Bytes[:len(Bytes)/2]
	secondB := Bytes[len(Bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstB)
	bigB.SetBytes(secondB)

	return &bigA, &bigB
}

func Verify(signature string, payload string, address string) bool {
	r, s := restoreBigInts(signature)
	x, y := restoreBigInts(address)
	PublicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	ok := ecdsa.Verify(&PublicKey, payloadAsB, r, s)

	return ok
}

func Wallet() *wallet {
	if w == nil {
		// has a wallet?
		w = &wallet{}
		if hasWalletFile() {
			w.privateKey = restoreKey()
			// yes -> restore the wallet
		} else {
			key := createKey()
			persistKey(key)    // 키를 DB에 저장
			w.privateKey = key // privateKey를 w에 저장
			// no -> create the wallet
		}
		w.Address = aFromK(w.privateKey)
	}
	//fmt.Println(w.Address)
	return w
}
