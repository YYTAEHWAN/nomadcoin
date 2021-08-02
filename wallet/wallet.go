package wallet

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"learngo/github.com/nomadcoders/utils"
	"math/big"
)

const (
	signature     string = "4f7f14ef3526228ea9fa7d6c2d9e6cbd587e37bd29fd2f86afa2a64dd788705a6e60785a358e61f869b5550726fe2c4917004e055ced4d700e5267b19d25804e"
	privateKey    string = "307702010104204e0a38a4d6104099c7c3e109b7c6c6758564cd613e00159544da5c78964d36a2a00a06082a8648ce3d030107a1440342000426377d2c4f912474e56ea1e467a9d2de9e0854fbb8d97386757716f463999322c36d28f9782bd052f21c7bf8391139562322e9e42c9f415e284cdb46d03c54ee"
	hashedMessage string = "1c5863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
	// 실제로는 Tx가 hash된 메세지일거야
	pri string = "100371616754094915234634347704562684656090136652171610345641000101660048327935"
)

func Start() {
	privBites, err := hex.DecodeString(privateKey)
	//hex.DecodeString는 16진수 string을 받아서 []byte형태로 바꿔준다
	// 그러므로 해당 데이터의 포맷이 16진수 string인지도 판별할 수 있다
	utils.HandleErr(err)

	restoreKey, err := x509.ParseECPrivateKey(privBites)
	utils.HandleErr(err)

	//fmt.Println(restoreKey)
	//fmt.Printf("\n\n")

	sigBytes, err := hex.DecodeString(signature)
	utils.HandleErr(err)

	rBytes := sigBytes[:len(sigBytes)/2] //처음부터 반까지
	sBytes := sigBytes[len(sigBytes)/2:] //반부터 끝까지

	var bigS, bigR = big.Int{}, big.Int{} //이 문장은 왜 이렇게 되는지 잘 모르겠다 big.Int라는 타입이 어디 패키지에 있겠찌 math/big 이네 그리고 := 를 안써주고 = 를 써준 이유는 var 를 썼기 때문
	// 밑에 SetByte() 초기화 함수를 쓰려면 변수가 초기화가 되어있어야 하나봐
	bigR.SetBytes(rBytes) // 리시버 함수이기 때문에 값을 바꿔줄 수 있으므로 대입 필요 없음
	bigS.SetBytes(sBytes)

	//fmt.Println(bigR, bigS)

	hashAsBytes, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)

	ok := ecdsa.Verify(&restoreKey.PublicKey, hashAsBytes, &bigR, &bigS)

	fmt.Println(ok)

}
