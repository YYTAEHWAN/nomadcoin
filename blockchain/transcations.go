package blockchain

import (
	"learngo/github.com/nomadcoders/utils"
	"time"
)

const (
	minerReward int = 50
)

type Tx struct { // Tx = transaction
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)

}

type TxIn struct { // TxIn = transaction Input
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type TxOut struct { // TxIn = transaction Output
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

func makeCoinbaseTx(address string) *Tx {
	// 채굴자를 주소로 삼는 코인베이스 거래내역을 생성해서 Tx 포인터를 리턴할 것이다
	txIns := []*TxIn{
		{"COINBASE", minerReward},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}
