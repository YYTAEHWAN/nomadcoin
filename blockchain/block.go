package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nomadcoders/utils"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"` // 한 번만 사용되는 숫자
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

func persistBlock(b *Block) {
	dbStorage.SaveBlock(b.Hash, utils.ToBytes(b))
	// 16진수로 된 해쉬 string이랑, b블록을 []byte로 변환한 값을 넣어준다
}

var ErrNotFound = errors.New("block not found")

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := dbStorage.GetBlockHashFromDb(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)

	return block, nil
}

func (b *Block) mine() {
	defer fmt.Println("채굴 완료!")

	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		//fmt.Printf("\nHash:%s\nTarget:%s\nNonce:%d\n\n", hash, target, b.Nonce)
		if strings.HasPrefix(hash, target) { // hash값과 target값이 같다면
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(prevHash string, height, diff int) *Block {
	block := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: diff,
		Nonce:      0,
		// Transactions: []*Tx{makeCoinbaseTx("taehwan")},
	}
	block.mine() // 조건에 맞는 해쉬값을 찾고,, 찾으면 그 값을 블록의 hash값으로 설정
	block.Transactions = Mempool().TxToComfirm()
	persistBlock(block)
	return block
	// dho dksehlfRKdy rltgjqmTL??ggg djfuqspdy~
}
