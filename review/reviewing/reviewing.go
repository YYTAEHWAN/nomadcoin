package reviewing

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	Data     string
	Hash     string
	Prevhash string
}

type blockchain struct {
	blocks []*block
	//4-4까지 있던 원래 코드는 blocks []*block 이었음, blockchain을 주소로 줄 수는 있지만
	// 특정한 함수를 만들기 전까지는 (b *blockchain) 만 가지고는 blocks 에 접근할 수 없었지
	// 하지만 Blocks []*block 로 바꾸면서 (b *blockchain)이라는 receiver or GetBlockchain() 같이 *blockchain을 받을 수만 있다면
	// ex) chain := GetBlockchain()
	//		chain.Blocks 을 통해서 []*blocks 를 보고 수정할 수 있는거지
	// 하지만 이걸 바꿀수 밖에 없어 왜나면 나중엔 데이터베이스로 blocks 를 만들 거고 함수로 데이터베이스를 가져와야 하기 때문이지
}

var b *blockchain  // 싱글턴 패턴을 사용하기 위한 블록체인 변수
var once sync.Once // 싱글턴 패턴을 사용하기 위한 전역변수

func (b *block) calculateHash() {
	// 블럭 복사본이 아닌 주소를 받아오기 때문에 값을 변경할 수 있다
	hash := sha256.Sum256([]byte(b.Data + b.Prevhash))
	b.Hash = fmt.Sprintf("%x", hash)
}
func getLastHash() string {
	// GetBlockchain() 함수로 블록체인 구조체 주소를 받아서 blocks를 사용한다
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].Hash
}
func createBlock(data string) *block {
	// 데이터와 이전 블럭의 해쉬값을 가지고있는 새로운 블럭을 생성한다
	newBlock := block{data, "", getLastHash()}
	newBlock.calculateHash() // 데이터와 이전 해쉬값을 가지고 해쉬를 생성한다
	return &newBlock         // 모두 채워진 새로운 블럭 반환
}
func (b *blockchain) AddBlock(data string) {
	// 4-5 에서 만든 함수인데 이렇게 만들면,, 1. 만든 블럭을 더해줄 때도 이 함수를 쓸 수 있고(재활용 가능)
	//2.  GetBlockchain 함수 밖에서도 사용 가능해
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockchain() *blockchain {
	// 이것이 싱글턴 패턴
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis Block")
			// 싱글턴 패턴을 실행할 때 초기 블럭을 만들어 넣는다
		})
	}
	return b
}

func (b *blockchain) AllBlock() []*block {
	return b.blocks
}
