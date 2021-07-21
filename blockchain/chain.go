package blockchain

import (
	"learngo/github.com/nomadcoders/db"
	"learngo/github.com/nomadcoders/utils"
	"sync"
)

type blockchain struct {
	NewstHash string `json:"newestHash"`
	Height    int    `json:"height"`
}

var b *blockchain // 이건 blockchain package 안에서만 사용할 수 있음!
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}
func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}
func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewstHash, b.Height+1)
	b.NewstHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) BlocksSlice() []*Block {
	var Blocks []*Block
	hashCursor := b.NewstHash

	for {
		block, _ := FindBlock(hashCursor) // FindBlock이 *Block, err를 리턴하기 때문에
		Blocks = append(Blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return Blocks
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			checkpoint := db.GetCheckPointFromDb()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				b.restore(checkpoint)
			}
		})
	}
	return b
}
