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
	TxID  string
	Index int
	Owner string `json:"owner"`
}

type TxOut struct { // TxOut = transaction Output2
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

func makeCoinbaseTx(address string) *Tx {
	// 채굴자를 주소로 삼는 코인베이스 거래내역을 생성해서 Tx 포인터를 리턴할 것이다
	txIns := []*TxIn{
		{"", -1, "COINBASE"}, // coinbase transaction이기 때문에
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
	if TotalBalanceByAddress(from, Blockchain()) < amount {
		return nil, errors.New("not enought 돈")
	}

	var TxIns []*TxIn
	var TxOuts []*TxOut
	uTxOuts := UTxOutsByAddress(from, Blockchain())
	total := 0

	for _, uTxOut := range uTxOuts { // 사용될 TxOut들을 TxIn으로 옮기는 과정  //// 근데 이거  복사야 아니면 참조에 의한 변경이야?
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from} // 이 from은 하드코딩이므로 나중에 바꿔버릴꺼야~
		TxIns = append(TxIns, txIn)                    // 의문 첫번째 - 위엣줄 코딩을 하게 되면 저 기록들이 옮겨져? 아니면 복사만 되어서 들어가는거야?
		total += uTxOut.Amount
	}
	if change := total - amount; change > 0 { // 거스름돈 받는 코드
		changeTxOut := &TxOut{from, change}
		TxOuts = append(TxOuts, changeTxOut)
	}
	TotxOut := &TxOut{to, amount}
	TxOuts = append(TxOuts, TotxOut)
	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     TxIns,
		TxOuts:    TxOuts,
	}
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

func IsOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			// mempool에 들어있는 inputs들은 모두 사용될 예정인 TxOut들이지 그러니 input과 같은 ID, Index가 있다면 삐삑 넌 걸렸어 UTxOuts로 다시 못들어가
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
}
