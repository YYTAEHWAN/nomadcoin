package blockchain

import (
	"errors"
	"learngo/github.com/nomadcoders/utils"
	"time"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct { // Tx = transaction
	Id        string   `json:"id"` // tx의 해시값?
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

type TxOut struct { // TxOut = transaction Output2
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

func makeTx(from string, to string, amount int) (*Tx, error) {
	if Blockchain().TotalBalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	oldTxOuts := Blockchain().TxOutsByAddress(from)
	for _, txOut := range oldTxOuts {
		txIn := &TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, txIn)
		total += txOut.Amount
		if total > amount {
			break
		}
	}
	changeTxOut := total - amount
	if changeTxOut != 0 {
		txOut := &TxOut{from, changeTxOut}
		txOuts = append(txOuts, txOut)
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		Id:        to,
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	//transaction을 mempool에 추가해줄 뿐, 거래를 만들진 않는다
	tx, err := makeTx("taehwan", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToComfirm() []*Tx {
	coinBase := makeCoinbaseTx("taehwan")
	txs := m.Txs
	txs = append(txs, coinBase)
	m.Txs = nil // 밈풀 초기화
	return txs
}
