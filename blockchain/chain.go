package blockchain

import (
	"sync"
)

type blockchain struct {
	NewstHash string `json:"newestHash"`
	Height    int    `json:"height"`
}

var b *blockchain // 이건 blockchain package 안에서만 사용할 수 있음!
var once sync.Once

func (b *blockchain) AddBlock(data string) {
	block := Block{data, "", b.NewstHash, b.Height + 1}
	b.NewstHash = block.Hash
	b.Height = block.Height
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			b.AddBlock("Genesis")
		})

	}
	return b
}
