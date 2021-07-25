package blockchain

import (
	"errors"
	"fmt"
	"learngo/github.com/nomadcoders/db"
	"learngo/github.com/nomadcoders/utils"
	"strings"
	"time"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"` // 한 번만 사용되는 숫자
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
	// 16진수로 된 해쉬 string이랑, b블록을 []byte로 변환한 값을 넣어준다
}

var ErrNotFound = errors.New("block not found")

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.GetBlockHashFromDb(hash)
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
		fmt.Printf("\nHash:%s\nTarget:%s\nNonce:%d\n\n", hash, target, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(prevHash string, height int) *Block {
	block := &Block{
		Hash:         "",
		PrevHash:     prevHash,
		Height:       height,
		Difficulty:   Blockchain().difficulty(),
		Nonce:        0,
		Transactions: []*Tx{makeCoinbaseTx("taehwan")},
	}
	block.mine()
	block.persist()
	return block
	// dho dksehlfRKdy rltgjqmTL??ggg djfuqspdy~
}
