package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"reflect"
	"testing"
)

const (
	testKey     string = "30770201010420bd2f27a160555c68dbe54f0293bbcf1ea7880e3546d00746ce28680a8b1672f4a00a06082a8648ce3d030107a14403420004dfa3dc9ad0b5657c8066f3fea5e952dd3e7f73b8cde2e438cf49c3e5738b1c6b9cfdb3f8d6e64f575ee8e7af8f9e918c2a9f32bb1c15dddda232f47d5217c91e"
	testPayload string = "a1c1fdd4b429ae5e4ee77cccdd308b59b5947d90fabe4e4a4ad04b5ea2d2f360"
	testSig     string = "e15ce853dd9ad3b7f812dfd0597226fb0eededd84a6c44d8f3638b6d321e81ebe42e2a77df171921a7559297c62e32172e6a08d31833d12a61fe623618286cb8"
)

type fakelayer struct {
	fakeHashWalletFile func() bool // 이게 무슨 형태의 함수일까...
}

func (f fakelayer) hasWalletFile() bool {
	return f.fakeHashWalletFile()
}

func (fakelayer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}

func (fakelayer) readFile(name string) ([]byte, error) {
	// 리턴값만 같도록 조작해준다
	return x509.MarshalECPrivateKey(makeTestwallet().privateKey)
}

func makeTestwallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)
	return w
}

func TestSign(t *testing.T) {
	signature := Sign(testPayload, *makeTestwallet())
	_, err := hex.DecodeString(signature)
	if err != nil {
		t.Errorf("Sign() should return a hex encoded string , but we got %s", signature)
	}
}

func TestVerify(t *testing.T) { // verify는 true 말고 false일 때도 체크해주어야한다 그러기 위해선 테이블테스트 만들기
	type testable struct {
		input string
		ok    bool
	}
	testT := []testable{
		{input: testPayload, ok: true},
		{input: "a3c1fdd4b429ae5e4ee77cccdd308b59b5947d90fabe4e4a4ad04b5ea2d2f360", ok: false},
	}
	w = makeTestwallet()
	for _, test := range testT {
		ok := Verify(testSig, test.input, w.Address)
		if ok != test.ok {
			t.Error("Verify() could not verify testSig and testPayload")
		}
	}
}

func TestWallet(t *testing.T) {
	// 이게 도데체 어떤 논리로 흘러가는 코드일까
	w = nil
	t.Run("New wallet is created", func(t *testing.T) {
		files = fakelayer{
			fakeHashWalletFile: func() bool {
				t.Log("i have benn called")
				return false
			},
		}
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("new Wallet should return a new wallet instance")
		}
	})
	t.Run("Wallet is restored", func(t *testing.T) {
		files = fakelayer{
			fakeHashWalletFile: func() bool {
				t.Log("i have benn called")
				return true
			},
		}
		w = nil
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("new Wallet should return a new wallet instance")
		}
	})
}
